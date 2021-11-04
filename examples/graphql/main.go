package main

import (
	"github.com/onichandame/gim"
	"github.com/onichandame/gim/pkg/graphql"
)

type MainModule struct{}

func (m *MainModule) Imports() []interface{}   { return []interface{}{&graphql.GraphqlModule{}} }
func (m *MainModule) Providers() []interface{} { return []interface{}{newMainResolver} }

type MainResolver struct{}

func newMainResolver(gsvc *graphql.GraphqlService) *MainResolver {
	var r MainResolver
	gsvc.AddResolver(&r)
	return &r
}
func (r *MainResolver) greet(ctx graphql.QueryContext) interface{} {
	return "hello world"
}

func main() {
	app := gim.Bootstrap(&MainModule{})
	app.Server().Run("0.0.0.0:80")
}
