package pkg

import (
	mapset "github.com/deckarep/golang-set/v2"
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
