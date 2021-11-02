package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
)

type MainModule struct{}

func (m *MainModule) Imports() []interface{}     { return []interface{}{&MWModule{}} }
func (m *MainModule) Controllers() []interface{} { return []interface{}{&MainController{}} }

type MainController struct{}

func (ctl *MainController) Get(c *gin.Context) interface{} {
	return `hello world`
}

type MWModule struct{}

func (m *MWModule) Server(eng *gin.Engine) *gin.Engine {
	eng.Use(func(c *gin.Context) {
		fmt.Println("before request")
		c.Next()
		fmt.Println("after request")
	})
	return eng
}

func main() {
	app := gim.Bootstrap(&MainModule{})
	app.Server().Run("0.0.0.0:80")
}
