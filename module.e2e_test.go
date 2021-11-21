package gim_test

import (
	"fmt"
	"testing"

	"github.com/onichandame/gim"
	"github.com/stretchr/testify/assert"
)

var MainModule = gim.Module{
	Name:      `MainModule`,
	Imports:   []*gim.Module{&SubModule, &SubDummyModule},
	Providers: []interface{}{newMainProvider},
}

type MainProvider struct {
	subprov SubProvider
}

func newMainProvider(subprov SubProvider) *MainProvider {
	var prov MainProvider
	prov.subprov = subprov
	return &prov
}

var SubModule = gim.Module{
	Name: `SubModule`,
	Providers: []interface{}{
		newSubPrivateService,
		newSubProvider,
	},
	Exports: []interface{}{newSubProvider},
}

type SubProvider interface {
	GetName() string
}

type SubService struct{}

func (*SubService) GetName() string { return `SubService` }

func newSubProvider() SubProvider {
	var svc SubService
	return &svc
}

type SubPrivateService struct {
	prov SubProvider
}

func newSubPrivateService(prov SubProvider) *SubPrivateService {
	var svc SubPrivateService
	svc.prov = prov
	return &svc
}

var SubDummyModule = gim.Module{Name: `SubDummyModule`}

func TestGimModule(t *testing.T) {
	MainModule.Bootstrap()
	getptr := func(ptr interface{}) string { return fmt.Sprintf("%p", ptr) }
	mainprov := MainModule.Get(new(MainProvider)).(*MainProvider)
	subprov := MainModule.Get(new(SubProvider)).(SubProvider)
	subprivprov := MainModule.Get(new(SubPrivateService)).(*SubPrivateService)
	assert.NotNil(t, mainprov)
	assert.NotNil(t, mainprov.subprov)
	assert.NotNil(t, subprov)
	assert.NotNil(t, subprivprov)
	assert.NotNil(t, subprivprov.prov)
	assert.True(t, mainprov.subprov == subprov)
	assert.True(t, getptr(mainprov.subprov) == getptr(subprov))
	assert.Equal(t, subprivprov.prov, subprov)
	assert.True(t, getptr(subprivprov.prov) == getptr(subprov))
	assert.Equal(t, subprivprov.prov, mainprov.subprov)
	assert.True(t, getptr(subprivprov.prov) == getptr(mainprov.subprov))
}
