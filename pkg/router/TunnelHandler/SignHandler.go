package TunnelHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/router/UserHandler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"time"
)

type Sign struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"varchar(255); not null" json:"username"`
	Signdate     int64  `gorm:"varchar(255); not null" json:"signdate"`
	Totalsign    int64  `json:"totalsign"`
	Totaltraffic int64  `json:"totaltraffic"`
}

func HandleSignGet(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户名
	username, _ := c.Get("username")

	var sign Sign
	result := db.Table("sign").Where("username = ?", username).First(&sign)
	if result.Error != nil {
		sign = Sign{
			Username:     username.(string),
			Signdate:     0,
			Totalsign:    0,
			Totaltraffic: 0,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": sign})
}
func HandleSignPost(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户名
	username, _ := c.Get("username")

	var sign Sign
	result := db.Table("sign").Where("username = ?", username).First(&sign)
	if result.Error != nil {
		// 如果签到记录不存在，则创建一个新的签到记录
		sign = Sign{
			Username:     username.(string),
			Signdate:     time.Now().Unix(),
			Totalsign:    1,
			Totaltraffic: 0,
		}
	} else {
		lastSignDate := sign.Signdate
		currentDate := time.Now().Unix()

		// 检查是否满足签到条件（每隔24小时才能签到）
		if currentDate-lastSignDate < 24*60*60 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "You can only sign once every 24 hours"})
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

	// 读取 users 表的 traffic 值 将获得流量 乘以 1024 相加写入
	var user UserHandler.User
	result = db.Table("users").Where("username = ?", username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	user.Traffic += randomTraffic * 1024
	db.Save(&user)
	// 将签到记录保存到数据库中
	if result.Error != nil {
		db.Table("sign").Create(&sign)
	} else {
		db.Table("sign").Save(&sign)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sign successfully", "traffic": randomTraffic})
}
