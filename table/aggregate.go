package table

import (
	"database/sql"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
)

func (this *Table) Count(columns ...string) int64 {
	queryStatement, bindings := this.WithCount(columns...).SelectSql()
	var num int64
	err := this.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && err != sql.ErrNoRows {
		exceptions.Throw(SelectException{exceptions.WithError(err, contracts.Fields{
			"columns":  columns,
			"sql":      queryStatement,
			"bindings": bindings,
		})})
	}
	return num
}

func (this *Table) Avg(column string, as ...string) int64 {
	queryStatement, bindings := this.WithAvg(column, as...).SelectSql()
	var num int64
	err := this.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && err != sql.ErrNoRows {
		exceptions.Throw(SelectException{exceptions.WithError(err, contracts.Fields{
			"column":   column,
			"sql":      queryStatement,
			"bindings": bindings,
		})})
	}
	return num
}

func (this *Table) Sum(column string, as ...string) int64 {
	queryStatement, bindings := this.WithSum(column, as...).SelectSql()
	var num int64
	err := this.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && err != sql.ErrNoRows {
		exceptions.Throw(SelectException{exceptions.WithError(err, contracts.Fields{
			"column":   column,
			"sql":      queryStatement,
			"bindings": bindings,
		})})
	}
	return num
}

func (this *Table) Max(column string, as ...string) int64 {
	queryStatement, bindings := this.WithMax(column, as...).SelectSql()
	var num int64
	err := this.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && err != sql.ErrNoRows {
		exceptions.Throw(SelectException{exceptions.WithError(err, contracts.Fields{
			"column":   column,
			"sql":      queryStatement,
			"bindings": bindings,
		})})
	}
	return num
}

func (this *Table) Min(column string, as ...string) int64 {
	queryStatement, bindings := this.WithMin(column, as...).SelectSql()
	var num int64
	err := this.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && err != sql.ErrNoRows {
		exceptions.Throw(SelectException{exceptions.WithError(err, contracts.Fields{
			"column":   column,
			"sql":      queryStatement,
			"bindings": bindings,
		})})
	}
	return num
}
