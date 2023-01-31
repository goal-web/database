package table

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/utils"
)

type Table struct {
	contracts.QueryBuilder
	executor contracts.SqlExecutor

	table      string
	primaryKey string
	class      contracts.Class
}

// SetConnection 参数要么是 contracts.DBConnection 要么是 string
func (table *Table) SetConnection(connection interface{}) *Table {
	if conn, ok := connection.(contracts.DBConnection); ok {
		table.executor = conn
	} else {
		table.executor = getFactory().Connection(utils.ConvertToString(connection, ""))
	}
	return table
}

// SetClass 设置类
func (table *Table) SetClass(class contracts.Class) *Table {
	table.class = class
	return table
}

// SetPrimaryKey 设置主键
func (table *Table) SetPrimaryKey(name string) *Table {
	table.primaryKey = name
	return table
}

// getExecutor 获取 sql 语句的执行者
func (table *Table) getExecutor() contracts.SqlExecutor {
	return table.executor
}

// SetExecutor 参数必须是 contracts.DBTx 实例
func (table *Table) SetExecutor(executor contracts.SqlExecutor) contracts.QueryBuilder {
	table.executor = executor
	return table
}

func (table *Table) Delete() int64 {
	sql, bindings := table.DeleteSql()
	result, err := table.getExecutor().Exec(sql, bindings...)
	if err != nil {
		panic(DeleteException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		})
	}
	num, err := result.RowsAffected()
	if err != nil {
		panic(DeleteException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		})
	}
	return num
}
