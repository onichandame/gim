package graphql

import "github.com/onichandame/gim"

var GraphqlModule = gim.Module{
	Name:      "GraphqlModule",
	Providers: []interface{}{newGraphqlService},
	Exports:   []interface{}{newGraphqlService},
}
