package middlewares

import (
	"another_node/internal/community/storage"
	"another_node/internal/web_server/pkg/response"
	"github.com/gin-gonic/gin"
)

type ApiKey struct {
	Key string `form:"apiKey" json:"apiKey" binding:"required"`
}

func ApiVerificationHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Query("apiKey")
		if apiKey == "" {
			response.BadRequest(c, "ApiKey is mandatory, visit to https://dashboard.aastar.io for more detail")
			return
		}
		apiModel, err := storage.GetApiInfoByApiKey(apiKey)
		if err != nil {
			response.BadRequest(c, "can Not Find Your Api Key")
			return
		}
		if apiModel.Disable {
			response.BadRequest(c, "api Key Is Disabled")
			return
		}
		if !apiModel.AirAccountEnable {
			response.BadRequest(c, "api Key Is Disabled AirAccount")
			return
		}

		if !VerifyRateLimit(*apiModel) {
			response.BadRequest(c, "too many requests")
			return
		}

		if apiModel.IPWhiteList != nil && apiModel.IPWhiteList.Cardinality() > 0 {
			clientIp := c.ClientIP()
			if !apiModel.IPWhiteList.Contains(clientIp) {
				response.BadRequest(c, "ip not in whitelist")
				return
			}
		}
		if apiModel.DomainWhitelist != nil && apiModel.DomainWhitelist.Cardinality() > 0 {
			domain := c.Request.Host
			if !apiModel.DomainWhitelist.Contains(domain) {
				response.BadRequest(c, "domain not in whitelist")
				return
			}
		}
	}
}
