package contract

import "database/sql"

type DatabaseDriverInterface interface {
	GetConnection() (*sql.DB, error)
	Execute(query string, args ...interface{}) (sql.Result, error)
	FetchOne(query string, args ...interface{}) (*sql.Row, error)
	FetchMany(query string, args ...interface{}) (*sql.Rows, error)
	Close()
}
