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
	db               storage.Db
	authSessionStore *seedworks.SessionStore
	txSessionStore   *seedworks.SessionStore
	node             *node.Community
	accountProvider  *account.Provider
}

func (r *RelayParty) RegisterRoutes(router *gin.Engine, community *node.Community) {

	router.POST("/api/passkey/v1/reg/prepare", r.regPrepareByEmail)
	router.POST("/api/passkey/v1/reg", r.beginRegistrationByEmail)
	router.POST("/api/passkey/v1/reg-account", r.beginRegistrationByAccount)
	router.POST("/api/passkey/v1/reg/verify", r.finishRegistrationByEmail)
	router.POST("/api/passkey/v1/reg-account/verify", r.finishRegistrationByAccount)
	router.POST("/api/passkey/v1/sign", r.beginSignIn)
	router.POST("/api/passkey/v1/sign/verify", r.finishSignIn)
	router.GET("/api/passkey/v1/chains/support", r.supportChains)

	_ = router.Use(AuthHandler())
	{
		// APIs needs signin here
		router.GET("/api/passkey/v1/imauthz", r.imauthz)
		router.GET("/api/passkey/v1/account/info", r.getAccountInfo)
		router.POST("/api/passkey/v1/tx/sign", r.beginTxSignature)
		router.POST("/api/passkey/v1/tx/sign/verify", r.finishTxSignature)
		router.POST("/api/passkey/v1/account/chain", r.createAA)
	}

	r.node = community
}

func NewRelay() *RelayParty {
	migrations.AutoMigrate()

	return &RelayParty{
		db:               storage.NewPgsqlStorage(),
		authSessionStore: seedworks.NewInMemorySessionStore(),
		txSessionStore:   seedworks.NewInMemorySessionStore(),
		accountProvider: func() *account.Provider {
			p, err := impl.NewAlchemyProvider(conf.GetProvider().Alchemy)
			if err != nil {
				panic(err)
			}
			return &p
		}(),
	}
}
