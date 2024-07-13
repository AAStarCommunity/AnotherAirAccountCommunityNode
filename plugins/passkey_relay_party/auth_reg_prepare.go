package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

// regPrepare
// @Summary sign up step1. send code
// @Tags Plugins Passkey
// @Description Send captcha to email for verifying belongs
// @Accept json
// @Product json
// @Param registrationBody body seedworks.RegistrationPrepare true "Send Captcha to Email"
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
