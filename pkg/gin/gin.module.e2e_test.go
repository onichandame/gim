package gin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
	gimgin "github.com/onichandame/gim/pkg/gin"
	"github.com/stretchr/testify/assert"
)

var MainModule =gim.Module{Imports: []*gim.Module{&gimgin.GinModule},Providers: []interface{}{newMainController}}

type MainController struct{}

func newMainController(ginsvc *gimgin.GinService) *MainController {
	var ctl MainController
	ginsvc.AddRoute(func(rg *gin.RouterGroup) {
		rg.GET("", gimgin.GetHTTPHandler(func(c *gin.Context) interface{} { return "hello world" }))
	})
	return &ctl
}

func TestGinModule(t *testing.T) {
	MainModule.Bootstrap()
	server := MainModule.Get(&gimgin.GinService{}).(*gimgin.GinService).Bootstrap()
	assert.NotNil(t, server)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	server.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "hello world", w.Body.String())
}
