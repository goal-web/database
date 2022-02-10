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
		Executor: support.NewExecutor(db, dispatcher),
		events:   dispatcher,
	}
}

func (this *Base) Begin() (contracts.DBTx, error) {
	sqlxTx, err := this.db.Beginx()
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
