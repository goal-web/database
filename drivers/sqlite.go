package drivers

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	*Base
}

func SqliteConnector(config contracts.Fields, events contracts.EventDispatcher) contracts.DBConnection {
	db, err := sqlx.Connect("sqlite3", utils.GetStringField(config, "database"))

	if err != nil {
		logs.WithError(err).WithField("config", config).Fatal("sqlite 数据库初始化失败")
	}

	return &Sqlite{NewDriver(db, events)}
}
