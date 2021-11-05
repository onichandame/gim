package graphql

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

type graphqlController struct {
	schema *graphql.Schema
}

func newGraphqlBuilderController(svc *GraphqlService) *graphqlController {
	var ctl graphqlController
	ctl.schema = svc.buildSchema()
	return &ctl
}

func (ctl *graphqlController) Get(c *gin.Context) interface{} {
	execQuery := func(query string) *graphql.Result {
		return graphql.Do(graphql.Params{
			Schema:        *ctl.schema,
			RequestString: query,
			Context:       context.WithValue(context.Background(), contextToken, c),
		})
	}
	query := c.Query("query")
	return execQuery(query)
}
