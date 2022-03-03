package onecserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/korableg/V8I.Manager/app/internal/sqlitedb"
)

const (
	insertServer = "INSERT INTO onecservers (name, lst_path, watch) VALUES (?, ?, ?);"

	selectServer = "SELECT id, name, lst_path, watch FROM onecservers"

	selectServerByID = selectServer + " WHERE id = ?;"

	selectServersList = selectServer + " ORDER BY name;"

	updateServer = "UPDATE onecservers SET name = ?, lst_path = ? WHERE id = ?;"

	updateWatch = "UPDATE onecservers SET watch = ? WHERE id = ?;"

	updateHash = "UPDATE onecservers SET lst_hash = ? WHERE id = ?;"

	deleteServerByID = "DELETE FROM onecservers WHERE id = ?;"
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

func (s *sqliteRepository) Add(ctx context.Context, server Server) (int64, error) {
	sqlResult, err := s.db.ExecContext(ctx, insertServer, server.Name, server.LSTPath, server.Watch)
	if err != nil {
		return 0, fmt.Errorf("inserting server to db: %w", err)
	}

	id, err := sqlResult.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last inserting id: %w", err)
	}

	return id, nil
}

func (s *sqliteRepository) Get(ctx context.Context, ID int64) (Server, error) {
	var server Server

	row := s.db.QueryRowContext(ctx, selectServerByID, ID)
	if err := row.Scan(&server.ID, &server.Name, &server.LSTPath, &server.Watch); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Server{}, ErrServerNotFound
		}

		return Server{}, fmt.Errorf("row scan: %w", err)
	}

	return server, nil
}

func (s *sqliteRepository) GetList(ctx context.Context) ([]Server, error) {
	servers := make([]Server, 0)

	rows, err := s.db.QueryContext(ctx, selectServersList)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	var server Server
	for rows.Next() {
		if err = rows.Scan(&server.ID, &server.Name, &server.LSTPath, &server.Watch); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		servers = append(servers, server)
	}

	return servers, nil
}

func (s *sqliteRepository) Update(ctx context.Context, server Server) error {
	if _, err := s.db.ExecContext(ctx, updateServer, server.Name, server.LSTPath); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (s *sqliteRepository) UpdateWatch(ctx context.Context, ID int64, watch bool) error {
	if _, err := s.db.ExecContext(ctx, updateWatch, watch, ID); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (s *sqliteRepository) UpdateHash(ctx context.Context, ID int64, hash string) error {
	if _, err := s.db.ExecContext(ctx, updateHash, hash, ID); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (s *sqliteRepository) Delete(ctx context.Context, ID int64) error {
	if _, err := s.db.ExecContext(ctx, deleteServerByID, ID); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
