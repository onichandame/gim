package gim

import (
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Path   string
	Routes []*Route
}

func (c *Controller) bootstrap(eng *gin.Engine) {
	g := eng.Group(c.Path)
	for _, r := range c.Routes {
		r.bootstrap(g)
	}
}
