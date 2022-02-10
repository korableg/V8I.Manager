package user

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	insertUser = "INSERT INTO users (name, password_hash, role, token) VALUES (?, ?, ?, ?);"
)

type (
	sqliteRepository struct {
		db *sql.DB
	}
)

func NewSqliteRepository(db *sql.DB) (Repository, error) {
	return &sqliteRepository{db: db}, nil
}

func (s *sqliteRepository) Add(ctx context.Context, u User) (int64, error) {
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

func (s *sqliteRepository) GetList(ctx context.Context) ([]User, error) {
	//TODO implement me
	panic("implement me")
}
