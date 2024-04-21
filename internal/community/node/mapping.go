package node

import "another_node/internal/community/storage"

// GetPublicKey get public key by hashed account
func GetPublicKey(hashedAccount *string) string {
	if member, err := storage.TryFindMember(*hashedAccount); err != nil {
		return ""
	} else {
		return member.PublicKey
	}
}
