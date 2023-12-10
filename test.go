package main

import (
	"encoding/base64"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Realname struct {
	ID     uint   `gorm:"column:id"`
	Name   string `gorm:"column:name"`
	IDCard string `gorm:"column:id_card"`
}

func main() {
	dsn := "mefrp:VW8B7k2jhKOD18dy@tcp(36.133.110.35:3306)/mefrp"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败:", err)
		return
	}

	var realnames []Realname
	err = db.Find(&realnames).Error
	if err != nil {
		fmt.Println("查询数据失败:", err)
		return
	}

	for _, realname := range realnames {
		encryptedName := base64.StdEncoding.EncodeToString([]byte(realname.Name))
		encryptedIDCard := base64.StdEncoding.EncodeToString([]byte(realname.IDCard))

		err = db.Model(&realname).Updates(map[string]interface{}{
			"name":    encryptedName,
			"id_card": encryptedIDCard,
		}).Error
		if err != nil {
			fmt.Println("更新数据失败:", err)
			// 可以选择在这里进行错误处理，或者跳过当前行继续处理下一行
			continue
		}
	}

	fmt.Println("数据加密并更新完成")
}
