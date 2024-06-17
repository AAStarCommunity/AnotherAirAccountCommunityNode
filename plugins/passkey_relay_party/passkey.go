package plugin_passkey_relay_party

import (
	"another_node/internal/community/node"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

type RelayParty struct {
	store *seedworks.SessionStore
}

func (r *RelayParty) RegisterRoutes(router *gin.Engine, community *node.Community) {
	router.POST("/api/passkey/v1/reg", r.beginRegistration)
	router.POST("/api/passkey/v1/reg/verify", r.finishRegistration)
	router.GET("/api/passkey/v1/sign", r.beginSignIn)
	router.POST("/api/passkey/v1/sign/verify", r.finishSignIn)
}

func NewPasskey() *RelayParty {
	return &RelayParty{
		store: seedworks.NewInMemorySessionStore(),
	}
}
