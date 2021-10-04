package graphql

import "github.com/graphql-go/graphql"

func GetSchemaFromStruct(instance interface{}) *graphql.Schema {
	rootQueryObject := GetObjectFromStruct(&struct{}{})
	rootMutationObject := GetObjectFromStruct(&struct{}{})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQueryObject,
		Mutation: rootMutationObject,
	})
	if err != nil {
		panic(err)
	}
	return &schema
}
