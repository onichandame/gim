package gim

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Module struct {
	Imports     []*Module
	Middlewares []*Middleware
	Controllers []*Controller
	Providers   []*Provider
	Jobs        []*Job
	app         context.Context
}

var eng = &Provider{
	Provide: gin.Default(),
}

func (m *Module) bootstrap(app context.Context) context.Context {
	// do not double-bootstrap a module
	if m := app.Value(m); m != nil {
		return app
	}
	app = context.WithValue(app, m, m)
	// init sub-modules
	for _, sub := range m.Imports {
		app = sub.bootstrap(app)
	}
	// init providers
	for _, p := range m.Providers {
		app = p.bootstrap(app)
	}
	// init jobs
	for _, job := range m.Jobs {
		job.bootstrap(app)
	}
	// init controllers
	eng := app.Value(eng).(*gin.Engine)
	for _, mw := range m.Middlewares {
		eng.Use(mw.Use)
	}
	for _, ctlr := range m.Controllers {
		ctlr.bootstrap(eng, app)
	}
	return app
}

func (m *Module) Bootstrap() *gin.Engine {
	app := context.Background()
	app = context.WithValue(app, eng, eng.Provide)
	app = m.bootstrap(app)
	m.app = app
	return m.app.Value(eng).(*gin.Engine)
}

func (m *Module) Get(prov *Provider) interface{} {
	key := prov.getToken()
	res := m.app.Value(key)
	return res
}
