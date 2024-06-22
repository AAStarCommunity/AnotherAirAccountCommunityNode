package node

import (
	member_storage "another_node/internal/community/storage/member"
)

// GetPublicKey get public key by hashed account
func GetPublicKey(hashedAccount *string) string {
	if member, err := member_storage.TryFindMember(*hashedAccount); err != nil {
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
			if members := member_storage.UnmarshalMembers(payload); len(members) > 0 {
				go member_storage.MergeRemoteAccounts(members)
			}

		case AddrStream:
			go member_storage.MergeRemoteAddr(payload)
		}
	}
}
