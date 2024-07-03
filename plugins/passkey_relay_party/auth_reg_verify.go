package plugin_passkey_relay_party

import (
	"another_node/internal/community/account"
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

	// body works for parser, the additional info appends to query
	stubReg := seedworks.Registration{
		Origin: ctx.Query("origin"),
		Email:  ctx.Query("email"),
	}

	if user, err := relay.store.FinishRegSession(&stubReg, ctx); err != nil {
		response.BadRequest(ctx, err)
		return
	} else {
		if err := createAA(*relay.accountProvider, user); err != nil {
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
func createAA(provider account.Provider, user *seedworks.User) error {
	if w, err := account.NewHdWallet(account.HierarchicalPath_Main_ETH_TestNet); err != nil {
		return err
	} else {
		address, err := provider.CreateAccount(w)
		if err != nil {
			return err
		}
		user.SetWallet(w, address)
		return nil
	}
}
