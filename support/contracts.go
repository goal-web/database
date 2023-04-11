package support

import (
	"database/sql"
	"github.com/goal-web/contracts"
	"github.com/jmoiron/sqlx"
)

type SqlxExecutor interface {
	DriverName() string
	Queryx(query string, args ...any) (*sqlx.Rows, error)
	Get(dest any, query string, args ...any) (err error)
	Select(dest any, query string, args ...any) (err error)
	Exec(query string, args ...any) (sql.Result, error)
}

type Executor interface {
	DriverName() string
	Query(query string, args ...any) (contracts.Collection[contracts.Fields], contracts.Exception)
	Get(dest any, query string, args ...any) (err contracts.Exception)
	Select(dest any, query string, args ...any) (err contracts.Exception)
	Exec(query string, args ...any) (contracts.Result, contracts.Exception)
}
