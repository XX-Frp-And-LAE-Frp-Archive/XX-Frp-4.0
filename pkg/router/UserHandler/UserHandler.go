package UserHandler

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleUser(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户名
	username, _ := c.Get("username")

	if username == "" {
		return
	}

	var user define.User
	result := db.Table("users").Where("username = ?", username).First(&user)
	if result.Error != nil {
		return
	}

	groupName := user.Group

	// 根据用户名查询 limit 表
	var limitUser define.Limit
	limitResult := db.Table("limits").Where("username = ?", user.Username).First(&limitUser)
	if limitResult.Error == nil {
		// 如果存在对应记录，则返回对应的 outbound 和 inbound proxies 值
		user.Proxies = limitUser.Proxies
		user.Outbound = int64(limitUser.Outbound)
		user.Inbound = limitUser.Inbound
	} else {
		// 如果不存在对应记录，则根据用户组查询 groups 表，找到对应的 outbound 和 inbound proxies 值
		var group define.Groups
		groupResult := db.Table("groups").Where("name = ?", groupName).First(&group)
		if groupResult.Error != nil {
			respond.Respond(c, 500, "用户组不存在，请联系管理员", 0)
			return
		}
		user.Proxies = group.Proxies
		user.Outbound = int64(group.Outbound)
		user.Inbound = group.Inbound
	}

	user.EmailMD5 = getMD5Hash(user.Email)

	c.JSON(200, user)
}

// 获取字符串的 MD5 值
func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
