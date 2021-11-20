package main

import (
	"fmt"

	"github.com/onichandame/gim"
)

var RootModule = gim.Module{
	Name:      `RootModule`,
	Imports:   []*gim.Module{&CommonModule, &ChildModule},
	Providers: []interface{}{newRootService},
}

type RootService struct{}

func newRootService(cmsvc *CommonService) *RootService {
	var svc RootService
	fmt.Printf("Count from root service: %v \n", cmsvc.Add(1))
	return &svc
}
