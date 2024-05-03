package table

import (
	"errors"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
)

var InsertError = errors.New("insert statement execution failed")

func (table *Table[T]) CreateE(fields contracts.Fields) (*T, contracts.Exception) {
	table.createdTime(fields)
	sql, bindings := table.CreateSql(fields)
	result, err := table.getExecutor().Exec(sql, bindings...)
	if err != nil {
		return nil, &CreateException{
			Sql:      sql,
			Bindings: bindings,
			Err:      InsertError,
			previous: err,
		}
	}
	var id, lastIdErr = result.LastInsertId()
	if lastIdErr != nil {
		logs.WithError(lastIdErr).Debug("Table.Create: get last insert id failed")
	}

	if _, existsPrimaryKey := fields[table.primaryKey]; !existsPrimaryKey && lastIdErr == nil {
		fields[table.primaryKey] = id
	}

	instance := table.class.New(fields)
	return &instance, nil
}

func (table *Table[T]) InsertE(values ...contracts.Fields) contracts.Exception {
	for i := range values {
		table.createdTime(values[i])
	}
	sql, bindings := table.InsertSql(values)
	_, exception := table.getExecutor().Exec(sql, bindings...)

	if exception != nil {
		return &InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      InsertError,
			previous: exception,
		}
	}

	return nil
}
func (table *Table[T]) Insert(values ...contracts.Fields) bool {
	return table.InsertE(values...) == nil
}

func (table *Table[T]) InsertGetId(values ...contracts.Fields) int64 {
	result, err := table.InsertGetIdE(values...)
	if err != nil {
		panic(err)
	}
	return result
}

func (table *Table[T]) InsertGetIdE(values ...contracts.Fields) (int64, contracts.Exception) {
	for i := range values {
		table.createdTime(values[i])
	}
	sql, bindings := table.InsertSql(values)
	result, exception := table.getExecutor().Exec(sql, bindings...)

	if exception != nil {
		return 0, &InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      InsertError,
			previous: exception,
		}
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, &InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      errors.New("failed to get last id"),
			previous: exceptions.WithError(err),
		}
	}

	return id, nil
}

func (table *Table[T]) InsertOrIgnoreE(values ...contracts.Fields) (int64, contracts.Exception) {
	for i := range values {
		table.createdTime(values[i])
	}
	sql, bindings := table.InsertIgnoreSql(values)
	result, exception := table.getExecutor().Exec(sql, bindings...)

	if exception != nil {
		return 0, &InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      InsertError,
			previous: exception,
		}
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return 0, &InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      errors.New("failed to get rows affected"),
			previous: exceptions.WithError(err),
		}
	}

	return rowsAffected, nil
}

func (table *Table[T]) InsertOrIgnore(values ...contracts.Fields) int64 {
	result, err := table.InsertOrIgnoreE(values...)
	if err != nil {
		panic(err)
	}
	return result
}

// InsertOrReplace 将新记录插入数据库，同时如果存在，则先删除此行数据，然后插入新的数据
// Insert a new record into the database, and if it exists, delete this row of data first, and then insert new data.
func (table *Table[T]) InsertOrReplace(values ...contracts.Fields) int64 {
	result, err := table.InsertOrReplaceE(values...)
	if err != nil {
		panic(err)
	}
	return result
}

func (table *Table[T]) InsertOrReplaceE(values ...contracts.Fields) (int64, contracts.Exception) {
	for i := range values {
		table.createdTime(values[i])
	}
	sql, bindings := table.InsertReplaceSql(values)
	result, exception := table.getExecutor().Exec(sql, bindings...)

	if exception != nil {
		return 0, &InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      InsertError,
			previous: exception,
		}
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return 0, &InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      errors.New("failed to get rows affected"),
			previous: exceptions.WithError(err),
		}
	}

	return rowsAffected, nil
}

func (table *Table[T]) FirstOrCreate(where contracts.Fields, values ...contracts.Fields) T {
	instance, err := table.FirstOrCreateE(where, values...)
	if err != nil {
		panic(err)
	}
	return *instance
}

func (table *Table[T]) FirstOrCreateE(where contracts.Fields, values ...contracts.Fields) (*T, contracts.Exception) {
	instance, _ := table.WhereFields(where).FirstE()
	if instance != nil {
		return instance, nil
	}
	if len(values) > 0 {
		utils.MergeFields(where, values[0])
	}
	return table.CreateE(where)
}
