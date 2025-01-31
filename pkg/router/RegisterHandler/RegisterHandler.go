package register

import (
	"encoding/hex"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/rand"
	"regexp"
	"time"
)

// 使用crypto/rand生成随机token
func GenerateToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func IsValidEmail(email string) bool {
	// 校验邮箱格式
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

func validatePassword(password string) bool {
	// 正则表达式
	lowerCaseLetter := regexp.MustCompile(`[a-z]`)
	upperCaseLetter := regexp.MustCompile(`[A-Z]`)
	number := regexp.MustCompile(`[0-9]`)
	specialChar := regexp.MustCompile(`[!@#\$%\^&\*]`)

	// 检查是否满足至少两个条件
	matches := 0
	if lowerCaseLetter.MatchString(password) || upperCaseLetter.MatchString(password) {
		matches++
	}
	if number.MatchString(password) {
		matches++
	}
	if specialChar.MatchString(password) {
		matches++
	}

	return matches >= 2
}

func isValidCode(code string) bool {
	// 校验验证码，验证码必须为六位数字
	codeRegex := `^\d{6}$`
	match, _ := regexp.MatchString(codeRegex, code)
	return match
}
func IsValidUsername(username string) bool {
	// 校验用户名格式，用户名长度必须为4-16个字符，且只能包含字母、数字和下划线
	usernameRegex := `^[a-zA-Z0-9_]{4,16}$`
	match, _ := regexp.MatchString(usernameRegex, username)
	return match
}

func checkUsername(username string, db *gorm.DB) bool {
	var count int64
	db.Model(&define.User{}).Where("LOWER(username) = LOWER(?)", username).Count(&count)
	// 如果count大于0，说明至少用户名已存在
	return count > 0
}
func checkEmail(email string, db *gorm.DB) bool {
	var count int64
	db.Model(&define.User{}).Where("LOWER(email) = LOWER(?)", email).Count(&count)
	// 如果count大于0，说明至少用户名已存在
	return count > 0
}

func HandleRegister(c *gin.Context, db *gorm.DB) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	code := c.PostForm("code")

	// 校验用户名 邮箱 密码 验证码格式
	if !IsValidUsername(username) {
		respond.Respond(c, 400, "用户名格式错误！", 0)
		return
	}
	if !IsValidEmail(email) {
		respond.Respond(c, 400, "邮箱格式错误！", 0)
		return
	}
	if !validatePassword(password) {
		respond.Respond(c, 400, "密码格式错误！", 0)
		return
	}
	if !isValidCode(code) {
		respond.Respond(c, 400, "验证码格式错误！", 0)
		return
	}
	usernameExists := checkUsername(username, db)
	if usernameExists {
		respond.Respond(c, 400, "注册失败，用户名已存在！", 0)
		return
	}
	// 使用email作为key，从数据库codes表中获取验证码 并与用户输入的验证码进行比对
	var codes define.Code
	db.Table("codes").Where("email = ?", email).First(&codes)
	if codes.Code != code {
		respond.Respond(c, 400, "注册失败，验证码错误！", 0)
		return
	} else {
		// 验证码正确，获取数据库对应行的time时间戳字段 与当前时间戳进行比对，判断是否超过 15 分钟
		if time.Now().Unix()-codes.Time > 900 {
			respond.Respond(c, 400, "注册失败，验证码已经过期！", 0)
			// 删除数据库中对应的验证码
			db.Table("codes").Where("email = ?", email).Delete(&codes)
			return
		}
		// 验证码正确，且未过期，删除数据库中对应的验证码
		// 异步删除验证码
		go func() {
			db.Table("codes").Where("email = ?", email).Delete(&define.Code{})
		}()
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		respond.Respond(c, 500, "注册失败，密码处理错误!", 0)
		return
	}

	// 生成注册时间戳
	regTime := time.Now().Unix()

	// 从数据库 groups 表中获取 name 为 default 组的 traffic 字段
	var groups define.Groups
	db.Table("groups").Where("name = ?", "default").First(&groups)
	traffic := groups.Traffic

	// 生成 token
	token := GenerateToken()

	// 创建用户对象
	user := define.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Status:   0,
		Group:    "default",
		Traffic:  traffic,
		Token:    token,
		Regtime:  regTime,
	}

	// 将用户对象写入数据库表 users
	result := db.Create(&user)
	if result.Error != nil {
		respond.Respond(c, 500, "注册失败，用户写入错误！", 0)
		return
	}

	respond.Respond(c, 200, "注册成功", gin.H{
		"token": token,
	})
}
