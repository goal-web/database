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

func (this *BaseExecutor) DriverName() string {
	return this.executor.DriverName()
}

func (this *BaseExecutor) getStatement(sql string) string {
	if this.wrapper != nil {
		return this.wrapper(sql)
	}
	return sql
}

func NewExecutor(executor SqlxExecutor, dispatcher contracts.EventDispatcher, wrapper func(sql string) string) Executor {
	return &BaseExecutor{
		executor: executor,
		events:   dispatcher,
		wrapper:  wrapper,
	}
}

func (this *BaseExecutor) Query(query string, args ...interface{}) (results contracts.Collection, err error) {
	query = this.getStatement(query)
	var timeConsuming time.Duration
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args, Time: timeConsuming})
		}
	}()
	var startAt = time.Now()
	rows, err := this.executor.Queryx(query, args...)
	timeConsuming = time.Now().Sub(startAt)
	if err != nil {
		return nil, err
	}

	return ParseRowsToCollection(rows)
}

func (this *BaseExecutor) Get(dest interface{}, query string, args ...interface{}) (err error) {
	query = this.getStatement(query)
	var startAt = time.Now()
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args, Time: time.Now().Sub(startAt)})
		}
	}()
	return this.executor.Get(dest, query, args...)
}

func (this *BaseExecutor) Select(dest interface{}, query string, args ...interface{}) (err error) {
	query = this.getStatement(query)
	var startAt = time.Now()
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args, Time: time.Now().Sub(startAt)})
		}
	}()
	return this.executor.Get(dest, query, args...)
}

func (this *BaseExecutor) Exec(query string, args ...interface{}) (result contracts.Result, err error) {
	query = this.getStatement(query)
	var startAt = time.Now()
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args, Time: time.Now().Sub(startAt)})
		}
	}()
	return this.executor.Exec(query, args...)
}
