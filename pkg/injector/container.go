package injector

import (
	"fmt"
	"reflect"

	goutils "github.com/onichandame/go-utils"
)

// Container is a simple DI container based on types where each type can have
// only one concret
type Container map[reflect.Type]interface{}

// NewContainer returns an initiated container
func NewContainer() Container {
	c := make(Container)
	return c
}

// Bind adds a concret to the container. the argument can be an instance or
// a constructor taking other concrets in the container as parameters.
// example:
//  container.Bind(`v0.1.0`)
//  container.Bind(func(version string) bool{ return version==`v0.1.0` })
func (c Container) Bind(entity interface{}) {
	v := goutils.UnwrapValue(reflect.ValueOf(entity))
	t := goutils.UnwrapType(reflect.TypeOf(entity))
	if _, ok := c[t]; ok {
		panic(fmt.Errorf("singleton %v already bound", t.Name()))
	}
	if t.Kind() == reflect.Func {
		inputs := make([]reflect.Value, 0)
		for i := 0; i < t.NumIn(); i++ {
			in := goutils.UnwrapType(t.In(i))
			inputs = append(inputs, reflect.ValueOf(c[in]))
		}
		t = goutils.UnwrapType(t.Out(0))
		v = goutils.UnwrapValue(v.Call(inputs)[0])
	}
	if v.Kind() == reflect.Struct {
		v = v.Addr()
	}
	c[t] = (v.Interface())
}

// Resolve returns the bound concret of the type of the argument. If not
// found, nil is returned
func (c Container) Resolve(entity interface{}) interface{} {
	return c[goutils.UnwrapType(reflect.TypeOf(entity))]
}

// ResolveOrPanic returns the resolved concret if found, panics if not found
func (c Container) ResolveOrPanic(entity interface{}) interface{} {
	res, ok := c[goutils.UnwrapType(reflect.TypeOf(entity))]
	if !ok {
		panic(fmt.Errorf("singleton %v not found", goutils.UnwrapType(reflect.TypeOf(entity)).Name()))
	}
	return res
}
