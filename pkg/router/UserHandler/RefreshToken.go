package UserHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	register "github.com/ahmr-bot/ME-Frp/pkg/router/RegisterHandler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleRefreshToken(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	NewToken := register.GenerateToken()
	// 在数据库tokens表中更新token
	db.Model(&define.User{}).Where("username = ?", user.Username).Update("token", NewToken)
	respond.Respond(c, 200, "Token更新成功", gin.H{
		"newToken": NewToken,
	})
}
