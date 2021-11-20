package main

import (
	"fmt"

	"github.com/onichandame/gim"
)

var ChildModule = gim.Module{
	Name:      "ChildModule",
	Imports:   []*gim.Module{&CommonModule},
	Providers: []interface{}{newChildService},
}

type ChildService struct{}

func newChildService(cmsvc *CommonService) *ChildService {
	var svc ChildService
	fmt.Printf("Count from child service: %v \n", cmsvc.Add(1))
	return &svc
}
