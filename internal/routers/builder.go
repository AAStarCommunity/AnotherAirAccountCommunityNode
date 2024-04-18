package routers

import (
	"another_node/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func buildRouters(router *gin.Engine) {

	api := router.Group("/api")
	api.POST("/login", nil)
	api.POST("/logout", nil)

	router.Use(middlewares.AuthHandler())
	{
	}
}
