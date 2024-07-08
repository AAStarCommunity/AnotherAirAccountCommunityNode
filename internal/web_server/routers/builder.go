package routers

import (
	account_v1 "another_node/internal/web_server/controllers/account/v1"
	dashboard_v1 "another_node/internal/web_server/controllers/dashboard/v1"
	"github.com/gin-gonic/gin"
)

func buildRouters(router *gin.Engine) {

	router.GET("/api/dashboard/v1/node", dashboard_v1.Node)
	router.POST("/api/account/v1/bind", account_v1.Bind)
}
