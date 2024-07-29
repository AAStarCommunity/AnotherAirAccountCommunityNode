package plugin_passkey_relay_party

import (
	jwt2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func CurrentUser(ctx *gin.Context) (exists bool, email string) {

	defer func() {
		if r := recover(); r != nil {
			exists = false
		}
	}()

	if j, ok := ctx.Keys["JWT_PAYLOAD"]; !ok {
		exists = false
		return exists, ""
	} else {
		mapping := j.(jwt2.MapClaims)
		email = mapping["sub"].(string)
		exists = true
	}

	return
}
