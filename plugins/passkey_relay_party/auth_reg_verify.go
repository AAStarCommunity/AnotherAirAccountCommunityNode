package plugin_passkey_relay_party

import (
	consts "another_node/internal/seedworks"

	"another_node/internal/community/account"
	"another_node/internal/community/chain"
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

// finishRegistration
// @Summary sign up step3. Finish Registration
// @Tags Plugins Passkey
// @Description Verify Passkey Registration
// @Accept json
// @Product json
// @Param registrationBody body seedworks.Registration true "Verify Registration"
// @Router /api/passkey/v1/reg/verify [post]
// @Success 200
func (relay *RelayParty) finishRegistration(ctx *gin.Context) {

	// TODO: for tokyo ONLY
	network := consts.Network(ctx.Query("network"))
	if len(network) > 0 && network != consts.OptimismSepolia {
		response.BadRequest(ctx, "network not supported")
		return
	} else {
		network = consts.OptimismSepolia
	}

	// body works for parser, the additional info appends to query
	stubReg := seedworks.Registration{
		RegistrationPrepare: seedworks.RegistrationPrepare{
			Email: ctx.Query("email"),
		},
		Origin:  ctx.Query("origin"),
		Network: consts.Network(ctx.Query("network")),
	}

	if user, err := relay.store.FinishRegSession(&stubReg, ctx); err != nil {
		response.BadRequest(ctx, err)
		return
	} else {
		if err := createAA(user, stubReg.Network); err != nil {
			response.InternalServerError(ctx, err.Error())
			return
		} else {
			relay.db.Save(user)
			response.GetResponse().WithDataSuccess(ctx, user)
			return
		}
	}
}

// createAA represents creating an Account Abstraction for the user
func createAA(user *seedworks.User, network consts.Network) error {
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
