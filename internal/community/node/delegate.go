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
	if ss, err := storage.GetSnapshot(); err == nil {
		if join {
			return ss
		} else {
			if s, err := storage.UnmarshalSnapshot(ss); err == nil {
				members := storage.GetAllMembers(s.TotalMembers)
				if members != nil {
					return storage.MarshalMembers(members)
				}
			}
		}
	}
	return nil
}

// MergeRemoteState merges the remote state while current node joins or sync
func (d *CommunityDelegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) > 0 {
		if join {
			if snap, err := storage.UnmarshalSnapshot(buf); err == nil {
				go MergeRemoteMembers(snap)
			}
		} else {
			go func() {
				if members, err := storage.Unmarshal(buf); err != nil {
					fmt.Print("sync: Failed to unmarshal remote state: ", err)
				} else {
					storage.MergeRemoteMember(members)
				}
			}()
		}
	}
}
