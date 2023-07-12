package router

import (
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
}

// Token 结构体
type Token struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null"`
	Token    string `gorm:"not null"`
	Status   string `gorm:"bool"`
}

func HandleLogin(c *gin.Context, db *gorm.DB) {
	// 解析请求参数
	var loginData struct {
		UsernameOrEmail string `json:"username_or_email" binding:"required"`
		Password        string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&loginData)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 根据用户名或邮箱查询用户
	user, err := findUserByUsernameOrEmail(loginData.UsernameOrEmail, db)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid username or email"})
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(loginData.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}

	// 查询数据库中的 Token
	token, err := findTokenByUserID(user.ID, db)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to query token"})
		return
	}

	// 检查 Token 是否存在
	if token == nil {
		c.JSON(401, gin.H{"error": "Token does not exist"})
		return
	}

	// 返回 Token
	c.JSON(200, gin.H{"access_token": token.Token})
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
