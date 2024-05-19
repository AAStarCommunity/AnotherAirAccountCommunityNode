package node

import (
	"another_node/internal/community/storage"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const syncAccountApi = "api/account/v1/sync?count="

func splitNodeAddr(nodeAddr string) string {
	u := strings.Split(nodeAddr, ":")
	return u.Scheme + "://" + u.Host + ":8080"
}

// MergeRemoteMembers merge remote members by calling the remote API /api/account/v{version}/sync
func MergeRemoteMembers(snap *storage.Snapshot) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s%d", splitNodeAddr(entrypointNodeAddr[0]), syncAccountApi, snap.TotalMembers), nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if members := storage.UnmarshalMembers(body); members != nil {
		for _, member := range members {
			storage.MergeRemoteMember(&member)
		}
	}
}
