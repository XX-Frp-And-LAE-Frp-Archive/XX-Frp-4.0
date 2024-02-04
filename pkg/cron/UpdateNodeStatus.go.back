package cron

import (
	"encoding/json"
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// 403 禁用+隐藏
// 200 正常
// 201 无状态数据
// 500 服务器异常

const (
	StatusDisabled    = 403
	StatusNormal      = 200
	StatusNoData      = 201
	StatusServerError = 500
)

// UpdateNodeStatus 更新节点状态
func UpdateNodeStatus(db *gorm.DB) {
	for {
		// 获取所有节点
		var nodes []define.Node
		db.Find(&nodes)

		// 发起 HTTP 请求获取 Kuma 状态数据
		resp, err := http.Get("https://kuma.mcserverx.top/api/status-page/heartbeat/frp")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		// 解析响应数据
		var kumaData define.KumaData
		if err := json.NewDecoder(resp.Body).Decode(&kumaData); err != nil {
			fmt.Println(err)
			return
		}

		// 更新节点状态
		for _, node := range nodes {
			key := fmt.Sprintf("%d_24", node.KumaId) // 构造 key
			uptime, exists := kumaData.UptimeList[key]
			heartbeat, exist2 := kumaData.HeartbeatList[strconv.Itoa(node.KumaId)]
			if node.Status != StatusDisabled {
				switch {
				case exists || exist2:
					var status int
					for i := len(heartbeat) - 1; i >= 0; i-- {
						if heartbeat[i].Status == 1 {
							status = StatusNormal
							break
						} else if heartbeat[i].Status == 0 {
							status = StatusServerError
							break
						}
					}
					db.Model(&node).Where("id = ?", node.ID).Update("status", status)
					db.Model(&node).Where("id = ?", node.ID).Update("health24", uptime)
				default:
					fmt.Printf("Kuma 不存在节点数据: %s\n", key)
					db.Model(&node).Where("id = ?", node.ID).Update("status", StatusNoData)
				}
			}
		}
		log.Print("Updated NodeStatus data from database.")
		time.Sleep(5 * time.Minute)
	}
}
