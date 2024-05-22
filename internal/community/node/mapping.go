package node

import (
	"another_node/internal/community/storage"
)

// GetPublicKey get public key by hashed account
func GetPublicKey(hashedAccount *string) string {
	if member, err := storage.TryFindMember(*hashedAccount); err != nil {
		return ""
	} else {
		return member.PublicKey
	}
}

func UpcomingHandler(buf []byte) {
	if len(buf) > 0 {
		protocol := buf[0]
		payload := buf[1:]

		switch protocol {
		case MemberStream:
			if members := storage.UnmarshalMembers(payload); len(members) > 0 {
				go storage.MergeRemoteAccounts(members)
			}

		case AddrStream:
			go storage.MergeRemoteAddr(payload)
		}
	}
}
