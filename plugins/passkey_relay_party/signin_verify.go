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

// finishSignIn
// @Summary sign in step 2. verify credentials
// @Description Finish the sign in process
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
		Registration: seedworks.Registration{
			RegistrationPrepare: seedworks.RegistrationPrepare{
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
