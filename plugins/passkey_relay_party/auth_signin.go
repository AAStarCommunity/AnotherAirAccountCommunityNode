package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

func (relay *RelayParty) beginSignIn(ctx *gin.Context) {
	var signIn seedworks.SiginIn
	if err := ctx.ShouldBindJSON(&signIn); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	if session := relay.store.Get(seedworks.GetSessionKey(signIn.Registration)); session != nil {
		response.BadRequest(ctx, "Already in SignIn")
		return
	} else {
		if options, err := relay.store.NewAuthSession(signIn.Origin, signIn.Email); err != nil {
			response.InternalServerError(ctx, err)
		} else {
			response.GetResponse().WithDataSuccess(ctx, options.Response)
		}
	}
}
