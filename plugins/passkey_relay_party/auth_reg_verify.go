package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
)

func (relay *RelayParty) finishRegistration(ctx *gin.Context) {
	parsedResponse, err := protocol.ParseCredentialCreationResponse(ctx.Request)
	if err != nil {
		response.BadRequest(ctx, err)
		return
	}

	// body works for parser, the additional info appends to query
	stubReg := seedworks.Registration{
		Origin: ctx.Query("origin"),
		Email:  ctx.Query("email"),
	}

	key := seedworks.GetSessionKey(stubReg)
	if session := relay.store.Get(key); session == nil {
		response.BadRequest(ctx, "Session not found")
		return
	} else {
		if credential, err := session.WebAuthn.CreateCredential(&session.User, session.Data, parsedResponse); err != nil {
			response.BadRequest(ctx, err)
			return
		} else {
			// TODO: save credential to user
			response.GetResponse().WithDataSuccess(ctx, credential)
		}
	}
}
