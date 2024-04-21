package account_v1

import "github.com/gin-gonic/gin"

// Bind a account to community node
// @Tags Account
// @Description bind a account to community node
// @Accept json
// @Produce json
// @Success 201
// @Param bind body request.Bind true "Account Binding"
// @Router /api/account/v1/bind [POST]
// @Security JWT
func Bind(ctx *gin.Context) {
	ctx.JSON(201, gin.H{})
}
