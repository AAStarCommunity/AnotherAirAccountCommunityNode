package plugin_passkey_relay_party

import (
	"another_node/conf"
	"another_node/internal/community/wallet"
	"another_node/internal/community/wallet/impl"
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

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
		if err := createByAlchemy(user); err != nil {
			response.InternalServerError(ctx, err)
			return
		} else {
			relay.db.Save(user)
				response.GetResponse().WithDataSuccess(ctx, user)
			return
		}
	}
}

func createByAlchemy(user *seedworks.User) error {
	apiKey := conf.GetProvider().Alchemy
	if porvider, err := impl.NewAlchemyProvider(apiKey); err != nil {
		return err
	} else {
		return createAA(porvider, user)
	}
}

// createAA represents creating an Account Abstraction for the user
func createAA(provider wallet.Provider, user *seedworks.User) error {
	if w, err := wallet.NewHdWallet(wallet.HierarchicalPath_Main_ETH_TestNet); err != nil {
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
