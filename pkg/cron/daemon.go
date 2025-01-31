package cron

import (
	"encoding/json"
	"errors"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"gorm.io/gorm"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// 403 禁用
// 200 正常
// 201 无状态数据
// 500 服务器异常

const (
	StatusDisabled    = 403
	StatusNormal      = 200
	StatusServerError = 500
)

var timeout = 5 * time.Second

// FetchData 发送HTTP GET请求到指定的节点和路径
func FetchData(node define.Node, path string) (*http.Response, error) {
	// 使用备案域 作为请求地址，防止防火墙拦截
	url := "http://admin:" + node.AdminPass + "@node.mefrp.com" + ":" + node.AdminPort + "/api" + path
	// 自定义Dial函数
	dialFunc := func(network, addr string) (net.Conn, error) {
		return net.Dial("tcp", node.Hostname+":"+node.AdminPort)
	}

	// 创建自定义Transport
	transport := &http.Transport{
		Dial: dialFunc,
	}

	// 创建自定义Client
	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	// 发送请求
	resp, err := client.Get(url)
	return resp, err
}

// FetchTraffic  发送HTTP GET请求到所有节点的server_info
func FetchTraffic(db *gorm.DB) {
	// nodes 表获取所有节点信息
	var nodes []define.Node
	result := db.Find(&nodes)
	if result.Error != nil {
		log.Println("获取节点信息失败:", result.Error)
		return
	}
	// 遍历所有节点
	for _, node := range nodes {
		// 跳过 禁用 的节点
		if node.Status == StatusDisabled {
			log.Println(node.Name + "节点已禁用")
			continue
		}
		log.Println("正在获取" + node.Name + "节点信息")
		// 使用 FetchData 发送请求 设置5s超时
		resp, err := FetchData(node, "/serverinfo")
		if err != nil {
			log.Println("请求失败:", err)
			// 更新节点状态
			node.Status = StatusServerError
			result := db.Save(&node)
			if result.Error != nil {
				log.Println("更新节点状态失败:", result.Error)
			}
			continue
		}
		// 关闭响应体
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println("关闭响应体失败:", err)
			}
		}(resp.Body)
		// 检查响应状态码是否为200
		if resp.StatusCode != StatusNormal {
			log.Println("请求失败, 状态码:", resp.StatusCode)
			// 更新节点状态
			node.Status = StatusServerError
			result := db.Save(&node)
			if result.Error != nil {
				log.Println("更新节点状态失败:", result.Error)
			}
			continue
		}
		// 打印响应
		log.Println("获取" + node.Name + "节点信息成功:")
		// 更新节点状态为正常
		result = db.Model(&node).Update("status", StatusNormal)
		if result.Error != nil {
			log.Println("更新"+node.Name+"节点状态失败", result.Error)
		}
		// 使用 define.FrpsData 结构体解析响应
		var frpsData define.FrpsData
		err = json.NewDecoder(resp.Body).Decode(&frpsData)
		if err != nil {
			log.Println("解析"+node.Name+"节点信息失败:", err)
			continue
		}
		// 判断 node.TotalTrafficIn 和 node.TotalTrafficOut 是否分别大于 frpsData.TotalTrafficIn 和 frpsData.TotalTrafficOut
		if node.TotalTrafficIn > frpsData.TotalTrafficIn || node.TotalTrafficOut > frpsData.TotalTrafficOut {
			// 如果大于 则代表 frps 重启过 需要给用户结算流量
			// 获取该节点的所有隧道信息
			var proxies []define.Proxies
			result := db.Find(&proxies, "node = ?", node.ID)
			if result.Error != nil {
				log.Println("查询数据库出错:", result.Error)
				continue
			}
			// 遍历所有隧道
			for _, proxy := range proxies {
				// 如果该隧道的 in out 流量为0 则跳过
				if proxy.TodayTrafficIn == 0 && proxy.TodayTrafficOut == 0 {
					continue
				}
				// 读取该隧道的用户名 和 in out 流量
				username := proxy.Username
				todayTrafficIn := proxy.TodayTrafficIn
				todayTrafficOut := proxy.TodayTrafficOut
				totalTrafficb := todayTrafficIn + todayTrafficOut
				// 计算mb
				totalTrafficmb := totalTrafficb / 1024 / 1024
				// 清空对应的数据
				proxy.TodayTrafficOut = 0
				proxy.TodayTrafficIn = 0
				result = db.Save(proxy)
				// 扣除用户流量
				var user define.User
				result := db.First(&user, "username = ?", username)
				if result.Error != nil {
					log.Println("查询数据库出错:", result.Error)
					continue
				}
				user.Traffic -= totalTrafficmb
				result = db.Save(&user)
			}
		}
		// 将 ProxyTypeCount 中tcp udp http https的数量相加
		node.OnlineCount = int64(frpsData.ProxyTypeCount.Tcp + frpsData.ProxyTypeCount.Udp + frpsData.ProxyTypeCount.Http + frpsData.ProxyTypeCount.Https)
		node.Version = frpsData.Version
		node.TotalTrafficIn = frpsData.TotalTrafficIn
		node.TotalTrafficOut = frpsData.TotalTrafficOut
		// 更新节点信息
		result = db.Save(&node)
		if result.Error != nil {
			log.Println("更新"+node.Name+"节点信息失败", result.Error)
		}
	}
	log.Printf("Updated Node data from database.")
}
func FetchUserTraffic(db *gorm.DB) {

	// 预加载所有用户信息到map中
	userMap := make(map[string]define.User)
	var users []define.User
	if err := db.Find(&users).Error; err != nil {
		log.Println("获取用户信息失败:", err)
		return
	}
	for _, user := range users {
		userMap[user.Token] = user
	}

	// 预加载所有隧道信息到map中，键为用户名+隧道名
	proxiesMap := make(map[string]define.Proxies)
	var proxiesList []define.Proxies
	if err := db.Find(&proxiesList).Error; err != nil {
		log.Println("获取隧道信息失败:", err)
		return
	}
	for _, proxy := range proxiesList {
		key := proxy.Username + "." + proxy.ProxyName
		proxiesMap[key] = proxy
	}

	// 遍历所有节点
	var updates []define.Proxies // 存储需要更新的隧道信息
	// 先将所有隧道的状态设置为离线 要记得写where条件
	if err := db.Model(&define.Proxies{}).Updates(map[string]interface{}{"online": 0}).Error; err != nil {
		log.Println("更新隧道状态失败:", err)
	}
	// 获取所有节点信息
	var nodes []define.Node
	if err := db.Find(&nodes).Error; err != nil {
		log.Println("获取节点信息失败:", err)
		return
	}
	for _, node := range nodes {
		if node.Status != StatusNormal {
			continue
		}

		allowTypeSlice := strings.Split(node.AllowType, ";")
		for _, allowType := range allowTypeSlice {
			if allowType == "" {
				continue
			}

			resp, err := FetchData(node, "/proxy/"+allowType)
			if err != nil {
				log.Println("请求失败:", err)
				continue
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Println("关闭响应体失败:", err)
				}
			}(resp.Body)

			if resp.StatusCode != StatusNormal {
				log.Println("请求失败, 状态码:", resp.StatusCode)
				continue
			}

			var frpsTraffic define.FrpsTrafficData
			if err = json.NewDecoder(resp.Body).Decode(&frpsTraffic); err != nil {
				log.Println("解析流量数据失败:", err)
				continue
			}

			for _, traffic := range frpsTraffic.Proxies {
				nameSlice := strings.Split(traffic.Name, ".")
				tunnelName, token := nameSlice[1], nameSlice[0]

				user, ok := userMap[token]
				if !ok {
					log.Println("未找到用户:", token)
					continue
				}
				key := user.Username + "." + tunnelName
				proxy, ok := proxiesMap[key]
				if !ok {
					log.Println("未找到隧道:", key)
					continue
				}

				// 更新流量数据
				proxy.TodayTrafficIn = traffic.TodayTrafficIn
				proxy.TodayTrafficOut = traffic.TodayTrafficOut
				proxy.CurConns = traffic.CurConns
				proxy.ClientVersion = traffic.ClientVersion
				proxy.Online = traffic.Status
				proxy.LastCloseTime = traffic.LastCloseTime
				proxy.LastStartTime = traffic.LastStartTime
				updates = append(updates, proxy)
			}
		}
	}

	// 执行批量更新
	for _, update := range updates {
		db.Save(&update)
	}
	log.Printf("Fetched User Traffic.")
}

