package routers

import (
	"another_node/conf"
	"another_node/internal/web_server/middlewares"
	"another_node/internal/web_server/pkg/response"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
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
		buildSwagger(routers)      // build swagger
	}

	// 加载中间件
	routers.Use(handlers...)

	buildRouters(routers) // 构建http路由

	routers.NoRoute(func(ctx *gin.Context) {
		response.GetResponse().SetHttpCode(http.StatusNotFound).FailCode(ctx, http.StatusNotFound)
	})

	return
}

// buildSwagger 创建swagger文档
func buildSwagger(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
