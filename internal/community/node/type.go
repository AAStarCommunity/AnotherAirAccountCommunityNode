package node

import "github.com/hashicorp/memberlist"

type Node struct {
	Members  *memberlist.Memberlist
	Delegate *CommunityDelegate
}
