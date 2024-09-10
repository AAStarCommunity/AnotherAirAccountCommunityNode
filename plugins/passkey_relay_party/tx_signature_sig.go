package plugin_passkey_relay_party

import (
	"another_node/internal/common_util"
	"another_node/plugins/passkey_relay_party/seedworks"
)

func sigTx(user *seedworks.User, signPayment *seedworks.TxSignature) (*seedworks.TxSignatureResult, error) {
	privateKey, err := user.GetPrivateKeyEcdsa()
	if err != nil {
		return nil, err
	}
	if signHexStr, err := common_util.EthereumSignHexStr(signPayment.TxData, privateKey); err != nil {
		return nil, err
	} else {
		txSigRlt := seedworks.TxSignatureResult{
			Code:   200,
			TxData: signPayment.TxData,
			Sign:   signHexStr,
			Address: func() string {
				_, aaAddr := user.GetChainAddresses("")
				return *aaAddr
			}(),
		}
		return &txSigRlt, nil
	}
}
