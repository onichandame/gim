package gim

import "github.com/gin-gonic/gin"

type HTTPModule struct{}

func (mod *HTTPModule) Providers() []interface{} {
	return []interface{}{ProviderConfig{Global: true, Provide: newEngine}}
}

func newEngine() *gin.Engine {
	return gin.Default()
}
