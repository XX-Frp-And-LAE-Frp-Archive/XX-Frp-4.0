package TunnelHandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetConfByNode(c *gin.Context, db *gorm.DB) {
	node := c.Param("node")
	username, _ := c.Get("username")
	user, _ := c.Get("token")
	// 通过nodeid 在nodes 表查询到对应的节点信息
	var nodeInfo Node
	result := db.First(&nodeInfo, node)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	//通过 username 和 nodeid 在proxies表中查询到对应的代理信息 并遍历
	var proxies []Proxy
	proxyResult := db.Find(&proxies, "username = ? AND node = ?", username, node)
	if proxyResult.Error != nil {
		c.JSON(500, gin.H{"error": proxyResult.Error.Error()})
		return
	}
	// 生成ini配置文件
	conf := "[common]\n"
	conf += "server_addr = " + nodeInfo.Hostname + "\n"
	conf += "server_port = " + nodeInfo.Port + "\n"
	conf += "token = " + nodeInfo.Token + "\n"
	conf += "tcp_mux = true\n"
	conf += "protocol = tcp\n"
	conf += "user = " + user.(string) + "\n"
	conf += "dns_server = 114.114.114.114\n"
	conf += "\n"

	for _, proxy := range proxies {
		conf += "[" + proxy.ProxyName + "]\n"
		conf += "privilege_mode = true\n"
		conf += "type = " + proxy.ProxyType + "\n"
		conf += "local_ip = " + proxy.LocalIP + "\n"
		conf += "local_port = " + proxy.LocalPort + "\n"
		if proxy.ProxyType == "http" || proxy.ProxyType == "https" {
			conf += "domain = " + proxy.Domain + "\n"
		} else {
			conf += "remote_port = " + proxy.RemotePort + "\n"
		}
		conf += "use_encryption = true\nuse_compression = true\n"
		conf += "\n"

	}
	c.String(200, conf)
}
func GetConfByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	username, _ := c.Get("username")
	user, _ := c.Get("token")
	// 通过 id 在proxies 表查询到对应的代理信息
	var proxy Proxy
	result := db.First(&proxy, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	// 判断是否有权限
	if proxy.Username != username {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	// 通过nodeid 在nodes 表查询到对应的节点信息
	var nodeInfo Node
	nodeResult := db.First(&nodeInfo, proxy.Node)
	if nodeResult.Error != nil {
		c.JSON(500, gin.H{"error": nodeResult.Error.Error()})
		return
	}
	// 生成ini配置文件
	conf := "[common]\n"
	conf += "server_addr = " + nodeInfo.Hostname + "\n"
	conf += "server_port = " + nodeInfo.Port + "\n"
	conf += "token = " + nodeInfo.Token + "\n"
	conf += "tcp_mux = true\n"
	conf += "protocol = tcp\n"
	conf += "user = " + user.(string) + "\n"
	conf += "dns_server = 114.114.114.114\n"
	conf += "\n"

	conf += "[" + proxy.ProxyName + "]\n"
	conf += "privilege_mode = true\n"
	conf += "type = " + proxy.ProxyType + "\n"
	conf += "local_ip = " + proxy.LocalIP + "\n"
	conf += "local_port = " + proxy.LocalPort + "\n"
	if proxy.ProxyType == "http" || proxy.ProxyType == "https" {
		conf += "domain = " + proxy.Domain + "\n"
	} else {
		conf += "remote_port = " + proxy.RemotePort + "\n"
	}
	conf += "use_encryption = true\nuse_compression = true\n"
	conf += "\n"
	c.String(200, conf)
}
