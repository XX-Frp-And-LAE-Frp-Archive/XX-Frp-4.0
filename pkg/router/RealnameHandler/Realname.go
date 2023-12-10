package RealnameHandler

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	_struct "github.com/ahmr-bot/ME-Frp/pkg/define"
	"github.com/ahmr-bot/ME-Frp/pkg/respond"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"regexp"
	"time"
)

type VerifyResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Result int `json:"result"`
	} `json:"data"`
}

type Realname struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	IDCard   string `json:"idcard"`
	Time     int64  `json:"time"`
}
type Failrealname struct {
	Username string `json:"username"`
	Time     int64  `json:"time"`
}

func IsValidName(name string) bool {
	//通过正则表达式判断是否为中文
	if m, _ := regexp.MatchString("^\\p{Han}+$", name); !m {
		return false
	}
	return true
}
func IsValidIDCard(idcard string) bool {
	//通过正则表达式判断是否为18位数字或包含X
	if m, _ := regexp.MatchString("^[0-9]{17}[0-9X]$", idcard); !m {
		return false
	}
	return true
}
func RealnameHandler(c *gin.Context, db *gorm.DB) {
	// 获取中间件传递的用户信息
	userInterface, _ := c.Get("user")
	user, _ := userInterface.(_struct.User)
	// 获取提交的表单中的名字和身份证号
	name := c.PostForm("name")
	idcard := c.PostForm("idcard")
	// 检查提交的表单中的名字和身份证号格式
	if !IsValidName(name) {
		respond.Respond(c, 400, "姓名格式错误!", 0)
		return
	}
	if !IsValidIDCard(idcard) {
		respond.Respond(c, 400, "身份证格式错误!", 0)
		return
	}
	// 检查是否已经实名认证
	var realname Realname
	exist := db.Where("username = ?", user.Username).First(&realname)
	if exist.RowsAffected != 0 {
		respond.Respond(c, 400, "已经实名认证!", 0)
		return
	}
	// 检查24小时内是否有失败记录
	var failrealname Failrealname
	fails := db.Where("username = ?", user.Username).First(&failrealname)
	if fails.RowsAffected != 0 {
		if time.Now().Unix()-failrealname.Time < 86400 {
			respond.Respond(c, 400, "24小时只能认证一次！", 0)
			return
		}
		// 超过24小时删除对应用户的失败记录
		db.Delete(&failrealname)
	}
	conf := config.GetConfig()
	url := "https://s.tji0.com/api/two?appid=" + conf.Realname.SecretID + "&appkey=" + conf.Realname.SecretKey + "&name=" + name + "&card=" + idcard
	// 发送请求
	resp, err := http.Get(url)
	if err != nil {
		respond.Respond(c, 500, "实名认证接口请求错误，请联系管理员", 0)
		panic(err)
	}
	// 关闭请求
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respond.Respond(c, 500, "实名认证接口响应错误，请联系管理员", 0)
		return
	}

	responseText := string(body)
	// 判断实名认证结果

	switch responseText {
	case "0":
		// 使用base64加密 name
		encodeName := base64.StdEncoding.EncodeToString([]byte(name))
		// 使用base64加密 idcard
		encodeIdcard := base64.StdEncoding.EncodeToString([]byte(idcard))

		// 在 realname 表的创建一行 写入 username name idcard time 信息
		db.Create(&Realname{
			Username: user.Username,
			Name:     encodeName,
			IDCard:   encodeIdcard,
			Time:     time.Now().Unix(),
		})
		// 更新 users 表的 group 为 realname
		db.Model(&User{}).Where("username = ?", user.Username).Update("group", "realname")

		// 返回更新成功的消息
		respond.Respond(c, 200, "认证成功", 0)
	case "1":
		// 将 username time 写入到 failrealname 表中
		db.Create(&Failrealname{
			Username: user.Username,
			Time:     time.Now().Unix(),
		})
		// 返回实名认证失败的消息
		respond.Respond(c, 400, "信息不匹配，请24小时后重试", 0)
	}
}

// RsaEncrypt rsa加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	publicKey := `-----BEGIN RSA PUBLIC KEY-----
	MEgCQQD3reUyiTiGapGOUcuSc66AtSHQRlDkMeYRDxX+FlbfUsZUrqf0tuVdrSaV
	hhL3RIxsB4Jc39slV9/xZC7ZgpStAgMBAAE=
	-----END RSA PUBLIC KEY-----`
	// 定义 pkcs1 的公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		panic("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	// 使用公钥加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}
