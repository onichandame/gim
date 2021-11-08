package gin

import "github.com/gin-gonic/gin"

type GinService struct {
	engine     *gin.Engine
	middewares []gin.HandlerFunc
	routes     []func(*gin.RouterGroup)
}

func newGinService() *GinService {
	var svc GinService
	svc.engine = gin.Default()
	svc.middewares = make([]gin.HandlerFunc, 0)
	svc.routes = make([]func(*gin.RouterGroup), 0)
	return &svc
}
func (svc *GinService) SetServer(s *gin.Engine) { svc.engine = s }

func (svc *GinService) Server() *gin.Engine { return svc.engine }

type Controller struct {
	Path string
}

func (svc *GinService) AddMiddleware(mw gin.HandlerFunc) {
	svc.middewares = append(svc.middewares, mw)
}

func (svc *GinService) AddRoute(fn func(*gin.RouterGroup)) {
	svc.routes = append(svc.routes, fn)
}

func (svc *GinService) Bootstrap() *gin.Engine {
	for _, mw := range svc.middewares {
		svc.engine.Use(mw)
	}
	for _, route := range svc.routes {
		route(svc.engine.Group(""))
	}
	return svc.engine
}
