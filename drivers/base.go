package drivers

import (
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/events"
	exceptions2 "github.com/goal-web/database/exceptions"
	"github.com/goal-web/database/table"
	"github.com/goal-web/database/tx"
	"github.com/goal-web/supports/exceptions"
	"github.com/jmoiron/sqlx"
)

type Base struct {
	*sqlx.DB
	events contracts.EventDispatcher
}

func (this *Base) Query(query string, args ...interface{}) (results contracts.Collection, err error) {
	defer func() {
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args})
		}
	}()
	rows, err := this.DB.Query(query, args...)

	if err != nil {
		return nil, err
	}

	data, err := table.ParseRows(rows)

	if err != nil {
		return nil, err
	}

	return collection.FromFieldsSlice(data), nil
}

func (this *Base) Get(dest interface{}, query string, args ...interface{}) (err error) {
	defer func() {
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args})
		}
	}()
	return this.DB.Get(dest, query, args...)
}
func (this *Base) Select(dest interface{}, query string, args ...interface{}) (err error) {
	defer func() {
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args})
		}
	}()
	return this.DB.Get(dest, query, args...)
}

func (this *Base) Begin() (contracts.DBTx, error) {
	sqlxTx, err := this.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return tx.New(sqlxTx, this.events), err
}

func (this *Base) Transaction(fn func(tx contracts.SqlExecutor) error) (err error) {
	sqlxTx, err := this.Begin()
	if err != nil {
		return exceptions2.BeginException{Exception: exceptions.WithError(err, nil)}
	}

	defer func() { // 处理 panic 情况
		if recoverErr := recover(); recoverErr != nil {
			rollbackErr := sqlxTx.Rollback()
			err = recoverErr.(error)
			if rollbackErr != nil {
				err = exceptions2.RollbackException{Exception: exceptions.WithPrevious(rollbackErr, nil, exceptions.WithError(err, nil))}
			} else {
				err = exceptions2.TransactionException{Exception: exceptions.WithError(err, nil)}
			}
		}
	}()

	err = fn(sqlxTx)

	if err != nil {
		rollbackErr := sqlxTx.Rollback()
		if rollbackErr != nil {
			return exceptions2.RollbackException{Exception: exceptions.WithPrevious(rollbackErr, nil, exceptions.WithError(err, nil))}
		}
		return exceptions2.TransactionException{Exception: exceptions.WithError(err, nil)}
	}

	return sqlxTx.Commit()
}

func (this *Base) Exec(query string, args ...interface{}) (result contracts.Result, err error) {
	defer func() {
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args})
		}
	}()
	return this.DB.Exec(query, args...)
}
