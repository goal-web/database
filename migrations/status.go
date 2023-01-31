package migrations

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/commands"
	"github.com/modood/table"
)

func Status(app contracts.Application) contracts.Command {
	return &status{
		production: app.IsProduction(),
		redis:      app.Get("redis").(contracts.RedisConnection),
		db:         app.Get("db.factory").(contracts.DBFactory),
		table:      app.Get("migrations.table").(string),
		migrations: app.Get("migrations").(contracts.Migrations),
		Command:    commands.Base("migrate:status", "Run the database migrations"),
	}
}

type status struct {
	commands.Command
	production bool
	table      string
	redis      contracts.RedisConnection
	db         contracts.DBFactory
	migrations contracts.Migrations
}

type MigrationStatus struct {
	Ran       string
	Migration string
	Batch     interface{}
}

func (cmd *status) Handle() interface{} {

	var (
		migrated = getMigrations(cmd.db.Connection(), cmd.table).Pluck("migration")
		data     = make([]MigrationStatus, 0)
	)

	for _, migration := range cmd.migrations {
		if migratedItem, exists := migrated[migration.Name].(contracts.Fields); exists {
			data = append(data, MigrationStatus{
				Ran:       "Yes",
				Migration: migration.Name,
				Batch:     migratedItem["batch"],
			})
		} else {
			data = append(data, MigrationStatus{
				Ran:       "No",
				Migration: migration.Name,
				Batch:     0,
			})
		}
	}

	table.Output(data)

	return nil
}
