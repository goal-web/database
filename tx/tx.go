package tx

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/support"
	"github.com/jmoiron/sqlx"
)

type Tx struct {
	tx *sqlx.Tx
	support.Executor
	events contracts.EventDispatcher
}

func (t *Tx) Commit() contracts.Exception {
	if err := t.tx.Commit(); err != nil {
		return &CommitException{Err: err}
	}
	return nil
}

func (t *Tx) Rollback() contracts.Exception {
	if err := t.tx.Rollback(); err != nil {
		return &RollbackException{Err: err}
	}
	return nil
}

func New(tx *sqlx.Tx, events contracts.EventDispatcher) contracts.DBTx {
	return &Tx{
		tx:       tx,
		Executor: support.NewExecutor(tx, events, nil),
		events:   events,
	}
}
