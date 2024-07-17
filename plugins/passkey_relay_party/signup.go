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

// regPrepare
// @Summary sign up step1. send code
// @Tags Plugins Passkey
// @Description Send captcha to email for verifying belongs
// @Accept json
// @Product json
// @Param registrationBody body seedworks.RegistrationPrepare true "Send Captcha to Email"
// @Router /api/passkey/v1/reg/prepare [post]
// @Success 200
func (relay *RelayParty) regPrepare(ctx *gin.Context) {
	var reg seedworks.Registration
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

// beginRegistration
// @Summary sign up step2. begin registration
// @Tags Plugins Passkey
// @Description Begin the registration process
// @Accept json
// @Product json
// @Param registrationBody body seedworks.Registration true "Begin Registration"
// @Router /api/passkey/v1/reg [post]
// @Success 200 {object} protocol.PublicKeyCredentialCreationOptions
func (relay *RelayParty) beginRegistration(ctx *gin.Context) {
	var reg seedworks.Registration
	if err := ctx.ShouldBindJSON(&reg); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	// TODO: special logic for align testing
	if !strings.HasSuffix(reg.Email, "@aastar.org") && reg.Captcha != "111111" {
		if err := relay.emailChallenge(reg.Email, reg.Captcha); err != nil {
			response.BadRequest(ctx, err.Error())
			return
		}
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

type finishRegistrationResponse struct {
	AccountInitCode string `json:"account_init_code"`
	AccountAddress  string `json:"account_address"`
	EoaAddress      string `json:"eoa_address"`
}

// func replayLoginResponse(ctx *gin.Context, append func(c *gin.Context)) {
// 	ctx.Writer.WriteString("{\"token\": ")
// 	ginJwtMiddleware().LoginHandler(ctx)
// 	ctx.Writer.WriteString(",\"append\": ")
// 	append(ctx)
// 	ctx.Writer.WriteString("}")
// }

// finishRegistration
// @Summary sign up step3. Finish Registration
// @Tags Plugins Passkey
// @Description Verify Passkey Registration
// @Accept json
// @Product json
// @Param email  query string true "user email" Format(email)
// @Param origin query string true "origin"
// @Param network query string false "network"
// @Param registrationBody body protocol.CredentialCreationResponse true "Verify Registration"
// @Router /api/passkey/v1/reg/verify [post]
// @Success 200 {object} SiginInResponse "OK"
func (relay *RelayParty) finishRegistration(ctx *gin.Context) {

	// TODO: for tokyo ONLY
	network := consts.Chain(ctx.Query("network"))
	if len(network) > 0 && network != consts.OptimismSepolia {
		response.BadRequest(ctx, "network not supported")
		return
	} else {
		network = consts.OptimismSepolia
	}

	// body works for parser, the additional info appends to query
	stubReg := seedworks.FinishRegistration{
		RegistrationPrepare: seedworks.RegistrationPrepare{
			Email: ctx.Query("email"),
		},
		Origin:  ctx.Query("origin"),
		Network: consts.Chain(ctx.Query("network")),
	}

	if user, err := relay.authSessionStore.FinishRegSession(&stubReg, ctx); err != nil {
		response.GetResponse().FailCode(ctx, 401, "SignUp failed: "+err.Error())
		return
	} else {
		// TODO: special logic for align testing
		if strings.HasSuffix(stubReg.Email, "@aastar.org") {
			response.GetResponse().WithDataSuccess(ctx, user)
			return
		}
		if initCode, address, eoaAddress, err := createAA(user, stubReg.Network); err != nil { // TODO: persistent initCode and address
			response.InternalServerError(ctx, err.Error())
			return
		} else {

			// TODO: special logic for tokyo
			relay.db.Save(user, false)
			relay.db.SaveAccounts(user, initCode, address, eoaAddress, string(network))

			ginJwtMiddleware().LoginHandler(ctx)

			return
		}
	}
}

// createAA represents creating an Account Abstraction for the user
func createAA(user *seedworks.User, network consts.Chain) (initCode, address, eoaAddress string, err error) {
	if w, err := account.NewHdWallet(account.HierarchicalPath_ETH); err != nil {
		return "", "", "", err
	} else {
		address, initCode, err := chain.CreateSmartAccount(w, network)
		if err != nil {
			return "", "", "", err
		}
		user.SetWallet(w, address, network)
		return initCode, address, w.Address(), nil
	}
}
