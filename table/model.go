package table

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"reflect"
	"strings"
)

type Model[T any] struct {
	Class      contracts.Class[T]
	Table      string
	connection      string
	PrimaryKeyField string
	Data  *T
	Value reflect.Value
}

func (model Model[T]) InitModel(class contracts.Class[T], table, PrimaryKeyField string, data *T, value reflect.Value) {
	model.Class = class
	model.Table = table
	model.PrimaryKeyField = PrimaryKeyField
	model.Data = data
	model.Value = value
}

func (model Model[T]) Exists() bool {
	return Query[T](model.Table).Where(model.PrimaryKeyField, model.GetPrimaryKey()).Count() > 0
}

func (model Model[T]) Update(fields contracts.Fields) contracts.Exception {
	_, err := Query[T](model.Table).Where(model.PrimaryKeyField, model.GetPrimaryKey()).UpdateE(fields)

	if err != nil {
		return err
	}

	data := model.Class.NewByTag(fields, "db")

	utils.EachStructField(reflect.ValueOf(data), data, func(field reflect.StructField, value reflect.Value) {
		if field.IsExported() {
			tag := field.Tag.Get("json")
			tag = strings.Split(tag, ",")[0]
			if v, exists := fields[tag]; exists {
				model.Value.Elem().FieldByName(field.Name).Set(reflect.ValueOf(v))
			}
		}
	})

	return err
}

func (model Model[T]) Save() contracts.Exception {
	var fields = contracts.Fields{}

	utils.EachStructField(model.Value.Elem(), *model.Data, func(field reflect.StructField, value reflect.Value) {
		if field.IsExported() {
			tag := field.Tag.Get("json")
			tag = strings.Split(tag, ",")[0]
			if tag != "" && tag != "-" && tag != model.PrimaryKeyField {
				fields[tag] = value.Interface()
			}
		}
	})

	_, err := Query[T](model.Table).Where(model.PrimaryKeyField, model.GetPrimaryKey()).UpdateE(fields)
	return err
}

func (model Model[T]) Refresh() contracts.Exception {
	data, err := Query[T](model.Table).Where(model.PrimaryKeyField, model.GetPrimaryKey()).FirstE()
	if err != nil {
		return err
	}

	dataValue := reflect.ValueOf(*data)
	utils.EachStructField(dataValue, *data, func(field reflect.StructField, value reflect.Value) {
		if _, isModel := value.Interface().(Model[T]); !isModel {
			model.Value.Elem().FieldByName(field.Name).Set(value)
		}
	})
	data = nil

	return err
}

func (model Model[T]) Delete() contracts.Exception {
	_, err := Query[T](model.Table).Where(model.PrimaryKeyField, model.GetPrimaryKey()).DeleteE()
	return err
}

func (model Model[T]) GetClass() contracts.Class[T] {
	return model.Class
}

func (model Model[T]) GetPrimaryKey() any {
	return model.Value.Elem().FieldByName(capitalize(model.PrimaryKeyField)).Interface()
}

func (model Model[T]) GetTable() string {
	return model.Table
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
