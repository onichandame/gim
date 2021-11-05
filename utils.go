package gim

import (
	"reflect"

	goutils "github.com/onichandame/go-utils"
)

func newEntity(entOrFunc interface{}) interface{} {
	t := goutils.UnwrapType(reflect.TypeOf(entOrFunc))
	if t.Kind() == reflect.Func {
		t = goutils.UnwrapType(t.Out(0))
	}
	return reflect.New(t).Interface()
}
