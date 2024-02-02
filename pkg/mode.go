package pkg

import (
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"log"
)

var LogMode logger.LogLevel

func SetMode() {
	// 初始化配置
	conf := config.GetConfig()

	// 模式切换
	if conf.Debug.Debug == true {
		gin.SetMode(gin.DebugMode)
		LogMode = logger.Info
		log.Printf("当前模式：Debug")
	} else {
		gin.SetMode(gin.ReleaseMode)
		LogMode = logger.Silent
		log.Printf("当前模式：Release")
	}
}
