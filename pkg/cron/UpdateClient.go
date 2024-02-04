package cron

import (
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"gorm.io/gorm"
	"log"
	"time"
)

var ClientCache []define.Client

// UpdateClientPeriodically 从数据库中更新数据到缓存
func UpdateClientPeriodically(db *gorm.DB) {
	for {
		var Client []define.Client
		result := db.Find(&Client)
		if result.Error != nil {
			fmt.Println("Failed to retrieve data from database:", result.Error)
			return
		}

		// 更新缓存
		ClientCache = Client

		log.Printf("Updated Client data from database.")
		time.Sleep(5 * time.Minute)
	}
}
