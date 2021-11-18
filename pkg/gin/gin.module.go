package gin

import "github.com/onichandame/gim"

var GinModule = gim.Module{
	Name:      "GinModule",
	Providers: []interface{}{newGinService},
	Exports:   []interface{}{newGinService},
}
