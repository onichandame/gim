package injector

import (
	"fmt"
	"reflect"

	goutils "github.com/onichandame/go-utils"
)

type Container map[reflect.Type]interface{}

func NewContainer() Container {
	c := make(Container)
	return c
}

func (c Container) Bind(entity interface{}) {
	v := goutils.UnwrapValue(reflect.ValueOf(entity))
	t := goutils.UnwrapType(reflect.TypeOf(entity))
	if t.Kind() == reflect.Func {
		inputs := make([]reflect.Value, 0)
		for i := 0; i < t.NumIn(); i++ {
			in := goutils.UnwrapType(t.In(i))
			inputs = append(inputs, reflect.ValueOf(c.Resolve(reflect.New(in).Interface())))
		}
		t = goutils.UnwrapType(t.Out(0))
		v = goutils.UnwrapValue(v.Call(inputs)[0])
	}
	if _, ok := c[t]; ok {
		panic(fmt.Errorf("singleton %v already bound", t.Name()))
	}
	c[t] = (v.Addr().Interface())
}

func (c Container) Resolve(entity interface{}) interface{} {
	return c[goutils.UnwrapType(reflect.TypeOf(entity))]
}

func (c Container) ResolveOrPanic(entity interface{}) interface{} {
	res, ok := c[goutils.UnwrapType(reflect.TypeOf(entity))]
	if !ok {
		panic(fmt.Errorf("singleton %v not found", goutils.UnwrapType(reflect.TypeOf(entity)).Name()))
	}
	return res
}
