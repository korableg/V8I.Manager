package onecdb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/korableg/V8I.Manager/app/internal/sqlitedb"
)

const ()

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

func (s sqliteRepository) Add(ctx context.Context, db DB) error {
	//TODO implement me
	panic("implement me")
}

func (s sqliteRepository) Get(ctx context.Context, id uuid.UUID) (DB, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqliteRepository) GetList(ctx context.Context) ([]DB, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqliteRepository) Update(ctx context.Context, db DB) error {
	//TODO implement me
	panic("implement me")
}

func (s sqliteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

//func (s *sqliteRepository) Add(ctx context.Context, u User) (int64, error) {
//	sqlResult, err := s.db.ExecContext(ctx, insertUser, u.Name, u.PasswordHash, u.Role, u.Token)
//	if err != nil {
//		sqliteErr := sqlite3.Error{}
//		if errors.As(err, &sqliteErr) && sqliteErr.Code == 19 {
//			return 0, ErrUserAlreadyCreated
//		}
//		return 0, fmt.Errorf("inserting user to db: %w", err)
//	}
//
//	id, err := sqlResult.LastInsertId()
//	if err != nil {
//		return 0, fmt.Errorf("getting last inserting id: %w", err)
//	}
//
//	return id, nil
//}
//
//func (s *sqliteRepository) Get(ctx context.Context, ID int64) (User, error) {
//	var u User
//
//	row := s.db.QueryRowContext(ctx, selectUserByID, ID)
//	if err := row.Scan(&u.ID, &u.Name, &u.PasswordHash, &u.Role, &u.Token); err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return User{}, ErrUserNotFound
//		}
//
//		return User{}, fmt.Errorf("row scan: %w", err)
//	}
//
//	return u, nil
//}
//
//func (s *sqliteRepository) GetByName(ctx context.Context, name string) (User, error) {
//	var u User
//
//	row := s.db.QueryRowContext(ctx, selectUserByName, name)
//	if err := row.Scan(&u.ID, &u.Name, &u.PasswordHash, &u.Role, &u.Token); err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return User{}, ErrUserNotFound
//		}
//
//		return User{}, fmt.Errorf("row scan: %w", err)
//	}
//
//	return u, nil
//}
//
//func (s *sqliteRepository) GetList(ctx context.Context) ([]User, error) {
//	users := make([]User, 0)
//
//	rows, err := s.db.QueryContext(ctx, selectUsersList)
//	if err != nil {
//		return nil, fmt.Errorf("exec query: %w", err)
//	}
//
//	var u User
//	for rows.Next() {
//		if err = rows.Scan(&u.ID, &u.Name, &u.PasswordHash, &u.Role, &u.Token); err != nil {
//			return nil, fmt.Errorf("rows scan: %w", err)
//		}
//
//		users = append(users, u)
//	}
//
//	return users, nil
//}
//
//func (s *sqliteRepository) Update(ctx context.Context, u User) error {
//	if _, err := s.db.ExecContext(ctx, updateUser, u.Name, u.Role, u.Token, u.ID); err != nil {
//		return fmt.Errorf("exec query: %w", err)
//	}
//
//	return nil
//}
//
//func (s *sqliteRepository) Delete(ctx context.Context, ID int64) error {
//	if _, err := s.db.ExecContext(ctx, deleteUserByID, ID); err != nil {
//		return fmt.Errorf("exec query: %w", err)
//	}
//
//	return nil
//}
