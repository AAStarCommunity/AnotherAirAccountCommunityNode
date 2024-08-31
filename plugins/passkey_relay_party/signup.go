package plugin_passkey_relay_party

import (
	"another_node/internal/community/account"
	"another_node/internal/community/chain"
	consts "another_node/internal/seedworks"
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// regPrepareByEmail
// @Summary Prepare SignUp
// @Tags Plugins Passkey
// @Description Send captcha to email for confirming ownership
// @Accept json
// @Product json
// @Param registrationBody body seedworks.RegistrationPrepare true "Send Captcha to Email"
// @Router /api/passkey/v1/reg/prepare [post]
// @Success 200
func (relay *RelayParty) regPrepareByEmail(ctx *gin.Context) {
	var reg seedworks.RegistrationByEmail
	if err := ctx.ShouldBindJSON(&reg); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	err := relay.emailStartChallenge(reg.Email, ctx.GetHeader("Accept-Language"))
	if err != nil {
		response.BadRequest(ctx, "challenge failed", err.Error())
		return
	}
	response.GetResponse().Success(ctx)
}

// beginRegistrationByEmail
// @Summary Begin SignUp
// @Tags Plugins Passkey
// @Description Send challenge for passkey
// @Accept json
// @Product json
// @Param registrationBody body seedworks.Registration true "Begin Registration"
// @Router /api/passkey/v1/reg [post]
// @Success 200 {object} protocol.PublicKeyCredentialCreationOptions
func (relay *RelayParty) beginRegistrationByEmail(ctx *gin.Context) {
	var reg seedworks.RegistrationByEmail
	if err := ctx.ShouldBindJSON(&reg); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if err := relay.emailChallenge(reg.Email, reg.Captcha); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if u, err := relay.findUserByEmail(reg.Email); err != nil && !errors.Is(err, seedworks.ErrUserNotFound{}) {
		response.InternalServerError(ctx, err.Error())
		return
	} else if u != nil {
		response.BadRequest(ctx, "User already exists")
		return
	}

	// TODO: if the user is not exists but in community, re-register the user

	sessionKey := seedworks.GetSessionKey(reg.Origin, reg.Email)
	if session := relay.authSessionStore.Get(sessionKey); session != nil {
		relay.authSessionStore.Remove(sessionKey)
	}

	if options, err := relay.authSessionStore.NewRegSession(&reg); err != nil {
		response.InternalServerError(ctx, err.Error())
	} else {
		response.GetResponse().WithDataSuccess(ctx, options.Response)
	}
}

//lint:ignore U1000 because it's used in the swagger
type finishRegistrationResponse struct {
	AccountInitCode string `json:"account_init_code"`
	AccountAddress  string `json:"account_address"`
	EoaAddress      string `json:"eoa_address"`
}

// finishRegistrationByEmail
// @Summary Finish SignUp By Email
// @Tags Plugins Passkey
// @Description Verify attestations, register user and return JWT
// @Accept json
// @Product json
// @Param email  query string true "user email" Format(email)
// @Param origin query string true "origin"
// @Param network query string false "network"
// @Param registrationBody body protocol.CredentialCreationResponse true "Verify Registration"
// @Router /api/passkey/v1/reg/verify [post]
// @Success 200 {object} SiginInResponse "OK"
func (relay *RelayParty) finishRegistrationByEmail(ctx *gin.Context) {

	network := consts.Chain(ctx.Query("network"))

	if !isSupportChain(network) {
		response.BadRequest(ctx, "network not supported, please specify a valid network, e.g.: optimism-mainnet, base-sepolia, optimism-sepolia")
		return
	}

	// body-stream works for parser, the additional info appends to query
	stubReg := seedworks.FinishRegistrationByEmail{
		RegistrationByEmailPrepare: seedworks.RegistrationByEmailPrepare{
			Email: ctx.Query("email"),
		},
		Origin:  ctx.Query("origin"),
		Network: consts.Chain(ctx.Query("network")),
	}

	if user, err := relay.authSessionStore.FinishRegSession(&stubReg, ctx); err != nil {
		response.GetResponse().FailCode(ctx, 401, "SignUp failed: "+err.Error())
	} else {
		signup(relay, ctx, &stubReg, user)
	}
}

func signup(relay *RelayParty, ctx *gin.Context, reg *seedworks.FinishRegistrationByEmail, user *seedworks.User) {
	// TODO: special logic for align testing
	if strings.HasSuffix(reg.Email, "@aastar.org") {
		response.GetResponse().WithDataSuccess(ctx, user)
		return
	}

	if err := initAA(user, reg.Network); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	} else {
		relay.db.SaveAccounts(user, reg.Network)

		ginJwtMiddleware().LoginHandler(ctx)

		return
	}
}

// initAA represents creating an Account Abstraction for the user
func initAA(user *seedworks.User, network consts.Chain) (err error) {
	if w, err := account.NewHdWallet(account.HierarchicalPath_ETH); err != nil {
		return err
	} else {
		aa_address, initCode, err := chain.CreateSmartAccount(w, network)
		if err != nil {
			return err
		}
		user.SetWallet(w, &initCode, &aa_address, &w.Address, network)
		return nil
	}
}
