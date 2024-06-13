package plugin

import "another_node/internal/community/node"

type CommunityPlugin interface {
	Initialize(*node.Community) error
	Start() error
	Stop() error
}
