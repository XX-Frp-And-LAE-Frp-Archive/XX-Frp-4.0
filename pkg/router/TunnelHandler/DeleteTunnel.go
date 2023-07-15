package TunnelHandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleDeleteTunnel(c *gin.Context, db *gorm.DB) {
	username, _ := c.Get("username")
	proxyname := c.PostForm("proxy_name")
	if proxyname == "" {
		c.JSON(400, gin.H{"error": "proxy_name is empty"})
		return
	}
	var proxy Proxies
	db.Where("username = ? AND proxy_name = ?", username, proxyname).First(&proxy)
	if proxy.ID == 0 {
		c.JSON(404, gin.H{"error": "proxy not found"})
		return
	}
	db.Delete(&proxy)
	c.JSON(200, gin.H{"message": "success"})
}
