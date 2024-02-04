package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"regexp"
)

func GetTunnelList(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	var proxies []define.Proxy
	result := db.Find(&proxies, "username = ?", user.Username)
	if result.Error != nil {
		respond.Respond(c, 403, "未找到该隧道", 0)
		return
	}
	// 读取每一个代理的node值 然后根据node值去查询node表 获取node的ip和port和name然后返回
	for i, proxy := range proxies {
		var node define.Node
		nodeResult := db.First(&node, proxy.Node)
		if nodeResult.Error != nil {
			respond.Respond(c, 403, "未找到该隧道", 0)
			return
		}
		re := regexp.MustCompile(`"([^"]+)"`)
		matches := re.FindStringSubmatch(proxies[i].Domain)
		if len(matches) > 1 {
			domain := matches[1]
			proxies[i].Domain = domain
		}
		proxies[i].NodeName = node.Name
		proxies[i].NodeHostname = node.Hostname
		proxies[i].NodePort = node.Port
		proxies[i].NodeToken = node.Token
		proxies[i].Status = proxy.Status
	}
	respond.Respond(c, 200, "Success!", proxies)
}
