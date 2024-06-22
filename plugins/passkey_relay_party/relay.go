package plugin_passkey_relay_party

import (
	"another_node/internal/community/node"
	"another_node/plugins/passkey_relay_party/seedworks"
	storage "another_node/plugins/passkey_relay_party/storage"

	"github.com/gin-gonic/gin"
)

type RelayParty struct {
	db    storage.Db
	store *seedworks.SessionStore
	node  *node.Community
}

func (r *RelayParty) RegisterRoutes(router *gin.Engine, community *node.Community) {
	router.POST("/api/passkey/v1/reg", r.beginRegistration)
	router.POST("/api/passkey/v1/reg/verify", r.finishRegistration)
	router.POST("/api/passkey/v1/sign", r.beginSignIn)
	router.POST("/api/passkey/v1/sign/verify", r.finishSignIn)

	r.node = community
}

func NewRelay() *RelayParty {
	return &RelayParty{
		db:    storage.NewInMemory(),
		store: seedworks.NewInMemorySessionStore(),
	}
}

// FindUserByEmail finds a user by email in relay storage
func (r *RelayParty) FindUserByEmail(email string) (*seedworks.User, error) {
	if email == "" {
		return nil, seedworks.EmailEmptyError{}
	}
	return r.db.Find(email)
}
