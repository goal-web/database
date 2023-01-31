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

func (base *Base) Begin() (contracts.DBTx, error) {
	sqlxTx, err := base.db.Beginx()
	if err != nil {
		return nil, err
	}
	return tx.New(sqlxTx, base.events), err
}

func (base *Base) Transaction(fn func(tx contracts.SqlExecutor) error) (err error) {
	sqlxTx, err := base.Begin()
	if err != nil {
		return &exceptions2.BeginException{Err: err}
	}

	defer func() { // 处理 panic 情况
		if recoverErr := recover(); recoverErr != nil {
			rollbackErr := sqlxTx.Rollback()
			err = recoverErr.(error)
			if rollbackErr != nil {
				err = &exceptions2.RollbackException{
					Err:      rollbackErr,
					Previous: exceptions.WithError(err),
				}
			} else {
				err = &exceptions2.TransactionException{Err: err}
			}
		}
	}()

	err = fn(sqlxTx)

	if err != nil {
		rollbackErr := sqlxTx.Rollback()
		if rollbackErr != nil {
			return &exceptions2.RollbackException{Err: rollbackErr, Previous: exceptions.WithError(err)}
		}
		return &exceptions2.TransactionException{Err: err}
	}

	return sqlxTx.Commit()
}
