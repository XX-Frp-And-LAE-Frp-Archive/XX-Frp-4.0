package UserHandler

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type UserInfo struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"varchar(255); not null" json:"username"`
	Password string `gorm:"varchar(255); not null" json:"-"`
	Email    string `gorm:"varchar(255); not null" json:"email"`
	EmailMD5 string `gorm:"-" json:"email_md5"`
	Traffic  int64  `json:"traffic"`
	Proxies  int    `json:"proxies"`
	Group    string `gorm:"varchar(255); not null" json:"group"`
	RegTime  int    `json:"reg_time"`
	Status   int    `gorm:"varchar(255); not null" json:"status"`
	Outbound int64  `json:"outbound"`
	Inbound  int    `json:"inbound"`
}
type Limit struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"varchar(255); not null" json:"username"`
	Outbound int64  `json:"outbound"`
	Inbound  int    `json:"inbound"`
	Proxies  int    `json:"proxies"`
}
type Group struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"varchar(255); not null" json:"name"`
	Outbound   int    `json:"outbound"`
	Inbound    int    `json:"inbound"`
	Proxies    int    `json:"proxies"`
	CreateTime string `gorm:"varchar(255); not null" json:"create_time"`
}

func HandleUser(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户名
	username, _ := c.Get("username")

	var user UserInfo
	result := db.Table("users").Where("username = ?", username).First(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	groupName := user.Group

	// 根据用户名查询 limit 表
	var limitUser Limit
	limitResult := db.Table("limits").Where("username = ?", user.Username).First(&limitUser)
	if limitResult.Error == nil {
		// 如果存在对应记录，则返回对应的 outbound 和 inbound proxies 值
		user.Proxies = limitUser.Proxies
		user.Outbound = limitUser.Outbound
		user.Inbound = limitUser.Inbound
	} else {
		// 如果不存在对应记录，则根据用户组查询 groups 表，找到对应的 outbound 和 inbound proxies 值
		var group Group
		groupResult := db.Table("groups").Where("name = ?", groupName).First(&group)
		if groupResult.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Group not found"})
			return
		}
		user.Proxies = group.Proxies
		user.Outbound = int64(group.Outbound)
		user.Inbound = group.Inbound
	}

	user.EmailMD5 = getMD5Hash(user.Email)

	c.JSON(200, user)
}

// 获取字符串的 MD5 值
func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
