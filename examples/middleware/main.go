package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
	gimgin "github.com/onichandame/gim/pkg/gin"
)

var MainModule = gim.Module{
	Imports:   []*gim.Module{&gimgin.GinModule, &MWModule},
	Providers: []interface{}{newMainController},
}

type MainController struct{}

func newMainController(ginsvc *gimgin.GinService) *MainController {
	var ctl MainController
	ginsvc.AddRoute(func(rg *gin.RouterGroup) {
		rg.GET("", gimgin.GetHTTPHandler(ctl.Get))
	})
	return &ctl
}

func (ctl *MainController) Get(c *gin.Context) interface{} {
	return `hello world`
}

var MWModule = gim.Module{
	Imports:   []*gim.Module{&gimgin.GinModule},
	Providers: []interface{}{newMWProvider},
}

type MWProvider struct{}

func newMWProvider(ginsvc *gimgin.GinService) *MWProvider {
	var prov MWProvider
	ginsvc.AddMiddleware(func(c *gin.Context) { fmt.Println("before request"); c.Next(); fmt.Println("after request") })
	return &prov
}

func main() {
	MainModule.Bootstrap()
	server := MainModule.Get(&gimgin.GinService{}).(*gimgin.GinService).Bootstrap()
	server.Run("0.0.0.0:80")
}
