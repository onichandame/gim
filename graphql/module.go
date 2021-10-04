package graphql

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/onichandame/gim/core"
)

func NewGraphqlModule(instance interface{}) (mod *core.Module) {
	mod = &core.Module{}
	schema := GetSchemaFromStruct(instance)
	execQuery := func(query string) (res *graphql.Result) {
		res = graphql.Do(graphql.Params{
			Schema:        *schema,
			RequestString: query,
		})
		if len(res.Errors) > 0 {
			panic(core.NewGimError(400, res))
		}
		return res
	}
	mod.Routes = []*core.Route{
		{
			Get: func(c *gin.Context) interface{} {
				return execQuery(c.Request.URL.Query().Get("query"))
			},
			Post: func(c *gin.Context) interface{} {
				if bodyStr, err := io.ReadAll(c.Request.Body); err != nil {
					panic(err)
				} else {
					return execQuery(string(bodyStr))
				}
			},
		},
	}
	return mod
}
