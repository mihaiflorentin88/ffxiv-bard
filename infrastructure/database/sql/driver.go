package database

import (
	"database/sql"
	"ffxvi-bard/config"
	"ffxvi-bard/port/contract"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteDriver struct {
	database string
	path     string
}

var (
	instance *sql.DB
	once     sync.Once
	mu       sync.Mutex
)

func NewSqlDriver(cfg *config.DatabaseConfig) (contract.DatabaseDriverInterface, error) {
	var err error
	once.Do(func() {
		instance, err = sql.Open(cfg.Database, cfg.Path)
		if err != nil {
			return
		}
	})
	if err != nil {
		return nil, err
	}
	return &sqliteDriver{database: cfg.Database, path: cfg.Path}, nil
}

func (d *sqliteDriver) getConnection() (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()
	err := instance.Ping()
	if err != nil {
		instance, err = sql.Open(d.database, d.path)
		if err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func (d *sqliteDriver) Execute(query string, args ...interface{}) (sql.Result, error) {
	db, err := d.getConnection()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, args...)
}

func (d *sqliteDriver) FetchOne(query string, args ...interface{}) (*sql.Row, error) {
	db, err := d.getConnection()
	if err != nil {
		return nil, err
	}
	return db.QueryRow(query, args...), nil
}

func (d *sqliteDriver) FetchMany(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := d.getConnection()
	if err != nil {
		return nil, err
	}
	return db.Query(query, args...)
}

func (d *sqliteDriver) Close() {
	if instance != nil {
		instance.Close()
	}
}
