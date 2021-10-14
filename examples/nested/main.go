package main

import "github.com/onichandame/gim"

func main() {
	mod := gim.Module{
		Imports: []*gim.Module{
			{
				Providers: []*gim.Provider{
					{
						Provide: "hello from deep", Key: "sub",
					},
				},
				Controllers: []*gim.Controller{
					{
						Path: "nested",
						Routes: []*gim.Route{
							{
								Get: func(args gim.RouteArgs) interface{} {
									return "greetings from depth"
								},
							},
						},
					},
				},
			},
		},
		Controllers: []*gim.Controller{
			{
				Path: "root",
				Routes: []*gim.Route{
					{
						Endpoint: "sub",
						Get: func(args gim.RouteArgs) interface{} {
							return args.App.Value("sub")
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
	}
	eng := mod.Bootstrap()
	eng.Run("0.0.0.0:3000")
}
