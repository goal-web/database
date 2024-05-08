package table

import (
	"database/sql"
	"errors"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
)

func (table *Table[T]) CountE(columns ...string) (int64, contracts.Exception) {
	queryStatement, bindings := table.WithCount(columns...).SelectSql()
	var num int64
	err := table.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) && err.Error() != sql.ErrNoRows.Error() {
		return 0, &SelectException{
			Sql:      queryStatement,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		}
	}
	return num, nil
}

func (table *Table[T]) AvgE(column string) (float64, contracts.Exception) {
	queryStatement, bindings := table.WithAvg(column).SelectSql()
	var num float64
	err := table.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) && err.Error() != sql.ErrNoRows.Error() {
		return 0, &SelectException{
			Sql:      queryStatement,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		}
	}
	return num, nil
}

func (table *Table[T]) SumE(column string) (float64, contracts.Exception) {
	queryStatement, bindings := table.WithSum(column).SelectSql()
	var num float64
	err := table.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) && err.Error() != sql.ErrNoRows.Error() {
		return 0, &SelectException{
			Sql:      queryStatement,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		}
	}
	return num, nil
}

func (table *Table[T]) MaxE(column string) (float64, contracts.Exception) {
	queryStatement, bindings := table.WithMax(column).SelectSql()
	var num float64
	err := table.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) && err.Error() != sql.ErrNoRows.Error() {
		return 0, &SelectException{
			Sql:      queryStatement,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		}
	}
	return num, nil
}

func (table *Table[T]) MinE(column string) (float64, contracts.Exception) {
	queryStatement, bindings := table.WithMin(column).SelectSql()
	var num float64
	err := table.getExecutor().Get(&num, queryStatement, bindings...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) && err.Error() != sql.ErrNoRows.Error() {
		return 0, &SelectException{
			Sql:      queryStatement,
			Bindings: bindings,
			Err:      err,
			previous: exceptions.WithError(err),
		}
	}
	return num, nil
}
