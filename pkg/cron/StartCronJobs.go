package cron

import (
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// StartCronJobs 启动所有定时任务
func StartCronJobs(db *gorm.DB) {
	// 加载统计数据
	go UpdateStatisticsPeriodically(db)
	// 加载赞助者
	go UpdateDataPeriodically(db)
	// 加载空闲端口
	go UpdateFreePort(db)
	// 每天更新 traffic 数据库为 0
	// go UpdateTrafficPeriodically(db)
	// 用不到了
	// 更新预设
	go UpdateSettingPeriodically(db)
	// 计算用户流量
	go CalculateUserTraffic(db)
	// 扣费与强制下线
	go UpdateUserTraffic(db)
	// 更新客户端
	go UpdateClientPeriodically(db)

	c := cron.New(cron.WithSeconds()) // 创建一个新的cron实例，支持秒级别的精度
	// 每天凌晨 0 点执行一次
	_, _ = c.AddFunc("0 0 0 * * *", func() {
		SettleUserTraffic(db)
	})
	// 每天凌晨 0：30 点执行一次
	_, _ = c.AddFunc("0 30 0 * * *", func() {
		// 此时 frps 用户流量清零 清空数据库
		ClearUserTraffic(db)
	})
	// 每10min执行一次

	_, _ = c.AddFunc("0 0/10 * * * *", func() {
		FetchTraffic(db)
	})
	// 启动定时任务
	c.Start()
	// 阻塞主线程
	select {}
}
