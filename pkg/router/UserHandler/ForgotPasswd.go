package UserHandler

import (
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	register "github.com/ahmr-bot/ME-Frp/pkg/router/RegisterHandler"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"time"
)

type Findpass struct {
	Username string `gorm:"primaryKey"`
	Link     string
	Time     int64
}

func HandleForgotPassword(c *gin.Context, db *gorm.DB) {
	// 获取表单数据中的邮箱
	email := c.PostForm("email")
	// 获取表单数据中的用户名
	username := c.PostForm("username")
	// 检查 邮箱 用户名 是否合法
	if register.IsValidEmail(email) == false || register.IsValidUsername(username) == false {
		c.JSON(400, gin.H{"error": "邮箱或用户名不合法"})
		return
	}
	// 检查邮箱是否存在
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "此邮箱未注册"})
		return
	}
	// 检测用户名和邮箱是否匹配
	if user.Username != username {
		c.JSON(400, gin.H{"error": "用户名和邮箱不匹配"})
		return
	}
	// 检测15分钟内是否已经发送过邮件
	var findpass Findpass
	if err := db.Table("findpass").Where("username = ?", username).First(&findpass).Error; err == nil {
		// 如果已经发送过邮件
		if time.Now().Unix()-findpass.Time < 900 {
			c.JSON(400, gin.H{"error": "15分钟内已经发送过邮件"})
			return
		}
		// 如果已经发送过邮件但是已经超过15分钟
		if time.Now().Unix()-findpass.Time > 900 {
			// 删除findpass表中的token
			if err := db.Table("findpass").Delete(&findpass).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to delete findpass"})
				return
			}
		}
	}
	// 生成重置密码的token
	token := register.GenerateToken()
	// 将username token time存入数据库 findpass 表中 而不是 findpasses 表中
	db.Table("findpass").Create(&Findpass{
		Username: username,
		Link:     token,
		Time:     time.Now().Unix(),
	})
	// 发送邮件
	err := SendEmail(email, token)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send email"})
		return
	}
	c.JSON(200, gin.H{"message": "邮件已发送"})
}

func SendEmail(email string, token string) error {
	conf := config.GetConfig()
	// 生成重置密码的链接
	passwordResetURL := conf.Server.Url + "/auth/reset_password/" + token
	// 发送邮件
	mailer := gomail.NewDialer(conf.Smtp.Addr, conf.Smtp.Port, conf.Smtp.From, conf.Smtp.Passwd)
	message := gomail.NewMessage()
	message.SetHeader("From", conf.Smtp.From)
	message.SetHeader("To", email)
	message.SetHeader("Subject", conf.Server.Name+"重置密码")
	message.SetBody("text/plain", fmt.Sprintf("你的密码重置链接是: %s 15分钟内有效", passwordResetURL))
	return mailer.DialAndSend(message)
}

func HandleResetPassword(c *gin.Context, db *gorm.DB) {
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
	var findpass Findpass
	if err := db.Table("findpass").Where("link = ?", token).First(&findpass).Error; err != nil {
		c.JSON(400, gin.H{"error": "token不存在"})
		return
	}
	// 检查token是否过期
	if findpass.Time+900 < time.Now().Unix() {
		c.JSON(400, gin.H{"error": "token已过期"})
		// 删除findpass表中的token
		if err := db.Table("findpass").Delete(&findpass).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete findpass"})
			return
		}
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"message": "密码加密失败"})
		return
	}
	// 更新对应用户的密码
	if err := db.Model(&User{}).Where("username = ?", findpass.Username).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update password"})
		return
	}
	// 删除findpass表中的token
	if err := db.Table("findpass").Delete(&findpass).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete findpass"})
		return
	}
	// 更新 token
	NewToken := register.GenerateToken()
	// 在数据库tokens表中更新token
	db.Model(&register.Token{}).Where("username = ?", findpass.Username).Update("token", NewToken)
	c.JSON(200, gin.H{
		"message": "密码重置成功，已为您自动/重新登录",
		"token":   NewToken,
	})
}
