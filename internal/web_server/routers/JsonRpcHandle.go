package routers

import (
	account_v1 "another_node/internal/web_server/controllers/account/v1"
	"another_node/internal/web_server/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func jsonrpcError(c *gin.Context, code pkg.JsonRpcError, message string, data any, id any) {
	c.JSON(http.StatusOK, gin.H{
		"jsonrpc": "2.0",
		"error": gin.H{
			"code":    code,
			"message": message,
			"data":    data,
		},
		"id": id,
	})
	c.Abort()
}

var JsonRpcMethods = map[string]pkg.RpcMethodFunctionFunc{}

func init() {
	JsonRpcMethods["aa_bind"] = account_v1.RpcBind()
	JsonRpcMethods["aa_sign"] = account_v1.RpcSign()

}

// JsonRpcHandle
// @Tags JsonRpcHandle
// @Description AirAccount JSON-RPC API
// @Accept json
// @Product json
// @param network path string true "Network"
// @Param rpcRequest body pkg.JsonRpcRequest true "JsonRpcRequest Model"
// @Param apiKey query string true "apiKey"
// @Router /api/v1/airaccount_rpc/{network}  [post]
// @Success 200
func JsonRpcHandle(c *gin.Context) {
	if c.Request.Body == nil {
		jsonrpcError(c, pkg.ParseError, "Parse error", nil, nil)
		return
	}
	jsonRpcRequest := pkg.JsonRpcRequest{}
	if err := c.ShouldBindJSON(&jsonRpcRequest); err != nil {
		jsonrpcError(c, pkg.ParseError, fmt.Sprintf(" Parse error"), err, nil)
		return
	}
	if jsonRpcRequest.JsonRpc != "2.0" {
		jsonrpcError(c, pkg.InvalidRequest, "Invalid Request", "Version of jsonrpc is not 2.0", nil)
		return
	}
	network := c.Param("network")
	if network == "" {
		jsonrpcError(c, pkg.InvalidRequest, "Invalid Request", "Network is empty", jsonRpcRequest.Id)
		return
	}
	jsonRpcRequest.Network = network

	if jsonRpcRequest.Method == "" {
		jsonrpcError(c, pkg.InvalidRequest, "Invalid Request", "Method is empty", jsonRpcRequest.Id)
		return
	}

	if methodFunc, ok := JsonRpcMethods[jsonRpcRequest.Method]; ok {
		defer func() {
			if r := recover(); r != nil {
				jsonrpcError(c, pkg.ServerError, "Server error", r, jsonRpcRequest.Id)
			}
		}()
		result, err := methodFunc(c, &jsonRpcRequest)
		if err != nil {
			jsonrpcError(c, pkg.ServerError, "Server error", err, jsonRpcRequest.Id)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"jsonrpc": "2.0",
			"result":  result,
			"id":      jsonRpcRequest.Id,
		})
		return
	} else {
		jsonrpcError(c, pkg.MethodNotFound, "Method not found", nil, jsonRpcRequest.Id)
		return
	}

}
