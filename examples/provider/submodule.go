package main

import (
	"github.com/gin-gonic/gin"
)

type SubModule struct{}

func (m *SubModule) Controllers() []interface{} { return []interface{}{newSubController} }
func (m *SubModule) Providers() []interface{}   { return []interface{}{newSubService} }
func (m *SubModule) Exports() []interface{}     { return []interface{}{&SubService{}} }

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

func newSubController(svc *SubService) *SubController {
	var ctl SubController
	ctl.svc = svc
	return &ctl
}

func (*SubController) Path() string { return "sub" }

func (c *SubController) Get(*gin.Context) interface{} {
	return c.svc.getGreeting()
}
