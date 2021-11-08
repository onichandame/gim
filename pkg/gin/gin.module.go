package gin

import "github.com/onichandame/gim"

var GinModule = gim.Module{Providers: []interface{}{newGinService}, Exports: []interface{}{newGinService}}
