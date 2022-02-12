package sqlitedb

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	_ "github.com/mattn/go-sqlite3"
)

const (
	schemaDir = "migrations/schema"
	dataDir   = "migrations/data"
)

//go:embed migrations
var migrationsDir embed.FS

type (
	SqliteConfig struct {
		Path string `yaml:"path"`
	}

	SqliteDB struct {
		db *sql.DB
	}
)

func NewSqliteDB(config SqliteConfig) (*SqliteDB, error) {
	db, err := sql.Open("sqlite3", config.Path)
	if err != nil {
		return nil, err
	}

	sdb := &SqliteDB{
		db: db,
	}

	if err = sdb.migrate(); err != nil {
		return nil, err
	}

	return sdb, nil
}

func (s *SqliteDB) Close() error {
	return s.db.Close()
}

func (s *SqliteDB) DB() *sql.DB {
	return s.db
}

func (s *SqliteDB) migrate() error {
	schemas, err := migrationsDir.ReadDir(schemaDir)
	if err != nil {
		return fmt.Errorf("read embed schema dir: %w", err)
	}

	datas, err := migrationsDir.ReadDir(dataDir)
	if err != nil {
		return fmt.Errorf("read embed data dir: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}

	if err = execMigrations(ctx, tx, schemas, schemaDir); err != nil {
		return fmt.Errorf("executing schema migrations: %w", err)
	}

	if err = execMigrations(ctx, tx, datas, dataDir); err != nil {
		return fmt.Errorf("executing data migrations: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func execMigrations(ctx context.Context, tx *sql.Tx, entries []fs.DirEntry, dir string) error {
	for _, entry := range entries {
		data, err := migrationsDir.ReadFile(fmt.Sprintf("%s/%s", dir, entry.Name()))
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		if _, err = tx.ExecContext(ctx, string(data)); err != nil {
			return fmt.Errorf("execute query: %w", err)
		}
	}

	return nil
}
