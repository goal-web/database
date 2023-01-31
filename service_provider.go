package database

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/drivers"
	"github.com/goal-web/database/migrations"
)

type ServiceProvider struct {
	migrations contracts.Migrations
}

func NewService(migrations contracts.Migrations) contracts.ServiceProvider {
	return &ServiceProvider{migrations: migrations}
}

func (provider *ServiceProvider) Register(application contracts.Application) {
	application.Instance("migrations", provider.migrations)
	application.Singleton("migrations.table", func(config contracts.Config) string {
		return config.Get("database").(Config).Migrations
	})

	application.Singleton("db.factory", func(config contracts.Config) contracts.DBFactory {
		events, _ := application.Get("events").(contracts.EventDispatcher)
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

	// 请确保 console 在 database 之前注册，否则迁移命令无法注册到 console 中
	if console, ok := application.Get("console").(contracts.Console); ok {
		console.RegisterCommand("migrate", migrations.Migrate)
		console.RegisterCommand("migrate:rollback", migrations.Rollback)
		console.RegisterCommand("migrate:refresh", migrations.Refresh)
		console.RegisterCommand("migrate:reset", migrations.Reset)
		console.RegisterCommand("migrate:status", migrations.Status)
	}
}

func (provider *ServiceProvider) Start() error {
	return nil
}

func (provider *ServiceProvider) Stop() {
}
