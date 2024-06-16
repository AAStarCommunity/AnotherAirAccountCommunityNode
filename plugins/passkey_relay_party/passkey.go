package plugin_passkey_relay_party

import (
	"another_node/internal/community/node"
	seedwork "another_node/plugins/passkey_relay_party/seedworks"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
)

const (
	RelayPartyDisplayName = "CommunityNode@AAStar"
	RelayPartyID          = "*"
	RelayPartyOrigins     = "*"
)

type Passkey struct {
	webAuthn *webauthn.WebAuthn
	store    *seedwork.SessionStore
}

func (p *Passkey) RegisterRoutes(router *gin.Engine, community *node.Community) {
	router.POST("/api/passkey/v1/reg", p.beginRegistration)
	router.POST("/api/passkey/v1/reg/verify", regAuthVerify)
	router.GET("/api/passkey/v1/sign", authSign)
	router.POST("/api/passkey/v1/sign/verify", authSignVerify)
}

func NewPasskeyByOrigin(id, origin string) *Passkey {
	wconfig := &webauthn.Config{
		RPDisplayName: RelayPartyDisplayName,
		RPID:          id,               // Generally the FQDN for your site
		RPOrigins:     []string{origin}, // The origin URLs allowed for WebAuthn requests
	}

	if webAuthn, err := webauthn.New(wconfig); err != nil {
		fmt.Println(err)
		return nil
	} else {
		return &Passkey{
			webAuthn: webAuthn,
			store:    seedwork.NewInMemorySessionStore(),
		}
	}
}

func NewPasskeyAuth() *Passkey {
	wconfig := &webauthn.Config{
		RPDisplayName: RelayPartyDisplayName,
		RPID:          RelayPartyID,                // Generally the FQDN for your site
		RPOrigins:     []string{RelayPartyOrigins}, // The origin URLs allowed for WebAuthn requests
	}

	if webAuthn, err := webauthn.New(wconfig); err != nil {
		fmt.Println(err)
		return nil
	} else {
		return &Passkey{
			webAuthn: webAuthn,
			store:    seedwork.NewInMemorySessionStore(),
		}
	}
}
