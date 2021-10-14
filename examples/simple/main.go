package main

import "github.com/onichandame/gim"

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
	}
	eng := mod.Bootstrap()
	eng.Run("0.0.0.0:3000")
}
