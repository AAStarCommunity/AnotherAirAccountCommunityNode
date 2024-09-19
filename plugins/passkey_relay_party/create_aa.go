package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// supportChains
// @Summary Get support chains in relay party
// @Tags Plugins Passkey
// @Description get support chains
// @Accept json
// @Product json
// @Router /api/passkey/v1/chains/support [get]
// @Success 200
func (relay *RelayParty) supportChains(ctx *gin.Context) {
	response.GetResponse().WithDataSuccess(ctx, supportChains)
}

// createAA
// @Summary Create AA with Alias, default empty
// @Tags Plugins Passkey
// @Description create aa by sepcify network(chain)
// @Accept json
// @Product json
// @Param createAABody body seedworks.CreateAARequest true "Create AA"
// @Router /api/passkey/v1/account/chain [post]
// @Success 200
// @Security JWT
func (relay *RelayParty) createAA(ctx *gin.Context) {
	if ok, email := CurrentUser(ctx); !ok {
		response.GetResponse().FailCode(ctx, http.StatusUnauthorized)
		return
	} else {
		if user, err := relay.db.FindUser(email); err != nil {
			response.NotFound(ctx, err.Error())
			return
		} else {
			createAA(relay, user, ctx)
			return
		}
	}
}

func createAA(relay *RelayParty, user *seedworks.User, ctx *gin.Context) {
	var req seedworks.CreateAARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if !isSupportChain(req.Network) {
		response.BadRequest(ctx, "network not supported, please specify a valid network, e.g.: optimism-mainnet, base-sepolia, optimism-sepolia")
		return
	}

	if err := user.TryCreateAA(req.Network, req.Alias); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	} else {
		if err := relay.db.SaveAccounts(user); err != nil {
			if errors.Is(err, seedworks.ErrUserAlreadyExists{}) {
				response.BadRequest(ctx, err.Error())
			} else {
				response.InternalServerError(ctx, err.Error())
			}
			return
		}

		response.Success(ctx)

		return
	}
}
