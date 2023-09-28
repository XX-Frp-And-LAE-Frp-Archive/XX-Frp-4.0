package InfoHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/cron"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
)

func HandleStatistics(c *gin.Context) {
	// 从缓存中读取数据
	data, ok := cron.CacheData.(map[string]int64)
	if !ok {
		respond.Respond(c, 500, "获取缓存数据出错", nil)
		return
	}
	
	respond.Respond(c, 200, "成功", gin.H{
		"userCount":  data["userCount"],
		"proxyCount": data["proxyCount"],
		"nodeCount":  data["nodeCount"],
		"trafficSum": data["trafficSum"],
	})
}
