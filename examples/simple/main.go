package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
	gimgin "github.com/onichandame/gim/pkg/gin"
)

var MainModule = gim.Module{Imports: []*gim.Module{&gimgin.GinModule}, Providers: []interface{}{newController}}

type MainController struct{}

func newController(ginsvc *gimgin.GinService) *MainController {
	var ctlr MainController
	ginsvc.AddRoute(func(rg *gin.RouterGroup) {
		rg.GET("", gimgin.GetHandler(func(c *gin.Context) interface{} {
			return `hello world`
		}))
	})
	return &ctlr
}

func main() {
	MainModule.Bootstrap()
	ginsvc := MainModule.Get(&gimgin.GinService{}).(*gimgin.GinService)
	ginsvc.Bootstrap().Run("0.0.0.0:80")
}
