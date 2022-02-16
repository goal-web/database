package migrations

import (
	"errors"
	"fmt"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/table"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
)

func Migrate(app contracts.Application) contracts.Command {
	return &migrate{
		production: app.IsProduction(),
		redis:      app.Get("redis").(contracts.RedisConnection),
		db:         app.Get("db.factory").(contracts.DBFactory),
		table:      app.Get("migrations.table").(string),
		migrations: app.Get("migrations").(contracts.Migrations),
		Command:    commands.Base("migrate {--force:是否在生产环境强制执行}", "Run the database migrations"),
	}
}

type migrate struct {
	commands.Command
	production bool
	table      string
	redis      contracts.RedisConnection
	db         contracts.DBFactory
	migrations contracts.Migrations
}

var MustForceErr = errors.New("use the force option in production")

func (this *migrate) Handle() interface{} {
	if this.production && !this.GetBool("force") {
		logs.WithError(MustForceErr).Error("refresh.Handle: ")
		return MustForceErr
	}

	var (
		raw           = getMigrations(this.db.Connection(), this.table)
		executedNum   = 0
		migratedItems = raw.Pluck("migration")
	)

	batch := raw.Max("batch") + 1

	migrations := collection.MustNew(this.migrations).Sort(func(migrate contracts.Migrate, next contracts.Migrate) bool {
		return migrate.CreatedAt.Before(next.CreatedAt)
	}).ToInterfaceArray()

	for _, item := range migrations {
		migration := item.(contracts.Migrate)
		if _, exists := migratedItems[migration.Name]; !exists {
			conn := this.db.Connection(migration.Connection)
			logs.Default().Info(fmt.Sprintf("migrate.Handle: %s migrating", migration.Name))
			if err := migration.Up(conn); err != nil {
				logs.Default().WithError(err).Error(fmt.Sprintf("migrate.Handle: %s failed to execute", migration.Name))
				panic(err)
			}
			logs.Default().Info(fmt.Sprintf("migrate.Handle: %s migrated", migration.Name))
			executedNum++
			table.WithConnection(this.table, conn).Insert(contracts.Fields{
				"batch":     batch,
				"migration": migration.Name,
			})
		}
	}

	if executedNum == 0 {
		logs.Default().Info("migrate.Handle: No migration was performed")
	}

	return nil
}
