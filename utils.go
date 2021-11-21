package gim

import (
	"reflect"

	"github.com/onichandame/gim/pkg/injector"
	goutils "github.com/onichandame/go-utils"
)

func getType(entOrFunc interface{}) reflect.Type {
	t := goutils.UnwrapType(reflect.TypeOf(entOrFunc))
	switch t.Kind() {
	case reflect.Func:
		t = goutils.UnwrapType(t.Out(0))
	}
	return t
}

func getTypeAndSingleton(container injector.Container, ent interface{}) (reflect.Type, interface{}) {
	t := getType(ent)
	return t, container[t]
}
