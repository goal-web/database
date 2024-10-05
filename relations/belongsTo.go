package relations

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/table"
)

type BelongsTo[T any, P contracts.ModelContext] struct {
	query           *table.Table[T]
	foreignKey      string
	ownerKey        string
	relation        contracts.RelationType
	foreignKeyValue any
}

func (b *BelongsTo[T, P]) GetRelationCollector() contracts.RelationCollector {
	return func(keys []any) []any {
		return b.query.WhereIn(b.ownerKey, keys).Get().ToAnyArray()
	}
}

func (b *BelongsTo[T, P]) GetForeignKeysCollector() contracts.ForeignKeysCollector[P] {
	return func(item *P) any {
		return (*item).Get(b.foreignKey)
	}
}

func (b *BelongsTo[T, P]) GetRelationSetter() contracts.RelationSetter[P] {
	return func(item *P, value any) {
		(*item).Set(contracts.Fields{
			string(b.relation): value,
		})
	}
}

func NewBelongsTo[T any, P contracts.ModelContext](query *table.Table[T], foreignKey, ownerKey string, relation contracts.RelationType) *BelongsTo[T, P] {
	return &BelongsTo[T, P]{
		query:      query,
		foreignKey: foreignKey,
		ownerKey:   ownerKey,
		relation:   relation,
	}
}

func (b *BelongsTo[T, P]) Where(q contracts.QueryFunc[T]) *BelongsTo[T, P] {
	b.query.WhereFunc(q)
	return b
}

func (b *BelongsTo[T, P]) Get() *T {
	return b.query.Where(b.ownerKey, b.foreignKeyValue).First()
}

func (b *BelongsTo[T, P]) SetForeignKey(value any) *BelongsTo[T, P] {
	b.foreignKeyValue = value
	return b
}

func (b *BelongsTo[T, P]) GetRelation() contracts.RelationType {
	return b.relation
}
