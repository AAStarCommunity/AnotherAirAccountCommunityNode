package routers

import (
	account_v1 "another_node/internal/web_server/controllers/account/v1"
	auth_v1 "another_node/internal/web_server/controllers/auth/v1"
	"another_node/internal/web_server/middlewares"

	"github.com/gin-gonic/gin"
)

func buildRouters(router *gin.Engine) {

	router.POST("/api/auth/v1/login", auth_v1.Login)

	router.Use(middlewares.AuthHandler())
	{
		account := router.Group("/api/account")
		account.POST("/v1/bind", account_v1.Bind)
	}
}
