package graphql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

const (
	QUERY_PREFIX    = "Query"
	MUTATION_PREFIX = "Mutation"
)

type ResolverFuncParam struct {
	Ctx  *gin.Context
	Args interface{}
}

type ResolverFunc func(*ResolverFuncParam) interface{}

type Resolvers struct {
	Query    graphql.Fields
	Mutation graphql.Fields
}

func GetQueryResolverFromStruct(instance interface{}, objType *graphql.Object) (resolvers *Resolvers) {
	resolvers = new(Resolvers)
	resolvers.Query = make(graphql.Fields)
	resolvers.Mutation = make(graphql.Fields)
	t := reflect.TypeOf(instance)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if strings.HasPrefix(m.Name, QUERY_PREFIX) {
			validateFunc(m.Func.Interface())
			name := strings.Replace(m.Name, QUERY_PREFIX, "", 1)
			args, fn := getResolverFunc(&m)
			resolvers.Query[name] = &graphql.Field{
				Type:    objType,
				Args:    *args,
				Resolve: fn,
			}
		} else if strings.HasPrefix(m.Name, MUTATION_PREFIX) {
			validateFunc(m.Func.Interface())
			validateFunc(m)
			name := strings.Replace(m.Name, MUTATION_PREFIX, "", 1)
			args, fn := getResolverFunc(&m)
			resolvers.Mutation[name] = &graphql.Field{
				Type:    objType,
				Args:    *args,
				Resolve: fn,
			}
		}
	}
	return resolvers
}

func validateFunc(fn interface{}) {
	if _, ok := fn.(ResolverFunc); !ok {
		panic(fmt.Errorf("resolver function must be of ResolverFunc type, %s is not", reflect.TypeOf(fn).Name()))
	}
}

func getResolverFunc(m *reflect.Method) (args *graphql.FieldConfigArgument, fn graphql.FieldResolveFn) {
	args = &graphql.FieldConfigArgument{}
	validateFunc(m.Func.Interface())
	paramType := m.Type.In(0)
	if paramType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("argument of a resolver function should be pointer to ResolverFuncParam"))
	}
	return args, fn
}
