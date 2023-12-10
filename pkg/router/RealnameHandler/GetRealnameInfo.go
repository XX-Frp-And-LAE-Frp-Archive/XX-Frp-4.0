package RealnameHandler

import (
	_struct "github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	Username string
	Group    string
}

func GetRealnameInfo(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(_struct.User)
	group := user.Group
	if group == "admin" {
		respond.Respond(c, 200, "Success!", gin.H{
			"code":     200,
			"realname": "管理员",
			"view":     "realname",
			"time":     0,
		})
	} else if group == "realname" {
		// 获取数据库 realname 表中的 time 字段 并将时间戳转换为时间
		var realname Realname
		db.Model(&Realname{}).Where("username = ?", user.Username).Select("time").Find(&realname)
		respond.Respond(c, 200, "Success!", gin.H{
			"code":     200,
			"realname": "已实名认证",
			"view":     "realname",
			"time":     realname.Time,
		})
	} else if group == "default" {
		respond.Respond(c, 200, "Success!", gin.H{
			"code":     200,
			"realname": "未实名认证",
			"view":     "default",
		})
	} else {
		respond.Respond(c, 500, "未知错误!", gin.H{
			"code":     500,
			"realname": "未知错误",
			"view":     "unknown",
		})
	}

}
