package gim

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	goutils "github.com/onichandame/go-utils"
)

type RouteFunc func(*gin.Context) interface{}
type Route struct {
	Endpoint string
	Get      RouteFunc
	Post     RouteFunc
	Put      RouteFunc
	Delete   RouteFunc
}

func (r *Route) bootstrap(g *gin.RouterGroup) {
	var loaded bool
	handleRequest := func(fn RouteFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			runHandler := func() (res interface{}, err error) {
				defer goutils.RecoverToErr(&err)
				res = fn(c)
				return res, err
			}
			res, err := runHandler()
			if err != nil {
				if c.Writer.Status() <= 200 {
					c.Status(400)
				}
				c.JSON(c.Writer.Status(), gin.H{"message": err.Error()})
				return
			}
			if res != nil {
				var body []byte
				var contentType string
				if r, ok := res.([]byte); ok {
					body = r
					contentType = "text/plain"
				} else if r, ok := res.(string); ok {
					body = []byte(r)
					contentType = "text/plain"
				} else {
					if body, err = json.Marshal(res); err != nil {
						c.JSON(500, gin.H{"message": err.Error()})
						return
					}
					contentType = "application/json"
				}
				c.Data(200, contentType, body)
			}
		}
	}
	if r.Post != nil {
		g.POST(r.Endpoint, handleRequest(r.Post))
		loaded = true
	}
	if r.Get != nil {
		g.GET(r.Endpoint, handleRequest(r.Get))
		loaded = true
	}
	if r.Put != nil {
		g.PUT(r.Endpoint, handleRequest(r.Put))
		loaded = true
	}
	if r.Delete != nil {
		g.DELETE(r.Endpoint, handleRequest(r.Delete))
		loaded = true
	}
	if !loaded {
		panic((fmt.Errorf("cannot load route %s as it does not implement any HTTP method",
			reflect.TypeOf(r).Elem().String())))
	}
}
