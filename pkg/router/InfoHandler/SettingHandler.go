package InfoHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/cron"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
)

type Ad struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Alert struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Response struct {
	Announce []map[string]string `json:"announce"`
	Ads      map[string]Ad       `json:"ads"`
	Alert    map[string]Alert    `json:"alert"`
}

func HandleSetting(c *gin.Context) {
	response := Response{
		Announce: nil,
		Ads:      nil,
		Alert:    nil,
	}

	adMap := make(map[string]Ad)
	alertMap := make(map[string]Alert)

	for _, setting := range cron.SettingCache {
		switch setting.Type {
		case "announce":
			if response.Announce == nil {
				response.Announce = make([]map[string]string, 0)
			}
			announceMap := map[string]string{
				"title":   setting.Title,
				"content": setting.Content,
			}
			response.Announce = append(response.Announce, announceMap)
		case "ad1", "ad2", "ad3", "ad4":
			ad := Ad{
				Title:   setting.Title,
				Content: setting.Content,
			}
			adMap[setting.Type] = ad
		case "error", "warning", "info", "success":
			alert := Alert{
				Title:   setting.Title,
				Content: setting.Content,
			}
			alertMap[setting.Type] = alert
		}
	}

	response.Ads = adMap
	response.Alert = alertMap

	respond.Respond(c, 200, "成功", response)
}
