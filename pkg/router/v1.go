package router

import (
	"github.com/gin-gonic/gin"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// reverseProxy 处理v1部分api代理
func reverseProxy(c *gin.Context) {
	path := "/api/v1"
	url := c.Request.URL.Path[len(path):]

	//去掉路径末尾的斜杠（如果有）
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	// 构建代理 URL
	proxyURL := "http://admin.mefrp.com:8123/api/" + url

	// 发送代理请求
	resp, err := http.Get(proxyURL)
	if err != nil {
		// 返回 502 错误 (Bad Gateway)
		_ = c.AbortWithError(http.StatusBadGateway, err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// handle the error
		}
	}(resp.Body)

	// 获取Content-Type
	ext := filepath.Ext(resp.Request.URL.Path)
	contentType := mime.TypeByExtension(ext)

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Data(resp.StatusCode, contentType, body)
}
