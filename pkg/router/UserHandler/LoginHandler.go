package UserHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HandleLogin(c *gin.Context, db *gorm.DB) {
	// 解析表单数据
	UsernameOrEmail := c.PostForm("username")
	Password := c.PostForm("password")

	//
	// 根据用户名或邮箱查询用户
	user, err := findUserByUsernameOrEmail(UsernameOrEmail, db)
	if err != nil {
		respond.Respond(c, 403, "账户不存在!", 0)
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password))
	if err != nil {
		respond.Respond(c, 403, "密码错误!", 0)
		return
	}

	// 返回 Token
	respond.Respond(c, 200, "登录成功!", gin.H{
		"access_token": user.Token,
	})
}

// 根据用户名或邮箱查询用户
func findUserByUsernameOrEmail(usernameOrEmail any, db *gorm.DB) (*define.User, error) {
	user := &define.User{}
	err := db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
