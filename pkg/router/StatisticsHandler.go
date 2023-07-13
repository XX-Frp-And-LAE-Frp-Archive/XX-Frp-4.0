package router

import (
	"github.com/ahmr-bot/ME-Frp/pkg/cron"
	"github.com/gin-gonic/gin"
)

func HandleStatistics(c *gin.Context) {
	// 从缓存中读取数据
	data, ok := cron.CacheData.(map[string]int64)
	if !ok {
		c.JSON(500, gin.H{
			"message": "Failed to fetch stats data",
		})
		return
	}

	c.JSON(200, gin.H{
		"userCount":  data["userCount"],
		"proxyCount": data["proxyCount"],
		"nodeCount":  data["nodeCount"],
		"trafficSum": data["trafficSum"],
	})
}
