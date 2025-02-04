package plugin_passkey_relay_party

import (
	consts "another_node/internal/seedworks"
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getAccountInfo represents for getting account info of user
// @Summary Get User Account Info
// @Description get user account info: eoa address, aa address, init code
// @Tags Plugins Passkey
// @Accept json
// @Produce json
// @Param network  query string true "network"
// @Param alias  query string false "alias"
// @Success 200 {object} seedworks.AccountInfo "OK"
// @Failure 400 {object} any "Bad Request"
// @Router /api/passkey/v1/account/info [get]
// @Security JWT
func (relay *RelayParty) getAccountInfo(ctx *gin.Context) {
	if ok, email := CurrentUser(ctx); !ok {
		response.GetResponse().FailCode(ctx, http.StatusUnauthorized)
		return
	} else {
		network := consts.Chain(ctx.Query("network"))
		alias := ctx.Query("alias")
		if !isSupportChain(network) {
			response.BadRequest(ctx, "network not supported, please specify a valid network, e.g.: optimism-mainnet, base-sepolia, optimism-sepolia")
			return
		}

		if user, err := relay.db.FindUser(email); err != nil {
			response.GetResponse().SuccessWithDataAndCode(http.StatusNotFound, ctx, err.Error())
			return
		} else {
			chain := user.GetSpecifiyChain(network, alias)
			if chain == nil {
				response.GetResponse().SuccessWithDataAndCode(http.StatusNotFound, ctx, seedworks.ErrChainNotFound{})
				return
			}

			response.GetResponse().SuccessWithData(ctx, seedworks.AccountInfo{
				InitCode: chain.InitCode,
				AA:       chain.AA_Addr,
				EOA:      user.GetEOA(chain),
				Email:    email,
			})
			return
		}
	}
}
