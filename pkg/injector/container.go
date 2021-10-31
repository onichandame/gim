package injector

import (
	"fmt"
	"reflect"

	goutils "github.com/onichandame/go-utils"
)

type Container struct {
	privateEntities map[reflect.Type]interface{}
	publicEntities  map[reflect.Type]interface{}
	children        map[*Container]interface{}
	parent          *Container
}

func NewContainer() *Container {
	c := new(Container)
	c.privateEntities = make(map[reflect.Type]interface{})
	c.publicEntities = make(map[reflect.Type]interface{})
	c.children = make(map[*Container]interface{})
	return c
}

func (c *Container) SetParent(parent *Container) {
	if c.parent != nil {
		delete(c.parent.children, c)
	}
	c.parent = parent
	c.parent.children[c] = nil
}

func (c *Container) Bind(entity interface{}, public bool) {
	v := goutils.UnwrapValue(reflect.ValueOf(entity))
	t := goutils.UnwrapType(reflect.TypeOf(entity))
	if t.Kind() == reflect.Func {
		inputs := make([]reflect.Value, 0)
		for i := 0; i < t.NumIn(); i++ {
			in := goutils.UnwrapType(t.In(i))
			inputs = append(inputs, reflect.ValueOf(c.resolve(reflect.New(in).Interface(), include_all)))
		}
		t = goutils.UnwrapType(t.Out(0))
		entity = v.Call(inputs)[0].Interface()
	}
	if public {
		if _, ok := c.publicEntities[t]; !ok {
			c.publicEntities[t] = entity
		}
	} else {
		if _, ok := c.privateEntities[t]; !ok {
			c.privateEntities[t] = entity
		}
	}
}

const (
	include_private  = 0b00000001
	include_public   = 0b00000010
	include_children = 0b00000100
	include_root     = 0b00001000
	include_all      = include_private ^ include_public ^ include_children ^ include_root
)

func (c *Container) GetRoot() *Container {
	if c.parent == nil {
		return c
	} else {
		return c.parent.GetRoot()
	}
}

func (c *Container) resolve(entity interface{}, scope int) interface{} {
	t := goutils.UnwrapType(reflect.TypeOf(entity))
	if scope&include_private > 0 {
		if ent, ok := c.privateEntities[t]; ok {
			return ent
		}
	}
	if scope&include_public > 0 {
		if ent, ok := c.publicEntities[t]; ok {
			return ent
		}
	}
	if scope&include_children > 0 {
		for child := range c.children {
			ent := child.resolve(entity, include_public)
			if ent != nil {
				return ent
			}
		}
	}
	if scope&include_root > 0 {
		root := c.GetRoot()
		ent := root.resolve(entity, include_public)
		if ent != nil {
			return ent
		}
	}
	return nil
}

func (c *Container) Resolve(entity interface{}) {
	ent := c.resolve(entity, include_all)
	if ent == nil {
		panic(fmt.Errorf("failed to resolve entity %v", goutils.UnwrapType(reflect.TypeOf(entity)).Name()))
	}
	v := reflect.ValueOf(entity)
	v.Elem().Set(reflect.ValueOf(ent).Elem())
}
