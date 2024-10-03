package plugin_passkey_relay_party

import "another_node/internal/seedworks"

var supportChains = map[seedworks.Chain]bool{
	seedworks.EthereumSepolia: true,
	seedworks.OptimismSepolia: true,
	seedworks.BaseSepolia:     true,
}

func isSupportChain(chain seedworks.Chain) bool {
	if len(chain) > 0 {
		if support, exists := supportChains[chain]; exists && support {
			return true
		}
	}
	return false
}
