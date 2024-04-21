package migrations

import (
	"another_node/conf"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var sortedUp, sortedDown []int

type migration001 struct {
}

func (m *migration001) Up(_ *gorm.DB) error {
	println("[up]001")
	sortedUp = append(sortedUp, 1)
	return nil
}

func (m *migration001) Down(_ *gorm.DB) error {
	println("[down]001")
	sortedDown = append(sortedDown, 1)
	return nil
}

type migration002 struct {
}

func (m *migration002) Up(_ *gorm.DB) error {
	println("[up]002")
	sortedUp = append(sortedUp, 2)
	return nil
}

func (m *migration002) Down(_ *gorm.DB) error {
	println("[down]002")
	sortedDown = append(sortedDown, 2)
	return nil
}

func initTestData() {
	sortedUp = make([]int, 0)
	sortedDown = make([]int, 0)

	migrations = []Migration{
		&migration001{},
		&migration002{},
	}
}
func TestMigrate(t *testing.T) {
	initTestData()
	migrate(nil)

	if testing.Short() {
		assert.Equal(t, 1, sortedUp[0])
		assert.Equal(t, 2, sortedUp[1])
	}
}

func TestRollback(t *testing.T) {
	var db *gorm.DB

	os.Setenv("Evn", "dev")

	if testing.Short() {
		initTestData()
	} else {
		db = conf.GetDbClient()
	}

	rollback(db)

	if testing.Short() {
		assert.Equal(t, 2, sortedDown[0])
		assert.Equal(t, 1, sortedDown[1])
	}
}

func TestAutoMigrate(t *testing.T) {
	var db *gorm.DB
	if err := os.Setenv("UnitTestEnv", "1"); err != nil {
		t.Skip("skipped due to CI")
	} else {
		defer os.Unsetenv("UnitTestEnv")
	}

	db = conf.GetDbClient()

	m1 := &Migration20240420{}

	assert.NoError(t, m1.Up(db))
	assert.NoError(t, m1.Down(db))
}
