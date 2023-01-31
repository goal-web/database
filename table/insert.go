package table

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
)

func (table *Table) Create(fields contracts.Fields) interface{} {
	sql, bindings := table.CreateSql(fields)
	result, err := table.getExecutor().Exec(sql, bindings...)
	if err != nil {
		panic(CreateException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}
	var id, lastIdErr = result.LastInsertId()
	if lastIdErr != nil {
		logs.WithError(lastIdErr).Debug("Table.Create: get last insert id failed")
	}

	if _, existsPrimaryKey := fields[table.primaryKey]; !existsPrimaryKey {
		fields[table.primaryKey] = id
	}

	if table.class != nil {
		return table.class.New(fields)
	}

	return fields
}

func (table *Table) Insert(values ...contracts.Fields) bool {
	sql, bindings := table.InsertSql(values)
	result, err := table.getExecutor().Exec(sql, bindings...)

	if err != nil {
		panic(InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		panic(InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}

	return rowsAffected > 0
}

func (table *Table) InsertGetId(values ...contracts.Fields) int64 {
	sql, bindings := table.InsertSql(values)
	result, err := table.getExecutor().Exec(sql, bindings...)

	if err != nil {
		panic(InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}

	id, err := result.LastInsertId()

	if err != nil {
		panic(InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}

	return id
}

func (table *Table) InsertOrIgnore(values ...contracts.Fields) int64 {
	sql, bindings := table.InsertIgnoreSql(values)
	result, err := table.getExecutor().Exec(sql, bindings...)

	if err != nil {
		panic(InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		panic(InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}

	return rowsAffected
}

func (table *Table) InsertOrReplace(values ...contracts.Fields) int64 {
	sql, bindings := table.InsertReplaceSql(values)
	result, err := table.getExecutor().Exec(sql, bindings...)

	if err != nil {
		panic(InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		panic(InsertException{
			Sql:      sql,
			Bindings: bindings,
			Err:      err,
		})
	}

	return rowsAffected
}

func (table *Table) FirstOrCreate(values ...contracts.Fields) interface{} {
	var attributes contracts.Fields
	argsLen := len(values)
	if argsLen > 0 {
		for field, value := range values[0] {
			attributes[field] = value
			table.Where(field, value)
		}
		if result := table.First(); result != nil {
			return result
		} else if argsLen > 1 {
			utils.MergeFields(attributes, values[1])
		}
		return table.Create(attributes)
	}

	return nil
}
