package table

import (
	"errors"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
)

func (this *Table) Get() contracts.Collection {
	sql, bindings := this.SelectSql()

	rows, err := this.getExecutor().Query(sql, bindings...)
	if err != nil {
		panic(SelectException{exceptions.WithError(err, contracts.Fields{"sql": sql, "bindings": bindings})})
	}

	// 返回 Collection<fields>
	if this.class == nil {
		return rows
	}

	return rows.Map(func(fields contracts.Fields) interface{} {
		return this.class.NewByTag(fields, "db")
	})
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
	panic(NotFoundException{exceptions.WithError(errors.New("未找到"), contracts.Fields{})})
}
