package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	seedwork "another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

func sessionKey(reg seedwork.Registration) string {
	return reg.Account + ":" + reg.Email
}
func (passkey Passkey) beginRegistration(ctx *gin.Context) {
	var reg seedwork.Registration
	if err := ctx.ShouldBindJSON(&reg); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	passkey.webAuthn = NewPasskeyByOrigin(reg.Origin, reg.Origin).webAuthn

	if passkey.store.Get(sessionKey(reg)) != nil {
		response.BadRequest(ctx, "Already in registration")
		return
	}

	user := seedwork.NewUser("ab", "de", "ff")
	if options, session, err := passkey.webAuthn.BeginRegistration(user); err != nil {
		response.InternalServerError(ctx, err)
	} else {
		passkey.store.Set(sessionKey(reg), session)
		response.GetResponse().WithDataSuccess(ctx, options.Response)
	}
}
