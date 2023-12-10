package cron

import (
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"gorm.io/gorm"
	"log"
	"time"
)

var Cache []define.Sponsor

// 从数据库中更新数据到缓存
func UpdateDataPeriodically(db *gorm.DB) {
	for {
		var sponsors []define.Sponsor
		result := db.Find(&sponsors)
		if result.Error != nil {
			fmt.Println("Failed to retrieve data from database:", result.Error)
			return
		}

		// 更新缓存
		Cache = sponsors

		log.Print("Updated sponsors data from database.")
		time.Sleep(5 * time.Minute)
	}
}
