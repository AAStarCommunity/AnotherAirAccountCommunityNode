package account_v1

import (
	"another_node/internal/community/node"
	"another_node/internal/web_server/pkg/request"
	"another_node/internal/web_server/pkg/response"

	"github.com/gin-gonic/gin"
)

// transfer a TX
// @Tags Account
// @Description transfer a TX
// @Accept json
// @Produce json
// @Success 201
// @Param tx body request.Transfer true "Transfer TX"
// @Router /api/account/v1/transfer [POST]
// @Security JWT
func Transfer(ctx *gin.Context) {
	var req request.Bind
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	if err := node.BindAccount(req.Account, &req.PublicKey); err != nil {
		response.InternalServerError(ctx, err)
	} else {
		response.Created(ctx, nil)
	}
}
