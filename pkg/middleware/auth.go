package middleware

import (
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
			respond.Respond(c, 401, "错误的 Token 格式!", 0)
			return
		} else {

			// 提取 token 值并去除前缀 "Bearer "
			token = strings.TrimPrefix(token, "Bearer ")

			// 查询Token对应的用户名
			var tokenData Token
			result := db.Where("token = ?", token).First(&tokenData)
			if result.Error != nil || result.RowsAffected == 0 {
				respond.Respond(c, 401, "Token 不存在!", 0)
				return
			} else {

				// 将用户名传递给下级路由
				c.Set("username", tokenData.Username)
				c.Set("token", token)

				c.Next()
			}
		}
	}
}
