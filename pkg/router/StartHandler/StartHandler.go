package StartHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/ahmr-bot/ME-Frp/pkg/router/TunnelHandler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

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
		respond.Respond(c, 503, "Token 错误", 0)
		return
	}
	// 读取node表id字段 判断id是否存在
	var node define.Node
	result := db.Where("id = ?", nodeid).First(&node)
	if result.Error != nil {
		respond.Respond(c, 503, "节点不存在", 0)
		return
	}
	// 根据 action 的值进行判断
	switch strings.ToLower(action) {
	case "checktoken":
		// 在users表中查询是否存在指定的token
		var User define.User
		result := db.Where("token = ?", user).First(&User)
		if result.Error != nil {
			respond.Respond(c, 403, "user 错误", 0)
			return
		}
		// 判断用户status是否为0
		if User.Status == 1 {
			respond.Respond(c, 403, "已被封禁", 0)
			return
		}
		// 判断 status 是否为2
		if User.Status == 2 {
			respond.Respond(c, 403, "流量已耗尽，隧道启动功能将在流量>0的5分钟内恢复", 0)
			return
		}
		// 判断用户是否有权限使用该节点
		var node define.Node
		nodeResult := db.First(&node, "id = ?", nodeid)
		if nodeResult.Error != nil {
			respond.Respond(c, 403, "无权使用此节点", 0)
			return
		}
		// 通过 User.Username 查询用户组
		var user define.User
		userResult := db.First(&user, "username = ?", User.Username)
		if userResult.Error != nil {
			respond.Respond(c, 403, "未找到用户", 0)
			return
		}
		// 判断用户是否有权限使用该节点
		if !TunnelHandler.CheckGroup(node.Group, user.Group) {
			respond.Respond(c, 403, "你无权使用此节点", 0)
			return
		}
		// 此处大抵是还不能改（）,直接写成规范格式罢
		c.JSON(200, gin.H{
			"status":  200,
			"success": true,
			"message": "登陆成功 欢迎使用 ME Frp 服务",
		})
	case "checkproxy":
		// 在tokens表中查询是否存在指定的token
		var User define.User
		result := db.Where("token = ?", user).First(&User)
		if result.Error != nil {
			respond.Respond(c, 403, "user 错误", 0)
			return
		}
		// 判断用户status是否为0
		if User.Status == 1 {
			respond.Respond(c, 403, "已被封禁", 0)
			return
		}
		ProxyType := c.Query("proxy_type")
		ProxyName := c.Query("proxy_name")
		// 直接获取的proxy_name 是 user.proxy_name 格式 将其分割 获取 proxy_name
		ProxyNameSlice := strings.Split(ProxyName, ".")
		ProxyName = ProxyNameSlice[1]

		if ProxyType == "tcp" || ProxyType == "udp" {
			RemotePort := c.Query("remote_port")
			// 通过 User.Username ProxyType ProxyName RemotePort 查询是否存在该条记录
			var proxy define.Proxies
			proxyResult := db.First(&proxy, "username = ? AND proxy_type = ? AND proxy_name = ? AND remote_port = ?", User.Username, ProxyType, ProxyName, RemotePort)
			if proxyResult.Error != nil {
				respond.Respond(c, 403, "隧道不存在", 0)
				return
			}
			// 判断用户是否有权限使用该隧道
			if proxy.Status == 1 {
				respond.Respond(c, 403, "隧道被禁用", 0)
				return
			}
			if proxy.Status == 2 {
				respond.Respond(c, 403, "隧道被封禁", 0)
				return
			}
		} else if ProxyType == "http" || ProxyType == "https" {
			Domain := c.Query("domain")
			// 通过 User.Username ProxyType ProxyName Domain 查询是否存在该条记录
			var proxy define.Proxies
			proxyResult := db.First(&proxy, "username = ? AND proxy_type = ? AND proxy_name = ? AND domain = ?", User.Username, ProxyType, ProxyName, Domain)
			if proxyResult.Error != nil {
				respond.Respond(c, 403, "隧道不存在 ", 0)
				return
			}
		} else {
			respond.Respond(c, 403, "隧道类型不支持", 0)
			return
		}
		run_id := c.Query("run_id")

		// 将该条记录的run_id更新为run_id
		var proxy define.Proxies
		proxyResult := db.First(&proxy, "username = ? AND proxy_type = ? AND proxy_name = ?", User.Username, ProxyType, ProxyName)
		if proxyResult.Error != nil {
			respond.Respond(c, 403, "隧道不存在", 0)
			return
		}
		// 获取该条记录的run_id
		// 存储到数据库中
		proxy.RunID = run_id
		db.Save(&proxy)
		// 依旧不能改，嘤嘤嘤……
		c.JSON(200, gin.H{
			"status":  200,
			"success": true,
			"message": "验证成功",
		})
	case "getlimit":
		// 在tokens表中查询是否存在指定的token
		var User define.User
		result := db.Where("token = ?", user).First(&User)
		if result.Error != nil {
			respond.Respond(c, 403, "user 错误", 0)
			return
		}
		// 判断用户status是否为0
		if User.Status == 1 {
			respond.Respond(c, 403, "已被封禁", 0)
			return
		}
		// 通过tokenRecord.Username查询 limit表中的记录
		var limit define.Limit
		limitResult := db.First(&limit, "username = ?", User.Username)
		if limitResult.Error != nil {
			// 通过 User.username 查询 user 表中的 group 字段
			var user define.User
			userResult := db.First(&user, "username = ?", User.Username)
			if userResult.Error != nil {
				respond.Respond(c, 403, "user 错误", 0)
				return
			}
			// 通过 user.group 查询 group 表中的记录
			var group define.Groups
			groupResult := db.First(&group, "name = ?", user.Group)
			if groupResult.Error != nil {
				respond.Respond(c, 403, "group 错误", 0)
				return
			}
			// 返回 group 表中的 limit 字段，还是不能改
			c.JSON(200, gin.H{
				"status":  200,
				"max-in":  group.Inbound,
				"max-out": group.Outbound,
			})
			return
		}
		// 返回 limit 表中的 inbound outbound 字段，还是不能改
		c.JSON(200, gin.H{
			"status":  200,
			"max-in":  limit.Inbound,
			"max-out": limit.Outbound,
		})
	}
}
