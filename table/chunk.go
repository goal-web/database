package table

import "github.com/goal-web/contracts"

func (table *Table) Chunk(size int, handler func(collection contracts.Collection, page int) error) (err error) {
	page := 1
	for err == nil {
		newCollection := table.SimplePaginate(int64(size), int64(page))
		err = handler(newCollection, page)
		page++
		if newCollection.Len() < size {
			break
		}
	}
	return
}

func (table *Table) ChunkById(size int, handler func(collection contracts.Collection, page int) error) error {
	return table.OrderBy("id").Chunk(size, handler)
}
