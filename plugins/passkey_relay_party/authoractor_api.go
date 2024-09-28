package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"

	"github.com/gin-gonic/gin"
)

// imauthz is a api should be called after signin
// @Tags Demo
// @Description a demo api to show api is authorized after signin
// @Success 200 {object} any "user is authorized"
// @Failure 401 {object} any "Unauthorized"
// @Router /api/passkey/v1/imauthz [get]
// @Security JWT
func (relay *RelayParty) imauthz(ctx *gin.Context) {
	response.GetResponse().SuccessWithData(ctx, "user is authorized")
}
