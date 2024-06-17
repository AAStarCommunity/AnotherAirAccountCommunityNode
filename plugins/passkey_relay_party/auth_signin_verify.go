package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

func (relay *RelayParty) finishSignIn(ctx *gin.Context) {
	var user seedworks.User

	// body works for SDK, the additional info appends to query
	stubReg := seedworks.Registration{
		Origin: ctx.Query("origin"),
		Email:  ctx.Query("email"),
	}

	session := relay.store.Get(seedworks.GetSessionKey(stubReg))

	if session == nil {
		response.BadRequest(ctx, "Session not found")
		return
	}

	credential, err := session.WebAuthn.FinishLogin(&user, session.Data, ctx.Request)
	if err != nil {
		response.BadRequest(ctx, "SignIn failed: "+err.Error())
		return
	}
	// TODO: save credential to user
	response.GetResponse().WithDataSuccess(ctx, credential)
}
