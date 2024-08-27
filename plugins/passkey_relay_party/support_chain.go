package plugin_passkey_relay_party

import "another_node/internal/seedworks"

func isSupportChain(chain seedworks.Chain) bool {
	if len(chain) > 0 {
		if chain != seedworks.OptimismSepolia && chain != seedworks.BaseSepolia && chain != seedworks.OptimismMainnet {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}
