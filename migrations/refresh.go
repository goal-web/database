package migrations

import (
	"fmt"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/table"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
)

func Refresh(app contracts.Application) contracts.Command {
	return &refresh{
		production: app.IsProduction(),
		redis:      app.Get("redis").(contracts.RedisConnection),
		db:         app.Get("db.factory").(contracts.DBFactory),
		table:      app.Get("migrations.table").(string),
		migrations: app.Get("migrations").(contracts.Migrations),
		Command:    commands.Base("migrate:refresh {--force:是否在生产环境强制执行}", "Run the database migrations"),
	}
}

type refresh struct {
	commands.Command
	production bool
	table      string
	redis      contracts.RedisConnection
	db         contracts.DBFactory
	migrations contracts.Migrations
}

func (cmd *refresh) Handle() interface{} {
	if cmd.production && !cmd.GetBool("force") {
		logs.WithError(MustForceErr).Error("refresh.Handle: ")
		return MustForceErr
	}

	// rollback all migrations
	if raw := getMigrations(cmd.db.Connection(), cmd.table); raw.Len() > 0 {
		var migrations = collection.MustNew(cmd.migrations).Pluck("name")

		raw.Map(func(item contracts.Fields) {
			migration, ok := migrations[item["migration"].(string)].(contracts.Migrate)
			if ok {
				conn := cmd.db.Connection(migration.Connection)
				logs.Default().Info(fmt.Sprintf("rollback.Handle: %s rollbacking", migration.Name))
				if err := migration.Down(conn); err != nil {
					logs.WithError(err).Error(fmt.Sprintf("rollback.Handle: %s failed to rollback", migration.Name))
					panic(err)
				}
				logs.Default().Info(fmt.Sprintf("rollback.Handle: %s rollbacked", migration.Name))
				table.WithConnection(cmd.table, conn).Where("id", item["id"]).Delete()
			} else {
				logs.Default().Warn(fmt.Sprintf("rollback.Handle: migration %s is not exists", migration.Name))
			}
		})
	}

	var (
		migrations = collection.MustNew(cmd.migrations).Sort(func(migrate contracts.Migrate, next contracts.Migrate) bool {
			return migrate.CreatedAt.Before(next.CreatedAt)
		}).ToInterfaceArray()
		executedNum = 0
	)

	for _, item := range migrations {
		migration := item.(contracts.Migrate)
		conn := cmd.db.Connection(migration.Connection)
		logs.Default().Info(fmt.Sprintf("migrate.Handle: %s migrating", migration.Name))
		if err := migration.Up(conn); err != nil {
			logs.Default().WithError(err).Error(fmt.Sprintf("migrate.Handle: %s failed to execute", migration.Name))
			return err
		}
		logs.Default().Info(fmt.Sprintf("migrate.Handle: %s migrated", migration.Name))
		executedNum++
		table.WithConnection(cmd.table, conn).Insert(contracts.Fields{
			"batch":     1,
			"migration": migration.Name,
		})
	}

	if executedNum == 0 {
		logs.Default().Info("migrate.Handle: No migration was performed")
	}

	return nil
}
