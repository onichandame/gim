package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
	gimgin "github.com/onichandame/gim/pkg/gin"
)

var SubModule = gim.Module{
	Imports:   []*gim.Module{&gimgin.GinModule},
	Providers: []interface{}{newSubService, newSubController},
	Exports:   []interface{}{newSubService},
}

type SubService struct {
	greetingIndex      int
	greetingcandidates []string
}

func newSubService() *SubService {
	var svc SubService
	svc.greetingcandidates = []string{"hello", "world"}
	return &svc
}
func (svc *SubService) getGreeting() string {
	res := svc.greetingcandidates[svc.greetingIndex]
	svc.greetingIndex++
	if svc.greetingIndex >= len(svc.greetingcandidates) {
		svc.greetingIndex = 0
	}
	return res
}

type SubController struct{ svc *SubService }

func newSubController(svc *SubService, ginsvc *gimgin.GinService) *SubController {
	var ctl SubController
	ctl.svc = svc
	ginsvc.AddRoute(func(rg *gin.RouterGroup) {
		rg.GET("sub", gimgin.GetHTTPHandler(ctl.Get))
	})
	return &ctl
}

func (c *SubController) Get(*gin.Context) interface{} {
	return c.svc.getGreeting()
}
