package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"

	"github.com/gin-gonic/gin"
)

// beginSignature
// @Summary request credential assertion for begin signature tx
// @Description Begins the signature process
// @Tags Plugins Passkey
// @Accept json
// @Produce json
// @Param dataSignature body seedworks.TxSignature true "send challenge to passkey for tx sign"
// @Success 200 {object} protocol.PublicKeyCredentialRequestOptions
// @Router /api/passkey/v1/tx/sign [post]
func (relay *RelayParty) beginSignature(ctx *gin.Context) {
	var tx seedworks.TxSignature
	if err := ctx.ShouldBindJSON(&tx); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	if session := relay.signSessionStore.Get(seedworks.GetSessionKey(tx.Origin, tx.Email, tx.Nonce)); session != nil {
		response.BadRequest(ctx, "Already in Signature Process")
		return
	} else {
		user, err := relay.db.Find(tx.Email)
		if err != nil {
			response.NotFound(ctx, err.Error())
		}
		if options, err := relay.authSessionStore.NewTxSignSession(user, &tx); err != nil {
			response.InternalServerError(ctx, err)
		} else {
			response.GetResponse().WithDataSuccess(ctx, options.Response)
		}
	}
}

// finishSignPayment
// @Summary finish sign payment request credential assertion
// @Description Finish the sign process for payment
// @Tags Plugins Passkey
// @Accept json
// @Produce json
// @Param paymentSign body protocol.CredentialAssertionResponse true "Verify SignIn"
// @Param email  query string true "user email" Format(email)
// @Param origin query string true "origin"
// @Param nonce query string true "nonce"
// @Success 200 {object} response.Response
// @Router /api/passkey/v1/payment/sign/verify [post]
func (relay *RelayParty) finishSignPayment(ctx *gin.Context) {
	signPayment := seedworks.TxSignature{
		Origin: ctx.Query("origin"),
		Email:  ctx.Query("email"),
		Nonce:  ctx.Query("nonce"),
	}

	user, err := relay.authSessionStore.FinishSignSession(&signPayment, ctx)
	if err != nil {
		response.BadRequest(ctx, "SignIn failed: "+err.Error())
		return
	}

	_ = user

	response.GetResponse().Success(ctx)
}
