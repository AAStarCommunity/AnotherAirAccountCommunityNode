package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

type siginInResponse struct {
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
// @Param origin query string true "Origin"
// @Param email query string true "Email"
// @Success 200 {object} siginInResponse "OK"
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

	user, _, err := relay.store.FinishAuthSession(&stubSignIn, ctx)
	if err != nil {
		response.BadRequest(ctx, "SignIn failed: "+err.Error())
		return
	}

	relay.db.Save(user, true)

	ginJwtMiddleware().LoginHandler(ctx)
}
