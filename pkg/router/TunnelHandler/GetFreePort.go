package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/cron"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"strconv"
)

func HandleGetFreePort(c *gin.Context) {
	nodeStr := c.Query("node")
	protocol := c.Query("protocol")

	// 检查node和protocol是否为空
	if nodeStr == "" || protocol == "" {
		respond.Respond(c, 403, "node 和 protocol 不能为空", 0)
		return
	}

	// 尝试将node从字符串转换为int
	node, err := strconv.Atoi(nodeStr)
	if err != nil {
		respond.Respond(c, 403, "node 必须是一个整数", 0)
		return
	}

	FreePort, err := cron.GetFreePort(node, protocol)
	if err != nil {
		respond.Respond(c, 403, err.Error(), 0)
		return
	}
	respond.Respond(c, 200, "获取成功", gin.H{
		"free_port": FreePort,
	})
}
