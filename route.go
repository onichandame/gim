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
			var contentType string
			var body interface{}
			var status int
			if err != nil {
				if e, ok := err.(*GimError); ok {
					if e.Status != 0 {
						status = e.Status
					}
					if e.Body != nil {
						body = e.Body
					}
				} else {
					body = err.Error()
				}
				if status < 400 {
					status = 400
				}
			} else if res != nil {
				body = res
			}
			c.Status(status)
			var responseBody []byte
			if r, ok := body.([]byte); ok {
				responseBody = r
				contentType = "text/plain"
			} else if r, ok := body.(string); ok {
				responseBody = []byte(r)
				contentType = "text/plain"
			} else {
				if responseBody, err = json.Marshal(body); err != nil {
					status = 500
					responseBody = []byte("failed to serialize response body")
					contentType = "text/plain"
					fmt.Print(err)
				} else {
					contentType = "application/json"
				}
			}
			c.Data(status, contentType, responseBody)
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
