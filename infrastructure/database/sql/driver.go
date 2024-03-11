package database

import (
	"database/sql"
	"ffxvi-bard/config"
	"sync"

	_ "modernc.org/sqlite"
)

type SqliteDriver struct {
	database string
	path     string
}

var (
	Instance *sql.DB
	once     sync.Once
	mu       sync.Mutex
)

func NewSqlDriver(cfg *config.DatabaseConfig) (*SqliteDriver, error) {
	var err error
	once.Do(func() {
		Instance, err = sql.Open(cfg.Database, cfg.Path)
		if err != nil {
			return
		}
	})
	if err != nil {
		return nil, err
	}
	return &SqliteDriver{database: cfg.Database, path: cfg.Path}, nil
}

func (d *SqliteDriver) GetConnection() (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()
	err := Instance.Ping()
	if err != nil {
		Instance, err = sql.Open(d.database, d.path)
		if err != nil {
			return nil, err
		}
	}
	return Instance, nil
}

func (d *SqliteDriver) Execute(query string, args ...interface{}) (sql.Result, error) {
	db, err := d.GetConnection()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, args...)
}

func (d *SqliteDriver) FetchOne(query string, args ...interface{}) (*sql.Row, error) {
	db, err := d.GetConnection()
	if err != nil {
		return nil, err
	}
	return db.QueryRow(query, args...), nil
}

func (d *SqliteDriver) FetchMany(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := d.GetConnection()
	if err != nil {
		return nil, err
	}
	return db.Query(query, args...)
}

func (d *SqliteDriver) Close() {
	if Instance != nil {
		Instance.Close()
	}
}
