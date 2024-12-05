package plugin_passkey_relay_party

import (
	"another_node/internal/community/account"
	consts "another_node/internal/seedworks"
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

//lint:ignore U1000 because it's used in the swagger
type finishRegistrationResponse struct {
	AccountInitCode string `json:"account_init_code"`
	AccountAddress  string `json:"account_address"`
	EoaAddress      string `json:"eoa_address"`
}

// regPrepareByEmail
// @Summary Prepare SignUp
// @Tags Plugins Passkey
// @Description Send captcha to email for confirming ownership
// @Accept json
// @Product json
// @Param registrationBody body seedworks.RegistrationByEmail true "Send Captcha to Email"
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
// @Param registrationBody body seedworks.RegistrationByEmail true "Begin Registration"
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

	sessionKey := seedworks.GetSessionKey(reg.Origin, reg.Email)
	if session := relay.authSessionStore.Get(sessionKey); session != nil {
		relay.authSessionStore.Remove(sessionKey)
	}

	if options, err := relay.authSessionStore.BeginRegSession(&reg); err != nil {
		response.InternalServerError(ctx, err.Error())
	} else {
		response.GetResponse().SuccessWithData(ctx, options.Response)
	}
}

// RegistrationByAccount
// @Summary Begin SignUp By EOA/UnionID/etc.,
// @Tags Plugins Passkey
// @Accept json
// @Product json
// @Param registrationBody body seedworks.RegistrationByAccount true "Begin Registration"
// @Router /api/passkey/v1/reg-account [post]
// @Success 200 {object} protocol.PublicKeyCredentialCreationOptions
func (relay *RelayParty) beginRegistrationByAccount(ctx *gin.Context) {
	var reg seedworks.RegistrationByAccount
	if err := ctx.ShouldBindJSON(&reg); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	sessionKey := seedworks.GetSessionKey(reg.Origin, reg.Account, reg.Type)
	if session := relay.authSessionStore.Get(sessionKey); session != nil {
		relay.authSessionStore.Remove(sessionKey)
	}

	if options, err := relay.authSessionStore.BeginRegSession(&reg); err != nil {
		response.InternalServerError(ctx, err.Error())
	} else {
		response.GetResponse().SuccessWithData(ctx, options.Response)
	}
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
// @Param alias query string false "network"
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
		Origin: ctx.Query("origin"),
	}

	if user, err := relay.authSessionStore.FinishRegSession(&stubReg, ctx); err != nil {
		response.GetResponse().FailCode(ctx, 401, "SignUp failed: "+err.Error())
	} else {
		signup(relay, ctx, user)
	}
}

const defaultWalletCount = 5

func createWalletsForNewUser(relay *RelayParty, ctx *gin.Context, user *seedworks.User) {
	paths := make([]account.HierarchicalPath, defaultWalletCount)
	for i := 0; i < defaultWalletCount && i < 10; i++ {
		paths[i] = account.HierarchicalPath(fmt.Sprintf(account.HierarchicalPath_ETH_FMT, i))
	}
	wallets, err := account.NewHdWallet(paths...)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	if err := relay.db.CreateAccount(user.GetEmail(), wallets); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
}

func signup(relay *RelayParty, ctx *gin.Context, user *seedworks.User) {

	createWalletsForNewUser(relay, ctx, user)

	// check if user passkey already exists
	if len(user.WebAuthnCredentials()) != 1 {
		response.BadRequest(ctx, errors.New("not support multiple credentials"))
		return
	}

	signupCredId := base64.URLEncoding.EncodeToString(user.WebAuthnCredentials()[0].ID)
	if u, _ := relay.db.FindUserByPasskey(user.GetEmail(), signupCredId); u != nil {
		user = u
	}

	if err := relay.db.SaveAccounts(user); err != nil {
		if errors.Is(err, seedworks.ErrUserAlreadyExists{}) {
			response.BadRequest(ctx, err.Error())
		} else {
			response.InternalServerError(ctx, err.Error())
		}
		return
	}

	ginJwtMiddleware().LoginHandler(ctx)
}
