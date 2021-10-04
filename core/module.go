package core

import (
	"github.com/gin-gonic/gin"
)

type Module struct {
	Imports     []*Module
	Path        string
	Routes      []*Route
	Middlewares []*Middleware
	Providers   []*Provider
}

func (m *Module) bootstrap(pg *gin.RouterGroup) {
	subGroup := pg.Group(m.Path)
	// middlewares have to be loaded before subGroups are added
	for _, mw := range m.Middlewares {
		subGroup.Use(mw.Use)
	}
	for _, p := range m.Providers {
		p.Inject(subGroup)
	}
	for _, ctlr := range m.Routes {
		ctlr.bootstrap(subGroup)
	}
	for _, sm := range m.Imports {
		sm.bootstrap(subGroup)
	}
}

func (m *Module) Bootstrap() *gin.Engine {
	eng := gin.Default()
	m.bootstrap(eng.Group(""))
	return eng
}
