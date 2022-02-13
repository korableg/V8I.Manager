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

	selectUser = "SELECT id, name, password_hash, role, token FROM users"

	selectUserByID = selectUser + " WHERE id = ?;"

	selectUsersList = selectUser + " ORDER BY name;"

	updateUser = "UPDATE users SET name = ?, password_hash = ?, role = ?, token = ? WHERE id = ?;"

	deleteUserByID = "DELETE FROM users WHERE id = ?;"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type (
	sqliteRepository struct {
		db *sql.DB
	}
)

func NewSqliteRepository(sdb *sqlitedb.SqliteDB) (*sqliteRepository, error) {
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

func (s *sqliteRepository) Get(ctx context.Context, ID int64) (user.User, error) {
	var u user.User

	row := s.db.QueryRowContext(ctx, selectUserByID, ID)
	if err := row.Scan(&u.ID, &u.Name, &u.PasswordHash, &u.Role, &u.Token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, ErrUserNotFound
		}

		return user.User{}, fmt.Errorf("row scan: %w", err)
	}

	return u, nil
}

func (s *sqliteRepository) GetList(ctx context.Context) ([]user.User, error) {
	users := make([]user.User, 0)

	rows, err := s.db.QueryContext(ctx, selectUsersList)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	var u user.User
	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.Name, &u.PasswordHash, &u.Role, &u.Token); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		users = append(users, u)
	}

	return users, nil
}

func (s *sqliteRepository) Update(ctx context.Context, u user.User) error {
	if _, err := s.db.ExecContext(ctx, updateUser, u.Name, u.PasswordHash, u.Role, u.Token, u.ID); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (s *sqliteRepository) Delete(ctx context.Context, ID int64) error {
	if _, err := s.db.ExecContext(ctx, deleteUserByID, ID); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
