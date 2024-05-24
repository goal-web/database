package tests

import (
	"fmt"
	"github.com/goal-web/application"
	"github.com/goal-web/config"
	"github.com/goal-web/contracts"
	"github.com/goal-web/database"
	"github.com/goal-web/database/table"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Age       int    `json:"age" db:"age"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var tableName string

func init() {
	app := application.Singleton()

	tableName = "test_" + utils.RandStr(10)

	app.RegisterServices(
		exceptions.NewService([]contracts.Exception{}),
		config.NewService(config.NewDotEnv(config.File("")), map[string]contracts.ConfigProvider{
			"app": func(env contracts.Env) any {
				return application.Config{
					Name:     env.GetString("app.name"),
					Debug:    env.GetBool("app.debug"),
					Timezone: env.GetString("app.timezone"),
					Env:      env.GetString("app.env"),
					Locale:   env.GetString("app.locale"),
					Key:      env.GetString("app.key"),
				}
			},
			"database": func(env contracts.Env) any {
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
							"charset":         env.StringOptional("db.charset", "utf8mb4"),
							"collation":       env.StringOptional("db.collation", "utf8mb4_unicode_ci"),
							"prefix":          env.GetString("db.prefix"),
							"strict":          env.GetBool("db.struct"),
							"max_connections": env.GetInt("db.max_connections"),
							"max_idles":       env.GetInt("db.max_idles"),
						},
					},
				}
			},
		}),
		database.NewService(),
	)

	app.Start()
	app.Call(func(connection contracts.DBConnection) {
		_, err := connection.Exec("CREATE TABLE IF NOT EXISTS " + tableName +
			"(" +
			"    `id`       INT UNSIGNED AUTO_INCREMENT," +
			"    name       varchar(20)," +
			"    age        int," +
			"    created_at timestamp," +
			"    updated_at timestamp," +
			"    PRIMARY KEY (`id`)" +
			") ENGINE = InnoDB" +
			"  DEFAULT CHARSET = utf8mb4;")
		if err != nil {
			panic(err)
		}
	})
}

func TestMysqlDatabaseService(t *testing.T) {

	assert.True(t, table.Query[User](tableName).Count() == 0)

	user := table.Query[User](tableName).Create(contracts.Fields{
		"name": "testing",
	})
	assert.NotNil(t, user)
	assert.True(t, user.Name == "testing")
	assert.True(t, table.Query[User](tableName).Count() == 1)
	assert.True(t, table.Query[User](tableName).Get().Count() == 1)
	table.Query[User](tableName).Where("name", "testing").Delete()
	assert.True(t, table.Query[User](tableName).Count() == 0)

}

func TestMysqlDatabaseWithoutApplication(t *testing.T) {
	// 实例化数据库工厂
	factory := database.NewFactory(
		database.Config{
			Default: "mysql",
			Connections: map[string]contracts.Fields{
				"mysql": {
					"driver":    "mysql",
					"host":      "localhost",
					"port":      "3306",
					"database":  "goal",
					"username":  "root",
					"password":  "123456",
					"charset":   "utf8mb4",
					"collation": "utf8mb4_unicode_ci",
				},
			},
		},
		nil, // 第二个参数是一个 goal 的事件实例，非 goal 环境的情况下，允许为 nil
	)

	// 为 table 包设置数据库工厂
	table.SetFactory(factory)

	assert.True(t, table.ArrayQuery(tableName).Count() == 0)

	user := *table.ArrayQuery(tableName).Create(contracts.Fields{
		"name": "testing",
	})
	assert.NotNil(t, user)
	assert.True(t, user["name"] == "testing")
	assert.True(t, table.ArrayQuery(tableName).Count() == 1)
	table.ArrayQuery(tableName).Where("name", "testing").Delete()
	assert.True(t, table.ArrayQuery(tableName).Count() == 0)

}

func TestMysqlDatabaseFeature(t *testing.T) {
	// 实例化数据库工厂
	factory := database.NewFactory(
		database.Config{
			Default: "mysql",
			Connections: map[string]contracts.Fields{
				"mysql": {
					"driver":    "mysql",
					"host":      "localhost",
					"port":      "3306",
					"database":  "goal",
					"username":  "root",
					"password":  "123456",
					"charset":   "utf8mb4",
					"collation": "utf8mb4_unicode_ci",
				},
			},
		},
		nil, // 第二个参数是一个 goal 的事件实例，非 goal 环境的情况下，允许为 nil
	)

	// 为 table 包设置数据库工厂
	table.SetFactory(factory)

	_, err := table.ArrayQuery(tableName).DeleteE()
	assert.NoError(t, err, err)

	count, err := table.ArrayQuery(tableName).CountE()
	assert.NoError(t, err, err)
	assert.True(t, count == 0)

	user, exception := table.Query[User](tableName).CreateE(contracts.Fields{"name": "testing", "age": 18})
	assert.NotNil(t, user)
	assert.NoError(t, exception, exception)
	assert.True(t, user.Name == "testing")
	assert.True(t, user.Age == 18)

	ageSum, err := table.ArrayQuery(tableName).SumE("age")
	assert.NoError(t, err, err)
	assert.True(t, ageSum == 18)

	ageAvg, err := table.ArrayQuery(tableName).AvgE("age")
	assert.NoError(t, err, err)
	assert.True(t, ageAvg == 18)

	ageMin, err := table.ArrayQuery(tableName).MinE("age")
	assert.NoError(t, err, err)
	assert.True(t, ageMin == 18)

	ageMax, err := table.ArrayQuery(tableName).MaxE("age")
	assert.NoError(t, err, err)
	assert.True(t, ageMax == 18)

	err = table.ArrayQuery(tableName).Chunk(10, func(collection contracts.Collection[*contracts.Fields], page int) contracts.Exception {
		assert.True(t, page == 1)
		assert.True(t, collection.Count() == 1)

		collection.Each(func(i int, fields *contracts.Fields) *contracts.Fields {
			fmt.Println(i, *fields)
			assert.True(t, (*fields)["name"] == "testing")
			return nil
		})

		return nil
	})
	assert.NoError(t, err, err)

	user, exception = table.Query[User](tableName).CreateE(contracts.Fields{"name": "testing2", "age": 18})
	assert.NotNil(t, user)
	assert.NoError(t, exception, exception)

	err = table.ArrayQuery(tableName).ChunkById(1, func(collection contracts.Collection[*contracts.Fields], page int) (any, contracts.Exception) {
		var p, _ = collection.First()
		result := *p
		switch page {
		case 1:
			assert.True(t, result["name"] == "testing")
		case 2:
			assert.True(t, result["name"] == "testing2")
		}

		return result["id"], nil
	})
	assert.NoError(t, err, err)

	err = table.ArrayQuery(tableName).ChunkByIdDesc(1, func(collection contracts.Collection[*contracts.Fields], page int) (any, contracts.Exception) {
		var p, _ = collection.First()
		result := *p
		switch page {
		case 2:
			assert.True(t, result["name"] == "testing")
		case 1:
			assert.True(t, result["name"] == "testing2")
		}

		return result["id"], nil
	})
	assert.NoError(t, err, err)

	assert.NoError(t, table.Query[User](tableName).InsertE(contracts.Fields{"name": "testing3", "age": 18}), err)

	id, err := table.Query[User](tableName).InsertGetIdE(contracts.Fields{"name": "testing4", "age": 18})
	assert.NoError(t, err, err)
	assert.True(t, id == int64(user.Id+2))

	num, err := table.Query[User](tableName).InsertOrIgnoreE(contracts.Fields{"name": "testing5", "age": 18})
	assert.NoError(t, err, err)
	assert.True(t, num > 0)

	num, err = table.Query[User](tableName).InsertOrReplaceE(contracts.Fields{
		"name": "testing6", "age": 18, "id": user.Id,
	})
	assert.NoError(t, err, err)
	assert.True(t, num > 0)
	user = table.Query[User](tableName).Find(user.Id)
	assert.NotNil(t, user)
	assert.True(t, user.Name == "testing6")

	num, err = table.Query[User](tableName).UpdateE(contracts.Fields{"age": 10})
	assert.NoError(t, err, err)
	assert.True(t, num > 0)
	user = table.Query[User](tableName).Find(user.Id)
	assert.NotNil(t, user)
	assert.True(t, user.Age == 10)

	var lastId = user.Id
	num, err = table.Query[User](tableName).Where("id", user.Id).DeleteE()
	assert.NoError(t, err, err)
	assert.True(t, num == 1)
	user = table.Query[User](tableName).Find(user.Id)
	assert.Nil(t, user)

	err = table.Query[User](tableName).UpdateOrInsertE(contracts.Fields{
		"id": lastId,
	}, contracts.Fields{
		"name": "testing6",
		"age":  8,
	})
	assert.NoError(t, err, err)
	user = table.Query[User](tableName).Find(lastId)
	assert.NotNil(t, user)
	assert.True(t, user.Id == lastId)
	assert.True(t, user.Age == 8)
	assert.True(t, user.Name == "testing6")

	user, err = table.Query[User](tableName).UpdateOrCreateE(contracts.Fields{
		"id": lastId,
	}, contracts.Fields{
		"name": "testing6",
		"age":  18,
	})
	assert.NoError(t, err, err)
	assert.NotNil(t, user)
	assert.True(t, user.Id == lastId)
	assert.True(t, user.Age == 18)
	assert.True(t, user.Name == "testing6")

	users, err := table.Query[User](tableName).SelectForUpdateE()
	assert.NoError(t, err, err)
	assert.NotNil(t, users)
	lastUser, _ := users.Last()
	assert.True(t, lastUser.Age == 10)

	assert.True(t, table.Query[User](tableName).Where("id", 0).FirstOr(func() *User {
		return user
	}).Name == "testing6")

	assert.Error(t, utils.NoPanic(func() {
		assert.True(t, table.Query[User](tableName).Find(user.Id).Name == user.Name)
		assert.True(t, table.Query[User](tableName).FirstWhere("id", user.Id).Name == user.Name)
		tmpUser, tmpErr := table.Query[User](tableName).FirstWhereE("id", user.Id)
		assert.NoError(t, tmpErr, tmpErr)
		assert.True(t, tmpUser.Name == user.Name)

		table.Query[User](tableName).Where("id", 0).FirstOrFail()
	}))

	assert.True(t, table.ArrayQuery(tableName).Count() == 5)
	list, total := table.Query[User](tableName).Paginate(2, 1)
	assert.True(t, table.ArrayQuery(tableName).Count() == total)
	assert.True(t, list.Len() == 2)
	assert.True(t, table.Query[User](tableName).SimplePaginate(2, 3).Len() == 1)

	table.ArrayQuery(tableName).Delete()
	assert.True(t, table.ArrayQuery(tableName).Count() == 0)

}
