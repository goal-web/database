package table

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
)

func (table *Table) UpdateOrInsert(attributes contracts.Fields, values ...contracts.Fields) bool {
	table.WhereFields(attributes)
	sql, bindings := table.UpdateSql(attributes)
	result, err := table.getExecutor().Exec(sql, bindings...)
	if err != nil {
		panic(UpdateException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}
	num, _ := result.RowsAffected()
	if num > 0 {
		return true
	}
	if len(values) > 0 {
		utils.MergeFields(attributes, values[0])
	}
	return table.Insert(attributes)
}

func (table *Table) UpdateOrCreate(attributes contracts.Fields, values ...contracts.Fields) interface{} {
	table.WhereFields(attributes)
	sql, bindings := table.UpdateSql(attributes)
	result, err := table.getExecutor().Exec(sql, bindings...)
	if err != nil {
		panic(UpdateException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}
	num, _ := result.RowsAffected()
	if num > 0 {
		return true
	}
	if len(values) > 0 {
		utils.MergeFields(attributes, values[0])
	}
	return table.Insert(attributes)
}

func (table *Table) Update(fields contracts.Fields) int64 {
	sql, bindings := table.UpdateSql(fields)
	result, err := table.getExecutor().Exec(sql, bindings...)
	if err != nil {
		panic(UpdateException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}
	num, err := result.RowsAffected()
	if err != nil {
		panic(UpdateException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}
	return num
}
