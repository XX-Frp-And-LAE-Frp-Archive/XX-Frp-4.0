package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleDeleteTunnel(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	// 获取动态路由的 tunnelid
	tunnelid := c.Param("tunnelid")
	var proxy define.Proxies
	db.Where("username = ? AND id = ?", user.Username, tunnelid).First(&proxy)
	if proxy.Status == 2 {
		respond.Respond(c, 403, "隧道已被封禁", 0)
		return
	}
	if proxy.ID == 0 {
		respond.Respond(c, 403, "隧道未找到，可能已被删除，请刷新", 0)
		return
	}
	db.Delete(&proxy)
	respond.Respond(c, 200, "删除成功", 0)
}
