package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/korableg/V8I.Manager/app/internal/sqlitedb"
	"github.com/korableg/V8I.Manager/app/user"
)

const (
	insertUser = "INSERT INTO users (name, password_hash, role, token) VALUES (?, ?, ?, ?);"
)

type (
	sqliteRepository struct {
		db *sql.DB
	}
)

func NewSqliteRepository(sdb *sqlitedb.SqliteDB) (Repository, error) {
	if sdb == nil {
		return nil, errors.New("SQLite db doesn't initialized")
	}

	return &sqliteRepository{db: sdb.DB()}, nil
}

func (s *sqliteRepository) Add(ctx context.Context, u user.User) (int64, error) {
	sqlResult, err := s.db.ExecContext(ctx, insertUser, u.Name, u.PasswordHash, u.Role, u.Token)
	if err != nil {
		return 0, fmt.Errorf("inserting user to db: %w", err)
	}

	id, err := sqlResult.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last inserting id: %w", err)
	}

	return id, nil
}

func (s *sqliteRepository) GetList(ctx context.Context) ([]user.User, error) {
	//TODO implement me
	panic("implement me")
}
