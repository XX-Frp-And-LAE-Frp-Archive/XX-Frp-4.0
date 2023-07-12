package main

import (
	"flag"
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg"
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	"github.com/ahmr-bot/ME-Frp/pkg/mysql"
	"github.com/ahmr-bot/ME-Frp/pkg/router"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

var (
	ConfigPath string
)

var (
	db *gorm.DB
)

func init() {
	flag.StringVar(&ConfigPath, "c", "config.toml", "配置文件路径")
	flag.Parse()
	err := config.LoadConfig(ConfigPath)
	if err != nil {
		panic(err)
	}
}

func main() {
	// 初始化配置
	conf := config.GetConfig()

	// 连接MySQL数据库
	db = mysql.ConnectMySQL()

	// 设置gin模式
	pkg.SetMode()

	// 创建 Gin 引擎
	engine := gin.Default()

	// 加载所有路由
	router.LoadRoutes(engine, db)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	log.Printf("服务器已启动，监听地址：" + conf.Server.Host + ":" + fmt.Sprint(conf.Server.Port))
	if err := engine.Run(addr); err != nil {
		panic(err)
	}
}
