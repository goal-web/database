package tests

import (
	"github.com/goal-web/application"
	"github.com/goal-web/application/exceptions"
	"github.com/goal-web/config"
	"github.com/goal-web/contracts"
	"github.com/goal-web/database"
	"github.com/goal-web/database/table"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMysqlDatabaseService(t *testing.T) {
	app := application.Singleton()
	hostname, _ := os.Hostname()
	userHome, _ := os.UserHomeDir()

	app.RegisterServices(
		exceptions.NewService([]contracts.Exception{}),
		config.NewService("testing", ".", map[string]contracts.ConfigProvider{
			"app": application.ConfigProvider(hostname, userHome),
			"database": func(env contracts.Env) interface{} {
				return database.Config{
					Default: "mysql",
					Connections: map[string]contracts.Fields{
						"mysql": {
							"driver":          "mysql",
							"host":            "localhost",
							"port":            "3306",
							"database":        "goal",
							"username":        "root",
							"password":        "123456",
							"charset":         env.StringOption("db.charset", "utf8mb4"),
							"collation":       env.StringOption("db.collation", "utf8mb4_unicode_ci"),
							"prefix":          env.GetString("db.prefix"),
							"strict":          env.GetBool("db.struct"),
							"max_connections": env.GetInt("db.max_connections"),
							"max_idles":       env.GetInt("db.max_idles"),
						},
					},
					Migrations: "migrations",
				}
			},
		}),
		database.NewService(contracts.Migrations{}),
	)

	app.Start()

	assert.True(t, table.Query("users").Count() == 0)

	user := table.Query("users").Create(contracts.Fields{
		"name": "testing",
	})
	assert.NotNil(t, user)
	assert.True(t, user.(contracts.Fields)["name"] == "testing")
	assert.True(t, table.Query("users").Count() == 1)
	table.Query("users").Where("name", "testing").Delete()
	assert.True(t, table.Query("users").Count() == 0)

}
