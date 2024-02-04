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

type Announce struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Alert struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Response struct {
	Announce map[string]Announce `json:"announce"`
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

	announceMap := make(map[string]Announce)
	for _, setting := range cron.SettingCache {
		switch setting.Type {
		case "announce1", "announce2", "announce3", "announce4":
			announce := Announce{
				Title:   setting.Title,
				Content: setting.Content,
			}
			announceMap[setting.Type] = announce
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
	response.Announce = announceMap
	response.Ads = adMap
	response.Alert = alertMap

	respond.Respond(c, 200, "成功", response)
}