// 下线该用户的所有隧道
func offlineUserProxies(db *gorm.DB, user define.User) {
	// 获取该用户的所有隧道
	var proxies []define.Proxies
	result := db.Find(&proxies, "username = ?", user.Username)
	if result.Error != nil {
		log.Println("查询数据库出错:", result.Error)
		return
	}

	// 遍历所有隧道
	for _, proxy := range proxies {
		// 下线该隧道
		path := "/client/close/" + proxy.RunID

		// 获取该隧道的节点信息
		var node define.Node
		nodeResult := db.First(&node, proxy.Node)
		if nodeResult.Error != nil {
			log.Println("未找到该节点:", nodeResult.Error)
			continue
		}
		resp, err := FetchData(node, path)
		if err != nil {
			log.Printf("关闭隧道失败，服务器错误")
			return
		}
		// 判断返回的状态码 如果不是200 则返回错误
		if resp.StatusCode != 200 {
			log.Printf("关闭隧道失败，客户端没有启动")
			return
		}
	}
}

func CalculateUserTraffic(db *gorm.DB) {
	// 获取所有用户的用户名
	var users []define.User
	if err := db.Model(&define.User{}).Find(&users).Error; err != nil {
		log.Println("获取用户信息失败:", err)
		return
	}

	// 获取所有隧道信息，并按照用户名分组累计流量
	type ProxyTraffic struct {
		Username        string
		TotalTrafficIn  int64
		TotalTrafficOut int64
	}
	var proxyTraffic []ProxyTraffic
	if err := db.Model(&define.Proxies{}).
		Select("username, SUM(today_traffic_in) AS total_traffic_in, SUM(today_traffic_out) AS total_traffic_out").
		Group("username").
		Find(&proxyTraffic).Error; err != nil {
		log.Println("查询隧道流量失败:", err)
		return
	}

	// 构建一个映射，方便查找每个用户的总流量
	userTrafficMap := make(map[string]int64)
	for _, pt := range proxyTraffic {
		userTrafficMap[pt.Username] = pt.TotalTrafficIn + pt.TotalTrafficOut
	}

	// 遍历所有用户，更新或创建todaytraffic记录
	for _, user := range users {
		totalTraffic, exists := userTrafficMap[user.Username]
		if totalTraffic == 0 || !exists {
			continue // 如果用户没有隧道流量，跳过
		}

		var todayTraffic define.TodayTraffic
		result := db.First(&todayTraffic, "user = ?", user.Username)
		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 创建新记录
			todayTraffic = define.TodayTraffic{
				User:    user.Username,
				Traffic: totalTraffic,
			}
			db.Create(&todayTraffic)
		} else if result.Error == nil {
			todayTraffic.Traffic = totalTraffic
			// 保存记录
			db.Where("user = ?", user.Username).Save(&todayTraffic)
		} else {
			log.Println("查询todaytraffic失败:", result.Error)
		}
	}
}

