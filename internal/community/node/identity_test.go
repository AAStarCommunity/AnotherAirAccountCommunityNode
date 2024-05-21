package node

import (
	"os"
	"testing"
)

func TestExtractAddress(t *testing.T) {
	t.Run("test extract address", func(t *testing.T) {
		str := "abcdefg"

		m := extractIdenty(&str)

		if len(*m) != 18 {
			t.Errorf("expect abcdefg, got %s", *m)
		}
	})
}

func TestGenerateIdentity(t *testing.T) {
	t.Run("test generate identity", func(t *testing.T) {
		os.Setenv("UnitTest", "1")
		defer os.Unsetenv("UnitTest")

		if m, err := generateIdentity(); err != nil {
			t.Errorf("expect nil, got %v", err)
		} else {
			if len(m) != 18 {
				t.Errorf("expect abcdefg, got %s", m)
			}
		}
	})
}
