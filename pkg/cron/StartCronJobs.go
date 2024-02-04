package cron

import (
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
	// 获取节点状态
	go FetchServerInfo(db)
	// 加载用户流量
	go FetchTraffic(db)
	// 计算用户流量
	go CalculateUserTraffic(db)
	// 扣费与强制下线
	go UpdateUserTraffic(db)
	// 更新客户端
	go UpdateClientPeriodically(db)
}
