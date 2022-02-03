package table

import (
	"github.com/goal-web/contracts"
)

type BaseModel struct {
	class      contracts.Class
	table      string
	connection string
	primaryKey string
}

func Model(class contracts.Class, table string, connection ...string) *Table {
	return FromModel(NewModel(class, table, connection...))
}

func NewModel(class contracts.Class, table string, connection ...string) BaseModel {
	conn := ""
	if len(connection) > 0 {
		conn = connection[0]
	}
	return BaseModel{class: class, table: table, connection: conn}
}

func (model BaseModel) GetClass() contracts.Class {
	return model.class
}

func (model BaseModel) GetPrimaryKey() string {
	if model.primaryKey == "" {
		return "id"
	}
	return model.primaryKey
}

func (model BaseModel) SetPrimaryKey(key string) {
	model.primaryKey = key
}

func (model BaseModel) GetTable() string {
	return model.table
}

func (model BaseModel) GetConnection() string {
	return model.connection
}
