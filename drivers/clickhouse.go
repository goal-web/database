package drivers

import (
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Clickhouse struct {
	*Base
}

func ClickHouseConnector(config contracts.Fields, events contracts.EventDispatcher) contracts.DBConnection {
	dsn := utils.GetStringField(config, "dsn")
	if dsn == "" {
		address := config["address"].([]string)
		dsn = fmt.Sprintf("tcp://%s?debug=%s",
			address[0],
			//config["database"].(string),
			config["debug"].(string),
		)
		if config["username"] != "" {
			dsn = fmt.Sprintf("%s&username=%s&password=%s",
				dsn,
				config["username"].(string),
				config["password"].(string),
			)
		}
		if len(address) > 1 {
			dsn = fmt.Sprintf("%s&alt_hosts=%s", dsn, strings.Join(address[1:], ","))
		}
	}
	db, err := sqlx.Connect("clickhouse", dsn)
	if err != nil {
		logs.WithError(err).WithField("dsn", dsn).WithField("config", config).Fatal("clickhouse 数据库初始化失败")
	}
	db.SetMaxOpenConns(utils.GetIntField(config, "max_connections"))
	db.SetMaxIdleConns(utils.GetIntField(config, "max_idles"))
	return &Clickhouse{NewDriver(db, events)}
}
