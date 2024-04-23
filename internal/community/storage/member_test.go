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
	assert.NoError(t, UpsertMember("A", &ptr, &ptr, &ptr, &prtN, 1))
	assert.NoError(t, UpsertMember("A", &ptr, &ptr, &ptr, &prtN, 2))
}

func TestFindMember(t *testing.T) {
	if err := os.Setenv("UnitTestEnv", "1"); err != nil {
		t.Skip("skipped due to CI")
	} else {
		defer os.Unsetenv("UnitTestEnv")
	}
	prepare_test()

	prtN := 1

	ptr := "B"
	assert.NoError(t, UpsertMember("A", &ptr, &ptr, &ptr, &prtN, 1))
	ptr = "C"
	assert.NoError(t, UpsertMember("A", &ptr, &ptr, &ptr, &prtN, 1))

	member, err := TryFindMember("A")
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, "B", member.PublicKey)
}
