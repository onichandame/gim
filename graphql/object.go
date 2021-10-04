package graphql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
	"github.com/graphql-go/graphql"
)

func GetObjectFromStruct(instance interface{}) (obj *graphql.Object) {
	var t reflect.Type
	if t = reflect.TypeOf(instance); t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	var fields graphql.Fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		tags := getTags(&field)
		if customName, ok := tags["name"]; ok {
			name = customName
		}
		// TODO: populate field type
		fields[name] = &graphql.Field{
			Name: name,
			Type: graphql.String,
		}
	}
	obj = graphql.NewObject(graphql.ObjectConfig{
		Name:   t.Name(),
		Fields: &fields,
	})
	return obj
}

func getTags(field *reflect.StructField) (res map[string]string) {
	res = make(map[string]string)
	if tags, err := structtag.Parse(string(field.Tag)); err != nil {
		panic(err)
	} else {
		if tag, err := tags.Get("gimgraphql"); err != nil {
			panic(err)
		} else {
			rawTags := []string{}
			rawTags[0] = tag.Name
			rawTags = append(rawTags, tag.Options...)
			for _, rawTag := range rawTags {
				tuple := strings.Split(rawTag, "=")
				if len(tuple) > 2 {
					panic(fmt.Errorf("tag value cannot include '=': %s", rawTag))
				}
				key := strings.TrimSpace(tuple[0])
				var value string
				if len(tuple) > 1 {
					value = tuple[1]
				}
				res[key] = value
			}
		}
	}
	return res
}
