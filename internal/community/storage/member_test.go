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
	assert.NoError(t, UpsertMember("A", &ptr, &ptr, &ptr, 1))
	assert.NoError(t, UpsertMember("A", &ptr, &ptr, &ptr, 2))
}
