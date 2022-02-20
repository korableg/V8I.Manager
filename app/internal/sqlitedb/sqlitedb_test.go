package sqlitedb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSqliteDB(t *testing.T) {
	config := Config{
		Path: "./db.db",
	}

	sdb, err := NewSqliteDB(config)
	require.Nil(t, err)

	err = sdb.Close()
	require.Nil(t, err)

	err = os.Remove(config.Path)
	require.Nil(t, err)
}
