package TunnelHandler

import (
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
		c.JSON(404, gin.H{"error": "proxy not found"})
		return
	}
	db.Delete(&proxy)
	c.JSON(200, gin.H{"message": "success"})
}
