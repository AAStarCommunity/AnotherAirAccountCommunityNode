package plugin_passkey_relay_party

import (
	"another_node/internal/community/node"
	seedwork "another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

const (
	RelayPartyDisplayName = "CommunityNode@AAStar"
	RelayPartyID          = "*"
	RelayPartyOrigins     = "*"
)

type RelayParty struct {
	store *seedwork.SessionStore
}

func (r *RelayParty) RegisterRoutes(router *gin.Engine, community *node.Community) {
	router.POST("/api/passkey/v1/reg", r.beginRegistration)
	router.POST("/api/passkey/v1/reg/verify", r.regAuthVerify)
	router.GET("/api/passkey/v1/sign", authSign)
	router.POST("/api/passkey/v1/sign/verify", authSignVerify)
}

func NewPasskey() *RelayParty {
	return &RelayParty{
		store: seedwork.NewInMemorySessionStore(),
	}
}
