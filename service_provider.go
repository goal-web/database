package database

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/table"
	"github.com/goal-web/supports/utils"
)

type ServiceProvider struct {
	app contracts.Application

	optimizeLoad bool // 是否尽快连数据库
}

func NewService(optimizeLoad ...bool) contracts.ServiceProvider {
	return &ServiceProvider{
		optimizeLoad: utils.DefaultBool(optimizeLoad, true),
	}
}

func (provider *ServiceProvider) Register(application contracts.Application) {
	provider.app = application
	application.Singleton("db.factory", func(config contracts.Config) contracts.DBFactory {
		events, _ := application.Get("events").(contracts.EventDispatcher)
		return NewFactory(config.Get("database").(Config), events)
	})
	application.Singleton("db", func(config contracts.Config, factory contracts.DBFactory) contracts.DBConnection {
		return factory.Connection()
	})
}

func (provider *ServiceProvider) Start() error {
	if provider.optimizeLoad {
		table.SetFactory(provider.app.Get("db.factory").(contracts.DBFactory))
	}
	return nil
}

func (provider *ServiceProvider) Stop() {
}
