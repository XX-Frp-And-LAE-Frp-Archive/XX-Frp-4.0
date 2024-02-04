package cron

import (
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"gorm.io/gorm"
	"log"
	"time"
)

// CacheData 创建一个缓存变量
var CacheData interface{}

func UpdateStatisticsPeriodically(db *gorm.DB) {
	for {

		// 查询 users 表的行数
		var userCount int64
		if err := db.Model(&define.User{}).Count(&userCount).Error; err != nil {
			fmt.Println("Failed to query user count:", err)
			continue
		}

		// 查询 proxies 表的行数
		var proxyCount int64
		if err := db.Model(&define.Proxies{}).Count(&proxyCount).Error; err != nil {
			fmt.Println("Failed to query proxy count:", err)
			continue
		}

		// 查询 nodes 表的行数
		var nodeCount int64
		if err := db.Model(&define.Node{}).Count(&nodeCount).Error; err != nil {
			fmt.Println("Failed to query node count:", err)
			continue
		}

		// 查询nodes表
		var nodes []define.Node
		if err := db.Find(&nodes).Error; err != nil {
			fmt.Println("Failed to query nodes:", err)
			continue
		}
		//对每一个节点的TotalTrafficIn TotalTrafficout进行求和
		var sum int64
		for _, node := range nodes {
			sum += node.TotalTrafficIn
			sum += node.TotalTrafficOut
		}

		// 将查询结果写入缓存
		CacheData = map[string]int64{
			"userCount":  userCount,
			"proxyCount": proxyCount,
			"nodeCount":  nodeCount,
			"trafficSum": sum,
		}
		log.Printf("Updated Statistics data from database.")
		time.Sleep(5 * time.Minute)
	}
}
