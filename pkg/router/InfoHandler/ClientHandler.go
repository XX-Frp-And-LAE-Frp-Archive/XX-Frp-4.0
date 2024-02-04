package InfoHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/cron"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
)

// HandleSponsor 获取赞助者
func HandleClient(c *gin.Context) {
	respond.Respond(c, 200, "Success!", cron.ClientCache)
}
