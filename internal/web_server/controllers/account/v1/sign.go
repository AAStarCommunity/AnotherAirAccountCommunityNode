package account_v1

import (
	"another_node/internal/community/node"
	"another_node/internal/web_server/pkg"
	"another_node/internal/web_server/pkg/request"
	"another_node/internal/web_server/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
)

// Sign a account to community node
// @Tags Account
// @Description sign a account to community node
// @Accept json
// @Produce json
// @Success 201
// @Router /api/account/v1/sign [POST]
func Sign(ctx *gin.Context) {
	var req request.Sign
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err)
		return
	}
	if err := node.Sign(req); err != nil {
		response.InternalServerError(ctx, err)
	} else {
		response.Created(ctx, struct {
			Msg string `json:"msg"`
		}{
			Msg: "Sign success",
		})
	}
}

func RpcSign() pkg.RpcMethodFunctionFunc {
	return func(ctx *gin.Context, jsonRpcRequest *pkg.JsonRpcRequest) (interface{}, error) {
		req := request.Sign{}
		if jsonRpcRequest.Params == nil || len(jsonRpcRequest.Params) <= 0 {
			return nil, errors.New("invalid request [params is empty]")
		}
		resp := struct {
			Msg string `json:"msg"`
		}{
			Msg: "Sign success",
		}
		//TODO convert Request Params to Sign Request
		if err := node.Sign(req); err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}
}
