package cron

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	ID   int
	Name string
	// 其他字段...
}

type Proxy struct {
	ID   int
	Name string
	// 其他字段...
}

type Node struct {
	ID   int
	Name string
	// 其他字段...
}

// 创建一个缓存变量
var CacheData interface{}

func UpdateStatisticsPeriodically(db *gorm.DB) {
	for {

		// 查询 users 表的行数
		var userCount int64
		if err := db.Model(&User{}).Count(&userCount).Error; err != nil {
			fmt.Println("Failed to query user count:", err)
			continue
		}

		// 查询 proxies 表的行数
		var proxyCount int64
		if err := db.Model(&Proxy{}).Count(&proxyCount).Error; err != nil {
			fmt.Println("Failed to query proxy count:", err)
			continue
		}

		// 查询 nodes 表的行数
		var nodeCount int64
		if err := db.Model(&Node{}).Count(&nodeCount).Error; err != nil {
			fmt.Println("Failed to query node count:", err)
			continue
		}

		// 查询 traffic 列的和
		var trafficSum int64
		if err := db.Model(&User{}).Select("SUM(traffic)").Scan(&trafficSum).Error; err != nil {
			fmt.Println("Failed to query traffic sum:", err)
			continue
		}

		// 将查询结果写入缓存
		CacheData = map[string]int64{
			"userCount":  userCount,
			"proxyCount": proxyCount,
			"nodeCount":  nodeCount,
			"trafficSum": trafficSum,
		}
		log.Print("Updated Statistics data from database.")
		time.Sleep(5 * time.Minute)
	}
}
