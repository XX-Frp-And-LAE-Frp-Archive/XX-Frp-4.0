package TunnelHandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Proxy struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	ProxyName    string `json:"proxy_name"`
	ProxyType    string `json:"proxy_type"`
	LocalIP      string `json:"local_ip"`
	LocalPort    string `json:"local_port"`
	Domain       string `json:"domain"`
	Node         string `json:"node"`
	NodeName     string `json:"node_name"`
	NodeHostname string `json:"node_hostname"`
	NodePort     string `json:"node_port"`
	NodeToken    string `json:"node_token"`
	RemotePort   string `json:"remote_port"`
}
type Node struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Token    string `json:"token"`
	Group    string `json:"group"`
}

func GetTunnelList(c *gin.Context, db *gorm.DB) {
	username, _ := c.Get("username")

	var proxies []Proxy
	result := db.Find(&proxies, "username = ?", username)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	// 读取每一个代理的node值 然后根据node值去查询node表 获取node的ip和port和name然后返回
	for i, proxy := range proxies {
		var node Node
		nodeResult := db.First(&node, proxy.Node)
		if nodeResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": nodeResult.Error.Error()})
			return
		}
		proxies[i].NodeName = node.Name
		proxies[i].NodeHostname = node.Hostname
		proxies[i].NodePort = node.Port
		proxies[i].NodeToken = node.Token
	}
	c.JSON(http.StatusOK, proxies)
}
