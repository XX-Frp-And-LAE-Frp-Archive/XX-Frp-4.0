package register

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"math/big"
	mailrand "math/rand"
	"strconv"
	"time"
)

func HandleEmail(c *gin.Context, db *gorm.DB) {
	// 获取表单数据
	email := c.PostForm("email")

	// 使用isValidEmail 函数检测邮箱
	if !IsValidEmail(email) {
		respond.Respond(c, 400, "邮箱不合法", 0)
		return
	}
	// 查询数据库中是否已经存在该邮箱
	if checkEmail(email, db) {
		respond.Respond(c, 400, "邮箱已经注册，请直接登录", 0)
		return
	}
	// 生成6位随机验证码
	code := generateRandomCode()
	// 发送邮件
	err := sendEmail(email, code)
	if err != nil {
		respond.Respond(c, 500, "邮件发送错误，请联系管理员", 0)
		fmt.Println(err)
		return
	}
	// 查询数据库中是否存在指定邮箱的记录
	var existingRecord define.Code

	result := db.Where("email = ?", email).First(&existingRecord)
	if result.Error == nil {
		// 更新已存在的记录
		existingRecord.Code = code
		existingRecord.Time = time.Now().Unix()
		err := db.Model(&existingRecord).Where("email = ?", email).Updates(define.Code{Code: existingRecord.Code, Time: existingRecord.Time}).Error
		if err != nil {
			respond.Respond(c, 500, "验证码暂存错误，已发送的邮件请忽略，请联系管理员或重试！", 0)
			return
		}
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 创建新的记录
		record := define.Code{
			Email: email,
			Code:  code,
			Time:  time.Now().Unix(),
		}
		err := db.Create(&record).Error
		if err != nil {
			respond.Respond(c, 500, "验证码暂存错误，已发送的邮件请忽略，请联系管理员或重试！", 0)
			return
		}
	} else {
		respond.Respond(c, 500, "验证码暂存错误，已发送的邮件请忽略，请联系管理员或重试！", 0)
		return
	}
	respond.Respond(c, 200, "邮件发送成功，15分钟内有效，邮件可能被误判为垃圾邮件，请检查垃圾箱！", 0)
}

func generateRandomCode() string {
	// 生成一个随机六位数
	var max int64 = 999999
	var min int64 = 100000
	randNum, _ := rand.Int(rand.Reader, big.NewInt(max-min))
	return strconv.FormatInt(randNum.Int64()+min, 10)
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
	messageId := fmt.Sprintf("<%d.%d@mefbi.com>", timestamp, randomNum)
	return messageId
}
