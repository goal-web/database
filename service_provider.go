package database

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/drivers"
	"github.com/goal-web/database/migrations"
)

type ServiceProvider struct {
	migrations contracts.Migrations
}

func Service(migrations contracts.Migrations) contracts.ServiceProvider {
	return &ServiceProvider{migrations: migrations}
}

func (this *ServiceProvider) Register(application contracts.Application) {
	application.Instance("migrations", this.migrations)
	application.Singleton("migrations.table", func(config contracts.Config) string {
		return config.Get("database").(Config).Migrations
	})
	application.Singleton("db.factory", func(config contracts.Config, events contracts.EventDispatcher) contracts.DBFactory {
		return &Factory{
			events:      events,
			config:      config,
			dbConfig:    config.Get("database").(Config),
			connections: make(map[string]contracts.DBConnection),
			drivers: map[string]contracts.DBConnector{
				"mysql":      drivers.MysqlConnector,
				"postgres":   drivers.PostgresSqlConnector,
				"sqlite":     drivers.SqliteConnector,
				"clickhouse": drivers.ClickHouseConnector,
			},
		}
	})
	application.Singleton("db", func(config contracts.Config, factory contracts.DBFactory) contracts.DBConnection {
		return factory.Connection()
	})

	// 一定要确保 console 在 database 之前注册
	application.Call(func(console contracts.Console) {
		console.RegisterCommand("migrate", migrations.Migrate)
		console.RegisterCommand("migrate:rollback", migrations.Rollback)
		console.RegisterCommand("migrate:refresh", migrations.Refresh)
		console.RegisterCommand("migrate:reset", migrations.Reset)
		console.RegisterCommand("migrate:status", migrations.Status)
	})
}

func (this *ServiceProvider) Start() error {
	return nil
}

func (this *ServiceProvider) Stop() {
}
