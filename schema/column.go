package schema

import "github.com/goal-web/querybuilder"

type ColumnDefinition struct {
}

// after Place the column "after" another column (MySQL)
func (Column *ColumnDefinition) after(column string) *ColumnDefinition {
	// TODO
	return Column
}

// always Used as a modifier for generatedAs()(PostgreSQL)
func (Column *ColumnDefinition) always(value ...bool) *ColumnDefinition {
	// TODO
	return Column
}

// autoIncrement() Set INTEGER columns as auto-increment (primary key)
func (Column *ColumnDefinition) autoIncrement() *ColumnDefinition {
	// TODO
	return Column
}

// change() Change the column
func (Column *ColumnDefinition) change() *ColumnDefinition {
	// TODO
	return Column
}

// charset Specify a character set for the column (MySQL)
func (Column *ColumnDefinition) charset(charset string) *ColumnDefinition {
	// TODO
	return Column
}

// collation Specify a collation for the column (MySQL/PostgreSQL/SQL Server)
func (Column *ColumnDefinition) collation(collation string) *ColumnDefinition {
	// TODO
	return Column
}

// comment Add a comment to the column (MySQL/PostgreSQL)
func (Column *ColumnDefinition) comment(comment string) *ColumnDefinition {
	// TODO
	return Column
}

// Default Specify a "default" value for the column
func (Column *ColumnDefinition) Default(value interface{}) *ColumnDefinition {
	// TODO
	return Column
}

// first Place the column "first" in the table (MySQL)
func (Column *ColumnDefinition) first() *ColumnDefinition {
	// TODO
	return Column
}

// from Set the starting value of an auto-incrementing field (MySQL / PostgreSQL)
func (Column *ColumnDefinition) from(startingValue int) *ColumnDefinition {
	// TODO
	return Column
}

// generatedAs Create a SQL compliant identity column (PostgreSQL)
func (Column *ColumnDefinition) generatedAs(expression ...querybuilder.Expression) *ColumnDefinition {
	// TODO
	return Column
}

// index Add an index
func (Column *ColumnDefinition) index(indexName ...string) *ColumnDefinition {
	// TODO
	return Column
}

// invisible Specify that the column should be invisible to "SELECT *" (MySQL)
func (Column *ColumnDefinition) invisible() *ColumnDefinition {
	// TODO
	return Column
}

// nullable Allow NULL values to be inserted into the column
func (Column *ColumnDefinition) nullable(value ...bool) *ColumnDefinition {
	// TODO
	return Column
}

// persisted() Mark the computed generated column as persistent (SQL Server)
func (Column *ColumnDefinition) persisted() *ColumnDefinition {
	// TODO
	return Column
}

// primary() Add a primary index
func (Column *ColumnDefinition) primary() *ColumnDefinition {
	// TODO
	return Column
}

// fulltext Add a fulltext index
func (Column *ColumnDefinition) fulltext(indexName ...string) *ColumnDefinition {
	// TODO
	return Column
}

// spatialIndex Add a spatial index
func (Column *ColumnDefinition) spatialIndex(indexName ...string) *ColumnDefinition {
	// TODO
	return Column
}

// storedAs Create a stored generated column (MySQL/PostgreSQL/SQLite)
func (Column *ColumnDefinition) startingValue(startingValue int) *ColumnDefinition {
	// TODO
	return Column
}

// storedAs Create a stored generated column (MySQL/PostgreSQL/SQLite)
func (Column *ColumnDefinition) storedAs(expression string) *ColumnDefinition {
	// TODO
	return Column
}

// Type Specify a type for the column
func (Column *ColumnDefinition) Type(typeValue string) *ColumnDefinition {
	// TODO
	return Column
}

// Unique Add a unique index
func (Column *ColumnDefinition) unique(indexName ...string) *ColumnDefinition {
	// TODO
	return Column
}

// unsigned() Set the INTEGER column as UNSIGNED (MySQL)
func (Column *ColumnDefinition) unsigned() *ColumnDefinition {
	// TODO
	return Column
}

// useCurrent() Set the TIMESTAMP column to use CURRENT_TIMESTAMP as default value
func (Column *ColumnDefinition) useCurrent() *ColumnDefinition {
	// TODO
	return Column
}

// useCurrentOnUpdate() Set the TIMESTAMP column to use CURRENT_TIMESTAMP when updating (MySQL)
func (Column *ColumnDefinition) useCurrentOnUpdate() *ColumnDefinition {
	// TODO
	return Column
}

// virtualAs Create a virtual generated column (MySQL/PostgreSQL/SQLite)
func (Column *ColumnDefinition) virtualAs(expression string) *ColumnDefinition {
	// TODO
	return Column
}
