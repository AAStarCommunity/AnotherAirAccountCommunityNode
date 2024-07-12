package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"
	"errors"

	"github.com/gin-gonic/gin"
)

// regPrepare
// @Summary sign up step1. send code
// @Tags Plugins Passkey
// @Description Send captcha to email for verifying belongs
// @Accept json
// @Product json
// @Param registrationBody body seedworks.Registration true "Send Captcha to Email"
// @Router /api/passkey/v1/reg/prepare [post]
// @Success 200
func (relay *RelayParty) regPrepare(ctx *gin.Context) {
	var reg seedworks.Registration
	if err := ctx.ShouldBindJSON(&reg); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	err := relay.emailStartChallenge(reg.Email, ctx.GetHeader("Accept-Language"))
	if err != nil {
		response.BadRequest(ctx, "challenge failed")
		return
	}
	response.GetResponse().Success(ctx)
}

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

	if u, err := relay.FindUserByEmail(reg.Email); err != nil && !errors.Is(err, seedworks.ErrUserNotFound{}) {
		response.InternalServerError(ctx, err.Error())
		return
	} else if u != nil {
		response.BadRequest(ctx, "User already exists")
		return
	}

	if relay.emailFinishChallenge(reg.Email, reg.Captcha) != nil {
		response.BadRequest(ctx, "Invalid captcha")
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
