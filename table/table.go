package table

import (
	"errors"
	"github.com/goal-web/contracts"
	"github.com/goal-web/querybuilder"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/utils"
)

type InstanceFactory[T any] func(fields contracts.Fields) *T

type Table[T any] struct {
	*querybuilder.Builder[T]
	executor contracts.SqlExecutor

	table             string
	primaryKeyField   string
	createdTimeColumn string
	UpdatedTimeColumn string

	instanceFactory InstanceFactory[T]

	relationCollectors   map[contracts.RelationType]contracts.RelationCollector
	foreignKeyCollectors map[contracts.RelationType]contracts.ForeignKeysCollector[T]
	relationSetters      map[contracts.RelationType]contracts.RelationSetter[T]
}

// SetConnection 参数要么是 contracts.DBConnection 要么是 string
func (table *Table[T]) SetConnection(connection any) *Table[T] {
	if conn, ok := connection.(contracts.DBConnection); ok {
		table.executor = conn
	} else {
		table.executor = getFactory().Connection(utils.ToString(connection, ""))
	}
	return table
}

// SetFactory 设置类
func (table *Table[T]) SetFactory(factory InstanceFactory[T]) *Table[T] {
	table.instanceFactory = factory
	return table
}

// SetRelation 设置关联关系
func (table *Table[T]) SetRelation(
	relation contracts.RelationType,
	collector contracts.ForeignKeysCollector[T],
	relationCollector contracts.RelationCollector,
	setter contracts.RelationSetter[T],
) *Table[T] {
	table.foreignKeyCollectors[relation] = collector
	table.relationCollectors[relation] = relationCollector
	table.relationSetters[relation] = setter
	return table
}

// SetByRelation 设置关联关系
func (table *Table[T]) SetByRelation(
	relation contracts.Relation[any, T],
) *Table[T] {
	table.foreignKeyCollectors[relation.GetRelation()] = relation.GetForeignKeysCollector()
	table.relationCollectors[relation.GetRelation()] = relation.GetRelationCollector()
	table.relationSetters[relation.GetRelation()] = relation.GetRelationSetter()
	return table
}

func (table *Table[T]) GetTable() string {
	return table.table
}

func (table *Table[T]) GetPrimaryKeyField() string {
	return table.primaryKeyField
}

// SetPrimaryKey 设置主键
func (table *Table[T]) SetPrimaryKey(name string) *Table[T] {
	table.primaryKeyField = name
	return table
}

// SetCreatedTimeColumn 设置创建时间字段
func (table *Table[T]) SetCreatedTimeColumn(column string) *Table[T] {
	table.createdTimeColumn = column
	return table
}

// SetUpdatedTimeColumn 设置更新时间字段
func (table *Table[T]) SetUpdatedTimeColumn(column string) *Table[T] {
	table.UpdatedTimeColumn = column
	return table
}

// getExecutor 获取 sql 语句的执行者
func (table *Table[T]) getExecutor() contracts.SqlExecutor {
	if table.executor != nil {
		return table.executor
	}
	return getFactory().Connection()
}

// SetExecutor 参数必须是 contracts.DBTx 实例
func (table *Table[T]) SetExecutor(executor contracts.SqlExecutor) *Table[T] {
	table.executor = executor
	return table
}

// SetWiths 初始化关联查询
func (table *Table[T]) SetWiths(withs ...contracts.RelationType) *Table[T] {
	table.With(withs...)
	return table
}

func (table *Table[T]) Delete() int64 {
	result, exception := table.DeleteE()
	if exception != nil {
		panic(exception)
	}
	return result
}

func (table *Table[T]) DeleteE() (int64, contracts.Exception) {
	sql, bindings := table.DeleteSql()
	result, exception := table.getExecutor().Exec(sql, bindings...)
	if exception != nil {
		return 0, &DeleteException{
			Sql:      sql,
			Bindings: bindings,
			Err:      errors.New("delete statement execution failed"),
			previous: exception,
		}
	}
	num, err := result.RowsAffected()
	if err != nil {
		return 0, &DeleteException{
			Sql:      sql,
			Bindings: bindings,
			Err:      errors.New("failed to get number of affected rows"),
			previous: exceptions.WithError(err),
		}
	}
	return num, nil
}
