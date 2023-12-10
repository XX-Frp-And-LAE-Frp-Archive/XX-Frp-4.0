package UserHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	register "github.com/ahmr-bot/ME-Frp/pkg/router/RegisterHandler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleRefreshToken(c *gin.Context, db *gorm.DB) {
	username, _ := c.Get("username")
	NewToken := register.GenerateToken()
	// 在数据库tokens表中更新token
	db.Model(&define.User{}).Where("username = ?", username).Update("token", NewToken)
	respond.Respond(c, 200, "Token更新成功", gin.H{
		"newToken": NewToken,
	})
}
