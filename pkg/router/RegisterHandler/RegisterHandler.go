package register

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// User 用户结构体
type User struct {
	Username string
	Password string
	Email    string
	Traffic  int64
	Group    string
	Status   int
	Regtime  int64
}

// Codes 验证码结构体
type Codes struct {
	Email string
	Code  string
	Time  int64
}
type Groups struct {
	name    string
	Traffic int64
}

// Token 定义Token结构体
type Token struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Status   int
}

// 生成随机token
func GenerateToken() string {
	// 生成当前时间戳
	timestamp := time.Now().Unix()
	// 生成随机数
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(1000000) // 随机范围可根据需要更改

	// 将时间戳和随机数拼接成字符串
	tokenString := strconv.FormatInt(timestamp, 10) + strconv.Itoa(random)

	// 计算 MD5 值
	hash := md5.Sum([]byte(tokenString))
	md5Str := hex.EncodeToString(hash[:])

	return md5Str
}

func IsValidEmail(email string) bool {
	// 校验邮箱格式
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

func IsValidPassword(password string) bool {
	passwordRegex := `^\d{6,16}$`
	match, _ := regexp.MatchString(passwordRegex, password)
	return match
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

// 检查邮箱是否已经存在,不使用where查询，因为where查询会返回一个数组
func checkEmail(email string, db *gorm.DB) bool {
	var user User
	db.First(&user, "email = ?", email)
	lower := strings.ToLower(user.Email)
	if lower == strings.ToLower(email) {
		return true
	}
	return false
}

// 检查用户名是否已经注册
func checkUsername(username string, db *gorm.DB) bool {
	var user User
	// 查找是否存在该用户
	db.First(&user, "username = ?", username)
	// fmt.Printf(user.Username)
	lower := strings.ToLower(user.Username)
	if lower == strings.ToLower(username) {
		return true
	}
	return false
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
	//if !IsValidPassword(password) {
	//	c.JSON(400, gin.H{
	//		"code": 400,
	//		"msg":  "密码格式错误",
	//	})
	//	return
	// }
	if !isValidCode(code) {
		respond.Respond(c, 400, "验证码格式错误！", 0)
		return
	}
	// 查询username是否被注册
	if checkUsername(username, db) {
		respond.Respond(c, 400, "注册失败，用户名已存在！", 0)
		return
	}
	// 查询邮箱是否被注册
	if checkEmail(email, db) {
		respond.Respond(c, 400, "注册失败，邮箱已存在！", 0)
		return
	}
	// 使用email作为key，从数据库codes表中获取验证码 并与用户输入的验证码进行比对
	var codes Codes
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
		db.Table("codes").Where("email = ?", email).Delete(&codes)
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
	var groups Groups
	db.Table("groups").Where("name = ?", "default").First(&groups)
	traffic := groups.Traffic

	// 创建用户对象
	user := User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Status:   0,
		Group:    "default",
		Traffic:  traffic,
		Regtime:  regTime,
	}

	// 将用户对象写入数据库表 users
	result := db.Create(&user)
	if result.Error != nil {
		respond.Respond(c, 500, "注册失败，用户写入错误！", 0)
		return
	}

	// 读取数据库中自增的id
	var id int
	db.Table("users").Where("email = ?", email).First(&user).Scan(&id)
	// 生成token
	token := GenerateToken()
	// 将token写入数据库表 tokens
	tokenObj := Token{
		ID:       id,
		Username: username,
		Token:    token,
		Status:   0,
	}
	// 将token写入数据库表 tokens
	result = db.Create(&tokenObj)

	respond.Respond(c, 200, "注册成功", gin.H{
		"token": token,
	})
}
