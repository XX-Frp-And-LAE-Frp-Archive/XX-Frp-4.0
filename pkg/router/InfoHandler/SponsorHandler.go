package InfoHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/cron"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
)

// 获取赞助者
func HandleSponsor(c *gin.Context) {
	respond.Respond(c, 200, "成功", cron.Cache)
}
