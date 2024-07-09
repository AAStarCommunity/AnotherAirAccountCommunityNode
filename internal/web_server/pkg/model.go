package pkg

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ApiKeyModel struct {
	Disable            bool               `json:"disable"`
	ApiKey             string             `json:"api_key"`
	RateLimit          rate.Limit         `json:"rate_limit"`
	UserId             int64              `json:"user_id"`
	NetWorkLimitEnable bool               `json:"network_limit_enable"`
	DomainWhitelist    mapset.Set[string] `json:"domain_whitelist"`
	IPWhiteList        mapset.Set[string] `json:"ip_white_list"`
	AirAccountEnable   bool               `json:"airaccount_enable"`
}

type RpcMethodFunctionFunc = func(ctx *gin.Context, jsonRpcRequest *JsonRpcRequest) (result interface{}, err error)

type JsonRpcRequest struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
	Network string        `json:"-"`
}
type JsonRpcError int

const (
	ParseError     JsonRpcError = -32700
	InvalidRequest JsonRpcError = -32600
	MethodNotFound JsonRpcError = -32601
	InvalidParams  JsonRpcError = -32602
	InternalError  JsonRpcError = -32603
	ServerError    JsonRpcError = -32099
)
