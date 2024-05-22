package node

import (
	"another_node/internal/community/storage"
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
		defer func() {
			os.Unsetenv("UnitTest")
			defer storage.Close()
		}()

		if db, err := storage.EnsureOpen(); err != nil {
			t.Errorf("expect nil, got %v", err)
		} else {
			if m, err := generateIdentity(db); err != nil {
				t.Errorf("expect nil, got %v", err)
			} else {
				if len(m) != 18 {
					t.Errorf("expect len(m) == 18, got %d", len(m))
				}
			}
		}
	})
}
