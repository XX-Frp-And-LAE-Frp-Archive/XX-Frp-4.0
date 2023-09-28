package cron

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

func UpdateTrafficPeriodically(db *gorm.DB) {
	// 启动定时器，在每天晚上24:00更新traffic为0
	ticker := time.NewTicker(24 * time.Hour)
	for {
		now := <-ticker.C
		if now.Hour() == 0 {
			// 更新todaytraffic表中的所有行的traffic为0
			result := db.Model(&TodayTraffic{}).Update("traffic", 0)
			if result.Error != nil {
				fmt.Printf("清空流量数据错误: %v\n", result.Error)
			} else {
				fmt.Println("流量数据清空成功")
			}
		}
	}
}
