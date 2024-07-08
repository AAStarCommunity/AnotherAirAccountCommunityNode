package middlewares

import (
	"another_node/internal/web_server/pkg"
	"golang.org/x/time/rate"
)

const (
	DefaultBurst int = 50 // burst size, for surge traffic
)

var limiter map[string]*rate.Limiter

func VerifyRateLimit(keyModel pkg.ApiKeyModel) bool {
	return limiting(&keyModel.ApiKey, keyModel.RateLimit)
}

func limiting(apiKey *string, defaultLimit rate.Limit) bool {

	var l *rate.Limiter
	if limit, ok := limiter[*apiKey]; ok {
		l = limit
	} else {
		l = rate.NewLimiter(defaultLimit, DefaultBurst)
		limiter[*apiKey] = l
	}

	return l.Allow()
}

func init() {
	limiter = make(map[string]*rate.Limiter, 100)
}
