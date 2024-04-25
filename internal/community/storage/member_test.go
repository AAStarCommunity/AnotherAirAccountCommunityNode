package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMember(t *testing.T) {
	if err := os.Setenv("UnitTestEnv", "1"); err != nil {
		t.Skip("skipped due to CI")
	} else {
		defer os.Unsetenv("UnitTestEnv")
	}
	prepare_test()

	ptr := "B"
	prtN := 1
	ver := 1
	assert.NoError(t, UpsertMember("A", ptr, ptr, ptr, prtN, &ver))
	assert.NoError(t, UpsertMember("A", ptr, ptr, ptr, prtN, &ver))
}

func TestFindMember(t *testing.T) {
	if err := os.Setenv("UnitTestEnv", "1"); err != nil {
		t.Skip("skipped due to CI")
	} else {
		defer os.Unsetenv("UnitTestEnv")
	}
	prepare_test()

	ptrN := 1
	ver := 1

	ptr := "B"
	assert.NoError(t, UpsertMember("A", ptr, ptr, ptr, ptrN, &ver))
	ptr = "C"
	assert.NoError(t, UpsertMember("A", ptr, ptr, ptr, ptrN, &ver))

	member, err := TryFindMember("A")
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, "B", member.PublicKey)
}
