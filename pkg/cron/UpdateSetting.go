package cron

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type Setting struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var SettingCache []Setting

// 从数据库中更新数据到缓存
func UpdateSettingPeriodically(db *gorm.DB) {
	for {
		var settings []Setting
		result := db.Find(&settings)
		if result.Error != nil {
			fmt.Println("Failed to retrieve data from database:", result.Error)
			return
		}

		// 更新缓存
		SettingCache = settings

		log.Print("Updated Setting data from database.")
		time.Sleep(5 * time.Minute)
	}
}
