package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

func HandleSignGet(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	var sign define.Sign
	result := db.Table("sign").Where("username = ?", user.Username).First(&sign)
	if result.Error != nil {
		sign = define.Sign{
			Username:     user.Username,
			Signdate:     0,
			Totalsign:    0,
			Totaltraffic: 0,
		}
	}

	respond.Respond(c, 200, "Success!", sign)
}
func HandleSignPost(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(define.User)

	var sign define.Sign
	result := db.Table("sign").Where("username = ?", user.Username).First(&sign)
	if result.Error != nil {
		// 如果签到记录不存在，则创建一个新的签到记录
		sign = define.Sign{
			Username:     user.Username,
			Signdate:     time.Now().Unix(),
			Totalsign:    1,
			Totaltraffic: 0,
		}
	} else {
		lastSignDate := sign.Signdate
		currentDate := time.Now().Unix()

		// 检查是否满足签到条件（每隔24小时才能签到）
		if currentDate-lastSignDate < 24*60*60 {
			respond.Respond(c, 403, "24小时只能签到一次", 0)

			return
		}

		// 更新签到表的相关字段
		sign.Signdate = currentDate
		sign.Totalsign += 1
	}

	// 生成1-10GB之间的随机数
	rand.Seed(time.Now().UnixNano())
	randomTraffic := rand.Int63n(10) + 1

	// 更新签到表的总流量字段
	sign.Totaltraffic += randomTraffic

	user.Traffic += randomTraffic * 1024
	db.Save(&user)
	// 将签到记录保存到数据库中
	if result.Error != nil {
		db.Table("sign").Create(&sign)
	} else {
		db.Table("sign").Save(&sign)
	}
	respond.Respond(c, 200, "签到成功", gin.H{
		"traffic": randomTraffic,
	})
}
