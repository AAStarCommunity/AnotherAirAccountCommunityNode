package routers

import (
	"another_node/internal/community"
	account_v1 "another_node/internal/web_server/controllers/account/v1"
	auth_v1 "another_node/internal/web_server/controllers/auth/v1"
	dashboard_v1 "another_node/internal/web_server/controllers/dashboard/v1"
	"another_node/internal/web_server/middlewares"

	"github.com/gin-gonic/gin"
)

func buildRouters(router *gin.Engine) {

	// testing router
	router.GET("/api/say", func(ctx *gin.Context) {
		community.SayHello()
	})

	router.POST("/api/auth/v1/login", auth_v1.Login)
	router.GET("api/dashboard/v1/node", dashboard_v1.Node)

	router.Use(middlewares.AuthHandler())
	{
		account := router.Group("/api/account")
		account.POST("/v1/bind", account_v1.Bind)
	}
}
