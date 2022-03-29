package support

import (
	"database/sql"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/jmoiron/sqlx"
	"strconv"
)

func ParseRows(rows *sqlx.Rows) (results []contracts.Fields, err error) {
	if err = rows.Err(); err != nil {
		return nil, err
	}
	columns, colErr := rows.Columns()
	if colErr != nil {
		return nil, colErr
	}
	colTypes, typeErr := rows.ColumnTypes()
	if typeErr != nil {
		return nil, typeErr
	}

	columnsLen := len(columns)
	for rows.Next() {
		var colVar = make([]interface{}, columnsLen)
		for i := 0; i < columnsLen; i++ {
			SetColVarType(&colVar, i, colTypes[i].DatabaseTypeName())
		}
		result := make(map[string]interface{})
		if scanErr := rows.Scan(colVar...); scanErr != nil {
			panic(scanErr)
		}
		for j := 0; j < columnsLen; j++ {
			SetResultValue(&result, columns[j], colVar[j], colTypes[j].DatabaseTypeName())
		}
		results = append(results, result)
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

func SetColVarType(colVar *[]interface{}, i int, typeName string) {
	switch typeName {
	case "INT", "TINYINT", "MEDIUMINT", "SMALLINT", "BIGINT":
		var s sql.NullInt64
		(*colVar)[i] = &s
	case "FLOAT", "DOUBLE":
		var s sql.NullFloat64
		(*colVar)[i] = &s
	case "DECIMAL":
		var s []uint8
		(*colVar)[i] = &s
	case "TEXT", "MEDIUMTEXT", "TINYTEXT", "LONGTEXT", "VARCHAR", "TIMESTAMP", "DATETIME", "YEAR", "DATE", "TIME", "JSON":
		var s sql.NullString
		(*colVar)[i] = &s
	default:
		var s interface{}
		(*colVar)[i] = &s
	}
}

func SetResultValue(result *map[string]interface{}, index string, colVar interface{}, typeName string) {
	switch typeName {
	case "INT", "TINYINT", "MEDIUMINT", "SMALLINT", "BIGINT":
		temp := *(colVar.(*sql.NullInt64))
		if temp.Valid {
			(*result)[index] = temp.Int64
		} else {
			(*result)[index] = nil
		}
	case "FLOAT", "DOUBLE":
		temp := *(colVar.(*sql.NullFloat64))
		if temp.Valid {
			(*result)[index] = temp.Float64
		} else {
			(*result)[index] = nil
		}
	case "DECIMAL":
		if len(*(colVar.(*[]uint8))) < 1 {
			(*result)[index] = nil
		} else {
			(*result)[index], _ = strconv.ParseFloat(string(*(colVar.(*[]uint8))), 64)
		}
	case "TEXT", "MEDIUMTEXT", "TINYTEXT", "LONGTEXT", "VARCHAR", "TIMESTAMP", "DATETIME", "YEAR", "DATE", "TIME", "JSON":
		temp := *(colVar.(*sql.NullString))
		if temp.Valid {
			(*result)[index] = temp.String
		} else {
			(*result)[index] = nil
		}
	default:
		(*result)[index] = colVar
	}
}
