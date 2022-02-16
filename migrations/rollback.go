package migrations

import (
	"fmt"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/table"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
)

func Rollback(app contracts.Application) contracts.Command {
	return &rollback{
		production: app.IsProduction(),
		redis:      app.Get("redis").(contracts.RedisConnection),
		db:         app.Get("db.factory").(contracts.DBFactory),
		table:      app.Get("migrations.table").(string),
		migrations: app.Get("migrations").(contracts.Migrations),
		Command:    commands.Base("migrate:rollback {--force:是否在生产环境强制执行}", "Run the database migrations"),
	}
}

type rollback struct {
	commands.Command
	production bool
	table      string
	redis      contracts.RedisConnection
	db         contracts.DBFactory
	migrations contracts.Migrations
}

func (this *rollback) Handle() interface{} {
	if this.production && !this.GetBool("force") {
		logs.WithError(MustForceErr).Error("refresh.Handle: ")
		return MustForceErr
	}

	var (
		raw        = getMigrations(this.db.Connection(), this.table)
		migrations = collection.MustNew(this.migrations).Pluck("name")
	)

	if raw.Len() == 0 {
		logs.Default().Info("rollback.Handle: No migrations need to be rolled back")
		return nil
	}

	raw.Where("batch", raw.Max("batch")).Map(func(item contracts.Fields) {
		migration, ok := migrations[item["migration"].(string)].(contracts.Migrate)
		if ok {
			conn := this.db.Connection(migration.Connection)
			logs.Default().Info(fmt.Sprintf("rollback.Handle: %s rollbacking", migration.Name))
			if err := migration.Down(conn); err != nil {
				logs.WithError(err).Error(fmt.Sprintf("rollback.Handle: %s failed to rollback", migration.Name))
				panic(err)
			}
			logs.Default().Info(fmt.Sprintf("rollback.Handle: %s rollbacked", migration.Name))
			table.WithConnection(this.table, conn).
				Where("migration", item["migration"]).
				Where("batch", item["batch"]).
				Delete()
		} else {
			logs.Default().Warn(fmt.Sprintf("rollback.Handle: migration %s is not exists", migration.Name))
		}
	})

	return nil
}
