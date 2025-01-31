package UserHandler

import (
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	register "github.com/ahmr-bot/ME-Frp/pkg/router/RegisterHandler"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

func HandleForgotPassword(c *gin.Context, db *gorm.DB) {
	// 获取表单数据中的邮箱
	email := c.PostForm("email")
	// 获取表单数据中的用户名
	username := c.PostForm("username")
	// 检查 邮箱 用户名 是否合法
	if register.IsValidEmail(email) == false || register.IsValidUsername(username) == false {
		respond.Respond(c, 403, "邮箱或用户名不合法!", 0)
		return
	}
	// 检查邮箱是否存在
	var user define.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		respond.Respond(c, 403, "该邮箱未注册!", 0)
		return
	}
	// 检测用户名和邮箱是否匹配
	if user.Username != username {
		respond.Respond(c, 403, "用户名和密码不匹配!", 0)
		return
	}
	// 检测15分钟内是否已经发送过邮件
	var findpass define.Findpass
	if err := db.Table("findpass").Where("username = ?", username).First(&findpass).Error; err == nil {
		// 如果已经发送过邮件
		if time.Now().Unix()-findpass.Time < 900 {
			respond.Respond(c, 403, "15分钟内已发送过邮件!", 0)
			return
		}
		// 如果已经发送过邮件但是已经超过15分钟
		if time.Now().Unix()-findpass.Time > 900 {
			// 删除findpass表中的token
			if err := db.Table("findpass").Delete(&findpass).Error; err != nil {
				respond.Respond(c, 500, "删除旧数据失败!", 0)
				return
			}
		}
	}
	// 生成重置密码的token
	token := register.GenerateToken()
	// 将username token time存入数据库 findpass 表中 而不是 findpasses 表中
	db.Table("findpass").Create(&define.Findpass{
		Username: username,
		Link:     token,
		Time:     time.Now().Unix(),
	})
	// 发送邮件
	err := SendEmail(email, token)
	if err != nil {
		respond.Respond(c, 500, "发送错误!", 0)
		return
	}
	respond.Respond(c, 200, "发送成功!", 0)
}

func SendEmail(email string, token string) error {
	conf := config.GetConfig()
	// 生成重置密码的链接
	passwordResetURL := conf.Server.Url + "/auth/reset/?token=" + token
	// 发送邮件
	mailer := gomail.NewDialer(conf.Smtp.Addr, conf.Smtp.Port, conf.Smtp.From, conf.Smtp.Passwd)
	message := gomail.NewMessage()
	message.SetHeader("From", conf.Smtp.From)
	message.SetHeader("To", email)
	message.SetHeader("Subject", conf.Server.Name+"重置密码")

	// 添加messageId头部信息
	messageId := GenerateMessageId()
	message.SetHeader("Message-ID", messageId)

	message.SetBody("text/plain", fmt.Sprintf("你的密码重置链接是: %s 15分钟内有效", passwordResetURL))
	return mailer.DialAndSend(message)
}

func GenerateMessageId() string {
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(1000)
	messageId := fmt.Sprintf("<%d.%d@mefbi.com>", timestamp, randomNum)
	return messageId
}

func HandleForgotResetPassword(c *gin.Context, db *gorm.DB) {
	// 获取动态路由中的token
	token := c.Param("link")
	// 获取表单数据中的新密码
	newPassword := c.PostForm("password")
	// 检查新密码是否合法
	//if register.IsValidPassword(newPassword) == false {
	//	c.JSON(400, gin.H{"error": "密码不合法"})
	//	return
	//}

	// 检查token是否存在 findpass 表中 并获取username
	var findpass define.Findpass
	if err := db.Table("findpass").Where("link = ?", token).First(&findpass).Error; err != nil {
		respond.Respond(c, 403, "Token不存在!", 0)
		return
	}
	// 检查token是否过期
	if findpass.Time+900 < time.Now().Unix() {
		respond.Respond(c, 403, "Token已过期!", 0)
		// 删除findpass表中的token
		if err := db.Table("findpass").Delete(&findpass).Error; err != nil {
			respond.Respond(c, 500, "Token删除失败!", 0)
			return
		}
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		respond.Respond(c, 500, "密码加密失败!", 0)
		return
	}
	// 更新对应用户的密码
	if err := db.Model(&define.User{}).Where("username = ?", findpass.Username).Update("password", string(hashedPassword)).Error; err != nil {
		respond.Respond(c, 500, "更新密码失败!", 0)
		return
	}
	// 删除findpass表中的token
	if err := db.Table("findpass").Delete(&findpass).Error; err != nil {
		respond.Respond(c, 500, "Token删除失败!", 0)
		return
	}
	// 更新 token
	NewToken := register.GenerateToken()
	// 在数据库users表中更新token
	db.Model(&define.User{}).Where("username = ?", findpass.Username).Update("token", NewToken)
	respond.Respond(c, 200, "密码重置成功，已为您自动重新登录!", gin.H{
		"token": NewToken,
	})
}
