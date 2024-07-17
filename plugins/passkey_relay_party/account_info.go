package plugin_passkey_relay_party

import (
	consts "another_node/internal/seedworks"
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"
	"net/http"

	"github.com/gin-gonic/gin"
)

// accountInfo user account info
// @Summary get user account info
// @Description get user account info: eoa address, aa address, init code
// @Tags Plugins Passkey
// @Accept json
// @Produce json
// @Param network  query string true "network"
// @Success 200 {object} seedworks.accountInfo "OK"
// @Failure 400 {object} any "Bad Request"
// @Router /api/passkey/v1/account/info [get]
// @Security JWT
func (relay *RelayParty) accountInfo(ctx *gin.Context) {
	if ok, email := CurrentUser(ctx); !ok {
		response.GetResponse().FailCode(ctx, http.StatusUnauthorized)
		return
	} else {
		// TODO: for tokyo ONLY
		chain := consts.Chain(ctx.Query("network"))
		if len(chain) > 0 && chain != consts.OptimismSepolia {
			response.BadRequest(ctx, "network not supported")
			return
		} else {
			chain = consts.OptimismSepolia
		}

		if initCode, addr, eoaAddr, err := relay.db.GetAccounts(email, string(chain)); err != nil {
			response.NotFound(ctx, err.Error())
		} else {
			response.GetResponse().WithDataSuccess(ctx, seedworks.AccountInfo{
				InitCode: initCode,
				AA:       addr,
				EOA:      eoaAddr,
			})
		}
	}
}
