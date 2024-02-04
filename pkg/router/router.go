package router

import (
	"github.com/ahmr-bot/ME-Frp/pkg/middleware"
	"github.com/ahmr-bot/ME-Frp/pkg/router/InfoHandler"
	"github.com/ahmr-bot/ME-Frp/pkg/router/RealnameHandler"
	register "github.com/ahmr-bot/ME-Frp/pkg/router/RegisterHandler"
	"github.com/ahmr-bot/ME-Frp/pkg/router/StartHandler"
	"github.com/ahmr-bot/ME-Frp/pkg/router/TunnelHandler"
	"github.com/ahmr-bot/ME-Frp/pkg/router/UserHandler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoadRoutes(r *gin.Engine, db *gorm.DB) {
	// 加载路由
	// v4 版本 API
	apiPublicRouter := r.Group("/api/v4/public")
	{
		apiPublicRouter.POST("/verify/login", func(c *gin.Context) {
			UserHandler.HandleLogin(c, db)
		})
		apiPublicRouter.GET("/info/sponsor", func(c *gin.Context) {
			InfoHandler.HandleSponsor(c)
		})
		apiPublicRouter.GET("/info/statistics", func(c *gin.Context) {
			InfoHandler.HandleStatistics(c)
		})
		apiPublicRouter.POST("/verify/register/email", func(c *gin.Context) {
			register.HandleEmail(c, db)
		})
		apiPublicRouter.POST("/verify/register", func(c *gin.Context) {
			register.HandleRegister(c, db)
		})
		apiPublicRouter.POST("/verify/forgot_password", func(c *gin.Context) {
			UserHandler.HandleForgotPassword(c, db)
		})
		apiPublicRouter.POST("/verify/reset_password/:link", func(c *gin.Context) {
			UserHandler.HandleForgotResetPassword(c, db)
		})
		apiPublicRouter.GET("/info/setting", func(c *gin.Context) {
			InfoHandler.HandleSetting(c)
		})
		apiPublicRouter.GET("/client/update", func(c *gin.Context) {
			InfoHandler.HandleClient(c)
		})
	}
	apiAuthRouter := r.Group("/api/v4/auth")
	apiAuthRouter.Use(middleware.AuthMiddleware(db))
	{
		apiAuthRouter.GET("/user", func(c *gin.Context) {
			UserHandler.HandleUser(c, db)
		})
		apiAuthRouter.GET("/user/sign", func(c *gin.Context) {
			TunnelHandler.HandleSignGet(c, db)
		})
		apiAuthRouter.POST("/user/sign", func(c *gin.Context) {
			TunnelHandler.HandleSignPost(c, db)
		})
		apiAuthRouter.POST("/user/refresh_token", func(c *gin.Context) {
			UserHandler.HandleRefreshToken(c, db)
		})
		apiAuthRouter.GET("/user/realname/get", func(c *gin.Context) {
			RealnameHandler.GetRealnameInfo(c, db)
		})
		apiAuthRouter.POST("/user/realname/post", func(c *gin.Context) {
			RealnameHandler.RealnameHandler(c, db)
		})
		apiAuthRouter.GET("/tunnel/list", func(c *gin.Context) {
			TunnelHandler.GetTunnelList(c, db)
		})
		apiAuthRouter.GET("/tunnel/conf/node/:node", func(c *gin.Context) {
			TunnelHandler.GetConfByNode(c, db)
		})
		apiAuthRouter.GET("/tunnel/conf/id/:id", func(c *gin.Context) {
			TunnelHandler.GetConfByID(c, db)
		})
		apiAuthRouter.POST("/tunnel/create", func(c *gin.Context) {
			TunnelHandler.HandleCreateTunnel(c, db)
		})
		apiAuthRouter.POST("/tunnel/delete/:tunnelid", func(c *gin.Context) {
			TunnelHandler.HandleDeleteTunnel(c, db)
		})
		apiAuthRouter.GET("/tunnel/info/:tunnelid", func(c *gin.Context) {
			TunnelHandler.GetTunnelByID(c, db)
		})
		apiAuthRouter.GET("/node/list", func(c *gin.Context) {
			TunnelHandler.HandleGetNodeList(c, db)
		})
		apiAuthRouter.GET("/node/list/all", func(c *gin.Context) {
			TunnelHandler.HandleGetAllNode(c, db)
		})
		apiAuthRouter.POST("/user/reset_password", func(c *gin.Context) {
			UserHandler.HandleResetPassword(c, db)
		})
		apiAuthRouter.GET("/tunnel/get_free_port", func(c *gin.Context) {
			TunnelHandler.HandleGetFreePort(c)
		})
		apiAuthRouter.POST("/tunnel/close_tunnel/:tunnelid", func(c *gin.Context) {
			TunnelHandler.CloseTunnel(c, db)
		})
		apiAuthRouter.POST("/tunnel/edit_tunnel", func(c *gin.Context) {
			TunnelHandler.HandleEditTunnel(c, db)
		})
	}
	//apiV4managerRouter := r.Group("/api/v4/manage")
	//{
	//	apiV4managerRouter.POST
	//}
	// 保留 v3 版本的 Frps 鉴权
	apiV3Router := r.Group("/api/v3")
	{
		apiV3Router.GET("/start", func(c *gin.Context) {
			StartHandler.HandleStart(c, db)
		})
	}
	// 考虑到地址更新后 简单启动受到影响，故保留 v2 版本简单启动的一部分
	apiV2Router := r.Group("/api/v2")
	apiV2Router.Use(middleware.AuthMiddleware(db))
	{
		apiV2Router.GET("/tunnel/conf/id/:id", func(c *gin.Context) {
			TunnelHandler.GetConfByID(c, db)
		})
	}
	// 原樱花面板（ME Frp 1.0） API（v1）
	// apiV1ProxyURL := "http://admin.mefrp.com:8123/api/"
	r.GET("/api/v1/*filepath", func(c *gin.Context) {
		reverseProxy(c)
	})

	// 监听默认路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "message": "ME Frp 4.0 API Server Is OK!"})
	})
}
