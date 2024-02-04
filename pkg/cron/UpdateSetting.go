package cron

import (
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"gorm.io/gorm"
	"log"
	"time"
)

var SettingCache []define.Setting

// UpdateSettingPeriodically 从数据库中更新数据到缓存
func UpdateSettingPeriodically(db *gorm.DB) {
	for {
		var settings []define.Setting
		result := db.Find(&settings)
		if result.Error != nil {
			fmt.Println("Failed to retrieve data from database:", result.Error)
			return
		}

		// 更新缓存
		SettingCache = settings

		log.Printf("Updated Setting data from database.")
		time.Sleep(5 * time.Minute)
	}
}
