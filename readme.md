# Goal-web/collection
goal 框架的数据库组件，当然你也可以在 goal 之外的框架使用他。
> 目前数据库组件暂时不能使用关联关系，你可以用 `WhereExists` 来代替

## 安装 - install
```bash
go get github.com/goal-web/database
```

## 使用
goal 的脚手架自带了绝大多数开发一个 web 应用的所需要的功能和组件，当然包括了数据库组件。一般情况下，我们只需要在 .env 修改自己的数据库配置即可，添加数据库连接可以 `config/database.go` 修改 `Connections` 属性。

### 配置
默认情况下，`config/database` 配置文件像下面那样，默认添加了 sqlite、MySQL、postgresSql 三个数据库连接的配置
> 和 `Laravel` 不同的是，goal 把 redis 配置独立出去了，因为 redis 也是一个独立的模块，不想让 redis 依赖 database
```go
package config

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/database"
)

func init() {
	configs["database"] = func(env contracts.Env) interface{} {
		return database.Config{
			Default: env.StringOption("db.connection", "mysql"),
			Connections: map[string]contracts.Fields{
				"sqlite": {
					"driver":   "sqlite",
					"database": env.GetString("sqlite.database"),
				},
				"mysql": {
					"driver":          "mysql",
					"host":            env.GetString("db.host"),
					"port":            env.GetString("db.port"),
					"database":        env.GetString("db.database"),
					"username":        env.GetString("db.username"),
					"password":        env.GetString("db.password"),
					"unix_socket":     env.GetString("db.unix_socket"),
					"charset":         env.StringOption("db.charset", "utf8mb4"),
					"collation":       env.StringOption("db.collation", "utf8mb4_unicode_ci"),
					"prefix":          env.GetString("db.prefix"),
					"strict":          env.GetBool("db.struct"),
					"max_connections": env.GetInt("db.max_connections"),
					"max_idles":       env.GetInt("db.max_idles"),
				},
				"pgsql": {
					"driver":          "postgres",
					"host":            env.GetString("db.pgsql.host"),
					"port":            env.GetString("db.pgsql.port"),
					"database":        env.GetString("db.pgsql.database"),
					"username":        env.GetString("db.pgsql.username"),
					"password":        env.GetString("db.pgsql.password"),
					"charset":         env.StringOption("db.pgsql.charset", "utf8mb4"),
					"prefix":          env.GetString("db.pgsql.prefix"),
					"schema":          env.StringOption("db.pgsql.schema", "public"),
					"sslmode":         env.StringOption("db.pgsql.sslmode", "disable"),
					"max_connections": env.GetInt("db.pgsql.max_connections"),
					"max_idles":       env.GetInt("db.pgsql.max_idles"),
				},
			},
		}
	}
}

```
`.env` 的数据库相关配置
```bash
# 默认连接
db.connection=sqlite

sqlite.database=/Users/qbhy/project/go/goal-web/goal/example/database/db.sqlite

db.host=localhost
db.port=3306
db.database=goal
db.username=root
db.password=password

db.pgsql.host=localhost
db.pgsql.port=55433
db.pgsql.database=postgres
db.pgsql.username=postgres
db.pgsql.password=123456
```

### 定义模型
`app/models/user.go` 文件

```go
package models

import (
	"github.com/goal-web/database/table"
	"github.com/goal-web/supports/class"
)

// UserClass 这个类变量，以后大有用处
var UserClass = class.Make(new(User))

// UserModel 返回 table 实例，继承自查询构造器并且实现了所有 future
func UserModel() *table.Table {
	return table.Model(UserClass, "users")
}

// User 模型结构体
type User struct {
	Id       int64  `json:"id"`
	NickName string `json:"name"`
}
```

```go

```

## 在 goal 之外的框架使用
略



[goal-web](https://github.com/goal-web/goal)  
[goal-web/collection](https://github.com/goal-web/collection)  
qbhy0715@qq.com
