package plugin_passkey_relay_party

import (
	"another_node/internal/web_server/pkg/response"
	"another_node/plugins/passkey_relay_party/seedworks"
	"crypto/sha256"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

// beginTxSignature
// @Summary request credential assertion for begin signature tx
// @Description Begins the signature process
// @Tags Plugins Passkey
// @Accept json
// @Produce json
// @Param dataSignature body seedworks.TxSignature true "send challenge to passkey for tx sign"
// @Success 200 {object} protocol.PublicKeyCredentialRequestOptions
// @Router /api/passkey/v1/tx/sign [post]
// @Security JWT
func (relay *RelayParty) beginTxSignature(ctx *gin.Context) {
	var tx seedworks.TxSignature

	if ok, email := CurrentUser(ctx); !ok {
		response.GetResponse().FailCode(ctx, http.StatusUnauthorized)
		return
	} else {
		if err := ctx.ShouldBindJSON(&tx); err != nil {
			response.BadRequest(ctx, err)
			return
		} else if len(tx.TxData) == 0 {
			response.BadRequest(ctx, "TxData is empty")
			return
		}
		tx.Email = email
	}

	if session := relay.txSessionStore.Get(seedworks.GetSessionKey(tx.Origin, tx.Email, tx.Nonce)); session != nil {
		response.BadRequest(ctx, "Already in Signature Process")
		return
	} else {
		user, err := relay.db.Find(tx.Email)
		if err != nil {
			response.NotFound(ctx, err.Error())
		}
		if options, err := relay.txSessionStore.NewTxSession(user, &tx); err != nil {
			response.InternalServerError(ctx, err)
		} else {
			response.GetResponse().WithDataSuccess(ctx, options.Response)
		}
	}
}

// finishTxSignature
// @Summary finish sign payment request credential assertion
// @Description Finish the sign process for payment
// @Tags Plugins Passkey
// @Accept json
// @Produce json
// @Param paymentSign body protocol.CredentialAssertionResponse true "Verify SignIn"
// @Param origin query string true "origin"
// @Param nonce query string true "nonce"
// @Success 200 {object} seedworks.TxSignatureResult
// @Router /api/passkey/v1/tx/sign/verify [post]
// @Security JWT
func (relay *RelayParty) finishTxSignature(ctx *gin.Context) {
	signPayment := seedworks.TxSignature{
		Origin: ctx.Query("origin"),
		Nonce:  ctx.Query("nonce"),
	}
	if ok, email := CurrentUser(ctx); !ok {
		response.GetResponse().FailCode(ctx, http.StatusUnauthorized)
		return
	} else {
		signPayment.Email = email
	}

	user, err := relay.txSessionStore.FinishSignSession(&signPayment, ctx)
	if err != nil {
		response.GetResponse().FailCode(ctx, 403, "SignIn failed: "+err.Error())
		return
	}

	_ = user

	txSigRlt := seedworks.TxSignatureResult{
		Code:   200,
		TxData: signPayment.TxData,
		Sign: func(opHash []byte) string {
			hash := sha256.New()
			hash.Write(opHash)
			hashByte := hash.Sum(nil)
			if privateKey, err := crypto.HexToECDSA(user.GetPrivateKey()); err != nil {
				return ""
			} else {
				if str, err := crypto.Sign(accounts.TextHash(hashByte), privateKey); err != nil {
					return ""
				} else {
					return string(str)
				}
			}
		}([]byte(signPayment.TxData)),
	}
	response.GetResponse().WithDataSuccess(ctx, &txSigRlt)
}
