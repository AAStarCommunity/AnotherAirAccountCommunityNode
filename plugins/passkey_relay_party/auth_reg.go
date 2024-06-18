package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

func (relay *RelayParty) beginRegistration(ctx *gin.Context) {
	var reg seedworks.Registration
	if err := ctx.ShouldBindJSON(&reg); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	if session := relay.store.Get(seedworks.GetSessionKey(&reg)); session != nil {
		response.BadRequest(ctx, "Already in registration")
		return
	} else {
		if options, err := relay.store.NewRegSession(&reg); err != nil {
			response.InternalServerError(ctx, err)
		} else {
			response.GetResponse().WithDataSuccess(ctx, options.Response)
		}
	}
}
