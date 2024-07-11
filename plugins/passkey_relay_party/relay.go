package plugin_passkey_relay_party

import (
	"another_node/conf"
	"another_node/internal/community/account"
	"another_node/internal/community/account/impl"
	"another_node/internal/community/node"
	"another_node/plugins/passkey_relay_party/seedworks"
	storage "another_node/plugins/passkey_relay_party/storage"
	"another_node/plugins/passkey_relay_party/storage/migrations"

	"github.com/gin-gonic/gin"
)

type RelayParty struct {
	db              storage.Db
	store           *seedworks.SessionStore
	node            *node.Community
	accountProvider *account.Provider
}

func (r *RelayParty) RegisterRoutes(router *gin.Engine, community *node.Community) {
	router.POST("/api/passkey/v1/reg/prepare", r.regPrepare)
	router.POST("/api/passkey/v1/reg", r.beginRegistration)
	router.POST("/api/passkey/v1/reg/verify", r.finishRegistration)
	router.POST("/api/passkey/v1/sign", r.beginSignIn)
	router.POST("/api/passkey/v1/sign/verify", r.finishSignIn)

	r.node = community
}

func NewRelay() *RelayParty {
	migrations.AutoMigrate()

	return &RelayParty{
		db:    storage.NewPgsqlStorage(),
		store: seedworks.NewInMemorySessionStore(),
		accountProvider: func() *account.Provider {
			p, err := impl.NewAlchemyProvider(conf.GetProvider().Alchemy)
			if err != nil {
				panic(err)
			}
			return &p
		}(),
	}
}

// FindUserByEmail finds a user by email in relay storage
func (r *RelayParty) FindUserByEmail(email string) (*seedworks.User, error) {
	if email == "" {
		return nil, seedworks.ErrEmailEmpty{}
	}
	return r.db.Find(email)
}
