package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/cron"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CloseTunnel(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	// 获取动态路由的 tunnelid
	tunnelid := c.Param("tunnelid")
	var proxy define.Proxies
	result := db.First(&proxy, "username = ? AND id = ?", user.Username, tunnelid)
	if result.Error != nil {
		respond.Respond(c, 403, "未找到隧道", 0)
		return
	}
	if proxy.Status == 2 {
		respond.Respond(c, 403, "隧道已被封禁", 0)
		return
	}
	if proxy.Status == 1 {
		respond.Respond(c, 403, "隧道处于禁用状态", 0)
		return
	}
	// 读取代理的node值 然后根据node值去查询node表 获取node的ip和port和name然后返回
	var node define.Node
	nodeResult := db.First(&node, proxy.Node)
	if nodeResult.Error != nil {
		respond.Respond(c, 403, "未找到该节点", 0)
		return
	}
	// 设置 proxy 的状态为 1
	proxy.Status = 1
	proxy.Online = "offline"
	db.Save(&proxy)
	// 使用 Basic Auth 认证方式 用户名 admin 密码为 node.AdminPass 地址为 node.Hostname:node.AdminPort 使用 node.mefrp.com 作为host域名
	// 发送请求到 /api/client/close/ proxy.RunID 代表关闭隧道
	path := "/client/close/" + proxy.RunID

	resp, err := cron.FetchData(node, path)
	if err != nil {
		respond.Respond(c, 500, "关闭隧道失败，服务器错误", 0)
		return
	}
	// 判断返回的状态码 如果不是200 则返回错误
	if resp.StatusCode != 200 {
		respond.Respond(c, 403, "关闭隧道失败，客户端没有启动", 0)
		return
	}

	respond.Respond(c, 200, "关闭成功", 0)
}
