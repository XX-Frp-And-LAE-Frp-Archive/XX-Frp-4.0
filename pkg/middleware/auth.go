package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type Token struct {
	ID       uint   `gorm:"primaryKey"`
	Token    string `gorm:"unique"`
	Username string
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求的Header中获取Bearer Token
		token := c.GetHeader("Authorization")

		// 检查 Authorization 头部是否为空或不包含前缀 "Bearer "
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			c.JSON(401, gin.H{"message": "Invalid authorization token"})
			return
		}

		// 提取 token 值并去除前缀 "Bearer "
		token = strings.TrimPrefix(token, "Bearer ")

		// 查询Token对应的用户名
		var tokenData Token
		result := db.Where("token = ?", token).First(&tokenData)
		if result.Error != nil || result.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			return
		}

		// 将用户名传递给下级路由
		c.Set("username", tokenData.Username)
		c.Set("token", token)

		// 继续处理后续的请求
		c.Next()
	}
}
