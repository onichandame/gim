package graphql

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/fatih/structtag"
	"github.com/graphql-go/graphql"
	goutils "github.com/onichandame/go-utils"
)

// query function: func(ctx QueryContext, params <args type>)<response type>
type QueryContext *graphql.ResolveParams

// mutation function: func(ctx MutationContext, params <args type>)<response type>
type MutationContext *graphql.ResolveParams

// subscription function: func(ctx SubscriptionContext, params <args type>)<response type, channel>
type SubscriptionContext *graphql.ResolveParams

type GraphqlService struct {
	enums     map[reflect.Type]*graphql.Enum
	scalars   map[reflect.Type]*graphql.Scalar
	resolvers []interface{}
	path      string
}

var svcsvc *GraphqlService

func newGraphqlService() *GraphqlService {
	var svc GraphqlService
	fmt.Println("svc")
	svc.enums = make(map[reflect.Type]*graphql.Enum)
	svc.scalars = make(map[reflect.Type]*graphql.Scalar)
	svc.resolvers = make([]interface{}, 0)
	svcsvc = &svc
	return &svc
}

func (svc *GraphqlService) SetPath(path string) { svc.path = path }

func (svc *GraphqlService) AddResolver(resolver interface{}) {
	fmt.Println(svcsvc == svc)
	fmt.Println("add")
	svc.resolvers = append(svc.resolvers, resolver)
}

func (svc *GraphqlService) AddEnum(ent interface{}, enum *graphql.Enum) {
	svc.enums[goutils.UnwrapType(reflect.TypeOf(ent))] = enum
}
func (svc *GraphqlService) AddScalar(ent interface{}, scalar *graphql.Scalar) {
	svc.scalars[goutils.UnwrapType(reflect.TypeOf(ent))] = scalar
}

