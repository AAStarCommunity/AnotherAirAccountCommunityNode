package plugin_passkey_relay_party

import (
	"another_node/internal/common_util"
	"another_node/plugins/dvt/signature"
	"another_node/plugins/passkey_relay_party/conf"
	"another_node/plugins/passkey_relay_party/seedworks"
)

func sigTx(user *seedworks.User, signPayment *seedworks.TxSignature) (*seedworks.TxSignatureResult, error) {
	if chain := user.GetSpecifiyChain(signPayment.Network, signPayment.NetworkAlias); chain == nil {
		return nil, &seedworks.ErrChainNotFound{}
	} else {

		threshold := conf.GetDVT().Threshold
		dvtNodes := conf.GetDVT().Nodes
		timeout := conf.GetDVT().Timeout
		if dvtResult, err := signature.Bls([]byte(signPayment.TxData), threshold, timeout, dvtNodes, signPayment.CA, signPayment.CAPublicKey); err != nil {
			return nil, err
		} else {
			privateKey, err := user.GetPrivateKeyEcdsa(chain)
			if err != nil {
				return nil, err
			}
			if signHexStr, err := common_util.EthereumSignHexStr(signPayment.TxData, privateKey); err != nil {
				return nil, err
			} else {
				txSigRlt := seedworks.TxSignatureResult{
					Code:    200,
					DVT:     dvtResult,
					TxData:  signPayment.TxData, // userOpHash
					Sign:    signHexStr,
					Address: user.GetEOA(chain),
				}
				return &txSigRlt, nil
			}
		}
	}
}
