package graphql

import (
	"errors"
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

type ResolverFuncArgs struct {
	Ctx  *gin.Context
	Args interface{}
}

type ResolverFunc func(ResolverFuncArgs) interface{}

type Resolvers struct {
	Query    graphql.Fields
	Mutation graphql.Fields
}

func GetQueryResolverFromStruct(instance interface{}, objType graphql.Object) (resolvers *Resolvers) {
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
			resolvers.Query[name] = &graphql.Field{
				Type: &objType,
				Args: graphql.FieldConfigArgument{},
			}
		} else if strings.HasPrefix(m.Name, MUTATION_PREFIX) {
			validateFunc(m.Func.Interface())
			validateFunc(m)
			name := strings.Replace(m.Name, MUTATION_PREFIX, "", 1)
			resolvers.Mutation[name] = &graphql.Field{}
		}
	}
	return resolvers
}

func validateFunc(fn interface{}) {
	if _, ok := fn.(ResolverFunc); !ok {
		panic(fmt.Errorf("resolver function must be of ResolverFunc type, %s is not", reflect.TypeOf(fn).Name()))
	}
}

func getResolverFunc(m *reflect.Method)(*graphql.FieldConfigArgument,graphql.FieldResolveFn) {
	validateFunc(m.Func.Interface())
	args:=m.Type.In(0)
}
