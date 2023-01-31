package table

import (
	"github.com/goal-web/application"
	"github.com/goal-web/contracts"
	"github.com/goal-web/querybuilder"
)

var factory contracts.DBFactory

func SetFactory(dbFactory contracts.DBFactory) {
	factory = dbFactory
}

func getFactory() contracts.DBFactory {
	if factory == nil {
		factory = application.Get("db.factory").(contracts.DBFactory)
	}
	return factory
}

func getTable(name string) *Table {
	builder := querybuilder.NewQuery(name)
	instance := &Table{
		QueryBuilder: builder,
		primaryKey:   "id",
		table:        name,
	}
	builder.Bind(instance)
	return instance
}

// Query 将使用默认 connection
func Query(name string) *Table {
	return getTable(name).SetConnection(factory.Connection())
}

func FromModel(model contracts.Model) *Table {
	return WithConnection(model.GetTable(), model.GetConnection()).SetClass(model.GetClass()).SetPrimaryKey(model.GetPrimaryKey())
}

// WithConnection 使用指定链接
func WithConnection(name string, connection interface{}) *Table {
	if connection == "" || connection == nil {
		return Query(name)
	}
	return getTable(name).SetConnection(connection)
}

// WithTX 使用TX
func WithTX(name string, tx contracts.DBTx) contracts.QueryBuilder {
	return getTable(name).SetExecutor(tx)
}
