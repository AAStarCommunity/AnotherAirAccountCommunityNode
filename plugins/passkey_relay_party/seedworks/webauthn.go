package seedworks

import (
	"another_node/plugins/passkey_relay_party/conf"
	"fmt"
	"net/url"

	"github.com/go-webauthn/webauthn/webauthn"
)

func newWebAuthn(origin string) (*webauthn.WebAuthn, error) {
	u, err := url.Parse(origin)
	if err != nil {
		return nil, err
	}
	hostname := u.Hostname()
	wconfig := &webauthn.Config{
		RPDisplayName: conf.Get().RelayParty.DisplayName,
		RPID:          conf.Get().RelayParty.Id,   // Generally the FQDN for your site
		RPOrigins:     []string{origin, hostname}, // The origin URLs allowed for WebAuthn requests
	}

	if webAuthn, err := webauthn.New(wconfig); err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		return webAuthn, nil
	}
}
