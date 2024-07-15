package plugin_passkey_relay_party

import (
	"another_node/internal/community/account"
	"another_node/internal/community/chain"
	"another_node/internal/global_const"
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"
	"github.com/gin-gonic/gin"
)

// var accountProvider account.Provider
//
//	func init() {
//		p, err := impl.NewAlchemyProvider(conf.GetProvider().Alchemy)
//		if err != nil {
//			panic(err)
//		}
//
//		accountProvider = p
//	}
type finishRegistrationResponse struct {
	AccountInitCode string          `json:"account_init_code"`
	AccountAddress  string          `json:"account_address"`
	User            *seedworks.User `json:"user"`
}

func (relay *RelayParty) finishRegistration(ctx *gin.Context) {

	// body works for parser, the additional info appends to query
	stubReg := seedworks.Registration{
		Origin:  ctx.Query("origin"),
		Email:   ctx.Query("email"),
		Network: global_const.Network(ctx.Query("network")),
	}

	if user, err := relay.store.FinishRegSession(&stubReg, ctx); err != nil {
		response.BadRequest(ctx, err)
		return
	} else {
		if initCodeStr, address, err := createAA(user, stubReg.Network); err != nil {
			response.InternalServerError(ctx, err.Error())
			return
		} else {
			relay.db.Save(user)

			response.GetResponse().WithDataSuccess(ctx, finishRegistrationResponse{
				AccountInitCode: initCodeStr,
				AccountAddress:  address,
				User:            user,
			})
			return
		}
	}
}

// createAA represents creating an Account Abstraction for the user
func createAA(user *seedworks.User, network global_const.Network) (initCode string, address string, err error) {
	if w, err := account.NewHdWallet(account.HierarchicalPath_Main_ETH_TestNet); err != nil {
		return "", "", err
	} else {
		address, initCode, err := chain.CreateSmartAccount(w, network)
		if err != nil {
			return "", "", err
		}
		user.SetWallet(w, address, network)
		return address, initCode, nil
	}
}
