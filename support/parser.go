package support

import (
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/jmoiron/sqlx"
)

func ParseRows(rows *sqlx.Rows) (results []contracts.Fields, err error) {
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var item = make(contracts.Fields)
		err = rows.MapScan(item)
		if err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	return
}

func ParseRowsToCollection(rows *sqlx.Rows) (contracts.Collection, error) {
	data, parseErr := ParseRows(rows)
	if parseErr != nil {
		return nil, parseErr
	}
	return collection.MustNew(data), nil
}
