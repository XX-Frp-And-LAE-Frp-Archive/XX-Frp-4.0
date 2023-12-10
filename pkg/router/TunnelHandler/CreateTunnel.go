package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func HandleCreateTunnel(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	// 从表单中获取数据
	var proxy define.Proxies
	proxy.Username = user.Username
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
	// 检查proxy_name 是否合法
	if !CheckProxyName(proxy.ProxyName) {
		respond.Respond(c, 403, "隧道名只能是 1-16为的字母数字", 0)
		return
	}
	// 检查proxytype是否合法
	if proxy.ProxyType != "tcp" && proxy.ProxyType != "udp" && proxy.ProxyType != "http" && proxy.ProxyType != "https" {
		respond.Respond(c, 403, "隧道类型不支持", 0)
		return
	}
	// 使用正则表达式检查 localip 是不是 ip地址
	// 使用正则表达式检查 localport 是不是 端口号
	// 使用正则表达式检查 domain 是不是域名
	// 使用正则表达式检查 remoteport 是不是 端口号
	if !IsvalidIP(proxy.LocalIP) {
		respond.Respond(c, 403, "本地IP不合法", 0)
		return
	}
	if !IsvalidPort(proxy.LocalPort) {
		respond.Respond(c, 403, "本地端口不合法", 0)
		return
	}
	if proxy.ProxyType == "http" || proxy.ProxyType == "https" {
		if !IsValidDomain() {
			respond.Respond(c, 403, "域名不合法", 0)
			return
		}
	} else {
		if !IsvalidPort(proxy.RemotePort) {
			respond.Respond(c, 403, "远程端口不合法", 0)
			return
		}
	}

	// 检查 单用户的代理名称是否重复
	var proxyCount int64
	db.Model(&define.Proxies{}).Where("username = ? AND proxy_name = ?", user.Username, proxy.ProxyName).Count(&proxyCount)
	if proxyCount > 0 {
		respond.Respond(c, 403, "隧道名重复", 0)
		return
	}

	// 获取用户的 group 值
	var group define.Groups
	groupResult := db.First(&group, "name = ?", user.Group)
	if groupResult.Error != nil {
		respond.Respond(c, 500, "无法查验组别", 0)
		return
	}
	//  limits 表中查询该用户是否有自定义的限制
	var limit define.Limit
	limitResult := db.First(&limit, "username = ?", user.Username)
	if limitResult.Error != nil {
		limit.Proxies = group.Proxies
	}
	var proxies []define.Proxies
	db.Find(&proxies, "username = ?", user.Username)
	if len(proxies) >= limit.Proxies {
		respond.Respond(c, 403, "您的隧道数已上限", 0)
		return
	}
	// 检查 node 是否存在
	var nodeCount int64
	db.Model(&define.Node{}).Where("id = ?", proxy.Node).Count(&nodeCount)
	if nodeCount == 0 {
		respond.Respond(c, 403, "节点不存在", 0)
		return
	}
	// 检查 nodes 表中对应的group值中是否包含该用户的group值
	var node define.Node
	nodeResult := db.First(&node, "id = ?", proxy.Node)
	if nodeResult.Error != nil {
		respond.Respond(c, 500, "无法查验节点", 0)
		return
	}
	// node 的 group 值是多个值组成的字符串 用分号分隔 例如 admin;default;realname;trustuser;
	// 检查 node 的 group 值中是否包含该用户的 group 值
	if !CheckGroup(node.Group, user.Group) {
		respond.Respond(c, 403, "你没有使用该节点的权限", 0)
		return
	}

	// 检查 选择的 proxy_type 是否在 nodes 表中的 allow_type 字段中 allow_type 字段是多个值组成的字符串 用分号分隔 例如 tcp;udp;http;https;
	if !CheckAllowType(node.AllowType, proxy.ProxyType) {
		respond.Respond(c, 403, "您提交的隧道类型在该节点不支持", 0)
		return
	}
	if proxy.ProxyType == "tcp" || proxy.ProxyType == "udp" {
		// 检查 选择的 remote_port 是否在 nodes 表中的 allow_port 字段中 allow_port 字段是 起始端口-结束端口 组成的字符串 用分号分隔 例如 1000-2000;3000-4000;
		// 将 proxy.RemotePort 转换成 int 类型
		remotePort, _ := strconv.Atoi(proxy.RemotePort)
		if !CheckAllowPort(node.AllowPort, remotePort) {
			respond.Respond(c, 403, "您提交的远程端口在该节点已被使用", 0)
			return
		}
		// 检查 proxies 表中 用户选择的 node 的对应端口是否已经被占用
		var proxyCount2 int64
		db.Model(&define.Proxies{}).Where("node = ? AND remote_port = ? proxy_type = ?", proxy.Node, proxy.RemotePort, proxy.ProxyType).Count(&proxyCount2)
		if proxyCount2 > 0 {
			respond.Respond(c, 403, "该远程端口已经被占用了哦", 0)
			return
		}
	} else {
		// 检查 proxies 表中 用户选择的 node 的对应域名是否已经被占用
		var proxyCount3 int64
		domain := "[\"" + proxy.Domain + "\"]"
		db.Model(&define.Proxies{}).Where("node = ? AND domain = ? AND proxy_type = ?", proxy.Node, domain, proxy.ProxyType).Count(&proxyCount3)
		if proxyCount3 > 0 {
			respond.Respond(c, 403, "该域名已被占用", 0)
			return
		}
		proxy.Domain = "[\"" + proxy.Domain + "\"]"
	}
	proxy.Lastupdate = time.Now().Unix()
	// 将 proxy 写入 proxies
	result := db.Table("proxies").Create(&proxy)
	if result.Error != nil {
		respond.Respond(c, 500, "写入失败", 0)
		return
	}
	respond.Respond(c, 200, "创建成功", 0)
}
func CheckGroup(nodeGroup string, userGroup string) bool {
	// 不去除分号 直接判断是否包含
	if strings.Contains(nodeGroup, userGroup) {
		return true
	} else {
		return false
	}
}

// 检查 proxyname 是不是字母数字组成
func CheckProxyName(proxyname string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{1,16}$", proxyname); ok {
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
func IsValidDomain() bool {
	// 使用正则表达式匹配域名格式
	// regex := `/^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(?:\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$/`
	// match, _ := regexp.MatchString(regex, domain)
	return true
}