// 更新用户流量
func UpdateUserTraffic(db *gorm.DB) {
	// 读取users表中的所有数据 并构建一个能用username作为键的map traffic status为值
	var userTraffic = make(map[string]int64)
	var users []define.User
	if err := db.Find(&users).Error; err != nil {
		log.Println("获取用户信息失败:", err)
		return
	}

	var userStatus = make(map[string]int)
	for _, user := range users {
		userTraffic[user.Username] = user.Traffic
		userStatus[user.Username] = user.Status
		// 如果剩余流量大于0 且status为2 则更新status为0 恢复隧道启动功能
		if user.Traffic > 0 && user.Status == 2 {
			if err := db.Model(&define.User{}).Where("username = ?", user.Username).Update("status", 0).Error; err != nil {
				log.Println("更新用户状态失败:", err)
			}
		}
	}

	// 读取todaytraffic表中的所有数据 并遍历
	var todayTraffics []define.TodayTraffic
	if err := db.Find(&todayTraffics).Error; err != nil {
		log.Println("获取todaytraffic信息失败:", err)
		return
	}
	for _, todayTraffic := range todayTraffics {
		// 计算剩余流量
		totalTrafficMB := float64(userTraffic[todayTraffic.User])
		usedTrafficMB := float64(todayTraffic.Traffic) / 1024.0 / 1024.0
		remainTraffic := int64(totalTrafficMB - usedTrafficMB)
		// 如果用户号被封 则跳过
		if userStatus[todayTraffic.User] == 1 {
			continue
		}
		if userStatus[todayTraffic.User] == 2 {
			continue
		}
		// 如果剩余流量小于<=0 则下线该用户的所有隧道
		if remainTraffic <= 0 {
			// 下线该用户的所有隧道
			offlineUserProxies(db, define.User{Username: todayTraffic.User})
			// 重置该用户的流量
			userTraffic[todayTraffic.User] = 0
			// 更新该用户的流量为0 status为2
			if err := db.Model(&define.User{}).Where("username = ?", todayTraffic.User).Update("traffic", 0).Error; err != nil {
				log.Println("更新用户流量失败:", err)
			}
			if err := db.Model(&define.User{}).Where("username = ?", todayTraffic.User).Update("status", 2).Error; err != nil {
				log.Println("更新用户状态失败:", err)
			}
			// 清空该用户的流量记录
			if err := db.Where("user = ?", todayTraffic.User).Delete(&todayTraffic).Error; err != nil {
				log.Println("清空用户流量记录失败:", err)
			}
			continue
		}
	}
	log.Printf("Updated User Traffic.")
}

