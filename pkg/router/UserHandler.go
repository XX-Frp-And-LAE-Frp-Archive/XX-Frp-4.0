package router

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	// 替换密码字段为 Gmail 的 MD5 值
	user.EmailMD5 = getMD5Hash(user.Email)

	c.JSON(200, user)
}

// 获取字符串的 MD5 值
func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
