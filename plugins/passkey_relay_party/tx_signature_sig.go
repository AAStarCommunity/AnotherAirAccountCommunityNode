package plugin_passkey_relay_party

import (
	"another_node/internal/common_util"
	"another_node/plugins/dvt"
	"another_node/plugins/passkey_relay_party/seedworks"
	"encoding/base64"
)

func sigTx(user *seedworks.User, signPayment *seedworks.TxSignature) (*seedworks.TxSignatureResult, error) {
	if chain := user.GetSpecifiyChain(signPayment.Network, signPayment.NetworkAlias); chain == nil {
		return nil, &seedworks.ErrChainNotFound{}
	} else {

		done := make(chan struct {
			signature []byte
			publickey []byte
			err       error
		})
		go func() {
			s, r, e := dvt.Signature(
				signPayment.CA,
				signPayment.CAPublicKey,
			)
			done <- struct {
				signature []byte
				publickey []byte
				err       error
			}{s, r, e}
		}()

		dvtSign := <-done

		if dvtSign.err != nil {
			return nil, dvtSign.err
		}

		privateKey, err := user.GetPrivateKeyEcdsa(chain)
		if err != nil {
			return nil, err
		}
		if signHexStr, err := common_util.EthereumSignHexStr(signPayment.TxData, privateKey); err != nil {
			return nil, err
		} else {
			txSigRlt := seedworks.TxSignatureResult{
				Code:      200,
				TxData:    signPayment.TxData, // userOpHash
				Sign:      signHexStr,
				BlsSign:   base64.URLEncoding.EncodeToString(dvtSign.signature),
				BlsPubKey: base64.URLEncoding.EncodeToString(dvtSign.publickey),
				Address:   user.GetEOA(chain),
			}
			return &txSigRlt, nil
		}
	}
}
