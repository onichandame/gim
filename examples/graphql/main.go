package main

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/onichandame/gim"
	gimgraphql "github.com/onichandame/gim/pkg/graphql"
	goutils "github.com/onichandame/go-utils"
)

var MainModule = gim.Module{
	Imports:   []*gim.Module{&gimgraphql.GraphqlModule},
	Providers: []interface{}{newMainResolver, newMainService},
}

type MainResolver struct{}

func newMainResolver(gsvc *gimgraphql.GraphqlService, svc *MainService) *MainResolver {
	var r MainResolver
	gsvc.AddQuery("greet", &graphql.Field{Type: graphql.NewNonNull(graphql.String), Resolve: func(p graphql.ResolveParams) (interface{}, error) { return "hello world", nil }})
	type EnlistResult struct {
		Name         string `json:"name"`
		IsSuccessful bool   `json:"isSuccessful"`
	}
	gsvc.AddMutation("enlist", &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "EnlistResult",
			Fields: graphql.Fields{
				"name":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
				"isSuccessful": &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
			},
		}),
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		},
		Resolve: func(p graphql.ResolveParams) (res interface{}, err error) {
			defer goutils.RecoverToErr(&err)
			name := p.Args["name"].(string)
			if _, ok := svc.list[name]; ok {
				panic(fmt.Errorf("%v already enlisted", name))
			}
			svc.list[name] = nil
			res = EnlistResult{Name: name, IsSuccessful: true}
			return res, err
		},
	},
	)
	return &r
}

type MainService struct {
	list map[string]interface{}
}

func newMainService() *MainService {
	var svc MainService
	svc.list = make(map[string]interface{})
	return &svc
}

func main() {
	MainModule.Bootstrap()
	schema := MainModule.Get(&gimgraphql.GraphqlService{}).(*gimgraphql.GraphqlService).BuildSchema()
	result := graphql.Do(graphql.Params{
		Schema:        *schema,
		RequestString: `{greet}`,
	})
	fmt.Println(result.Data.(map[string]interface{})["greet"])
	result = graphql.Do(graphql.Params{
		Schema:         *schema,
		RequestString:  `mutation($name:String!){enlist(name:$name){name isSuccessful}}`,
		VariableValues: map[string]interface{}{"name": "Tim"},
	})
	fmt.Println(result.Data.(map[string]interface{})["enlist"])
	result = graphql.Do(graphql.Params{
		Schema:         *schema,
		RequestString:  `mutation($name:String!){enlist(name:$name){name isSuccessful}}`,
		VariableValues: map[string]interface{}{"name": "Tim"},
	})
	fmt.Println(result.Errors)
}
