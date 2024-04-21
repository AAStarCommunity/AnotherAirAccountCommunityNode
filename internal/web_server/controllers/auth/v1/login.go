package auth_v1

import (
	"another_node/internal/web_server/middlewares"

	"github.com/gin-gonic/gin"
)

// Login
// @Tags Auth
// @Description Get AccessToken By ApiKey
// @Accept json
// @Product json
// @Param credential body request.ClientCredential true "AccessToken Model"
// @Router /api/auth/v1/login [post]
// @Success 201
func Login(ctx *gin.Context) {
	middlewares.GinJwtMiddleware().LoginHandler(ctx)
}
