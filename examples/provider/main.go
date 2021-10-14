package main

import (
	"log"

	"github.com/onichandame/gim"
)

func main() {
	prov := gim.Provider{
		Provide: "hello world",
	}
	mod := gim.Module{
		Providers: []*gim.Provider{&prov, {Provide: "hello world(alt)", Token: "alt"}},
		Controllers: []*gim.Controller{
			{
				Routes: []*gim.Route{
					{
						Get: func(args gim.RouteArgs) interface{} {
							return args.App.Value(&prov)
						},
					},
					{
						Endpoint: "alt",
						Get: func(args gim.RouteArgs) interface{} {
							return args.App.Value("alt")
						},
					},
				},
			},
		},
	}
	eng := mod.Bootstrap()
	str := mod.Get(&prov).(string)
	log.Printf("response content: %s", str)
	eng.Run("0.0.0.0:3000")
}
