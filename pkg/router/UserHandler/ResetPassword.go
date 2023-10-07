package UserHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	register "github.com/ahmr-bot/ME-Frp/pkg/router/RegisterHandler"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HandleResetPassword(c *gin.Context, db *gorm.DB) {
	// 获取表单数据中的新密码
	newPassword := c.PostForm("password")
	username, _ := c.Get("username")

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		respond.Respond(c, 500, "密码加密失败", 0)
		return
	}
	// 更新对应用户的密码
	if err := db.Model(&User{}).Where("username = ?", username).Update("password", string(hashedPassword)).Error; err != nil {
		respond.Respond(c, 500, "更新密码失败", 0)
		return
	}
	// 更新 token
	NewToken := register.GenerateToken()
	// 在数据库tokens表中更新token
	db.Model(&register.Token{}).Where("username = ?", username).Update("token", NewToken)
	respond.Respond(c, 200, "密码重置成功，已为您自动重新登录", gin.H{
		"token": NewToken,
	})
}
