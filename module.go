package gim

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Module struct {
	Imports     []*Module
	Middlewares []*Middleware
	Controllers []*Controller
	Jobs        []*Job
	app         context.Context
	rg          *gin.RouterGroup
}

func (m *Module) bootstrap(app context.Context) context.Context {
	// do not double-bootstrap a module
	if m := app.Value(m); m != nil {
		return app
	}
	app = context.WithValue(app, m, m)
	// load middlewares
	for _, mw := range m.Middlewares {
		m.rg.Use(mw.Use)
	}
	// init sub-modules
	for _, sub := range m.Imports {
		sub.rg = m.rg.Group("")
		app = sub.bootstrap(app)
	}
	// init jobs
	for _, job := range m.Jobs {
		job.bootstrap(app)
	}
	// init controllers
	for _, ctlr := range m.Controllers {
		ctlr.bootstrap(m.rg, app)
	}
	return app
}

func (m *Module) Bootstrap() *gin.Engine {
	eng := gin.Default()
	m.rg = eng.Group("")
	app := context.Background()
	app = m.bootstrap(app)
	m.app = app
	return eng
}
