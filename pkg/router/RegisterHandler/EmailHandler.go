package register

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"log"
	"math/big"
	mailrand "math/rand"
	"net/http"
	"strconv"
	"time"
)

type CodeData struct {
	Email string `json:"email"`
	Code  string `json:"code"`
	Time  string `json:"time"`
}

type Code struct {
	Email string
	Code  string
	Time  int64
}

func HandleEmail(c *gin.Context, db *gorm.DB) {
	// 获取表单数据
	email := c.PostForm("email")

	// 使用isValidEmail 函数检测邮箱
	if !IsValidEmail(email) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "邮箱格式错误",
		})
	}
	// 查询数据库中是否已经存在该邮箱
	if checkEmail(email, db) {
		c.JSON(200, gin.H{
			"message": "注册失败，邮箱已存在",
		})
		return
	}
	// 生成6位随机验证码
	code := generateRandomCode()
	// 发送邮件
	// 发送邮件
	err := sendEmail(email, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send email" + err.Error()})
		return
	}
	// 查询数据库中是否存在指定邮箱的记录
	var existingRecord Code

	result := db.Where("email = ?", email).First(&existingRecord)
	if result.Error == nil {
		// 更新已存在的记录
		existingRecord.Code = code
		existingRecord.Time = time.Now().Unix()
		err := db.Model(&existingRecord).Where("email = ?", email).Updates(Code{Code: existingRecord.Code, Time: existingRecord.Time}).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update record"})
			return
		}
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 创建新的记录
		record := Code{
			Email: email,
			Code:  code,
			Time:  time.Now().Unix(),
		}
		err := db.Create(&record).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save record"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to query database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Email sent successfully"})
}

func generateRandomCode() string {
	max := big.NewInt(999999)
	code, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
	}
	return strconv.Itoa(int(code.Int64()))
}

func sendEmail(email, code string) error {
	conf := config.GetConfig()
	mailer := gomail.NewDialer(conf.Smtp.Addr, conf.Smtp.Port, conf.Smtp.From, conf.Smtp.Passwd)
	message := gomail.NewMessage()
	message.SetHeader("From", conf.Smtp.From)
	message.SetHeader("To", email)
	message.SetHeader("Subject", conf.Server.Name+"注册验证码")
	// 添加messageId头部信息
	messageId := GenerateMessageId()
	message.SetHeader("Message-ID", messageId)

	message.SetBody("text/plain", fmt.Sprintf("你的验证码是: %s 15分钟内有效", code))
	return mailer.DialAndSend(message)
}

func GenerateMessageId() string {
	timestamp := time.Now().Unix()
	randomNum := mailrand.Intn(1000)
	messageId := fmt.Sprintf("<%d.%d@yourdomain.com>", timestamp, randomNum)
	return messageId
}
