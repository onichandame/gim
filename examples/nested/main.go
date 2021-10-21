package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
)

func main() {
	mod := gim.Module{
		Imports: []*gim.Module{
			{
				Middlewares: []*gim.Middleware{
					{
						Use: func(c *gin.Context) {
							fmt.Println("nested middleware in")
							c.Next()
							fmt.Println("nested middleware out")
						},
					},
				},
				Controllers: []*gim.Controller{
					{
						Path: "nested",
						Routes: []*gim.Route{
							{
								Get: func(args gim.RouteArgs) interface{} {
									return "greetings from nested"
								},
							},
						},
					},
				},
			},
		},
		Controllers: []*gim.Controller{
			{
				Routes: []*gim.Route{
					{
						Endpoint: "sub",
						Get: func(args gim.RouteArgs) interface{} {
							return "greetings from sub"
						},
					},
					{
						Get: func(args gim.RouteArgs) interface{} {
							return "greetings from root"
						},
					},
				},
			},
		},
		Middlewares: []*gim.Middleware{
			{
				Use: func(c *gin.Context) {
					fmt.Println("root in")
					c.Next()
					fmt.Println("root out")
				},
			},
		},
	}
	eng := mod.Bootstrap()
	eng.Run("0.0.0.0:3000")
}
