package gim

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Path   string
	Routes []*Route
}

func (c *Controller) bootstrap(eng *gin.RouterGroup, app context.Context) {
	g := eng.Group(c.Path)
	for _, r := range c.Routes {
		r.bootstrap(g, app)
	}
}
