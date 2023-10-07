package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"regexp"
)

func GetTunnelByID(c *gin.Context, db *gorm.DB) {
	username, _ := c.Get("username")
	// 获取动态路由的 tunnelid
	tunnelid := c.Param("tunnelid")
	var proxy Proxy
	result := db.First(&proxy, "username = ? AND id = ?", username, tunnelid)
	if result.Error != nil {
		respond.Respond(c, 403, "未找到隧道", 0)
		return
	}
	// 读取代理的node值 然后根据node值去查询node表 获取node的ip和port和name然后返回
	var node Node
	nodeResult := db.First(&node, proxy.Node)
	if nodeResult.Error != nil {
		respond.Respond(c, 403, "未找到该节点", 0)
		return
	}
	re := regexp.MustCompile(`"([^"]+)"`)
	matches := re.FindStringSubmatch(proxy.Domain)
	if len(matches) > 1 {
		domain := matches[1]
		proxy.Domain = domain
	}
	proxy.NodeName = node.Name
	proxy.NodeHostname = node.Hostname
	proxy.NodePort = node.Port
	proxy.NodeToken = node.Token
	respond.Respond(c, 200, "Success!", proxy)
}