// ClearUserTraffic 清空用户流量
func ClearUserTraffic(db *gorm.DB) {
	// 清空 today_traffic 表
	if err := db.Exec("TRUNCATE TABLE today_traffic").Error; err != nil {
		log.Println("清空用户流量失败:", err)
	}
	// 清空 proxies 表 today_traffic_in today_traffic_out
	if err := db.Model(&define.Proxies{}).Updates(map[string]interface{}{"today_traffic_in": 0, "today_traffic_out": 0}).Error; err != nil {
		log.Println("清空隧道流量失败:", err)
	}
}

// SettleUserTraffic 结算用户流量
func SettleUserTraffic(db *gorm.DB) {
	// 读取users表中的所有数据 并构建一个能用username作为键的map traffic status为值
	var userTraffic = make(map[string]int64)
	var users []define.User
	if err := db.Find(&users).Error; err != nil {
		log.Println("获取用户信息失败:", err)
		return
	}

	var userStatus = make(map[string]int)
	for _, user := range users {
		userTraffic[user.Username] = user.Traffic
		userStatus[user.Username] = user.Status
	}

	// 读取todaytraffic表中的所有数据 并遍历
	var todayTraffics []define.TodayTraffic
	if err := db.Find(&todayTraffics).Error; err != nil {
		log.Println("获取todaytraffic信息失败:", err)
		return
	}
	for _, todayTraffic := range todayTraffics {
		// 计算剩余流量
		totalTrafficMB := float64(userTraffic[todayTraffic.User])
		usedTrafficMB := float64(todayTraffic.Traffic) / 1024.0 / 1024.0
		remainTraffic := int64(totalTrafficMB - usedTrafficMB)

		// 如果用户号被封 则跳过
		if userStatus[todayTraffic.User] == 1 {
			continue
		}

		if userStatus[todayTraffic.User] == 2 {
			continue
		}

		if remainTraffic <= 0 {
			if err := db.Model(&define.User{}).Where("username = ?", todayTraffic.User).Update("traffic", 0).Error; err != nil {
				log.Println("更新用户流量失败:", err)
			}
		}
		if err := db.Model(&define.User{}).Where("username = ?", todayTraffic.User).Update("traffic", remainTraffic).Error; err != nil {
			log.Println("更新用户流量失败:", err)
		}
	}
}
