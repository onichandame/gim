package gim_test

import (
	"fmt"
	"testing"

	"github.com/onichandame/gim"
	"github.com/stretchr/testify/assert"
)

var MainModule = gim.Module{
	Imports:   []*gim.Module{&SubModule, &SubDummyModule},
	Providers: []interface{}{newMainProvider},
}

type MainProvider struct{}

var mainSubService *SubService

func newMainProvider(subsvc *SubService) *MainProvider {
	mainSubService = subsvc
	var prov MainProvider
	return &prov
}

var SubModule = gim.Module{
	Providers: []interface{}{
		newSubPrivateService,
		newSubService,
	},
	Exports: []interface{}{newSubService},
}

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

var SubDummyModule = gim.Module{}

func TestGimModule(t *testing.T) {
	MainModule.Bootstrap()
	getptr := func(ptr interface{}) string { return fmt.Sprintf("%p", ptr) }
	assert.NotNil(t, subService)
	assert.NotNil(t, mainSubService)
	assert.NotNil(t, spsvc)
	assert.True(t, mainSubService == subService)
	assert.True(t, getptr(mainSubService) == getptr(subService))
	assert.Equal(t, subService, mainSubService)
	assert.Equal(t, subService, spsvc)
	assert.Equal(t, subService, MainModule.Get(&SubService{}))
}
