package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// HandleEditTunnel 处理编辑隧道请求
func HandleEditTunnel(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	// 获取请求参数
	tunnelID := c.PostForm("tunnel_id")
	tunnelName := c.PostForm("tunnel_name")
	LocalPort := c.PostForm("local_port")
	LocalIP := c.PostForm("local_ip")
	Status := c.PostForm("status")

	var proxy define.Proxies
	result := db.First(&proxy, "username = ? AND id = ?", user.Username, tunnelID)
	if result.Error != nil {
		respond.Respond(c, 403, "未找到隧道", 0)
		return
	}
	if proxy.Status == 2 {
		respond.Respond(c, 403, "隧道已被封禁", 0)
		return
	}
	// 判断是否相同
	if tunnelName != "" {
		// 判断是否重名
		var count int64
		db.Model(&define.Proxies{}).Where("username = ? AND name = ?", user.Username, tunnelName).Count(&count)
		if count > 0 {
			respond.Respond(c, 403, "隧道名称已存在", 0)
			return
		}
		proxy.ProxyName = tunnelName
	}
	if LocalPort != "" {
		// 判断 port 是否合法
		port, err := strconv.Atoi(LocalPort)
		if err != nil || port < 1 || port > 65535 {
			respond.Respond(c, 403, "端口不合法", 0)
			return
		}
		proxy.LocalPort = LocalPort
	}
	if LocalIP != "" {
		// 判断 ip 是否为合法 ip
		if IsvalidIP(LocalIP) == false {
			respond.Respond(c, 403, "IP不合法", 0)
			return
		}
		proxy.LocalIP = LocalIP
	}
	if Status != "" {
		// 将字符串转换为int
		status, err := strconv.Atoi(Status)
		if err != nil {
			respond.Respond(c, 403, "参数错误", 0)
			return
		}
		proxy.Status = status
	}
	// 更新数据
	db.Save(&proxy)
	respond.Respond(c, 200, "更新成功!", 0)
}
