package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	seedwork "another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

func sessionKey(reg seedwork.Registration) string {
	return reg.Origin + ":" + reg.Email
}
func (relay *RelayParty) beginRegistration(ctx *gin.Context) {
	var reg seedwork.Registration
	if err := ctx.ShouldBindJSON(&reg); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	if session := relay.store.Get(sessionKey(reg)); session != nil {
		response.BadRequest(ctx, "Already in registration")
		return
	} else {
		if options, err := relay.store.NewSession(reg.Origin, reg.Email); err != nil {
			response.InternalServerError(ctx, err)
		} else {
			response.GetResponse().WithDataSuccess(ctx, options.Response)
		}
	}
}
