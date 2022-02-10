package support

import (
	"database/sql"
	"github.com/goal-web/contracts"
	"github.com/jmoiron/sqlx"
)

type SqlxExecutor interface {
	DriverName() string
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Get(dest interface{}, query string, args ...interface{}) (err error)
	Select(dest interface{}, query string, args ...interface{}) (err error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Executor interface {
	DriverName() string
	Query(query string, args ...interface{}) (contracts.Collection, error)
	Get(dest interface{}, query string, args ...interface{}) (err error)
	Select(dest interface{}, query string, args ...interface{}) (err error)
	Exec(query string, args ...interface{}) (contracts.Result, error)
}
