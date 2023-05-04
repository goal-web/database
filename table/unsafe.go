package table

import (
	"database/sql"
	"github.com/goal-web/contracts"
)

// Count 检索查询的“count”结果
// Retrieve the "count" result of the query.
func (table *Table[T]) Count(columns ...string) int64 {
	result, err := table.CountE(columns...)
	if err != nil {
		panic(err)
	}
	return result
}

// Avg 检索给定列的平均值
// Retrieve the average of the values of a given column.
func (table *Table[T]) Avg(column string) float64 {
	result, err := table.AvgE(column)
	if err != nil {
		panic(err)
	}
	return result
}

// Sum 检索给定列的值的总和
// Retrieve the sum of the values of a given column.
func (table *Table[T]) Sum(column string) float64 {
	result, err := table.SumE(column)
	if err != nil {
		panic(err)
	}
	return result
}

// Max 检索给定列的最大值
// Retrieve the maximum value of a given column.
func (table *Table[T]) Max(column string) float64 {
	result, err := table.MaxE(column)
	if err != nil {
		panic(err)
	}
	return result
}

// Min 检索给定列的最小值
// Retrieve the minimum value of a given column.
func (table *Table[T]) Min(column string) float64 {
	result, err := table.MinE(column)
	if err != nil {
		panic(err)
	}
	return result
}

// Create 保存新模型并返回实例
// Save a new model and return the instance.
func (table *Table[T]) Create(fields contracts.Fields) T {
	result, err := table.CreateE(fields)
	if err != nil {
		panic(err)
	}
	return *result
}

// Update 更新数据库中的记录
// update records in the database.
func (table *Table[T]) Update(fields contracts.Fields) int64 {
	result, _ := table.UpdateE(fields)
	return result
}

// FindOrFail 按 ID 对单个记录执行查询
// Execute a query for a single record by ID.
func (table *Table[T]) FindOrFail(key any) T {
	result := table.Find(key)
	if result == nil {
		panic(NotFoundException{Err: sql.ErrNoRows})
	}
	return *result
}

// FirstOr 执行查询并获得第一个结果或调用回调
// Execute the query and get the first result or call a callback.
func (table *Table[T]) FirstOr(provider contracts.InstanceProvider[T]) T {
	if result := table.First(); result != nil {
		return *result
	}
	return provider()
}

// FirstWhere 向查询添加基本 where 子句，并返回第一个结果
// Add a basic where clause to the query, and return the first result.
func (table *Table[T]) FirstWhere(column string, args ...any) *T {
	return table.Where(column, args...).First()
}

// FirstWhereE 向查询添加基本 where 子句，并返回第一个结果
// Add a basic where clause to the query, and return the first result.
func (table *Table[T]) FirstWhereE(column string, args ...any) (*T, contracts.Exception) {
	return table.Where(column, args...).FirstE()
}

// Paginate 对给定的查询进行分页。
// paginate the given query.
func (table *Table[T]) Paginate(perPage int64, current ...int64) (contracts.Collection[T], int64) {
	return table.WithPagination(perPage, current...).Get(), table.Count()
}

// SimplePaginate 将给定的查询分页成一个简单的分页器
// paginate the given query into a simple paginator.
func (table *Table[T]) SimplePaginate(perPage int64, current ...int64) contracts.Collection[T] {
	return table.WithPagination(perPage, current...).Get()
}
