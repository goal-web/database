package support

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/events"
	"github.com/goal-web/supports/exceptions"
)

type BaseExecutor struct {
	executor SqlxExecutor
	events   contracts.EventDispatcher
}

func (this *BaseExecutor) DriverName() string {
	return this.executor.DriverName()
}

func NewExecutor(executor SqlxExecutor, dispatcher contracts.EventDispatcher) Executor {
	return &BaseExecutor{
		executor: executor,
		events:   dispatcher,
	}
}

func (this *BaseExecutor) Query(query string, args ...interface{}) (results contracts.Collection, err error) {
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args})
		}
	}()

	rows, err := this.executor.Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	return ParseRowsToCollection(rows)
}

func (this *BaseExecutor) Get(dest interface{}, query string, args ...interface{}) (err error) {
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args})
		}
	}()
	return this.executor.Get(dest, query, args...)
}

func (this *BaseExecutor) Select(dest interface{}, query string, args ...interface{}) (err error) {
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args})
		}
	}()
	return this.executor.Get(dest, query, args...)
}

func (this *BaseExecutor) Exec(query string, args ...interface{}) (result contracts.Result, err error) {
	defer func() {
		if err == nil {
			err = exceptions.ResolveException(recover())
		}
		if err == nil {
			this.events.Dispatch(&events.QueryExecuted{Sql: query, Bindings: args})
		}
	}()
	return this.executor.Exec(query, args...)
}
