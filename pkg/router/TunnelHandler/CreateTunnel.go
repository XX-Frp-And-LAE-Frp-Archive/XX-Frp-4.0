package TunnelHandler

import (
	register "github.com/ahmr-bot/ME-Frp/pkg/router/RegisterHandler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net"
	"strconv"
	"strings"
	"time"
)

type Group struct {
	ID      int    `json:"id"`
	Group   string `json:"group"`
	Proxies int    `json:"proxies"`
}
type Limit struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Inbound  int    `json:"inbound"`
	Outbound int    `json:"outbound"`
	Proxies  int    `json:"proxies"`
}

type Proxies struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	ProxyName  string `json:"proxy_name"`
	ProxyType  string `json:"proxy_type"`
	LocalIP    string `json:"local_ip"`
	LocalPort  string `json:"local_port"`
	Domain     string `json:"domain"`
	Node       string `json:"node"`
	RemotePort string `json:"remote_port"`
	Lastupdate int64  `json:"lastupdate"`
}

func HandleCreateTunnel(c *gin.Context, db *gorm.DB) {
	username, _ := c.Get("username")
	// 从表单中获取数据
	var proxy Proxies
	proxy.Username = username.(string)
	proxy.ProxyName = c.PostForm("proxy_name")
	proxy.ProxyType = c.PostForm("proxy_type")
	proxy.LocalIP = c.PostForm("local_ip")
	proxy.LocalPort = c.PostForm("local_port")
	proxy.Node = c.PostForm("node")
	if proxy.ProxyType == "http" || proxy.ProxyType == "https" {
		proxy.Domain = c.PostForm("domain")
	} else {
		proxy.RemotePort = c.PostForm("remote_port")
	}
	//检查proxyname是否合法
	if proxy.ProxyName == "" {
		c.JSON(400, gin.H{"message": "proxy_name is empty"})
		return
	}
	// 检查proxytype是否合法
	if proxy.ProxyType != "tcp" && proxy.ProxyType != "udp" && proxy.ProxyType != "http" && proxy.ProxyType != "https" {
		c.JSON(400, gin.H{"message": "proxy_type is not valid"})
		return
	}
	// 使用正则表达式检查 localip 是不是 ip地址
	// 使用正则表达式检查 localport 是不是 端口号
	// 使用正则表达式检查 domain 是不是域名
	// 使用正则表达式检查 remoteport 是不是 端口号
	if !IsvalidIP(proxy.LocalIP) {
		c.JSON(400, gin.H{"message": "local_ip is not valid"})
		return
	}
	if !IsvalidPort(proxy.LocalPort) {
		c.JSON(400, gin.H{"message": "local_port is not valid"})
		return
	}
	if proxy.ProxyType == "http" || proxy.ProxyType == "https" {
		if !IsValidDomain(proxy.Domain) {
			c.JSON(400, gin.H{"message": "domain is not valid"})
			return
		}
	} else {
		if !IsvalidPort(proxy.RemotePort) {
			c.JSON(400, gin.H{"message": "remote_port is not valid"})
			return
		}
	}

	// 检查 单用户的代理名称是否重复
	var proxyCount int64
	db.Model(&Proxy{}).Where("username = ? AND proxy_name = ?", username, proxy.ProxyName).Count(&proxyCount)
	if proxyCount > 0 {
		c.JSON(400, gin.H{"message": "proxy name already exists"})
		return
	}
	//在 users 表使用 username 获取 group 值 然后根据 group 值去查询 group 表中的 proxies 字段 如果 该用户已经拥有的代理数量大于等于 group 表中的 proxies 字段的值 则不允许创建
	var user register.User
	userResult := db.First(&user, "username = ?", username)
	if userResult.Error != nil {
		c.JSON(500, gin.H{"message": userResult.Error.Error()})
		return
	}
	// 获取用户的 group 值
	var group Group
	groupResult := db.First(&group, "name = ?", user.Group)
	if groupResult.Error != nil {
		c.JSON(500, gin.H{"message": groupResult.Error.Error()})
		return
	}
	//  limits 表中查询该用户是否有自定义的限制
	var limit Limit
	limitResult := db.First(&limit, "username = ?", username)
	if limitResult.Error != nil {
		limit.Proxies = group.Proxies
	}
	var proxies []Proxy
	db.Find(&proxies, "username = ?", username)
	if len(proxies) >= limit.Proxies {
		c.JSON(400, gin.H{"message": "proxy count limit"})
		return
	}
	// 检查 node 是否存在
	var nodeCount int64
	db.Model(&Node{}).Where("id = ?", proxy.Node).Count(&nodeCount)
	if nodeCount == 0 {
		c.JSON(400, gin.H{"message": "node not exists"})
		return
	}
	// 检查 nodes 表中对应的group值中是否包含该用户的group值
	var node Node
	nodeResult := db.First(&node, "id = ?", proxy.Node)
	if nodeResult.Error != nil {
		c.JSON(500, gin.H{"message": nodeResult.Error.Error()})
		return
	}
	// node 的 group 值是多个值组成的字符串 用分号分隔 例如 admin;default;realname;trustuser;
	// 检查 node 的 group 值中是否包含该用户的 group 值
	if !CheckGroup(node.Group, user.Group) {
		c.JSON(400, gin.H{"message": "You are not allowed to create a proxy on this node"})
		return
	}

	// 检查 选择的 proxy_type 是否在 nodes 表中的 allow_type 字段中 allow_type 字段是多个值组成的字符串 用分号分隔 例如 tcp;udp;http;https;
	if !CheckAllowType(node.AllowType, proxy.ProxyType) {
		c.JSON(400, gin.H{"message": "proxy_type is not allowed on this node"})
		return
	}
	if proxy.ProxyType == "tcp" || proxy.ProxyType == "udp" {
		// 检查 选择的 remote_port 是否在 nodes 表中的 allow_port 字段中 allow_port 字段是 起始端口-结束端口 组成的字符串 用分号分隔 例如 1000-2000;3000-4000;
		// 将 proxy.RemotePort 转换成 int 类型
		remotePort, _ := strconv.Atoi(proxy.RemotePort)
		if !CheckAllowPort(node.AllowPort, remotePort) {
			c.JSON(400, gin.H{"message": "remote_port is not allowed on this node"})
			return
		}
		// 检查 proxies 表中 用户选择的 node 的对应端口是否已经被占用
		var proxyCount2 int64
		db.Model(&Proxy{}).Where("node = ? AND remote_port = ?", proxy.Node, proxy.RemotePort).Count(&proxyCount2)
		if proxyCount2 > 0 {
			c.JSON(400, gin.H{"message": "remote port already exists"})
			return
		}
	} else {
		// 检查 proxies 表中 用户选择的 node 的对应域名是否已经被占用
		var proxyCount3 int64
		domain := "[\"" + proxy.Domain + "\"]"
		db.Model(&Proxy{}).Where("node = ? AND domain = ?", proxy.Node, domain).Count(&proxyCount3)
		if proxyCount3 > 0 {
			c.JSON(400, gin.H{"message": "domain already exists"})
			return
		}
		proxy.Domain = "[\"" + proxy.Domain + "\"]"
	}
	proxy.Lastupdate = time.Now().Unix()
	// 将 proxy 写入 proxies
	result := db.Table("proxies").Create(&proxy)
	if result.Error != nil {
		c.JSON(500, gin.H{"message": result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
func CheckGroup(nodeGroup string, userGroup string) bool {
	// 不去除分号 直接判断是否包含
	if strings.Contains(nodeGroup, userGroup) {
		return true
	} else {
		return false
	}
}
func CheckAllowType(nodeAllowType string, proxyType string) bool {
	// 不去除分号 直接判断是否包含
	if strings.Contains(nodeAllowType, proxyType) {
		return true
	} else {
		return false
	}
}
func CheckAllowPort(nodeAllowPort string, remotePort int) bool {
	ports := strings.Split(nodeAllowPort, ";")

	for _, port := range ports {
		ranges := strings.Split(port, "-")

		startPort, err := strconv.Atoi(ranges[0])
		if err != nil {
			continue
		}

		endPort := startPort
		if len(ranges) > 1 {
			endPort, err = strconv.Atoi(ranges[1])
			if err != nil {
				continue
			}
		}

		if remotePort >= startPort && remotePort <= endPort {
			return true
		}
	}

	return false
}
func IsvalidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
func IsvalidPort(port string) bool {
	if _, err := strconv.Atoi(port); err != nil {
		return false
	}
	return true
}
func IsValidDomain(domain string) bool {
	// 使用正则表达式匹配域名格式
	// regex := `/^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(?:\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$/`
	// match, _ := regexp.MatchString(regex, domain)
	return true
}
