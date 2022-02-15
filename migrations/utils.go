package migrations

import (
	"fmt"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/table"
	"github.com/goal-web/supports/logs"
)

func Transaction(sql string) contracts.MigrateHandler {
	return func(db contracts.DBConnection) error {
		return db.Transaction(func(executor contracts.SqlExecutor) error {
			_, err := executor.Exec(sql)
			return err
		})
	}
}

func Exec(sql string) contracts.MigrateHandler {
	return func(db contracts.DBConnection) error {
		_, err := db.Exec(sql)
		return err
	}
}

func getMigrations(db contracts.DBConnection, tableName string) contracts.Collection {
	query := table.WithConnection(tableName, db)
	ddl := fmt.Sprintf("create table %s\n(\n    id        int unsigned auto_increment\n        primary key,\n    migration varchar(191) not null,\n    batch     int          not null\n)\n", tableName)
	_, err := db.Exec(ddl)

	if err == nil {
		logs.Default().Info("migrations.getMigrations: Migration table has been generated")
		return collection.MustNew([]contracts.Fields{})
	}

	return query.Get()
}
