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
	apiV1Router := r.Group("/api/v1")
	{
		apiV1Router.POST("/auth/login", func(c *gin.Context) {
			UserHandler.HandleLogin(c, db)
		})
		apiV1Router.GET("/info/sponsor", func(c *gin.Context) {
			InfoHandler.HandleSponsor(c)
		})
		apiV1Router.GET("/info/statistics", func(c *gin.Context) {
			InfoHandler.HandleStatistics(c)
		})
		apiV1Router.POST("/auth/reg/email", func(c *gin.Context) {
			register.HandleEmail(c, db)
		})
		apiV1Router.POST("/auth/register", func(c *gin.Context) {
			register.HandleRegister(c, db)
		})
		apiV1Router.POST("/auth/forgot_password", func(c *gin.Context) {
			UserHandler.HandleForgotPassword(c, db)
		})
		apiV1Router.POST("/auth/reset_password/:link", func(c *gin.Context) {
			UserHandler.HandleResetPassword(c, db)
		})

	}
	apiV2Router := r.Group("/api/v2")
	apiV2Router.Use(middleware.AuthMiddleware(db))
	{
		apiV2Router.GET("/user", func(c *gin.Context) {
			UserHandler.HandleUser(c, db)
		})
		apiV2Router.GET("/sign", func(c *gin.Context) {
			TunnelHandler.HandleSignGet(c, db)
		})
		apiV2Router.POST("/sign", func(c *gin.Context) {
			TunnelHandler.HandleSignPost(c, db)
		})
		apiV2Router.POST("/refresh_token", func(c *gin.Context) {
			UserHandler.HandleRefreshToken(c, db)
		})
		apiV2Router.GET("/realname/get", func(c *gin.Context) {
			RealnameHandler.GetRealnameInfo(c, db)
		})
		apiV2Router.POST("/realname/post", func(c *gin.Context) {
			RealnameHandler.RealnameHandler(c, db)
		})
		apiV2Router.GET("/tunnel/list", func(c *gin.Context) {
			TunnelHandler.GetTunnelList(c, db)
		})
		apiV2Router.GET("/tunnel/conf/node/:node", func(c *gin.Context) {
			TunnelHandler.GetConfByNode(c, db)
		})
		apiV2Router.GET("/tunnel/conf/id/:id", func(c *gin.Context) {
			TunnelHandler.GetConfByID(c, db)
		})
		apiV2Router.POST("/tunnel/create", func(c *gin.Context) {
			TunnelHandler.HandleCreateTunnel(c, db)
		})
		apiV2Router.POST("/tunnel/delete/:tunnelid", func(c *gin.Context) {
			TunnelHandler.HandleDeleteTunnel(c, db)
		})
		apiV2Router.GET("/tunnel/info/:tunnelid", func(c *gin.Context) {
			TunnelHandler.GetTunnelByID(c, db)
		})
		apiV2Router.GET("/node/list", func(c *gin.Context) {
			TunnelHandler.HandleGetNodeList(c, db)
		})
	}
	apiV3Router := r.Group("/api/v3")
	{
		apiV3Router.GET("/start", func(c *gin.Context) {
			StartHandler.HandleStart(c, db)
		})
	}
}