func (svc *GraphqlService) buildSchema() *graphql.Schema {
	queries := make(graphql.Fields)
	mutations := make(graphql.Fields)
	objects := make(map[reflect.Type]graphql.Type)
	inputs := make(map[reflect.Type]*graphql.ArgumentConfig)
	ints := []interface{}{int(0), int8(0), int16(0), int32(0), int64(0), uint(0), uint8(0), uint16(0), uint32(0), uint64(0)}
	floats := []interface{}{float32(0), float64(0)}
	strings := []interface{}{string(``), []byte(``)}
	for _, i := range ints {
		objects[reflect.TypeOf(i)] = graphql.Int
	}
	for _, f := range floats {
		objects[reflect.TypeOf(f)] = graphql.Float
	}
	for _, s := range strings {
		objects[reflect.TypeOf(s)] = graphql.String
	}
	objects[reflect.TypeOf(true)] = graphql.Boolean
	objects[reflect.TypeOf(time.Time{})] = graphql.DateTime
	for t, scalar := range svc.scalars {
		objects[t] = scalar
	}
	for t, enum := range svc.enums {
		objects[t] = enum
	}
	loadResolver := func(resolver interface{}) {
		v := goutils.UnwrapValue(reflect.ValueOf(resolver))
		t := goutils.UnwrapType(reflect.TypeOf(resolver))
		pt := reflect.TypeOf(reflect.New(t))
		pv := v.Addr()
		var loadObjectType func(reflect.Type)
		loadObjectType = func(t reflect.Type) {
			t = goutils.UnwrapType(t) // object type
			if _, ok := objects[t]; ok {
				return
			}
			name := t.Name()
			if n, ok := reflect.New(t).Interface().(interface{ Name() string }); ok {
				name = n.Name()
			}
			fields := make(graphql.Fields)
			getFieldType := func(f *reflect.StructField) graphql.Type {
				ft := goutils.UnwrapType(f.Type)
				tags, err := structtag.Parse(string(f.Tag))
				goutils.Assert(err)
				tag, _ := tags.Get(TAG)
				res, ok := objects[ft]
				if !ok {
					loadObjectType(ft)
					res = objects[ft]
				}
				if !tag.HasOption(TAG_NULLABLE) {
					res = graphql.NewNonNull(res)
				}
				return res
			}
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i) // field
				name := f.Name
				tags, err := structtag.Parse(string(f.Tag))
				goutils.Assert(err)
				tag, _ := tags.Get(TAG)
				if tag != nil {
					name = tag.Name
				}
				fields[name] = &graphql.Field{
					Type: getFieldType(&f),
				}
			}
			objects[t] = graphql.NewObject(graphql.ObjectConfig{
				Name:   name,
				Fields: fields,
			})
		}
		loadInputType := func(t reflect.Type) {} //TODO
		isQueryFunc := func(t reflect.Type) bool {
			ctxType := t.In(0)
			return ctxType == reflect.TypeOf(new(QueryContext)).Elem() && t.NumOut() == 1
		}
		isMutationFunc := func(t reflect.Type) bool {
			ctxType := t.In(0)
			return ctxType == reflect.TypeOf(new(MutationContext)).Elem() && t.NumOut() == 1
		}
		isSubscriptionFunc := func(t reflect.Type) bool {
			ctxType := t.In(0)
			return ctxType == reflect.TypeOf(new(SubscriptionContext)).Elem() && t.NumOut() == 2 && t.Out(1).Kind() == reflect.Chan
		}
		isResolverFunc := func(t reflect.Type) bool {
			if t.NumIn() < 1 || t.NumIn() > 2 {
				return false
			}
			return isQueryFunc(t) || isMutationFunc(t) || isSubscriptionFunc(t)
		}
		loadResolverFunc := func(name string, fn interface{}) {
			t := goutils.UnwrapType(reflect.TypeOf(fn))
			v := goutils.UnwrapValue(reflect.ValueOf(fn))
			list := &queries
			if isQueryFunc(t) {
				list = &queries
			} else if isMutationFunc(t) {
				list = &mutations
			} else {
				panic(fmt.Errorf("only query and mutation are supported"))
			}
			loadObjectType(t.Out(0))
			args := make(graphql.FieldConfigArgument)
			if t.NumIn() > 1 {
				it := goutils.UnwrapType(t.In(1))
				if it != nil {
					for i := 0; i < it.NumField(); i++ {
						f := it.Field(i)
						loadInputType(f.Type)
						args[f.Name] = inputs[goutils.UnwrapType(f.Type)]
					}
				}
			}
			(*list)[name] = &graphql.Field{
				Type: objects[t.Out(0)],
				Args: args,
			}
			if isQueryFunc(t) || isMutationFunc(t) {
				(*list)[name].Resolve = func(p graphql.ResolveParams) (res interface{}, err error) {
					defer goutils.RecoverToErr(&err)
					args := reflect.New(goutils.UnwrapType(t.In(1)))
					argsMarshalled, err := json.Marshal(p.Args)
					goutils.Assert(err)
					json.Unmarshal(argsMarshalled, args.Interface())
					res = v.Call([]reflect.Value{reflect.ValueOf(p), args})[0].Interface()
					return res, err
				}
			}
		}
		for i := 0; i < pt.NumMethod(); i++ {
			mt := pt.Method(i)
			mv := pv.Method(i)
			if isResolverFunc(mt.Type) {
				loadResolverFunc(mt.Name, mv.Interface())
			}
		}
		for i := 0; i < t.NumMethod(); i++ {
			mt := t.Method(i)
			mv := v.Method(i)
			if isResolverFunc(mt.Type) {
				loadResolverFunc(mt.Name, mv.Interface())
			}
		}
	}
	for _, resolver := range svc.resolvers {
		loadResolver(resolver)
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   `Query`,
			Fields: queries,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   `Mutation`,
			Fields: mutations,
		}),
	})
	goutils.Assert(err)
	return &schema
}
