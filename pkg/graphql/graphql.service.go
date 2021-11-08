package graphql

import (
	"github.com/graphql-go/graphql"
	goutils "github.com/onichandame/go-utils"
)

type GraphqlService struct {
	queries   map[string]*graphql.Field
	mutations map[string]*graphql.Field
}

func newGraphqlService() *GraphqlService {
	var svc GraphqlService
	svc.queries = make(map[string]*graphql.Field)
	svc.mutations = make(map[string]*graphql.Field)
	return &svc
}

func (svc *GraphqlService) AddQuery(name string, query *graphql.Field) {
	svc.queries[name] = query
}
func (svc *GraphqlService) AddMutation(name string, mutation *graphql.Field) {
	svc.mutations[name] = mutation
}

func (svc *GraphqlService) BuildSchema() *graphql.Schema {
	queries := make(graphql.Fields)
	mutations := make(graphql.Fields)
	for name, query := range svc.queries {
		queries[name] = query
	}
	for name, mutation := range svc.mutations {
		mutations[name] = mutation
	}
	var schemaConf graphql.SchemaConfig
	if len(queries) > 0 {
		schemaConf.Query = graphql.NewObject(graphql.ObjectConfig{
			Name:   `Query`,
			Fields: queries,
		})
	}
	if len(mutations) > 0 {
		schemaConf.Mutation = graphql.NewObject(graphql.ObjectConfig{
			Name:   `Mutation`,
			Fields: mutations,
		})
	}
	schema, err := graphql.NewSchema(schemaConf)
	goutils.Assert(err)
	return &schema
}
