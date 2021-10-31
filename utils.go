package gim

import (
	"reflect"

	goutils "github.com/onichandame/go-utils"
)

func newEntity(ent interface{})interface{}{
	t:=goutils.UnwrapType(reflect.TypeOf(ent))
	if t.Kind()==reflect.Func{
		t=goutils.UnwrapType( t.Out(0))
	}
	return reflect.New(t).Interface()
}