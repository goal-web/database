package table

import (
	"github.com/goal-web/application"
	"github.com/goal-web/contracts"
	"github.com/goal-web/querybuilder"
	class2 "github.com/goal-web/supports/class"
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

func getTable[T any](name string) *Table[T] {
	builder := querybuilder.NewBuilder[T](name)
	instance := &Table[T]{
		QueryBuilder:      builder,
		primaryKey:        "id",
		table:             name,
		createdTimeColumn: "created_at",
		UpdatedTimeColumn: "updated_at",
	}
	builder.Bind(instance)
	return instance
}

// Query 将使用默认 connection
func Query[T any](name string) *Table[T] {
	return getTable[T](name).SetClass(class2.Make[T]())
}

func Class[T any](class contracts.Class[T], table string) *Table[T] {
	return Query[T](table).SetClass(class)
}

func Auth(class contracts.Class[contracts.Authenticatable], table, primaryKey string) contracts.Query[contracts.Authenticatable] {
	return Query[contracts.Authenticatable](table).SetClass(class).SetPrimaryKey(primaryKey)
}

func ArrayQuery(name string) *Table[contracts.Fields] {
	return getTable[contracts.Fields](name).SetClass(arrayClass[contracts.Fields]{})
}

func FromModel[T any](model contracts.Model[T]) *Table[T] {
	return WithConnection[T](model.GetTable(), model.GetConnection()).SetClass(model.GetClass()).SetPrimaryKey(model.GetPrimaryKey())
}

// WithConnection 使用指定链接
func WithConnection[T any](name string, connection any) *Table[T] {
	if connection == "" || connection == nil {
		return Query[T](name)
	}
	return getTable[T](name).SetConnection(connection)
}

// WithTX 使用TX
func WithTX[T any](name string, tx contracts.DBTx) contracts.Query[T] {
	return getTable[T](name).SetExecutor(tx)
}
