package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshotMarshal(t *testing.T) {
	snapshot := &Snapshot{
		Version:      1,
		TotalMembers: 1,
		HashedDigest: []byte{1, 2, 3},
	}

	expected := []byte{1, 1, 0, 0, 0, 1, 2, 3}
	actual := snapshot.Digest().Marshal()

	assert.Equal(t, expected, actual)
}

func TestSnapshotUnmarshal(t *testing.T) {
	data := []byte{1, 1, 0, 0, 0, 1, 2, 3}

	snapshot, err := UnmarshalSnapshot(data)
	assert.NoError(t, err)
	assert.Equal(t, uint8(1), snapshot.Version)
	assert.Equal(t, uint32(1), snapshot.TotalMembers)
	assert.Equal(t, []byte{1, 2, 3}, snapshot.HashedDigest)
}
