package account_v1

import (
	"another_node/internal/community"
	"another_node/internal/web_server/pkg/request"
	"another_node/internal/web_server/pkg/response"

	"github.com/gin-gonic/gin"
)

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
	var req request.Bind
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	community.BindAccount(req.Account, &req.PublicKey)
}
