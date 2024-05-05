package node

import (
	"another_node/internal/community/storage"
	"fmt"
)

type CommunityDelegate struct {
	DataChannel  chan []byte
	Broadcasts   [][]byte
	BroadcastCap int
}

// NotifyMsg receives a message from the network
func (d *CommunityDelegate) NotifyMsg(data []byte) {
	d.DataChannel <- data
}

// NodeMeta returns the current node metadata
func (d *CommunityDelegate) NodeMeta(limit int) []byte {
	return nil
}

// GetBroadcasts returns the broadcast messages
func (d *CommunityDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	if len(d.Broadcasts) > 0 {
		broadcasts := d.Broadcasts
		d.Broadcasts = nil
		return broadcasts
	}
	return nil
}

// LocalState return the local state data while a remote node joins or sync
func (d *CommunityDelegate) LocalState(join bool) []byte {
	if join {
		if _, err := storage.GetAllMembers(); err != nil {
			return nil
		} else {

			return []byte{1, 2, 3}
		}
	} else {
		//TODO: retrive partial data by non-init sync policy form storage and return to joiner
		return []byte{4, 5, 6}
	}
	return nil
}

// MergeRemoteState merges the remote state while current node joins or sync
func (d *CommunityDelegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) > 0 {
		if join {
			// TODO: merge partial data by init sync policy from remote
			fmt.Print("join: Merge remote state: ", buf)
		} else {
			// TODO: merge partial data by non-init sync policy from remote
			fmt.Print("sync: Merge remote state: ", buf)
		}
	}
}
