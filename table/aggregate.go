package table

import (
	"database/sql"
	"github.com/goal-web/supports/exceptions"
)

func (table *Table) Count(columns ...string) int64 {
	queryStatement, bindings := table.WithCount(columns...).SelectSql()
	var num int64
	err := table.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && err != sql.ErrNoRows {
		exceptions.Throw(SelectException{
			Sql:      queryStatement,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		})
	}
	return num
}

func (table *Table) Avg(column string, as ...string) int64 {
	queryStatement, bindings := table.WithAvg(column, as...).SelectSql()
	var num int64
	err := table.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && err != sql.ErrNoRows {
		exceptions.Throw(SelectException{
			Sql:      queryStatement,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		})
	}
	return num
}

func (table *Table) Sum(column string, as ...string) int64 {
	queryStatement, bindings := table.WithSum(column, as...).SelectSql()
	var num int64
	err := table.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && err != sql.ErrNoRows {
		exceptions.Throw(SelectException{
			Sql:      queryStatement,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		})
	}
	return num
}

func (table *Table) Max(column string, as ...string) int64 {
	queryStatement, bindings := table.WithMax(column, as...).SelectSql()
	var num int64
	err := table.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && err != sql.ErrNoRows {
		exceptions.Throw(SelectException{
			Sql:      queryStatement,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		})
	}
	return num
}

func (table *Table) Min(column string, as ...string) int64 {
	queryStatement, bindings := table.WithMin(column, as...).SelectSql()
	var num int64
	err := table.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && err != sql.ErrNoRows {
		exceptions.Throw(SelectException{
			Sql:      queryStatement,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		})
	}
	return num
}
