package plugin_passkey_relay_party

import (
	"another_node/conf"
	"another_node/internal/community/account"
	"another_node/internal/community/account/impl"
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

var accountProvider account.Provider

func init() {
	p, err := impl.NewAlchemyProvider(conf.GetProvider().Alchemy)
	if err != nil {
		panic(err)
	}

	accountProvider = p
}

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
		if err := createAA(accountProvider, user); err != nil {
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
