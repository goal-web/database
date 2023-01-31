package database

import (
	"github.com/goal-web/contracts"
)

type Config struct {
	// 默认数据库连接
	Default string
	
	// 数据库连接配置
	Connections map[string]contracts.Fields

	// 迁移表名字
	Migrations string
}
