package table

import (
	"github.com/goal-web/contracts"
)

type BaseModel[T any] struct {
	class      contracts.Class[T]
	table      string
	connection string
	primaryKey string
}

func Model[T any](class contracts.Class[T], table string, connection ...string) *Table[T] {
	model := NewModel[T](class, table, connection...)
	return FromModel[T](model)
}

func NewModel[T any](class contracts.Class[T], table string, connection ...string) BaseModel[T] {
	conn := ""
	if len(connection) > 0 {
		conn = connection[0]
	}
	return BaseModel[T]{class: class, table: table, connection: conn}
}

func (model BaseModel[T]) GetClass() contracts.Class[T] {
	return model.class
}

func (model BaseModel[T]) GetPrimaryKey() string {
	if model.primaryKey == "" {
		return "id"
	}
	return model.primaryKey
}

func (model BaseModel[T]) SetPrimaryKey(key string) {
	model.primaryKey = key
}

func (model BaseModel[T]) GetTable() string {
	return model.table
}

func (model BaseModel[T]) GetConnection() string {
	return model.connection
}
