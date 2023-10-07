package UserHandler

import (
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户结构体
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password []byte
	Traffic  int64
}

// Token 结构体
type Token struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null"`
	Token    string `gorm:"not null"`
	Status   string `gorm:"bool"`
}

func HandleLogin(c *gin.Context, db *gorm.DB) {
	// 解析表单数据
	UsernameOrEmail := c.PostForm("username")
	Password := c.PostForm("password")

	//
	// 根据用户名或邮箱查询用户
	user, err := findUserByUsernameOrEmail(UsernameOrEmail, db)
	if err != nil {
		respond.Respond(c, 403, "账户不存在", 0)
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(Password))
	if err != nil {
		respond.Respond(c, 403, "密码错误", 0)
		return
	}

	// 查询数据库中的 Token
	token, err := findTokenByUserID(user.ID, db)
	if err != nil {
		respond.Respond(c, 500, "Token获取失败，请联系管理员", 0)
		return
	}

	// 检查 Token 是否存在
	if token == nil {
		respond.Respond(c, 500, "Token获取失败（2类），请联系管理员", 0)
		return
	}

	// 返回 Token
	respond.Respond(c, 200, "登录成功", gin.H{
		"access_token": token.Token,
	})
}

// 根据用户名或邮箱查询用户
func findUserByUsernameOrEmail(usernameOrEmail string, db *gorm.DB) (*User, error) {
	user := &User{}
	err := db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 根据用户ID查询 Token
func findTokenByUserID(userID uint, db *gorm.DB) (*Token, error) {
	token := &Token{}
	err := db.Where("id = ?", userID).First(token).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return token, nil
}
