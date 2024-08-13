package plugin_passkey_relay_party

import (
	"another_node/conf"
	"another_node/internal/community/account"
	"another_node/internal/community/account/impl"
	"another_node/internal/community/node"
	"another_node/plugins/passkey_relay_party/seedworks"
	storage "another_node/plugins/passkey_relay_party/storage"
	"another_node/plugins/passkey_relay_party/storage/migrations"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RelayParty struct {
	db               storage.Db
	authSessionStore *seedworks.SessionStore
	txSessionStore   *seedworks.SessionStore
	node             *node.Community
	accountProvider  *account.Provider
}

func (r *RelayParty) RegisterRoutes(router *gin.Engine, community *node.Community) {

	router.POST("/api/passkey/v1/reg/prepare", r.regPrepare)
	router.POST("/api/passkey/v1/reg", r.beginRegistration)
	router.POST("/api/passkey/v1/reg/verify", r.finishRegistration)
	router.POST("/api/passkey/v1/sign", r.beginSignIn)
	router.POST("/api/passkey/v1/sign/verify", r.finishSignIn)

	_ = router.Use(AuthHandler())
	{
		// APIs needs signin here
		router.GET("/api/passkey/v1/imauthz", r.imauthz)
		router.GET("/api/passkey/v1/account/info", r.accountInfo)
		router.POST("/api/passkey/v1/tx/sign", r.beginTxSignature)
		router.POST("/api/passkey/v1/tx/sign/verify", r.finishTxSignature)
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

// findUserByEmail finds a user by email in relay storage
func (r *RelayParty) findUserByEmail(email string) (*seedworks.User, error) {
	if email == "" {
		return nil, seedworks.ErrEmailEmpty{}
	}

	u, err := r.db.Find(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, seedworks.ErrUserNotFound{}
		}
	}

	return u, err
}
