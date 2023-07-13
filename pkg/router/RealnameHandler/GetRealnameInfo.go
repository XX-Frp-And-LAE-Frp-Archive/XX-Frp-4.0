package RealnameHandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	Username string
	Group    string
}

func GetRealnameInfo(c *gin.Context, db *gorm.DB) {
	username, _ := c.Get("username")
	// 获取数据库 users 表中的 group 字段
	var user User
	db.Model(&User{}).Where("username = ?", username).Select("group").Find(&user)
	// 获取数据库 groups 表中的 realname 字段
	group := user.Group
	if group == "admin" {
		c.JSON(200, gin.H{
			"code":     200,
			"realname": "管理员",
			"view":     "realname",
		})
	} else if group == "realname" {
		// 获取数据库 realname 表中的 time 字段 并将时间戳转换为时间
		var realname Realname
		db.Model(&Realname{}).Where("username = ?", username).Select("time").Find(&realname)
		c.JSON(200, gin.H{
			"code":     200,
			"realname": "已实名认证",
			"view":     "realname",
			"time":     realname.Time,
		})
	} else if group == "default" {
		c.JSON(200, gin.H{
			"code":     200,
			"realname": "未实名认证",
			"view":     "default",
		})
	} else {
		c.JSON(200, gin.H{
			"code":     200,
			"realname": "未知错误",
			"view":     "unknown",
		})
	}

}
