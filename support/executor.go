package support

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/events"
	"github.com/goal-web/supports/exceptions"
	"time"
)

type BaseExecutor struct {
	executor SqlxExecutor
	events   contracts.EventDispatcher
	wrapper  func(sql string) string
}

func NewExecutor(executor SqlxExecutor, dispatcher contracts.EventDispatcher, wrapper func(sql string) string) Executor {
	return &BaseExecutor{
		executor: executor,
		events:   dispatcher,
		wrapper:  wrapper,
	}
}

func (base *BaseExecutor) DriverName() string {
	return base.executor.DriverName()
}
func (base *BaseExecutor) dispatchEvent(event contracts.Event) {
	if base.events != nil {
		base.events.Dispatch(event)
	}
}

func (base *BaseExecutor) getStatement(sql string) string {
	if base.wrapper != nil {
		return base.wrapper(sql)
	}
	return sql
}

func (base *BaseExecutor) Query(query string, args ...interface{}) (results contracts.Collection, err error) {
	query = base.getStatement(query)
	var timeConsuming time.Duration
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}
		base.dispatchEvent(&events.QueryExecuted{Sql: query, Bindings: args, Time: timeConsuming, Error: err})
	}()
	var startAt = time.Now()
	rows, err := base.executor.Queryx(query, args...)
	timeConsuming = time.Now().Sub(startAt)
	if err != nil {
		return nil, err
	}

	return ParseRowsToCollection(rows)
}

func (base *BaseExecutor) Get(dest interface{}, query string, args ...interface{}) (err error) {
	query = base.getStatement(query)
	var startAt = time.Now()
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}

		base.dispatchEvent(&events.QueryExecuted{Sql: query, Bindings: args, Time: time.Now().Sub(startAt), Error: err})
	}()
	return base.executor.Get(dest, query, args...)
}

func (base *BaseExecutor) Select(dest interface{}, query string, args ...interface{}) (err error) {
	query = base.getStatement(query)
	var startAt = time.Now()
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}
		base.dispatchEvent(&events.QueryExecuted{Sql: query, Bindings: args, Time: time.Now().Sub(startAt), Error: err})
	}()
	return base.executor.Get(dest, query, args...)
}

func (base *BaseExecutor) Exec(query string, args ...interface{}) (result contracts.Result, err error) {
	query = base.getStatement(query)
	var startAt = time.Now()
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}
		base.dispatchEvent(&events.QueryExecuted{Sql: query, Bindings: args, Time: time.Now().Sub(startAt), Error: err})
	}()
	return base.executor.Exec(query, args...)
}
