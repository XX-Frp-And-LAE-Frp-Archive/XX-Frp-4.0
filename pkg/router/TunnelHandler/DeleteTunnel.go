package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleDeleteTunnel(c *gin.Context, db *gorm.DB) {
	username, _ := c.Get("username")
	// 获取动态路由的 tunnelid
	tunnelid := c.Param("tunnelid")
	var proxy Proxies
	db.Where("username = ? AND id = ?", username, tunnelid).First(&proxy)
	if proxy.ID == 0 {
		respond.Respond(c, 403, "隧道未找到，可能已被删除，请刷新", 0)
		return
	}
	db.Delete(&proxy)
	respond.Respond(c, 200, "删除成功", 0)
}
