package respond

import (
	"github.com/gin-gonic/gin"
)

// Response 是所有 API 响应的基本结构体
type Response struct {
	Code    int
	Message string
	Data    interface{} `json:"data,omitempty"`
}

// Response 返回的 API 响应
func Respond(c *gin.Context, code int, message string, data interface{}) {
	res := gin.H{
		"status":    code,
		"message": message,
	}

	// 判断 data 是否为空，不为空则添加到响应中
	if data != nil {
		res["data"] = data
	}

	c.JSON(code, res)
}
