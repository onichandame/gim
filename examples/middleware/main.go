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

func (m *MWModule) Middlewares() []interface{} { return []interface{}{&LoggerMW{}} }

type LoggerMW struct{}

func (l *LoggerMW) Use() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("before request")
		c.Next()
		c.Status(500)
		fmt.Println("after request")
	}
}

func main() {
	root := gim.Bootstrap(&MainModule{})
	var eng gin.Engine
	root.Resolve(&eng)
	eng.Run("0.0.0.0:80")
}
