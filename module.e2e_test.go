package gim_test

import (
	"fmt"
	"testing"

	"github.com/onichandame/gim"
	"github.com/stretchr/testify/assert"
)

type MainModule struct{}

func (m *MainModule) Imports() []interface{}     { return []interface{}{&SubModule{}} }
func (m *MainModule) Controllers() []interface{} { return []interface{}{newMainController} }

type MainController struct {
}

var ctlSubService *SubService

func newMainController(subSvc *SubService) *MainController {
	var ctl MainController
	ctlSubService = subSvc
	return &ctl
}

type SubModule struct{}

func (m *SubModule) Providers() []interface{} {
	return []interface{}{newSubService, newSubPrivateService}
}
func (m *SubModule) Exports() []interface{} { return []interface{}{newSubService} }

type SubService struct{}

var subService *SubService

func newSubService() *SubService {
	var svc SubService
	subService = &svc
	return &svc
}

type SubPrivateService struct{}

var spsvc *SubService

func newSubPrivateService(s *SubService) *SubPrivateService {
	var svc SubPrivateService
	spsvc = s
	return &svc
}

func TestGimModule(t *testing.T) {
	gim.Bootstrap(&MainModule{})
	getptr := func(ptr interface{}) string { return fmt.Sprintf("%p", ptr) }
	assert.NotNil(t, subService)
	assert.NotNil(t, ctlSubService)
	assert.NotNil(t, spsvc)
	assert.True(t, ctlSubService == subService)
	assert.True(t, getptr(ctlSubService) == getptr(subService))
	assert.Equal(t, subService, ctlSubService)
	assert.Equal(t, subService, spsvc)
}
