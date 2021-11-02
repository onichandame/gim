package graphql

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/onichandame/gim/pkg/injector"
	goutils "github.com/onichandame/go-utils"
)

type GraphqlModule struct{}

func (m *GraphqlModule) Providers() []interface{}   { return []interface{}{newGraphqlService} }
func (m *GraphqlModule) Controllers() []interface{} { return []interface{}{newGraphqlController} }

type ResolverArgs struct {
	Context *gin.Context
	Params  interface{}
}
type Resolver func(ResolverArgs) interface{}
type GraphqlService struct {
	queries    []Resolver
	mutations  []Resolver
	path       string
	outputs    injector.Container
	outputsmap map[interface{}]*graphql.Object
	inputs     injector.Container
	inputsmap  map[interface{}]*graphql.Object
}

func newGraphqlService() *GraphqlService {
	var svc GraphqlService
	svc.outputs = injector.NewContainer()
	svc.outputsmap = make(map[interface{}]*graphql.Object)
	svc.inputs = injector.NewContainer()
	svc.inputsmap = make(map[interface{}]*graphql.Object)
	return &svc
}

func (svc *GraphqlService) SetPath(path string) {
	svc.path = path
}

func (svc *GraphqlService) AddQuery(resolver Resolver) {
	svc.queries = append(svc.queries, resolver)
}
func (svc *GraphqlService) AddMutation(resolver Resolver) {
	svc.queries = append(svc.mutations, resolver)
}

func (svc *GraphqlService) buildSchema() *graphql.Schema {
	buildResolver := func(fn Resolver) (name string, resolver *graphql.Field) {
		t := goutils.UnwrapType(reflect.TypeOf(fn))
		return name, resolver
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{})
	goutils.Assert(err)
	return &schema
}

type GraphqlController struct {
	schema *graphql.Schema
}

func newGraphqlController(svc *GraphqlService) *GraphqlController {
	var ctl GraphqlController
	ctl.schema = svc.buildSchema()
	return &ctl
}

func (ctl *GraphqlController) Get(c *gin.Context) interface{} {
	execQuery := func(query string) *graphql.Result {
		return graphql.Do(graphql.Params{
			Schema:        *ctl.schema,
			RequestString: query,
		})
	}
	query := c.Query("query")
	return execQuery(query)
}
