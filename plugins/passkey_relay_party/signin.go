package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

// beginSignIn
// @Summary sign in step 1. request credential assertion
// @Description Begins the sign in process
// @Tags Plugins Passkey
// @Accept json
// @Produce json
// @Param signIn body seedworks.SiginIn true "Sign in details"
// @Success 200 {object} protocol.PublicKeyCredentialRequestOptions
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/passkey/v1/sign [post]
func (relay *RelayParty) beginSignIn(ctx *gin.Context) {
	var signIn seedworks.SiginIn
	if err := ctx.ShouldBindJSON(&signIn); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	if session := relay.authSessionStore.Get(seedworks.GetSessionKey(signIn.Origin, signIn.Email)); session != nil {
		response.BadRequest(ctx, "Already in SignIn")
		return
	} else {
		user, err := relay.db.Find(signIn.Email)
		if err != nil {
			response.NotFound(ctx, err.Error())
		}
		if options, err := relay.authSessionStore.NewAuthSession(user, &signIn); err != nil {
			response.InternalServerError(ctx, err)
		} else {
			response.GetResponse().WithDataSuccess(ctx, options.Response)
		}
	}
}
