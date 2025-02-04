package routers

import (
	"another_node/conf"
	"another_node/docs"
	"another_node/internal/web_server/middlewares"
	"another_node/internal/web_server/pkg/response"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetRouters setting the routers
func SetRouters() (routers *gin.Engine) {
	routers = gin.New()

	// middlewares
	handlers := make([]gin.HandlerFunc, 0)
	handlers = append(handlers, middlewares.GenericRecovery())
	if conf.Environment.IsDevelopment() {
		handlers = append(handlers, gin.Logger())
	}
	handlers = append(handlers, middlewares.CorsHandler())

	// product mode
	if conf.Environment.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard // forbidden gin log output
	}

	// development mode
	if conf.Environment.IsDevelopment() {
		gin.SetMode(gin.DebugMode) // debugger mode
	}
	buildSwagger(routers) // build swagger

	// use middlewares
	routers.Use(handlers...)
	routers.GET("/api/healthz", Healthz)
	if conf.IsApiKeyAccessEnable() {
		fmt.Println("[AirAccount] Api Key Access Enable")
		routers.Use(middlewares.ApiVerificationHandler())
	}
	buildRouters(routers)

	routers.NoRoute(func(ctx *gin.Context) {
		response.GetResponse().SetHttpCode(http.StatusNotFound).FailCode(ctx, http.StatusNotFound)
	})

	return
}

// buildSwagger build swagger document
func buildSwagger(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Healthz
// @Tags Healthz
// @Description Get Healthz
// @Accept json
// @Product json
// @Router /api/healthz [get]
// @Success 200
func Healthz(c *gin.Context) {
	resp := response.GetResponse()
	resp.SuccessWithData(c, gin.H{
		"hello":   "AAStar AirAccount",
		"time":    time.Now(),
		"version": "v1.0.0",
	})
}
