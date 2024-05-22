package storage

import (
	"encoding/binary"
	"fmt"
)

type NodeAddr struct {
	Addr     string
	Endpoint string
}

const addrSize = 32
const endpointSize = 128
const nodeAddrPrefix = "nodes:"

func NodeKey(node *NodeAddr) string {
	return fmt.Sprintf("%s%s", nodeAddrPrefix, node.Addr)
}

func unmarshalNodes(b []byte) []NodeAddr {
	ret := []NodeAddr{}
	for len(b) > 0 {
		sz := binary.LittleEndian.Uint16(b[:2])
		b = b[2:]
		if m := unmarshalNode(b[:sz]); m == nil {
			return nil
		} else {
			ret = append(ret, *m)
			b = b[sz:]
		}
	}
	return ret
}

func unmarshalNode(b []byte) *NodeAddr {
	if len(b) < addrSize+endpointSize {
		return nil
	}
	addr := string(b[:addrSize])
	endpoint := string(b[addrSize:])
	return &NodeAddr{addr, endpoint}
}

func (n *NodeAddr) Marshal() []byte {
	buf := make([]byte, addrSize+endpointSize)
	copy(buf[:addrSize], []byte(n.Addr))
	copy(buf[addrSize:], []byte(n.Endpoint))
	return buf
}
