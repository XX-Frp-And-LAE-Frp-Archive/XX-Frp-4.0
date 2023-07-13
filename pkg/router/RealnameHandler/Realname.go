package RealnameHandler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ahmr-bot/ME-Frp/pkg/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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
	//通过正则表达式判断是否为18位数字
	if m, _ := regexp.MatchString("^[0-9]{18}$", idcard); !m {
		return false
	}
	return true
}
func calcAuthorization(source string, secretId string, secretKey string) (auth string, datetime string, err error) {
	timeLocation, _ := time.LoadLocation("Etc/GMT")
	datetime = time.Now().In(timeLocation).Format("Mon, 02 Jan 2006 15:04:05 GMT")
	signStr := fmt.Sprintf("x-date: %s\nx-source: %s", datetime, source)

	// hmac-sha1
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	auth = fmt.Sprintf("hmac id=\"%s\", algorithm=\"hmac-sha1\", headers=\"x-date x-source\", signature=\"%s\"",
		secretId, sign)

	return auth, datetime, nil
}
func RealnameHandler(c *gin.Context, db *gorm.DB) {
	conf := config.GetConfig()
	username, _ := c.Get("username")
	// 获取提交的表单中的名字和身份证号
	name := c.PostForm("name")
	idcard := c.PostForm("idcard")
	// 检查提交的表单中的名字和身份证号格式
	if !IsValidName(name) {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "姓名格式错误",
		})
		return
	}
	if !IsValidIDCard(idcard) {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "身份证号格式错误",
		})
		return
	}
	// 检查是否已经实名认证
	var realname Realname
	chachong := db.Where("username = ?", username).First(&realname)
	if chachong.RowsAffected != 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "已经实名认证",
		})
		return
	}
	// 检查24小时内是否有失败记录
	var failrealname Failrealname
	fangshua := db.Where("username = ?", username).First(&failrealname)
	if fangshua.RowsAffected != 0 {
		if time.Now().Unix()-failrealname.Time < 86400 {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "24小时内只能认证一次",
			})
			return
		}
		// 超过24小时删除对应用户的失败记录
		db.Delete(&failrealname)
	}
	// 通过 api 调用 https://service-lbiior1h-1307960160.sh.apigw.tencentcs.com/release/idcard/verify 进行实名认证
	// post请求参数为 name 和 idcard
	// 返回值为 {"code":200,"msg":"success","data":{"result":1}}
	// result为1表示实名认证成功 result为2表示不匹配 result为3表示暂无
	// 通过 result 判断是否实名认证成功
	// 通过 username 更新数据库中的 realname 和 idcard 字段
	// 返回 {"code":200,"msg":"success"}

	source := "market"
	// 进行签名
	auth, datetime, err := calcAuthorization(source, conf.Realname.SecretID, conf.Realname.SecretKey)
	// 拼接 body 数据
	// body参数

	// 使用 x-www-form-urlencoded 格式
	postData := url.Values{}
	postData.Add("name", name)
	postData.Add("idcard", idcard)
	// 发送 POST 请求 并使用 bodyParams 作为 body headers 作为请求头
	// 创建一个请求对象，并设置请求方法、URL 和请求体
	req, err := http.NewRequest("POST", "https://service-lbiior1h-1307960160.sh.apigw.tencentcs.com/release/idcard/verify", strings.NewReader(postData.Encode()))
	if err != nil {
		panic(err)
	}
	// 设置请求头
	req.Header.Set("X-Source", source)
	req.Header.Set("X-Date", datetime)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// debug
	// fmt.Println(auth)
	// fmt.Println(datetime)
	// fmt.Println(source)
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// 关闭请求
	defer resp.Body.Close()

	// 解析返回的 JSON 数据
	var result VerifyResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("解析响应异常:", err)
		return
	}

	// 判断实名认证结果
	if result.Code == 200 && result.Data.Result == 1 {
		// 在 realname 表的创建一行 写入 username name idcard time 信息
		db.Create(&Realname{
			Username: username.(string),
			Name:     name,
			IDCard:   idcard,
			Time:     time.Now().Unix(),
		})
		// 更新 users 表的 group 为 realname
		db.Model(&User{}).Where("username = ?", username).Update("group", "realname")

		// 返回更新成功的消息

		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "success",
		})
	} else {
		// 将 username time 写入到 failrealname 表中
		db.Create(&Failrealname{
			Username: username.(string),
			Time:     time.Now().Unix(),
		})
		// 返回实名认证失败的消息
		c.JSON(200, gin.H{
			"code":        400,
			"msg":         "实名认证失败",
			"result":      resp.StatusCode,
			"result-data": result.Data.Result,
		})
	}
}
