package table

import (
	"errors"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
)

var UpdateErr = errors.New("update statement execution failed")

func (table *Table[T]) UpdateOrInsert(attributes contracts.Fields, values contracts.Fields) bool {
	return table.UpdateOrInsertE(attributes, values) == nil
}

func (table *Table[T]) UpdateOrInsertE(attributes contracts.Fields, values contracts.Fields) contracts.Exception {
	table.WhereFields(attributes)
	sql, bindings := table.UpdateSql(attributes)
	result, exception := table.getExecutor().Exec(sql, bindings...)
	if exception != nil {
		return &UpdateException{
			Sql:      sql,
			Bindings: bindings,
			Err:      UpdateErr,
			previous: exception,
		}
	}
	num, _ := result.RowsAffected()
	if num > 0 {
		return nil
	}
	utils.MergeFields(attributes, values)
	return table.InsertE(attributes)
}

func (table *Table[T]) UpdateOrCreateE(attributes contracts.Fields, values contracts.Fields) (*T, contracts.Exception) {
	sql, bindings := table.WhereFields(attributes).UpdateSql(values)
	result, exception := table.getExecutor().Exec(sql, bindings...)
	if exception != nil {
		return nil, &UpdateException{
			Sql:      sql,
			Bindings: bindings,
			Err:      UpdateErr,
			previous: exception,
		}
	}
	num, _ := result.RowsAffected()
	if num > 0 {
		return table.WhereFields(attributes).FirstE()
	}
	utils.MergeFields(attributes, values)
	return table.CreateE(attributes)
}

func (table *Table[T]) UpdateOrCreate(attributes contracts.Fields, values contracts.Fields) T {
	result, e := table.UpdateOrCreateE(attributes, values)
	if e != nil {
		panic(e)
	}
	return *result
}

func (table *Table[T]) UpdateE(fields contracts.Fields) (int64, contracts.Exception) {
	sql, bindings := table.UpdateSql(fields)
	result, exception := table.getExecutor().Exec(sql, bindings...)
	if exception != nil {
		return 0, &UpdateException{
			Sql:      sql,
			Bindings: bindings,
			Err:      UpdateErr,
			previous: exception,
		}
	}
	num, err := result.RowsAffected()
	if err != nil {
		return 0, &UpdateException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		}
	}
	return num, nil
}
