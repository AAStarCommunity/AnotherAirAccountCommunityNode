package account_v1

import (
	"another_node/internal/community/node"
	"another_node/internal/web_server/pkg"
	"another_node/internal/web_server/pkg/request"
	"another_node/internal/web_server/pkg/response"
	"errors"

	"github.com/gin-gonic/gin"
)

// Bind a account to community node
// @Tags Account
// @Description bind a account to community node
// @Accept json
// @Produce json
// @Success 201
// @Param bind body request.Bind true "Account Binding"
// @Param apiKey query string true "apiKey"
// @Router /api/account/v1/bind [POST]
func Bind(ctx *gin.Context) {
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

func RpcBind() pkg.RpcMethodFunctionFunc {
	return func(ctx *gin.Context, jsonRpcRequest *pkg.JsonRpcRequest) (*interface{}, error) {
		req := request.Bind{}
		errors.Is(ctx.ShouldBindJSON(&req), nil)
		if jsonRpcRequest.Params == nil || len(jsonRpcRequest.Params) <= 0 {
			return nil, errors.New("invalid request [params is empty]")
		}

		accountParam := jsonRpcRequest.Params[0]
		if accountParam == nil {
			return nil, errors.New("invalid request [account is empty]")
		}
		req.Account = accountParam.(string)

		publicKeyParam := jsonRpcRequest.Params[1]
		if publicKeyParam == nil {
			return nil, errors.New("invalid request [publicKey is empty]")
		}

		if err := node.BindAccount(req.Account, &req.PublicKey); err != nil {
			return nil, err
		} else {
			return nil, nil
		}
	}
}
