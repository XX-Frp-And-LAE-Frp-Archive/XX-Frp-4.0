package mysql

import (
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectMySQL() *gorm.DB {
	// 初始化配置
	conf := config.GetConfig()

	dsn := conf.Mysql.User + ":" + conf.Mysql.Password + "@tcp(" + conf.Mysql.Host + ":" + conf.Mysql.Port + ")/" + conf.Mysql.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to connect to MySQL database")
	}
	return db
}
