package table

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"reflect"
	"strings"
)

type Model[T any] struct {
	class           contracts.Class[T]
	table           string
	connection      string
	primaryKeyField string
	data            *T
	value           reflect.Value
}

func NewModel[T any](class contracts.Class[T], table string, connection ...string) Model[T] {
	conn := ""
	if len(connection) > 0 {
		conn = connection[0]
	}
	return Model[T]{class: class, table: table, connection: conn}
}

func (model Model[T]) InitModel(class contracts.Class[T], table, primaryKeyField string, data *T, value reflect.Value) {
	model.class = class
	model.table = table
	model.primaryKeyField = primaryKeyField
	model.data = data
	model.value = value
}

func (model Model[T]) Exists() bool {
	return Query[T](model.table).Where(model.primaryKeyField, model.GetPrimaryKey()).Count() > 0
}

func (model Model[T]) Update(fields contracts.Fields) contracts.Exception {
	_, err := Query[T](model.table).Where(model.primaryKeyField, model.GetPrimaryKey()).UpdateE(fields)

	if err != nil {
		return err
	}

	data := model.class.NewByTag(fields, "db")

	utils.EachStructField(reflect.ValueOf(data), data, func(field reflect.StructField, value reflect.Value) {
		if field.IsExported() {
			tag := field.Tag.Get("json")
			tag = strings.Split(tag, ",")[0]
			if v, exists := fields[tag]; exists {
				model.value.Elem().FieldByName(field.Name).Set(reflect.ValueOf(v))
			}
		}
	})

	return err
}

func (model Model[T]) Save() contracts.Exception {
	var fields = contracts.Fields{}

	utils.EachStructField(model.value.Elem(), *model.data, func(field reflect.StructField, value reflect.Value) {
		if field.IsExported() {
			tag := field.Tag.Get("json")
			tag = strings.Split(tag, ",")[0]
			if tag != "" && tag != "-" && tag != model.primaryKeyField {
				fields[tag] = value.Interface()
			}
		}
	})

	_, err := Query[T](model.table).Where(model.primaryKeyField, model.GetPrimaryKey()).UpdateE(fields)
	return err
}

func (model Model[T]) Refresh() contracts.Exception {
	data, err := Query[T](model.table).Where(model.primaryKeyField, model.GetPrimaryKey()).FirstE()
	if err != nil {
		return err
	}

	dataValue := reflect.ValueOf(*data)
	utils.EachStructField(dataValue, *data, func(field reflect.StructField, value reflect.Value) {
		if _, isModel := value.Interface().(Model[T]); !isModel {
			model.value.Elem().FieldByName(field.Name).Set(value)
		}
	})
	data = nil

	return err
}

func (model Model[T]) Delete() contracts.Exception {
	_, err := Query[T](model.table).Where(model.primaryKeyField, model.GetPrimaryKey()).DeleteE()
	return err
}

func (model Model[T]) GetClass() contracts.Class[T] {
	return model.class
}

func (model Model[T]) GetPrimaryKey() any {
	return model.value.Elem().FieldByName(capitalize(model.primaryKeyField)).Interface()
}

func (model Model[T]) GetTable() string {
	return model.table
}

func (model Model[T]) GetConnection() string {
	return model.connection
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-'a'+'A') + s[1:]
}
