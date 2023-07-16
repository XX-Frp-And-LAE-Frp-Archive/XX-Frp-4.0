package StartHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	"github.com/ahmr-bot/ME-Frp/pkg/router/TunnelHandler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

type Node struct {
	ID    uint
	Group string
}
type Token struct {
	ID       uint
	Username string
	Token    string
	status   bool
}
type User struct {
	ID       uint
	Username string
	Group    string
}

type Proxy struct {
	ID         uint
	Username   string
	ProxyName  string
	ProxyType  string
	Domain     string
	RemotePort string
	RunID      string
	Status     bool
	Lastupdate int64
	Node       string
}
type Limit struct {
	ID       uint
	Username string
	Inbound  int
	Outbound int
}

type Group struct {
	ID       uint
	Name     string
	Inbound  int
	Outbound int
}

func HandleStart(c *gin.Context, db *gorm.DB) {
	conf := config.GetConfig()
	// 获取 url 参数 action 的值
	action := c.Query("action")
	// 获取 url 参数 token 的值
	apitoken := c.Query("apitoken")
	user := c.Query("user")

	//将 apitoken 分为 | 前后两部分
	apitokenSlice := strings.Split(apitoken, "|")
	//apitoken 为 | 前面的部分
	token := apitokenSlice[0]
	//id 为 | 后面的部分
	nodeid := apitokenSlice[1]
	if token != conf.Server.Token {
		c.JSON(403, gin.H{
			"status":  403,
			"message": "token 验证失败",
		})
		return
	}
	// 读取node表id字段 判断id是否存在
	var node Node
	result := db.Where("id = ?", nodeid).First(&node)
	if result.Error != nil {
		c.JSON(403, gin.H{
			"status":  403,
			"message": "nodeid 不存在",
		})
		return
	}
	// 根据 action 的值进行判断
	switch strings.ToLower(action) {
	case "checktoken":
		// 在tokens表中查询是否存在指定的token
		var tokenRecord Token
		result := db.Where("token = ?", user).First(&tokenRecord)
		if result.Error != nil {
			c.JSON(403, gin.H{
				"status":  403,
				"message": "user  错误",
			})
			return
		}
		// 判断用户status是否为0
		if tokenRecord.status == true {
			c.JSON(403, gin.H{
				"status":  403,
				"message": "token  已被禁用",
			})
			return
		}
		// 判断用户是否有权限使用该节点
		var node Node
		nodeResult := db.First(&node, "id = ?", nodeid)
		if nodeResult.Error != nil {
			c.JSON(500, gin.H{
				"status":  500,
				"message": nodeResult.Error.Error()},
			)
			return
		}
		// 通过 tokenRecord.Username 查询用户组
		var user User
		userResult := db.First(&user, "username = ?", tokenRecord.Username)
		if userResult.Error != nil {
			c.JSON(500, gin.H{
				"status":  500,
				"message": userResult.Error.Error(),
			})
			return
		}
		// 判断用户是否有权限使用该节点
		if !TunnelHandler.CheckGroup(node.Group, user.Group) {
			c.JSON(403, gin.H{
				"status":  403,
				"message": "你没有权限使用该节点",
			})
			return
		}
		c.JSON(200, gin.H{
			"status":  200,
			"success": true,
			"message": "登陆成功 欢迎使用 ME Frp 服务",
		})
	case "checkproxy":
		// 在tokens表中查询是否存在指定的token
		var tokenRecord Token
		result := db.Where("token = ?", user).First(&tokenRecord)
		if result.Error != nil {
			c.JSON(403, gin.H{
				"status":  403,
				"message": "user  错误",
			})
			return
		}
		// 判断用户status是否为0
		if tokenRecord.status == true {
			c.JSON(403, gin.H{
				"status":  403,
				"message": "token  已被禁用",
			})
			return
		}
		ProxyType := c.Query("proxy_type")
		ProxyName := c.Query("proxy_name")
		// 直接获取的proxy_name 是 user.proxy_name 格式 将其分割 获取 proxy_name
		ProxyNameSlice := strings.Split(ProxyName, ".")
		ProxyName = ProxyNameSlice[1]

		if ProxyType == "tcp" || ProxyType == "udp" {
			RemotePort := c.Query("remote_port")
			// 通过 tokenRecord.Username ProxyType ProxyName RemotePort 查询是否存在该条记录
			var proxy Proxy
			proxyResult := db.First(&proxy, "username = ? AND proxy_type = ? AND proxy_name = ? AND remote_port = ?", tokenRecord.Username, ProxyType, ProxyName, RemotePort)
			if proxyResult.Error != nil {
				c.JSON(403, gin.H{
					"status":  403,
					"message": "该条记录不存在",
				})
				return
			}
			// 判断用户是否有权限使用该隧道
			if proxy.Status == true {
				c.JSON(403, gin.H{
					"status":  403,
					"message": "该隧道已被禁用",
				})
				return
			}
		} else if ProxyType == "http" || ProxyType == "https" {
			Domain := c.Query("domain")
			// 通过 tokenRecord.Username ProxyType ProxyName Domain 查询是否存在该条记录
			var proxy Proxy
			proxyResult := db.First(&proxy, "username = ? AND proxy_type = ? AND proxy_name = ? AND domain = ?", tokenRecord.Username, ProxyType, ProxyName, Domain)
			if proxyResult.Error != nil {
				c.JSON(403, gin.H{
					"status":  403,
					"message": "该条记录不存在",
				})
				return
			}
			// 判断用户是否有权限使用该隧道
			if proxy.Status == true {
				c.JSON(403, gin.H{
					"status":  403,
					"message": "该隧道已被禁用",
				})
				return
			}
		} else {
			c.JSON(403, gin.H{
				"status":  403,
				"message": "proxy_type 错误",
			})
			return
		}
		run_id := c.Query("run_id")

		// 将该条记录的run_id更新为run_id
		var proxy Proxy
		proxyResult := db.First(&proxy, "username = ? AND proxy_type = ? AND proxy_name = ?", tokenRecord.Username, ProxyType, ProxyName)
		if proxyResult.Error != nil {
			c.JSON(403, gin.H{
				"status":  403,
				"message": "该条记录不存在",
			})
			return
		}
		// 获取该条记录的run_id
		// 存储到数据库中
		proxy.RunID = run_id
		db.Save(&proxy)
		c.JSON(200, gin.H{
			"status":  200,
			"success": true,
			"message": "验证成功",
		})
	case "getlimit":
		// 在tokens表中查询是否存在指定的token
		var tokenRecord Token
		result := db.Where("token = ?", user).First(&tokenRecord)
		if result.Error != nil {
			c.JSON(403, gin.H{
				"status":  403,
				"message": "user  错误",
			})
			return
		}
		// 判断用户status是否为0
		if tokenRecord.status == true {
			c.JSON(403, gin.H{
				"status":  403,
				"message": "token  已被禁用",
			})
			return
		}
		// 通过tokenRecord.Username查询 limit表中的记录
		var limit Limit
		limitResult := db.First(&limit, "username = ?", tokenRecord.Username)
		if limitResult.Error != nil {
			// 通过 tokenRecord.username 查询 user 表中的 group 字段
			var user User
			userResult := db.First(&user, "username = ?", tokenRecord.Username)
			if userResult.Error != nil {
				c.JSON(403, gin.H{
					"status":  403,
					"message": "user  错误",
				})
				return
			}
			// 通过 user.group 查询 group 表中的记录
			var group Group
			groupResult := db.First(&group, "name = ?", user.Group)
			if groupResult.Error != nil {
				c.JSON(403, gin.H{
					"status":  403,
					"message": "group  错误",
				})
				return
			}
			// 返回 group 表中的 limit 字段
			c.JSON(200, gin.H{
				"status":  200,
				"max-in":  group.Inbound,
				"max-out": group.Outbound,
			})
			return
		}
		// 返回 limit 表中的 inbound outbound 字段
		c.JSON(200, gin.H{
			"status":  200,
			"max-in":  limit.Inbound,
			"max-out": limit.Outbound,
		})
	}
}
