package InfoHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/cron"
	"github.com/gin-gonic/gin"
)

// 获取赞助者
func HandleSponsor(c *gin.Context) {
	c.JSON(200, cron.Cache)
}
