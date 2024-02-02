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
	go UpdateTrafficPeriodically(db)
	// 更新预设
	go UpdateSettingPeriodically(db)
	// 获取节点状态
	go UpdateNodeStatus(db)
}
