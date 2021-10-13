package gim

import "github.com/gin-gonic/gin"

type Provider struct {
	Inject func(g *gin.RouterGroup)
}
