package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
)

type SiginInResponse struct {
	Code   int    `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

// beginSignIn
// @Summary Begins SignIn
// @Description Send challenge for passkey
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

	if options, err := relay.authSessionStore.BeginDiscoverableAuthSession(&signIn); err != nil {
		response.InternalServerError(ctx, err)
	} else {
		response.GetResponse().WithDataSuccess(ctx, options.Response)
	}
}

// finishSignIn
// @Summary Finish SingIn
// @Description Verify attestations and return JWT
// @Tags Plugins Passkey
// @Accept json
// @Produce json
// @Param origin query string true "origin"
// @Param signinBody body protocol.CredentialAssertionResponse true "Verify SignIn"
// @Success 200 {object} SiginInResponse "OK"
// @Failure 400 {object} any "Bad Request"
// @Router /api/passkey/v1/sign/verify [post]
func (relay *RelayParty) finishSignIn(ctx *gin.Context) {
	// body works for SDK, the additional info appends to query
	stubSignIn := seedworks.SiginIn{
		Origin: ctx.Query("origin"),
	}

	var user *seedworks.User
	_, err := relay.authSessionStore.FinishDiscoverableAuthSession(&stubSignIn, ctx, func(rawID, userHandle []byte) (webauthn.User, error) {
		var err error
		user, err = relay.db.FindUser(string(userHandle))
		return user, err
	})
	if err != nil {
		response.GetResponse().FailCode(ctx, 401, "SignIn failed: "+err.Error())
		return
	}

	ctx.Set("user", user)
	ginJwtMiddleware().LoginHandler(ctx)
}
