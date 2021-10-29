package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
)

func main() {
	mod := gim.Module{
		Imports: []*gim.Module{
			{
				Jobs: []*gim.Job{
					{
						Immediate: &gim.ImmediateJobConfig{
							Blocking: true,
						},
						Run: func() {
							fmt.Println("blocking job started")
							time.Sleep(time.Second * 2)
							fmt.Println("blocking job done")
						},
					},
					{
						Immediate: &gim.ImmediateJobConfig{},
						Run: func() {
							fmt.Println("non-blocking job started")
							time.Sleep(time.Second * 2)
							fmt.Println("non-blocking job done")
						},
					},
				},
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
