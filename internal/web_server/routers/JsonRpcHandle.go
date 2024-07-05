package routers

import (
	account_v1 "another_node/internal/web_server/controllers/account/v1"
	"another_node/internal/web_server/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func jsonrpcError(c *gin.Context, code int, message string, data any, id any) {
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
	//JsonRpcMethods["aa_sgin"] = Sign()

}

// JsonRpcHandle
// @Tags JsonRpcHandle
// @Description Airaccount JSON-RPC API
// @Accept json
// @Product json
// @param network path string true "Network"
// @Param rpcRequest body model.JsonRpcRequest true "JsonRpcRequest Model"
// @Param apiKey query string true "apiKey"
// @Router /api/v1/airaccount/{network}  [post]
// @Success 200
func JsonRpcHandle(c *gin.Context) {
	if c.Request.Body == nil {
		jsonrpcError(c, -32700, "Parse error", nil, nil)
		return
	}
	jsonRpcRequest := pkg.JsonRpcRequest{}
	if err := c.ShouldBindJSON(&jsonRpcRequest); err != nil {
		jsonrpcError(c, -32700, fmt.Sprintf(" Parse error"), err, nil)
		return
	}
	if jsonRpcRequest.JsonRpc != "2.0" {
		jsonrpcError(c, -32600, "Invalid Request", "Version of jsonrpc is not 2.0", nil)
		return
	}
	network := c.Param("network")
	if network == "" {
		jsonrpcError(c, -32600, "Invalid Request", "Network is empty", jsonRpcRequest.Id)
		return
	}
	jsonRpcRequest.Network = network

	if jsonRpcRequest.Method == "" {
		jsonrpcError(c, -32600, "Invalid Request", "Method is empty", jsonRpcRequest.Id)
		return
	}

	if methodFunc, ok := JsonRpcMethods[jsonRpcRequest.Method]; ok {
		defer func() {
			if r := recover(); r != nil {
				jsonrpcError(c, -32603, "Internal error", r, jsonRpcRequest.Id)
			}
		}()
		result, err := methodFunc(c, &jsonRpcRequest)
		if err != nil {
			jsonrpcError(c, -32603, "Internal error", err, jsonRpcRequest.Id)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"jsonrpc": "2.0",
			"result":  result,
			"id":      jsonRpcRequest.Id,
		})
		return
	} else {
		jsonrpcError(c, -32601, "Method not found", nil, jsonRpcRequest.Id)
		return
	}

}
