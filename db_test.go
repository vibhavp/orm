package orm

import (
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	os.Remove("./test.db")
	db, err := Open("sqlite3", "test.db")
	require.NoError(t, err)
	require.NoError(t, db.DB.Ping())
}

func TestAutoMigrate(t *testing.T) {
	os.Remove("./test.db")
	db, _ := Open("sqlite3", "test.db")
	type Test struct {
		A int
		B string
		C bool
	}
	err := db.Migrate(Test{})
	assert.NoError(t, err)
}
