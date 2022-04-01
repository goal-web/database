package table

import (
	"database/sql"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
)

func (this *Table) fetch(query string, bindings ...interface{}) contracts.Collection {
	rows, err := this.getExecutor().Query(query, bindings...)
	if err != nil && err != sql.ErrNoRows {
		panic(SelectException{exceptions.WithError(err, contracts.Fields{"sql": query, "bindings": bindings})})
	}

	// 返回 Collection<fields>
	if this.class == nil {
		return rows
	}

	return rows.Map(func(fields contracts.Fields) interface{} {
		return this.class.NewByTag(fields, "db")
	})
}

func (this *Table) Get() contracts.Collection {
	queryStatement, bindings := this.SelectSql()

	return this.fetch(queryStatement, bindings...)
}

func (this *Table) SelectForUpdate() contracts.Collection {
	queryStatement, bindings := this.SelectForUpdateSql()

	return this.fetch(queryStatement, bindings...)
}

func (this *Table) Find(key interface{}) interface{} {
	return this.Where(this.primaryKey, key).First()
}

func (this *Table) First() interface{} {
	return this.Take(1).Get().First()
}

func (this *Table) FirstOrFail() interface{} {
	if result := this.First(); result != nil {
		return result
	}
	panic(NotFoundException{exceptions.WithError(sql.ErrNoRows, nil)})
}
