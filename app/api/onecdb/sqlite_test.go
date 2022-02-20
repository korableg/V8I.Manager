package onecdb

import (
	"context"
	"os"
	"testing"

	"github.com/korableg/V8I.Manager/app/internal/sqlitedb"
	"github.com/stretchr/testify/require"
)

func TestSqliteRepository(t *testing.T) {
	var r Repository

	_ = r

	sdbcfg := sqlitedb.Config{Path: "./db.db"}

	sdb, err := sqlitedb.NewSqliteDB(sdbcfg)
	require.Nil(t, err)

	defer func() {
		sdb.Close()

		err = os.Remove(sdbcfg.Path)
		require.Nil(t, err)
	}()

	r, err = NewSqliteRepository(sdb)
	require.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = ctx

	//u := User{
	//	Name:         "testuser2",
	//	PasswordHash: "123456",
	//	Token:        "555",
	//	Role:         "reader",
	//}
	//
	//id, err := r.Add(ctx, u)
	//require.Nil(t, err)
	//assert.NotEqual(t, id, 0)
	//
	//u.ID = id
	//u.Name = "test222"
	//u.Role = "admin"
	//u.Token = "333"
	//
	//err = r.Update(ctx, u)
	//assert.Nil(t, err)
	//
	//u1, err := r.Get(ctx, id)
	//assert.Nil(t, err)
	//assert.EqualValues(t, u, u1)
	//
	//u1, err = r.GetByName(ctx, u.Name)
	//assert.Nil(t, err)
	//assert.EqualValues(t, u, u1)
	//
	//users, err := r.GetList(ctx)
	//assert.Nil(t, err)
	//assert.Equal(t, len(users), 2)
	//
	//err = r.Delete(ctx, id)
	//assert.Nil(t, err)
	//
	//u1, err = r.Get(ctx, id)
	//assert.Error(t, ErrUserNotFound)
}
