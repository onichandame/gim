package gim

import (
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Path   string
	Routes []*Route
}

func (c *Controller) bootstrap(g *gin.RouterGroup) {
	sg := g.Group(c.Path)
	for _, route := range c.Routes {
		route.bootstrap(sg)
	}
}
