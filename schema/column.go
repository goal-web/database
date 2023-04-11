package schema

import (
	"fmt"
	"github.com/goal-web/querybuilder"
)

type ColumnDefinition struct {
	Name                 string
	Table                string
	AfterColumn          string
	IsAutoIncrement      bool
	IsChange             bool
	CharsetValue         string
	CollationValue       string
	CommentValue         string
	DefaultValue         any
	IsFirst              bool
	StartValue           int
	IndexName            string
	FulltextIndexName    string
	SpatialIndexName     string
	UniqueIndexName      string
	IsInvisible          bool
	IsNullable           bool
	IsPersisted          bool
	IsPrimary            bool
	StoredAsExpression   string
	TypeValue            string
	IsUnsigned           bool
	IsUseCurrent         bool
	IsUseCurrentOnUpdate bool
	VirtualAsExpression  string
	// float
	TotalValue  int // 总长度
	PlacesValue int // 小数点长度

	// time
	PrecisionValue int
}

// After Place the column "after" another column (MySQL)
func (column *ColumnDefinition) After(columnName string) *ColumnDefinition {
	column.AfterColumn = columnName
	return column
}

// Always Used as a modifier for generatedAs()(PostgreSQL)
func (column *ColumnDefinition) Always(value ...bool) *ColumnDefinition {
	// TODO
	return column
}

// AutoIncrement Set INTEGER columns as auto-increment (primary key)
func (column *ColumnDefinition) AutoIncrement() *ColumnDefinition {
	column.IsAutoIncrement = true
	return column
}

// Change the column
func (column *ColumnDefinition) Change() *ColumnDefinition {
	column.IsChange = true
	return column
}

// Charset Specify a character set for the column (MySQL)
func (column *ColumnDefinition) Charset(charset string) *ColumnDefinition {
	column.CharsetValue = charset
	return column
}

// Collation Specify a collation for the column (MySQL/PostgreSQL/SQL Server)
func (column *ColumnDefinition) Collation(collation string) *ColumnDefinition {
	column.CollationValue = collation
	return column
}

// Comment Add a comment to the column (MySQL/PostgreSQL)
func (column *ColumnDefinition) Comment(comment string) *ColumnDefinition {
	column.CommentValue = comment
	return column
}

// Default Specify a "default" value for the column
func (column *ColumnDefinition) Default(value any) *ColumnDefinition {
	column.DefaultValue = value
	return column
}

// First Place the column "first" in the table (MySQL)
func (column *ColumnDefinition) First() *ColumnDefinition {
	column.IsFirst = true
	return column
}

// From Set the starting value of an auto-incrementing field (MySQL / PostgreSQL)
func (column *ColumnDefinition) From(startingValue int) *ColumnDefinition {
	column.StartValue = startingValue
	return column
}

// GeneratedAs Create a SQL compliant identity column (PostgreSQL)
func (column *ColumnDefinition) GeneratedAs(expression ...querybuilder.Expression) *ColumnDefinition {
	// TODO
	return column
}

// Index Add an index
func (column *ColumnDefinition) Index(indexName ...string) *ColumnDefinition {
	if len(indexName) > 0 {
		column.IndexName = indexName[0]
	} else {
		column.IndexName = fmt.Sprintf("%s_%s_index", column.Table, column.Name)
	}
	return column
}

// Invisible Specify that the column should be invisible to "SELECT *" (MySQL)
func (column *ColumnDefinition) Invisible() *ColumnDefinition {
	column.IsInvisible = true
	return column
}

// Nullable Allow NULL values to be inserted into the column
func (column *ColumnDefinition) Nullable(value ...bool) *ColumnDefinition {
	if len(value) > 0 {
		column.IsNullable = value[0]
	} else {
		column.IsNullable = true
	}
	return column
}

// Persisted Mark the computed generated column as persistent (SQL Server)
func (column *ColumnDefinition) Persisted() *ColumnDefinition {
	column.IsPersisted = true
	return column
}

// Primary Add a primary index
func (column *ColumnDefinition) Primary() *ColumnDefinition {
	column.IsPrimary = true
	return column
}

// Fulltext Add a fulltext index
func (column *ColumnDefinition) Fulltext(indexName ...string) *ColumnDefinition {
	if len(indexName) > 0 {
		column.FulltextIndexName = indexName[0]
	} else {
		column.FulltextIndexName = fmt.Sprintf("%s_%s_fulltext_index", column.Table, column.Name)
	}
	return column
}

// SpatialIndex Add a spatial index
func (column *ColumnDefinition) SpatialIndex(indexName ...string) *ColumnDefinition {
	if len(indexName) > 0 {
		column.SpatialIndexName = indexName[0]
	} else {
		column.SpatialIndexName = fmt.Sprintf("%s_%s_spatial_index", column.Table, column.Name)
	}
	return column
}

// StartingValue create a stored generated column (MySQL/PostgreSQL/SQLite)
func (column *ColumnDefinition) StartingValue(startingValue int) *ColumnDefinition {
	// TODO
	return column
}

// StoredAs Create a stored generated column (MySQL/PostgreSQL/SQLite)
func (column *ColumnDefinition) StoredAs(expression string) *ColumnDefinition {
	column.StoredAsExpression = expression
	return column
}

// Type Specify a type for the column
func (column *ColumnDefinition) Type(typeValue string) *ColumnDefinition {
	column.TypeValue = typeValue
	return column
}

// Unique Add a unique index
func (column *ColumnDefinition) Unique(indexName ...string) *ColumnDefinition {
	if len(indexName) > 0 {
		column.UniqueIndexName = indexName[0]
	} else {
		column.UniqueIndexName = fmt.Sprintf("%s_%s_unique_index", column.Table, column.Name)
	}
	return column
}

// Unsigned Set the INTEGER column as UNSIGNED (MySQL)
func (column *ColumnDefinition) Unsigned() *ColumnDefinition {
	column.IsUnsigned = true
	return column
}

func (column *ColumnDefinition) Float(total, places int) *ColumnDefinition {
	column.TotalValue = total
	column.PlacesValue = places
	return column
}

// UseCurrent Set the TIMESTAMP column to use CURRENT_TIMESTAMP as default value
func (column *ColumnDefinition) UseCurrent() *ColumnDefinition {
	column.IsUseCurrent = true
	return column
}

// UseCurrentOnUpdate Set the TIMESTAMP column to use CURRENT_TIMESTAMP when updating (MySQL)
func (column *ColumnDefinition) UseCurrentOnUpdate() *ColumnDefinition {
	column.IsUseCurrentOnUpdate = true
	return column
}

// VirtualAs Create a virtual generated column (MySQL/PostgreSQL/SQLite)
func (column *ColumnDefinition) VirtualAs(expression string) *ColumnDefinition {
	column.VirtualAsExpression = expression
	return column
}

func (column *ColumnDefinition) Precision(precision int) *ColumnDefinition {
	column.PrecisionValue = precision
	return column
}
