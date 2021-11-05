package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
)

type MainModule struct{}

func (*MainModule) Imports() []interface{}     { return []interface{}{&SubModule{}} }
func (*MainModule) Controllers() []interface{} { return []interface{}{newMainController} }

type MainController struct{ svc *SubService }

func newMainController(svc *SubService) *MainController {
	var c MainController
	c.svc = svc
	return &c
}
func (c *MainController) Get(*gin.Context) interface{} {
	return c.svc.getGreeting()
}

func main() {
	app := gim.Bootstrap(&MainModule{})
	app.Server().Run("0.0.0.0:80")
}
