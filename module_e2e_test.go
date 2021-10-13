package gim_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	gim "github.com/onichandame/gim"
	"github.com/stretchr/testify/assert"
)

func TestModule(t *testing.T) {
	mod := gim.Module{
		Imports: []*gim.Module{
			{
				Path: "1", Routes: []*gim.Route{
					{
						Endpoint: "",
						Get: func(c *gin.Context) interface{} {
							return fmt.Sprintf("%s1", c.GetString("response"))
						},
					},
				},
				Middlewares: []*gim.Middleware{
					{
						Use: func(c *gin.Context) {
							c.Set("response", fmt.Sprintf("%smid2", c.GetString("response")))
						},
					},
				},
				Providers: []*gim.Provider{
					{
						Inject: func(g *gin.RouterGroup) {
							g.Use(func(c *gin.Context) {
								c.Set("response", fmt.Sprintf("%sprov2", c.GetString("response")))
								c.Next()
							})
						},
					},
				},
			},
			{
				Path: "/errors",
				Routes: []*gim.Route{
					{
						Endpoint: "/err",
						Get: func(c *gin.Context) interface{} {
							panic(errors.New("err"))
						},
					},
					{
						Endpoint: "/structErr",
						Get: func(c *gin.Context) interface{} {
							type Error struct {
								Msg string `json:"msg"`
							}
							panic(gim.NewGimError(500, &Error{Msg: "err"}))
						},
					},
				},
			},
		},
		Middlewares: []*gim.Middleware{
			{
				Use: func(c *gin.Context) {
					c.Set("response", "mid1")
					c.Next()
				},
			},
		},
		Providers: []*gim.Provider{
			{
				Inject: func(g *gin.RouterGroup) {
					g.Use(func(c *gin.Context) {
						c.Set("response", fmt.Sprintf("%sprov1", c.GetString("response")))
						c.Next()
					})
				},
			},
		},
	}
	r := mod.Bootstrap()
	t.Run("parent middleware", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/1", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, 200, rec.Code)
		assert.True(t, strings.Contains(rec.Body.String(), "mid1"))
	})
	t.Run("parent provider", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/1", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, 200, rec.Code)
		assert.True(t, strings.Contains(rec.Body.String(), "prov1"))
	})
	t.Run("sibling middleware", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/1", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, 200, rec.Code)
		assert.True(t, strings.Contains(rec.Body.String(), "mid2"))
		assert.True(t, strings.Index(rec.Body.String(), "mid1") < strings.Index(rec.Body.String(), "mid2"))
	})
	t.Run("sibling provider", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/1", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, 200, rec.Code)
		assert.True(t, strings.Contains(rec.Body.String(), "prov2"))
		assert.True(t, strings.Index(rec.Body.String(), "prov1") < strings.Index(rec.Body.String(), "prov2"))
	})
	t.Run("string error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/errors/err", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, 400, rec.Code)
		assert.Equal(t, "err", rec.Body.String())
	})
	t.Run("structured error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/errors/structErr", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, 500, rec.Code)
	})
}
