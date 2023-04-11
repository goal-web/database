package table

import (
	"database/sql"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
)

func (table *Table[T]) fetch(query string, bindings ...any) (contracts.Collection[T], contracts.Exception) {
	rows, err := table.getExecutor().Query(query, bindings...)
	if err != nil && err != sql.ErrNoRows {
		return collection.New[T](nil), &SelectException{
			Sql:      query,
			Bindings: bindings,
			Err:      err,
		}
	}

	var list = make([]T, rows.Len())

	rows.Foreach(func(i int, fields contracts.Fields) {
		list[i] = table.class.NewByTag(fields, "db")
	})

	return collection.New(list), nil
}

func (table *Table[T]) Get() contracts.Collection[T] {
	data, exception := table.GetE()
	if exception != nil {
		panic(exception)
	}
	return data
}
func (table *Table[T]) GetE() (contracts.Collection[T], contracts.Exception) {
	queryStatement, bindings := table.SelectSql()
	return table.fetch(queryStatement, bindings...)
}

func (table *Table[T]) SelectForUpdateE() (contracts.Collection[T], contracts.Exception) {
	queryStatement, bindings := table.SelectForUpdateSql()
	return table.fetch(queryStatement, bindings...)
}

func (table *Table[T]) SelectForUpdate() contracts.Collection[T] {
	data, _ := table.SelectForUpdateE()
	return data
}

func (table *Table[T]) Find(key any) *T {
	result, _ := table.Where(table.primaryKey, key).FirstE()
	return result
}

func (table *Table[T]) First() *T {
	result, _ := table.FirstE()
	return result
}

func (table *Table[T]) FirstE() (*T, contracts.Exception) {
	list, e := table.Take(1).GetE()
	if e != nil {
		return list.First(), e
	}
	if list.IsEmpty() {
		statement, bindings := table.SelectSql()
		e = &NotFoundException{Sql: statement, Bindings: bindings, Err: sql.ErrNoRows}
	}
	return list.First(), e
}

func (table *Table[T]) FirstOrFail() T {
	result, err := table.FirstE()
	if err != nil {
		panic(err)
	}
	return *result
}
