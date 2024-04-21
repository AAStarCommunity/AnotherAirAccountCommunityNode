package global

import "another_node/internal/community/node"

type Community struct {
	Node *node.Node
}

var CommunityInstance *Community
