package router

import (
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
)

// reverseProxy 是一个自定义的中间件函数，用于实现反向代理
func reverseProxy(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
