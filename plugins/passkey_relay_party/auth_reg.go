package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// beginRegistration
// @Summary sign up step2. begin registration
// @Tags Plugins Passkey
// @Description Begin the registration process
// @Accept json
// @Product json
// @Param registrationBody body seedworks.Registration true "Begin Registration"
// @Router /api/passkey/v1/reg [post]
// @Success 200
func (relay *RelayParty) beginRegistration(ctx *gin.Context) {
	var reg seedworks.Registration
	if err := ctx.ShouldBindJSON(&reg); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// TODO: special logic for align testing
	if !strings.HasSuffix(reg.Email, "@aastar.org") && reg.Captcha != "111111" {
		if err := relay.emailChallenge(reg.Email, reg.Captcha); err != nil {
			response.BadRequest(ctx, err.Error())
			return
		}
	}

	if u, err := relay.findUserByEmail(reg.Email); err != nil && !errors.Is(err, seedworks.ErrUserNotFound{}) {
		response.InternalServerError(ctx, err.Error())
		return
	} else if u != nil {
		response.BadRequest(ctx, "User already exists")
		return
	}

	// TODO: if the user is not exists but in community, re-register the user

	if session := relay.store.Get(seedworks.GetSessionKey(reg.Origin, reg.Email)); session != nil {
		response.BadRequest(ctx, "Already in registration process")
		return
	} else {
		if options, err := relay.store.NewRegSession(&reg); err != nil {
			response.InternalServerError(ctx, err.Error())
		} else {
			response.GetResponse().WithDataSuccess(ctx, options.Response)
		}
	}
}
