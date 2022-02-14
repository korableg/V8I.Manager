package service

import (
	"context"
	"github.com/korableg/V8I.Manager/app/internal/sqlitedb"
	"github.com/korableg/V8I.Manager/app/user"
	"github.com/korableg/V8I.Manager/app/user/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestService(t *testing.T) {
	var s Service

	sdbcfg := sqlitedb.SqliteConfig{Path: "./db.db"}

	sdb, err := sqlitedb.NewSqliteDB(sdbcfg)
	require.Nil(t, err)

	defer func() {
		sdb.Close()

		err = os.Remove(sdbcfg.Path)
		assert.Nil(t, err)
	}()

	r, err := repository.NewSqliteRepository(sdb)
	require.Nil(t, err)

	s, err = NewService(r)
	require.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	u := user.AddUserRequest{
		Name:     "Test",
		Password: "111",
		Role:     "admin",
		Token:    "1122",
	}

	id, err := s.Add(ctx, u)
	assert.Nil(t, err)
	assert.NotEqual(t, id, 0)

}
