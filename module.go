package gim

import (
	"github.com/gin-gonic/gin"
)

type Module struct {
	Imports     []*Module
	Controllers []*Controller
	Middlewares []*Middleware
}

func (m *Module) bootstrap(pg *gin.RouterGroup) {
	subGroup := pg.Group("")
	// middlewares have to be loaded before subGroups are added
	for _, mw := range m.Middlewares {
		subGroup.Use(mw.Use)
	}
	for _, ctlr := range m.Controllers {
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
