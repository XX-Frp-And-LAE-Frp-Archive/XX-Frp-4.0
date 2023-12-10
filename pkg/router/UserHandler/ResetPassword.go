package UserHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	register "github.com/ahmr-bot/ME-Frp/pkg/router/RegisterHandler"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HandleResetPassword(c *gin.Context, db *gorm.DB) {
	// 获取表单数据中的新密码
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("password")
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		respond.Respond(c, 403, "旧密码错误!", 0)
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		respond.Respond(c, 500, "密码加密失败", 0)
		return
	}
	// 更新对应用户的密码
	if err := db.Model(&define.User{}).Where("username = ?", user.Username).Update("password", string(hashedPassword)).Error; err != nil {
		respond.Respond(c, 500, "更新密码失败", 0)
		return
	}
	// 更新 token
	NewToken := register.GenerateToken()
	// 在数据库users表中更新token
	db.Model(&define.User{}).Where("username = ?", user.Username).Update("token", NewToken)
	respond.Respond(c, 200, "密码重置成功，已为您自动重新登录", gin.H{
		"token": NewToken,
	})
}
