package TunnelHandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

type User struct {
	ID       uint
	Username string
	Group    string
}

func HandleGetNodeList(c *gin.Context, db *gorm.DB) {
	username, _ := c.Get("username")
	var user User
	db.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	// 获取所有节点
	var nodes []Node
	db.Find(&nodes)
	// 遍历所有节点 获取节点支持的 group 列表 然后判断用户是否在 group 列表中 如果在则添加到返回列表中
	var nodeList []Node
	for _, node := range nodes {
		groups := strings.Split(node.Group, ";")
		for _, group := range groups {
			if group == user.Group {
				nodeList = append(nodeList, node)
			}
		}
		// 如果节点支持的 group 列表中包含 all 则添加到返回列表中
		if node.Group == "all" {
			nodeList = append(nodeList, node)
		}
		// 返回节点列表
		c.JSON(200, gin.H{"nodes": nodeList})
	}
}
