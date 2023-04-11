package table

import (
	"github.com/goal-web/contracts"
	class2 "github.com/goal-web/supports/class"
	"github.com/goal-web/supports/utils"
	"reflect"
)

type arrayClass[T contracts.Fields] struct {
	reflect.Type
}

func (class arrayClass[T]) NewByTag(data contracts.Fields, tag string) T {
	return T(data)
}

func (class arrayClass[T]) New(data contracts.Fields) T {
	return T(data)
}

func (class arrayClass[T]) ClassName() string {
	return utils.GetTypeKey(class)
}

func (class arrayClass[T]) GetType() reflect.Type {
	return class.Type
}

func (class arrayClass[T]) IsSubClass(subclass any) bool {
	if value, ok := subclass.(reflect.Type); ok {
		return value.ConvertibleTo(class.Type)
	}

	return reflect.TypeOf(subclass).ConvertibleTo(class.Type)
}

func (class arrayClass[T]) Implements(classType reflect.Type) bool {
	switch value := classType.(type) {
	case *class2.Interface:
		return class.Type.Implements(value.Type)
	case arrayClass[T]:
		return class.Type.Implements(value.Type)
	}

	return class.Type.Implements(classType)
}
