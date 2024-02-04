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
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	groupName := user.Group

	// 根据用户名查询 limit 表

	var userRes define.UserRes
	var limitUser define.Limit
	limitResult := db.Table("limits").Where("username = ?", user.Username).First(&limitUser)
	if limitResult.Error == nil {
		// 如果存在对应记录，则返回对应的 outbound 和 inbound proxies 值
		userRes.Proxies = limitUser.Proxies
		userRes.Outbound = int64(limitUser.Outbound)
		userRes.Inbound = limitUser.Inbound
	} else {
		// 如果不存在对应记录，则根据用户组查询 groups 表，找到对应的 outbound 和 inbound proxies 值
		var group define.Groups
		groupResult := db.Table("groups").Where("name = ?", groupName).First(&group)
		if groupResult.Error != nil {
			respond.Respond(c, 500, "用户组不存在，请联系管理员", 0)
			return
		}
		userRes.Proxies = group.Proxies
		userRes.Outbound = int64(group.Outbound)
		userRes.Inbound = group.Inbound
	}
	// 读取 todaytraffic 表中的记录
	var todayTraffic define.TodayTraffic
	todayTrafficResult := db.Table("todaytraffics").Where("username = ?", user.Username).First(&todayTraffic)
	if todayTrafficResult.Error != nil {
		respond.Respond(c, 500, "读取 todaytraffic 表失败", 0)
		return
	}

	userRes.EmailMD5 = getMD5Hash(user.Email)
	userRes.ID = user.ID
	userRes.Username = user.Username
	userRes.Password = user.Password
	userRes.Email = user.Email
	userRes.Traffic = user.Traffic
	userRes.RegTime = user.Regtime
	userRes.Token = user.Token
	userRes.Group = user.Group
	userRes.Status = user.Status
	userRes.TodayTraffic = todayTraffic.Traffic

	respond.Respond(c, 200, "Success!", userRes)
}

// 获取字符串的 MD5 值
func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
