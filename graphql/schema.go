package graphql

import "github.com/graphql-go/graphql"

func GetSchemaFromStruct(instance interface{}) *graphql.Schema {
	objectsSchema := GetObjectFromStruct(instance)
	resolver := GetQueryResolverFromStruct(instance, objectsSchema)
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: &resolver.Query,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: &resolver.Query,
		}),
	})
	if err != nil {
		panic(err)
	}
	return &schema
}
