package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

func HandleGetNodeList(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	// 获取所有节点
	var nodes []define.Nodes
	db.Find(&nodes)
	// 遍历所有节点 获取节点支持的 group 列表 然后判断用户是否在 group 列表中 如果在则添加到返回列表中
	var nodeList []define.Nodes
	for _, node := range nodes {
		groups := strings.Split(node.Group, ";")
		for _, group := range groups {
			if group == user.Group && node.Status != 403 {
				nodeList = append(nodeList, node)
			}
		}
		// 如果节点支持的 group 列表中包含 all 则添加到返回列表中
		// 403 禁用+隐藏
		// 200 正常
		// 201 无状态数据
		// 500 服务器异常
		if node.Group == "all" && node.Status != 403 {
			nodeList = append(nodeList, node)
		}
	}
	// 返回节点列表
	respond.Respond(c, 200, "获取成功", nodeList)
}
func HandleGetAllNode(c *gin.Context, db *gorm.DB) {
	// 获取所有节点
	var nodes []define.Nodes
	// 指定查询 nodes 表
	db.Find(&nodes)
	respond.Respond(c, 200, "获取成功", nodes)
}
