package router

import (
	"github.com/ahmr-bot/ME-Frp/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoadRoutes(r *gin.Engine, db *gorm.DB) {
	// 加载路由
	apiv1Router := r.Group("/api/v1")
	{
		apiv1Router.POST("/login", func(c *gin.Context) {
			HandleLogin(c, db)
		})
		apiv1Router.GET("/sponsor", func(c *gin.Context) {
			HandleSponsor(c)
		})
		apiv1Router.GET("/statistics", func(c *gin.Context) {
			HandleStatistics(c)
		})

	}
	apiv2Router := r.Group("/api/v2")
	apiv2Router.Use(middleware.AuthMiddleware(db))
	{
		apiv2Router.GET("/user", func(c *gin.Context) {
			HandleUser(c, db)
		})
		apiv2Router.GET("/sign", func(c *gin.Context) {
			HandleSignGet(c, db)
		})
		apiv2Router.POST("/sign", func(c *gin.Context) {
			HandleSignPost(c, db)
		})
		apiv2Router.GET("")
	}
}
