package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
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

// finishSignIn
// @Summary Finish SingIn
// @Description Verify attestations and return JWT
// @Tags Plugins Passkey
// @Accept json
// @Produce json
// @Param email  query string true "user email" Format(email)
// @Param origin query string true "origin"
// @Param signinBody body protocol.CredentialAssertionResponse true "Verify SignIn"
// @Success 200 {object} SiginInResponse "OK"
// @Failure 400 {object} any "Bad Request"
// @Router /api/passkey/v1/sign/verify [post]
func (relay *RelayParty) finishSignIn(ctx *gin.Context) {
	// body works for SDK, the additional info appends to query
	stubSignIn := seedworks.SiginIn{
		RegistrationByEmail: seedworks.RegistrationByEmail{
			RegistrationByEmailPrepare: seedworks.RegistrationByEmailPrepare{
				Email: ctx.Query("email"),
			},
			Origin: ctx.Query("origin"),
		},
	}

	user, _, err := relay.authSessionStore.FinishAuthSession(&stubSignIn, ctx)
	if err != nil {
		response.GetResponse().FailCode(ctx, 401, "SignIn failed: "+err.Error())
		return
	}

	relay.db.Save(user, true)

	ginJwtMiddleware().LoginHandler(ctx)
}
