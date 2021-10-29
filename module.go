package gim

import (
	"github.com/gin-gonic/gin"
)

type Module struct {
	Imports     []*Module
	Middlewares []*Middleware
	Controllers []*Controller
	Jobs        []*Job
	Engine      *gin.Engine
	modules     map[*Module]interface{}
}

func (m *Module) bootstrap() {
	// do not double-bootstrap a module
	if _, ok := m.modules[m]; ok {
		return
	}
	m.modules[m] = nil
	// load middlewares
	for _, mw := range m.Middlewares {
		m.Engine.Use(mw.Use)
	}
	// init sub-modules
	for _, sub := range m.Imports {
		sub.Engine = m.Engine
		sub.bootstrap()
	}
	// init jobs
	for _, job := range m.Jobs {
		job.bootstrap()
	}
	// init controllers
	for _, ctlr := range m.Controllers {
		ctlr.bootstrap(m.Engine)
	}
}

func (m *Module) Bootstrap() *gin.Engine {
	eng := gin.Default()
	m.Engine = eng
	m.bootstrap()
	return eng
}
