package node

import (
	"another_node/internal/community/storage"
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
	if addr, err := getAddr(); err != nil {
		return nil
	} else {
		// TODO: merge []NodeAddr
		_ = addr
		skip := uint32(0)
		members := storage.GetMembers(skip, ^uint32(0))
		if len(members) > 0 {
			m := []byte{MemberStream}
			m = append(m, members.Marshal()...)
			return m
		}
		return nil
	}
}

const (
	MemberStream uint8 = 0x01
	AddrStream   uint8 = 0x02
	All          uint8 = MemberStream | AddrStream
)

// MergeRemoteState merges the remote state while current node joins or sync
func (d *CommunityDelegate) MergeRemoteState(buf []byte, join bool) {
	// TODO: merge []NodeAddr
	UpcomingHandler(buf)
}
