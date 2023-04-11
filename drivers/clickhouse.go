package drivers

import (
	"errors"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/support"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Clickhouse struct {
	support.Executor
}

type TransactionException = exceptions.Exception

var exception = &TransactionException{
	Err: errors.New("begin is not supported"),
}

func (c Clickhouse) Begin() (contracts.DBTx, contracts.Exception) {
	return nil, exception
}

func (c Clickhouse) Transaction(f func(executor contracts.SqlExecutor) contracts.Exception) contracts.Exception {
	return exception
}

func ClickHouseConnector(config contracts.Fields, events contracts.EventDispatcher) contracts.DBConnection {
	var dsn = utils.GetStringField(config, "dsn")
	if dsn == "" {
		address := config["address"].([]string)
		dsn = fmt.Sprintf("tcp://%s?debug=%s",
			address[0],
			//config["database"].(string),
			config["debug"].(string),
		)
		if config["username"] != "" {
			dsn = fmt.Sprintf("%s&username=%s",
				dsn,
				config["username"].(string),
			)
		}
		if config["database"] != "" {
			dsn = fmt.Sprintf("%s&database=%s",
				dsn,
				config["database"].(string),
			)
		}
		if config["password"] != "" {
			dsn = fmt.Sprintf("%s&password=%s",
				dsn,
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
	return &Clickhouse{
		support.NewExecutor(db, events, DollarNParamBindWrapper),
	}
}
