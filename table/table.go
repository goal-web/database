package table

import (
	"errors"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/utils"
)

type Table[T any] struct {
	contracts.QueryBuilder[T]
	executor contracts.SqlExecutor

	table             string
	primaryKey        string
	class             contracts.Class[T]
	createdTimeColumn string
	UpdatedTimeColumn string
}

// SetConnection 参数要么是 contracts.DBConnection 要么是 string
func (table *Table[T]) SetConnection(connection any) *Table[T] {
	if conn, ok := connection.(contracts.DBConnection); ok {
		table.executor = conn
	} else {
		table.executor = getFactory().Connection(utils.ToString(connection, ""))
	}
	return table
}

// SetClass 设置类
func (table *Table[T]) SetClass(class contracts.Class[T]) *Table[T] {
	table.class = class
	return table
}

// SetPrimaryKey 设置主键
func (table *Table[T]) SetPrimaryKey(name string) *Table[T] {
	table.primaryKey = name
	return table
}

// SetCreatedTimeColumn 设置创建时间字段
func (table *Table[T]) SetCreatedTimeColumn(column string) *Table[T] {
	table.createdTimeColumn = column
	return table
}

// SetUpdatedTimeColumn 设置更新时间字段
func (table *Table[T]) SetUpdatedTimeColumn(column string) *Table[T] {
	table.UpdatedTimeColumn = column
	return table
}

// getExecutor 获取 sql 语句的执行者
func (table *Table[T]) getExecutor() contracts.SqlExecutor {
	if table.executor != nil {
		return table.executor
	}
	return getFactory().Connection()
}

// SetExecutor 参数必须是 contracts.DBTx 实例
func (table *Table[T]) SetExecutor(executor contracts.SqlExecutor) contracts.Query[T] {
	table.executor = executor
	return table
}

func (table *Table[T]) Delete() int64 {
	result, exception := table.DeleteE()
	if exception != nil {
		panic(exception)
	}
	return result
}

func (table *Table[T]) DeleteE() (int64, contracts.Exception) {
	sql, bindings := table.DeleteSql()
	result, exception := table.getExecutor().Exec(sql, bindings...)
	if exception != nil {
		return 0, &DeleteException{
			Sql:      sql,
			Bindings: bindings,
			Err:      errors.New("delete statement execution failed"),
			previous: exception,
		}
	}
	num, err := result.RowsAffected()
	if err != nil {
		return 0, &DeleteException{
			Sql:      sql,
			Bindings: bindings,
			Err:      errors.New("failed to get number of affected rows"),
			previous: exceptions.WithError(err),
		}
	}
	return num, nil
}
