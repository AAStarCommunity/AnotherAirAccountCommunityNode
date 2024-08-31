package plugin_passkey_relay_party

import "another_node/internal/seedworks"

var supportChains = map[seedworks.Chain]byte{
	seedworks.EthereumSepolia: 0,
	seedworks.OptimismSepolia: 0,
	seedworks.BaseSepolia:     0,
}

func isSupportChain(chain seedworks.Chain) bool {
	if len(chain) > 0 {
		if _, ok := supportChains[chain]; !ok {
			return false
		}

		return true
	} else {
		return false
	}
}
