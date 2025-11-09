// Package database provides database functionality and singleton access.
// 包 database 提供数据库功能和单例访问。
package database

import (
	"sync"

	"github.com/goal-web/application"
	"github.com/goal-web/contracts"
)

// singleton holds the singleton instance of the database factory.
// singleton 保存数据库工厂的单例实例。
var singleton contracts.DBFactory

// once ensures the singleton is initialized only once.
// once 确保单例只被初始化一次。
var once sync.Once

// Default returns the singleton instance of the database factory.
// Default 返回数据库工厂的单例实例。
func Default() contracts.DBFactory {
	once.Do(func() {
		singleton = application.Get("db.factory").(contracts.DBFactory)
	})

	return singleton
}

// Connection returns a database connection instance with the given name.
// Connection 返回指定名称的数据库连接实例。
func Connection(key ...string) contracts.DBConnection {
	return Default().Connection(key...)
}

// Extend registers a custom database connector.
// Extend 注册自定义数据库连接器。
func Extend(name string, driver contracts.DBConnector) {
	Default().Extend(name, driver)
}

// Query executes a query and returns a collection of results.
// Query 执行查询并返回结果集合。
func Query(query string, args ...any) (contracts.Collection[contracts.Fields], contracts.Exception) {
	return Connection().Query(query, args...)
}

// GetQuery executes a query to get a single result.
// GetQuery 执行查询以获取单个结果。
func GetQuery(dest any, query string, args ...any) contracts.Exception {
	return Connection().Get(dest, query, args...)
}

// Select executes a select query and fills the destination with results.
// Select 执行选择查询并将结果填充到目标中。
func Select(dest any, query string, args ...any) contracts.Exception {
	return Connection().Select(dest, query, args...)
}

// Exec executes a query without returning results.
// Exec 执行不返回结果的查询。
func Exec(query string, args ...any) (contracts.Result, contracts.Exception) {
	return Connection().Exec(query, args...)
}

// Transaction executes a callback within a database transaction.
// Transaction 在数据库事务中执行回调。
func Transaction(callback func(executor contracts.SqlExecutor) contracts.Exception) contracts.Exception {
	return Connection().Transaction(callback)
}

// DriverName returns the name of the database driver.
// DriverName 返回数据库驱动的名称。
func DriverName() string {
	return Connection().DriverName()
}

// Begin starts a new database transaction.
// Begin 开始新的数据库事务。
func Begin() (contracts.DBTx, contracts.Exception) {
	return Connection().Begin()
}
