package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// createAA
// @Summary Create AA with Purpose
// @Tags Plugins Passkey
// @Description create aa by sepcify network(chain)
// @Accept json
// @Product json
// @Param createAABody body seedworks.CreateAARequest true "Create AA"
// @Router /api/passkey/v1/account/chain [post]
// @Success 200
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

	// TODO: support memo
	if err := user.TryCreateAA(req.Network); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	} else {
		if err := relay.db.SaveAccounts(user, req.Network); err != nil {
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
