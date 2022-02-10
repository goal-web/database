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

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t Tx) Rollback() error {
	return t.tx.Rollback()
}

func New(tx *sqlx.Tx, events contracts.EventDispatcher) contracts.DBTx {
	return &Tx{
		tx:       tx,
		Executor: support.NewExecutor(tx, events),
		events:   events,
	}
}
