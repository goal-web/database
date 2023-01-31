package database

import (
	"errors"
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/drivers"
	"github.com/goal-web/supports/utils"
)

type Factory struct {
	events      contracts.EventDispatcher
	connections map[string]contracts.DBConnection
	drivers     map[string]contracts.DBConnector
	dbConfig    Config
}

func NewFactory(config Config, events contracts.EventDispatcher) contracts.DBFactory {
	return &Factory{
		events:      events,
		dbConfig:    config,
		connections: make(map[string]contracts.DBConnection),
		drivers: map[string]contracts.DBConnector{
			"mysql":      drivers.MysqlConnector,
			"postgres":   drivers.PostgresSqlConnector,
			"sqlite":     drivers.SqliteConnector,
			"clickhouse": drivers.ClickHouseConnector,
		},
	}
}

func (factory *Factory) Connection(name ...string) contracts.DBConnection {
	connection := factory.dbConfig.Default
	if len(name) > 0 && name[0] != "" {
		connection = name[0]
	}
	if conn, existsConnection := factory.connections[connection]; existsConnection {
		return conn
	}

	factory.connections[connection] = factory.make(connection)

	return factory.connections[connection]
}

func (factory *Factory) Extend(name string, driver contracts.DBConnector) {
	factory.drivers[name] = driver
}

func (factory *Factory) make(name string) contracts.DBConnection {
	config := factory.dbConfig

	if connectionConfig, existsConnection := config.Connections[name]; existsConnection {
		driverName := utils.GetStringField(connectionConfig, "driver")
		if driver, existsDriver := factory.drivers[driverName]; existsDriver {
			return driver(connectionConfig, factory.events)
		}

		panic(DBConnectionException{
			Err:    errors.New("该数据库驱动不存在：" + driverName),
			Code:   DbDriverDontExist,
			Config: connectionConfig,
		})
	}

	panic(DBConnectionException{
		Err:        errors.New("数据库连接不存在：" + name),
		Code:       DbConnectionDontExist,
		Connection: name,
	})
}
