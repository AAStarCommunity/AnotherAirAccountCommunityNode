package node

import "another_node/internal/community/storage"

func GetPublicKey(hashedAccount *string) string {
	if member, err := storage.TryFindMember(*hashedAccount); err != nil {
		return ""
	} else {
		return member.PublicKey
	}
}
