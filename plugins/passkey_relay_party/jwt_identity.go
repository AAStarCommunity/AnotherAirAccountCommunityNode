package plugin_passkey_relay_party

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CurrentUser(ctx *gin.Context) (exists bool, email string) {

	defer func() {
		if r := recover(); r != nil {
			exists = false
		}
	}()

	mapping := ctx.MustGet("JWT_PAYLOAD").(jwt.MapClaims)

	email = mapping["sub"].(string)
	exists = true

	return
}
