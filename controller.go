package gim

import (
	"github.com/gin-gonic/gin"
)

type pathed interface {
	Path() string
}

type getRouted interface {
	Get(*gin.Context) interface{}
}
type postRouted interface {
	Post(*gin.Context) interface{}
}
type putRouted interface {
	Put(*gin.Context) interface{}
}
type deleteRouted interface {
	Delete(*gin.Context) interface{}
}
