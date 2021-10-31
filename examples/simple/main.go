package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
)

type MainModule struct{}

func (mod *MainModule) Controllers() []interface{} {
	return []interface{}{newController}
}

type MainController struct{}

func newController() *MainController {
	var ctlr MainController
	return &ctlr
}
func (ctlr *MainController) Get(c *gin.Context) interface{} {
	return `hello world`
}

func main() {
	root := gim.Bootstrap(&MainModule{})
	var eng gin.Engine
	root.Resolve(&eng)
	eng.Run("0.0.0.0:80")
}
