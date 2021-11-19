package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
	gimgin "github.com/onichandame/gim/pkg/gin"
)

var MainModule = gim.Module{
	Imports:   []*gim.Module{&SubModule, &gimgin.GinModule, &CommonModule},
	Providers: []interface{}{newMainController},
}

type MainController struct{ svc *SubService }

func newMainController(svc *SubService, ginsvc *gimgin.GinService, cmsvc *CommonService) *MainController {
	var c MainController
	c.svc = svc
	ginsvc.AddRoute(func(rg *gin.RouterGroup) {
		rg.GET("", gimgin.GetHTTPHandler(c.Get))
		rg.GET("tick", gimgin.GetHTTPHandler(func(c *gin.Context) interface{} { return cmsvc.GetNextTick() }))
	})
	return &c
}
func (c *MainController) Get(*gin.Context) interface{} {
	return c.svc.getGreeting()
}

func main() {
	MainModule.Bootstrap()
	server := MainModule.Get(&gimgin.GinService{}).(*gimgin.GinService).Bootstrap()
	server.Run("0.0.0.0:80")
}
