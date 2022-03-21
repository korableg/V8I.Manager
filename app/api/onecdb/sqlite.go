package onecdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/korableg/V8I.Manager/app/internal/sqlitedb"
)

const (
	insertDB = `INSERT OR IGNORE INTO onecdbs
	(uuid, name, connect, order_in_list,
	 order_in_tree, folder, client_connection_speed, app,
	 wa, version, additional_parameters)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	selectDB = `SELECT
       id, uuid, name, connect, order_in_list,
       order_in_tree, folder, client_connection_speed, app,
       wa, version, additional_parameters
	FROM onecdbs`

	selectDBByID = selectDB + ` WHERE id = ?;`

	selectDBList = selectDB + ` ORDER BY id`

	updateDB = `UPDATE onecdbs SET
		uuid = ?, name = ?, connect = ?, order_in_list = ?,
		order_in_tree = ?, folder = ?, client_connection_speed = ?, app = ?,
		wa = ?, version = ?, additional_parameters = ?
	WHERE id = ?;`

	deleteDBByID = `DELETE FROM onecdbs WHERE id = ?`
)

type (
	sqliteRepository struct {
		db *sql.DB
	}
)

var (
	ErrDBNotFound = errors.New("db not found")
)

func NewSqliteRepository(sdb *sqlitedb.SqliteDB) (*sqliteRepository, error) {
	if sdb == nil {
		return nil, errors.New("SQLite db doesn't initialized")
	}

	return &sqliteRepository{db: sdb.DB()}, nil
}

func (s sqliteRepository) Add(ctx context.Context, db DB) (int64, error) {
	sqlResult, err := s.db.ExecContext(
		ctx, insertDB,
		db.UUID.String(), db.Name, db.Connect, db.OrderInList,
		db.OrderInTree, db.Folder, db.ClientConnectionSpeed, db.App,
		db.WA, db.Version, db.AdditionalParameters)
	if err != nil {
		return 0, fmt.Errorf("inserting db to onecdbs: %w", err)
	}

	id, err := sqlResult.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last inserting id: %w", err)
	}

	return id, nil
}

func (s sqliteRepository) Get(ctx context.Context, id int64) (DB, error) {
	var db DB

	row := s.db.QueryRowContext(ctx, selectDBByID, id)
	if err := row.Scan(
		&db.ID, &db.UUID, &db.Name, &db.Connect, &db.OrderInList,
		&db.OrderInTree, &db.Folder, &db.ClientConnectionSpeed, &db.App,
		&db.WA, &db.Version, &db.AdditionalParameters); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DB{}, ErrDBNotFound
		}

		return DB{}, fmt.Errorf("row scan: %w", err)
	}

	return db, nil
}

func (s sqliteRepository) GetList(ctx context.Context) ([]DB, error) {
	dbs := make([]DB, 0)

	rows, err := s.db.QueryContext(ctx, selectDBList)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	var db DB
	for rows.Next() {
		if err = rows.Scan(&db.ID, &db.UUID, &db.Name, &db.Connect, &db.OrderInList,
			&db.OrderInTree, &db.Folder, &db.ClientConnectionSpeed, &db.App,
			&db.WA, &db.Version, &db.AdditionalParameters); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		dbs = append(dbs, db)
	}

	return dbs, nil
}

func (s sqliteRepository) Update(ctx context.Context, db DB) error {
	if _, err := s.db.ExecContext(
		ctx, updateDB,
		db.UUID, db.Name, db.Connect, db.OrderInList,
		db.OrderInTree, db.Folder, db.ClientConnectionSpeed, db.App,
		db.WA, db.Version, db.AdditionalParameters, db.ID); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (s sqliteRepository) Delete(ctx context.Context, id int64) error {
	if _, err := s.db.ExecContext(ctx, deleteDBByID, id); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
