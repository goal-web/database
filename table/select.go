package table

import (
	"database/sql"
	"github.com/goal-web/contracts"
)

func (table *Table) fetch(query string, bindings ...interface{}) (contracts.Collection, contracts.Exception) {
	rows, err := table.getExecutor().Query(query, bindings...)
	if err != nil && err != sql.ErrNoRows {
		return nil, &SelectException{
			Sql:      query,
			Bindings: bindings,
			Err:      err,
		}
	}

	// 返回 Collection<fields>
	if table.class == nil {
		return rows, nil
	}

	return rows.Map(func(fields contracts.Fields) interface{} {
		return table.class.NewByTag(fields, "db")
	}), nil
}

func (table *Table) Get() contracts.Collection {
	queryStatement, bindings := table.SelectSql()

	collection, exception := table.fetch(queryStatement, bindings...)
	if exception == nil {
		return collection
	}
	panic(exception)
}

func (table *Table) SelectForUpdate() contracts.Collection {
	queryStatement, bindings := table.SelectForUpdateSql()
	collection, exception := table.fetch(queryStatement, bindings...)
	if exception == nil {
		return collection
	}
	panic(exception)
}

func (table *Table) Find(key interface{}) interface{} {
	return table.Where(table.primaryKey, key).First()
}

func (table *Table) First() interface{} {
	return table.Take(1).Get().First()
}

func (table *Table) FirstOrFail() interface{} {
	if result := table.First(); result != nil {
		return result
	}
	queryStatement, bindings := table.SelectSql()
	panic(NotFoundException{
		Sql:      queryStatement,
		Bindings: bindings,
		Err:      sql.ErrNoRows,
	})
}
