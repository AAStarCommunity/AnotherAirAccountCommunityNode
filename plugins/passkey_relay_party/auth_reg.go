package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"
	"errors"

	"github.com/gin-gonic/gin"
)

func (relay *RelayParty) beginRegistration(ctx *gin.Context) {
	var reg seedworks.Registration
	if err := ctx.ShouldBindJSON(&reg); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if u, err := relay.FindUserByEmail(reg.Email); err != nil && !errors.Is(err, seedworks.UserNotFoundError{}) {
		response.InternalServerError(ctx, err.Error())
		return
	} else if u != nil {
		response.BadRequest(ctx, "User already exists")
		return
	}

	// TODO: if the user is not exists but in community, re-register the user

	if session := relay.store.Get(seedworks.GetSessionKey(&reg)); session != nil {
		response.BadRequest(ctx, "Already in registration")
		return
	} else {
		if options, err := relay.store.NewRegSession(&reg); err != nil {
			response.InternalServerError(ctx, err.Error())
		} else {
			response.GetResponse().WithDataSuccess(ctx, options.Response)
		}
	}
}
