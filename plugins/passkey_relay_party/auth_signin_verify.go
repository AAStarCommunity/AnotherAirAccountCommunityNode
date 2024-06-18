package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

func (relay *RelayParty) finishSignIn(ctx *gin.Context) {
	// body works for SDK, the additional info appends to query
	stubSignIn := seedworks.SiginIn{
		Registration: seedworks.Registration{
			Origin: ctx.Query("origin"),
			Email:  ctx.Query("email"),
		},
	}

	user, credential, err := relay.store.FinishAuthSession(&stubSignIn, ctx)
	if err != nil {
		response.BadRequest(ctx, "SignIn failed: "+err.Error())
		return
	}
	relay.db.Save(user)
	response.GetResponse().WithDataSuccess(ctx, credential)
}
