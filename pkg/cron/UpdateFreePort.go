package cron

import (
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 初始化一个全局的缓存实例，设置默认的过期时间和清理间隔时间
var c = cache.New(15*time.Minute, 20*time.Minute)

func UpdateFreePort(db *gorm.DB) {
	for {
		var nodes []define.Node
		db.Find(&nodes)
		for _, node := range nodes {
			protocols := strings.Split(node.AllowType, ";") // 假设 AllowType 是以逗号分隔的字符串
			for _, protocol := range protocols {
				updateNodePorts(db, node, protocol)
			}
		}
		log.Print("Updated Free Port.")
		time.Sleep(10 * time.Minute)
	}
}

func updateNodePorts(db *gorm.DB, node define.Node, protocol string) {
	var proxies []define.Proxies
	result := db.Find(&proxies, "node = ? AND proxy_type = ?", node.ID, protocol)
	if result.Error != nil {
		fmt.Println("查询数据库出错:", result.Error)
		return
	}
	availablePorts := parseAllowPort(node.AllowPort)
	availablePorts, _ = filterUsedPorts(availablePorts, proxies)

	// 构造一个键名，格式为 "nodeID:protocol"，例如 "123:tcp"
	cacheKey := fmt.Sprintf("%d:%s", node.ID, protocol)
	// 将数据写入缓存
	c.Set(cacheKey, availablePorts, cache.NoExpiration)
}

// 解析 allow_port 字符串，生成所有可用端口的列表
func parseAllowPort(allowPort string) []int {
	var ports []int
	ranges := strings.Split(allowPort, ";")
	for _, r := range ranges {
		if r == "" {
			continue
		}
		bounds := strings.Split(r, "-")
		if len(bounds) != 2 {
			continue
		}
		start, err1 := strconv.Atoi(bounds[0])
		end, err2 := strconv.Atoi(bounds[1])
		if err1 != nil || err2 != nil {
			continue // 出错则跳过此范围
		}
		for i := start; i <= end; i++ {
			ports = append(ports, i)
		}
	}
	return ports
}

// getUsedPorts 遍历 proxies 列表，提取并返回所有已使用的端口号
func getUsedPorts(proxies []define.Proxies) ([]int, error) {
	var usedPorts []int // 用于存储已使用端口号的切片

	for _, proxy := range proxies {
		port, err := strconv.Atoi(proxy.RemotePort) // 将字符串类型的端口号转换为整数
		if err != nil {
			return nil, err // 如果转换失败，返回错误
		}
		usedPorts = append(usedPorts, port) // 将转换后的端口号添加到结果切片中
	}

	return usedPorts, nil // 返回已使用端口号的切片和 nil 错误
}

type Proxies struct {
	RemotePort string // 假设 RemotePort 是一个字符串类型
}

// filterUsedPorts 移除已使用的端口
func filterUsedPorts(availablePorts []int, proxies []define.Proxies) ([]int, error) {
	usedPortsSlice, err := getUsedPorts(proxies) // 获取已使用的端口列表

	if err != nil {
		return nil, err // 如果获取过程中出现错误，直接返回错误
	}

	usedPortsMap := make(map[int]bool) // 创建映射以快速检查端口是否已使用
	for _, port := range usedPortsSlice {
		usedPortsMap[port] = true
	}

	var result []int // 存储未被使用的端口
	for _, port := range availablePorts {
		if !usedPortsMap[port] { // 检查端口是否未被使用
			result = append(result, port)
		}
	}

	return result, nil // 返回未被使用的端口列表和 nil 错误
}

// GetFreePort 从缓存中随机获取一个可用端口并返回
func GetFreePort(nodeID int, protocol string) (int, error) {
	cacheKey := fmt.Sprintf("%d:%s", nodeID, protocol)
	value, found := c.Get(cacheKey)
	if !found {
		return 0, fmt.Errorf("没有找到与节点ID %d 和协议 %s 相关的可用端口", nodeID, protocol)
	}

	availablePorts, ok := value.([]int)
	if !ok || len(availablePorts) == 0 {
		return 0, fmt.Errorf("节点ID %d 和协议 %s 的可用端口列表为空或格式错误", nodeID, protocol)
	}

	// 随机选择一个端口的索引
	// 使用当前时间的纳秒作为种子值创建一个新的随机数生成器
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	index := r.Intn(len(availablePorts))
	selectedPort := availablePorts[index]

	return selectedPort, nil
}

// RemovePortFromCache 从缓存中删除指定节点和类型对应的特定端口
func RemovePortFromCache(nodeID int, protocol string, port int) error {
	cacheKey := fmt.Sprintf("%d:%s", nodeID, protocol)
	value, found := c.Get(cacheKey)
	if !found {
		return fmt.Errorf("没有找到与节点ID %d 和协议 %s 相关的可用端口", nodeID, protocol)
	}

	availablePorts, ok := value.(map[int]bool)
	if !ok {
		return fmt.Errorf("节点ID %d 和协议 %s 的可用端口列表格式错误", nodeID, protocol)
	}

	// 直接删除指定的端口
	if _, exists := availablePorts[port]; exists {
		delete(availablePorts, port)
		// 更新缓存中的可用端口列表
		c.Set(cacheKey, availablePorts, cache.NoExpiration)
	} else {
		return fmt.Errorf("端口 %d 不在节点ID %d 和协议 %s 的可用端口列表中", port, nodeID, protocol)
	}

	return nil
}
