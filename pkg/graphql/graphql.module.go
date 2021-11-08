package graphql

import "github.com/onichandame/gim"

var GraphqlModule = gim.Module{
	Providers: []interface{}{newGraphqlService},
	Exports:   []interface{}{newGraphqlService},
}
