package drivers

import (
	"github.com/goal-web/contracts"
	exceptions2 "github.com/goal-web/database/exceptions"
	"github.com/goal-web/database/support"
	"github.com/goal-web/database/tx"
	"github.com/goal-web/supports/exceptions"
	"github.com/jmoiron/sqlx"
)

type Base struct {
	support.Executor
	db     *sqlx.DB
	events contracts.EventDispatcher
}

func NewDriver(db *sqlx.DB, dispatcher contracts.EventDispatcher) *Base {
	return &Base{
		db:       db,
		Executor: support.NewExecutor(db, dispatcher, nil),
		events:   dispatcher,
	}
}

func WithWrapper(db *sqlx.DB, dispatcher contracts.EventDispatcher, wrapper func(string) string) *Base {
	return &Base{
		db:       db,
		Executor: support.NewExecutor(db, dispatcher, wrapper),
		events:   dispatcher,
	}
}

func (base *Base) Begin() (contracts.DBTx, contracts.Exception) {
	sqlxTx, err := base.db.Beginx()
	if err != nil {
		return nil, exceptions.WithError(err)
	}
	return tx.New(sqlxTx, base.events), nil
}

func (base *Base) Transaction(fn func(tx contracts.SqlExecutor) contracts.Exception) (e contracts.Exception) {
	sqlxTx, e := base.Begin()
	if e != nil {
		return &exceptions2.BeginException{Err: e}
	}

	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			rollbackErr := sqlxTx.Rollback()
			e = exceptions.WithRecover(recoverErr)
			if rollbackErr != nil {
				e = &exceptions2.RollbackException{Err: rollbackErr, Previous: e}
			} else {
				e = &exceptions2.TransactionException{Err: e}
			}
		}
	}()

	e = fn(sqlxTx)

	if e != nil {
		rollbackErr := sqlxTx.Rollback()
		if rollbackErr != nil {
			return &exceptions2.RollbackException{Err: rollbackErr, Previous: e}
		}
		return &exceptions2.TransactionException{Err: e}
	}

	return sqlxTx.Commit()
}
