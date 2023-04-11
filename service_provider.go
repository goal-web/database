package database

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/table"
)

type ServiceProvider struct {
	app contracts.Application
}

func NewService() contracts.ServiceProvider {
	return &ServiceProvider{}
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
	table.SetFactory(provider.app.Get("db.factory").(contracts.DBFactory))
	return nil
}

func (provider *ServiceProvider) Stop() {
}
