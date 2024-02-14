package database

import (
	"database/sql"
	port "ffxvi-bard/port/contract"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteDriver struct {
	database string
	path     string
}

func NewSqlDriver(database string) port.SqlDriverInterface {
	return sqliteDriver{database: database}
}

func (d sqliteDriver) connection() (*sql.DB, error) {
	db, err := sql.Open(d.database, d.path)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db, nil
}

func (d sqliteDriver) Execute(query string, args string) (sql.Result, error) {
	db, err := d.connection()
	if err != nil {
		return nil, err
	}
	result, err := db.Exec(query, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}
