package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
)

func (relay *RelayParty) regAuthVerify(ctx *gin.Context) {
	parsedResponse, err := protocol.ParseCredentialCreationResponse(ctx.Request)
	if err != nil {
		response.BadRequest(ctx, err)
		return
	}

	if session := relay.store.Get(parsedResponse.ID); session == nil {
		response.BadRequest(ctx, "Session not found")
		return
	} else {
		if k, err := session.WebAuthn.CreateCredential(&session.User, session.Data, parsedResponse); err != nil {
			response.BadRequest(ctx, err)
			return
		} else {
			response.GetResponse().WithDataSuccess(ctx, k)
		}
	}
}
