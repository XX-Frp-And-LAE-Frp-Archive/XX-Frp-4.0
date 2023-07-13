package UserHandler

import (
	register "github.com/ahmr-bot/ME-Frp/pkg/router/RegisterHandler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleRefreshToken(c *gin.Context, db *gorm.DB) {
	username, _ := c.Get("username")
	NewToken := register.GenerateToken()
	// 在数据库tokens表中更新token
	db.Model(&register.Token{}).Where("username = ?", username).Update("token", NewToken)
	c.JSON(200, gin.H{
		"code":     200,
		"msg":      "刷新成功",
		"newToken": NewToken,
	})
}
