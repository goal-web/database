package drivers

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"github.com/jmoiron/sqlx"
)

type Mysql struct {
	*Base
}

func MysqlConnector(config contracts.Fields, events contracts.EventDispatcher) contracts.DBConnection {
	dsn := utils.GetStringField(config, "unix_socket")
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s%s",
			utils.GetStringField(config, "username"),
			utils.GetStringField(config, "password"),
			utils.GetStringField(config, "host"),
			utils.GetStringField(config, "port"),
			utils.GetStringField(config, "database"),
			utils.GetStringField(config, "charset"),
			utils.GetStringField(config, "other"),
		)
	}
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		logs.WithError(err).WithField("config", config).Fatal("mysql 数据库初始化失败")
	}
	db.SetMaxOpenConns(utils.GetIntField(config, "max_connections"))
	db.SetMaxIdleConns(utils.GetIntField(config, "max_idles"))

	return &Mysql{NewDriver(db, events)}
}
