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
// @Success 200 {object} seedworks.AccountInfo "OK"
// @Failure 400 {object} any "Bad Request"
// @Router /api/passkey/v1/account/info [get]
// @Security JWT
func (relay *RelayParty) getAccountInfo(ctx *gin.Context) {
	if ok, email := CurrentUser(ctx); !ok {
		response.GetResponse().FailCode(ctx, http.StatusUnauthorized)
		return
	} else {
		chain := consts.Chain(ctx.Query("network"))
		if len(chain) > 0 {
			if chain != consts.OptimismSepolia && chain != consts.BaseSepolia {
				response.BadRequest(ctx, "network not supported")
				return
			}
		} else {
			chain = consts.OptimismSepolia
		}

		if user, err := relay.db.FindUser(email); err != nil {
			response.NotFound(ctx, err.Error())
		} else {
			initCode, aaAddr, eoaAddr := user.GetChainAddresses(chain)
			if aaAddr == nil || eoaAddr == nil || initCode == nil {
				response.NotFound(ctx, seedworks.ErrChainNotFound{})
				return
			}

			response.GetResponse().WithDataSuccess(ctx, seedworks.AccountInfo{
				InitCode: *initCode,
				AA:       *aaAddr,
				EOA:      *eoaAddr,
				Email:    email,
			})
		}
	}
}
