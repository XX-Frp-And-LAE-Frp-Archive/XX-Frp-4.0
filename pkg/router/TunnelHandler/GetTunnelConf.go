package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"regexp"
)

func GetConfByNode(c *gin.Context, db *gorm.DB) {
	node := c.Param("node")
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	// 通过nodeid 在nodes 表查询到对应的节点信息
	var nodeInfo define.Node
	result := db.First(&nodeInfo, node)
	if result.Error != nil {
		respond.Respond(c, 500, "未找到该节点", 0)
		return
	}
	//通过 username 和 nodeid 在proxies表中查询到对应的代理信息 并遍历
	var proxies []define.Proxies
	proxyResult := db.Find(&proxies, "username = ? AND node = ?", user.Username, node)
	if proxyResult.Error != nil {
		respond.Respond(c, 403, "未找到该隧道", 0)

		return
	}
	// 生成ini配置文件
	conf := "[common]\n"
	conf += "server_addr = " + nodeInfo.Hostname + "\n"
	conf += "server_port = " + nodeInfo.Port + "\n"
	conf += "token = " + nodeInfo.Token + "\n"
	conf += "tcp_mux = true\n"
	conf += "protocol = tcp\n"
	conf += "user = " + user.Token + "\n"
	conf += "dns_server = 114.114.114.114\n"
	conf += "\n"

	for _, proxy := range proxies {
		conf += "[" + proxy.ProxyName + "]\n"
		conf += "privilege_mode = true\n"
		conf += "type = " + proxy.ProxyType + "\n"
		conf += "local_ip = " + proxy.LocalIP + "\n"
		conf += "local_port = " + proxy.LocalPort + "\n"
		if proxy.ProxyType == "http" || proxy.ProxyType == "https" {
			// Proxy.Domain 是 ["domain"]  格式 需要提取其中的 domain
			re := regexp.MustCompile(`"([^"]+)"`)
			matches := re.FindStringSubmatch(proxy.Domain)
			if len(matches) > 1 {
				domain := matches[1]
				conf += "custom_domains = " + domain + "\n"
			}
		} else {
			conf += "remote_port = " + proxy.RemotePort + "\n"
		}
		conf += "use_encryption = true\nuse_compression = true\n"
		conf += "\n"

	}
	// 直接ini
	c.String(200, conf)
}
func GetConfByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	// 通过 id 在proxies 表查询到对应的代理信息
	var proxy define.Proxy
	result := db.First(&proxy, id)
	if result.Error != nil {
		respond.Respond(c, 403, "未找到该隧道", 0)
		return
	}
	// 判断是否有权限
	if proxy.Username != user.Username {
		respond.Respond(c, 403, "这不是你的隧道", 0)
		return
	}

	// 通过nodeid 在nodes 表查询到对应的节点信息
	var nodeInfo define.Node
	nodeResult := db.First(&nodeInfo, proxy.Node)
	if nodeResult.Error != nil {
		respond.Respond(c, 500, "节点信息获取错误", 0)
		return
	}
	// 生成ini配置文件
	conf := "[common]\n"
	conf += "server_addr = " + nodeInfo.Hostname + "\n"
	conf += "server_port = " + nodeInfo.Port + "\n"
	conf += "token = " + nodeInfo.Token + "\n"
	conf += "tcp_mux = true\n"
	conf += "protocol = tcp\n"
	conf += "user = " + user.Token + "\n"
	conf += "dns_server = 114.114.114.114\n"
	conf += "\n"

	conf += "[" + proxy.ProxyName + "]\n"
	conf += "privilege_mode = true\n"
	conf += "type = " + proxy.ProxyType + "\n"
	conf += "local_ip = " + proxy.LocalIP + "\n"
	conf += "local_port = " + proxy.LocalPort + "\n"
	if proxy.ProxyType == "http" || proxy.ProxyType == "https" {
		re := regexp.MustCompile(`"([^"]+)"`)
		matches := re.FindStringSubmatch(proxy.Domain)
		if len(matches) > 1 {
			domain := matches[1]
			conf += "custom_domains = " + domain + "\n"
		}
	} else {
		conf += "remote_port = " + proxy.RemotePort + "\n"
	}
	conf += "use_encryption = true\nuse_compression = true\n"
	conf += "\n"
	// 直接返回ini
	c.String(200, conf)
}
