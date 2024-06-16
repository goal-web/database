package table

import "github.com/goal-web/contracts"

func (table *Table[T]) Chunk(size int, handler func(collection contracts.Collection[*T], page int) contracts.Exception) (err contracts.Exception) {
	page := 1
	for err == nil {
		list, listErr := table.WithPagination(int64(size), int64(page)).GetE()
		if listErr != nil {
			return listErr
		}
		if list.IsEmpty() {
			return nil
		}
		err = handler(list, page)
		page++
		if list.Len() < size {
			break
		}
	}
	return
}

func (table *Table[T]) ChunkById(size int, handler func(collection contracts.Collection[*T], page int) (any, contracts.Exception)) contracts.Exception {
	page := 1
	var err contracts.Exception
	var id any
	for err == nil {
		list, listErr := table.When(id != nil, func(q contracts.Query[T]) contracts.Query[T] {
			return q.Where(table.primaryKeyField, ">", id)
		}).Take(int64(size)).GetE()
		if listErr != nil {
			return listErr
		}
		if list.IsEmpty() {
			return nil
		}
		id, err = handler(list, page)
		page++
		if list.Len() < size {
			break
		}
	}
	return err
}

// ChunkByIdDesc 通过比较 ID 对查询结果进行分块
// chunk the results of a query by comparing IDs.
func (table *Table[T]) ChunkByIdDesc(size int, handler func(collection contracts.Collection[*T], page int) (any, contracts.Exception)) contracts.Exception {
	page := 1
	var err contracts.Exception
	var id any
	for err == nil {
		list, listErr := table.OrderByDesc(table.primaryKeyField).
			When(id != nil, func(q contracts.Query[T]) contracts.Query[T] {
				return q.Where(table.primaryKeyField, "<", id)
			}).Take(int64(size)).GetE()
		if listErr != nil {
			return listErr
		}
		if list.IsEmpty() {
			return nil
		}
		id, err = handler(list, page)
		page++
		if list.Len() < size {
			break
		}
	}
	return err
}
