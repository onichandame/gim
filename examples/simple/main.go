package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
)

func main() {
	mod := &gim.Module{
		Controllers: []*gim.Controller{
			{
				Path: "",
				Routes: []*gim.Route{
					{
						Get: func(args gim.RouteArgs) interface{} {
							return "hello world"
						},
					},
				},
			},
		},
		Middlewares: []*gim.Middleware{
			{
				Use: func(c *gin.Context) {
					fmt.Println("received request")
					c.Next()
					fmt.Println("responded request")
				},
			},
		},
	}
	eng := mod.Bootstrap()
	eng.Run("0.0.0.0:3000")
}
