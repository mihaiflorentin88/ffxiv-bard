package contract

import "database/sql"

type SqlDriverInterface interface {
	Execute(query string, args string) (sql.Result, error)
}
