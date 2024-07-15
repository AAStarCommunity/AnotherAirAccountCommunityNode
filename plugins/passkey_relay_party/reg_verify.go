package plugin_passkey_relay_party

import (
	consts "another_node/internal/seedworks"
	"strings"

	"another_node/internal/community/account"
	"another_node/internal/community/chain"
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

type regVerifyResponse struct {
	Code   int    `json:"code"`
	Token  string `json:"token"`
	Expire string `json:"expire"`
}

// finishRegistration
// @Summary sign up step3. Finish Registration
// @Tags Plugins Passkey
// @Description Verify Passkey Registration
// @Accept json
// @Product json
// @Param email  query string true "user email" Format(email)
// @Param origin query string true "origin"
// @Param network query string false "network"
// @Param registrationBody body protocol.CredentialCreationResponse true "Verify Registration"
// @Router /api/passkey/v1/reg/verify [post]
// @Success 200 {object} regVerifyResponse
func (relay *RelayParty) finishRegistration(ctx *gin.Context) {

	// TODO: for tokyo ONLY
	network := consts.Chain(ctx.Query("network"))
	if len(network) > 0 && network != consts.OptimismSepolia {
		response.BadRequest(ctx, "network not supported")
		return
	} else {
		network = consts.OptimismSepolia
	}

	// body works for parser, the additional info appends to query
	stubReg := seedworks.FinishRegistration{
		RegistrationPrepare: seedworks.RegistrationPrepare{
			Email: ctx.Query("email"),
		},
		Origin:  ctx.Query("origin"),
		Network: consts.Chain(ctx.Query("network")),
	}

	if user, err := relay.authSessionStore.FinishRegSession(&stubReg, ctx); err != nil {
		response.BadRequest(ctx, err)
		return
	} else {
		// TODO: special logic for align testing
		if strings.HasSuffix(stubReg.Email, "@aastar.org") {
			response.GetResponse().WithDataSuccess(ctx, user)
			return
		}
		if err := createAA(user, stubReg.Network); err != nil {
			response.InternalServerError(ctx, err.Error())
			return
		} else {
			relay.db.Save(user, false)

			ginJwtMiddleware().LoginHandler(ctx)
			return
		}
	}
}

// createAA represents creating an Account Abstraction for the user
func createAA(user *seedworks.User, network consts.Chain) error {
	if w, err := account.NewHdWallet(account.HierarchicalPath_ETH); err != nil {
		return err
	} else {
		address, err := chain.CreateSmartAccount(w, network)
		if err != nil {
			return err
		}
		user.SetWallet(w, address, network)
		return nil
	}
}
