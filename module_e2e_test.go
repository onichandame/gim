package gim_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
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
			},
			{
				Path: "2", Routes: []*gim.Route{
					{
						Endpoint: "",
						Post: func(c *gin.Context) interface{} {
							return fmt.Sprintf("%s2", c.GetString("response"))
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
	}
	r := mod.Bootstrap()
	rec1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/1", nil)
	r.ServeHTTP(rec1, req1)
	assert.Equal(t, 200, rec1.Code)
	assert.Equal(t, "mid11", rec1.Body.String())
	rec2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/2", nil)
	r.ServeHTTP(rec2, req2)
	assert.Equal(t, 200, rec2.Code)
	assert.Equal(t, "mid1mid22", rec2.Body.String())
}
