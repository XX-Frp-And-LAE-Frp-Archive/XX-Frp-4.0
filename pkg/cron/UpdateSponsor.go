package cron

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type Sponsor struct {
	ID      int    `gorm:"primaryKey" json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Thing   string `json:"thing"`
	Comment string `json:"comment"`
}

var Cache []Sponsor

// 从数据库中更新数据到缓存
func UpdateDataPeriodically(db *gorm.DB) {
	for {
		var sponsors []Sponsor
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
