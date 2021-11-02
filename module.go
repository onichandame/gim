package gim

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim/pkg/injector"
	goutils "github.com/onichandame/go-utils"
)

type withImports interface {
	Imports() []interface{}
}
type withServer interface {
	Server(*gin.Engine) *gin.Engine
}
type withExports interface {
	Exports() []interface{}
}
type withControllers interface {
	Controllers() []interface{}
}
type withMiddlewares interface {
	Middlewares() []interface{}
}
type withProviders interface {
	Providers() []interface{}
}

type App struct {
	modules       injector.Container
	modcontainers map[interface{}]injector.Container
	eng           *gin.Engine
}

func (a App) Server() *gin.Engine {
	return a.eng
}

func Bootstrap(main interface{}) *App {
	var app App
	app.modules = injector.NewContainer()
	app.modcontainers = make(map[interface{}]injector.Container)
	app.eng = gin.Default()
	getSingleton := func(container injector.Container, ent interface{}) interface{} {
		sing := newEntity(ent)
		container.Resolve(sing)
		return sing
	}
	var loadModule func(mod interface{}, visited map[interface{}]interface{})
	loadModule = func(mod interface{}, visited map[interface{}]interface{}) {
		newVisited := func(self interface{}) map[interface{}]interface{} {
			res := make(map[interface{}]interface{})
			for k, v := range visited {
				res[k] = v
			}
			res[self] = nil
			return res
		}

		if sing := newEntity(mod); goutils.Try(func() { app.modules.Resolve(sing) }) == nil {
			if _, ok := visited[sing]; ok {
				panic(fmt.Errorf("circular module dependency detected for module %v", goutils.UnwrapType(reflect.TypeOf(sing)).Name()))
			}
		} else {
			app.modules.Bind(mod)
			sing := getSingleton(app.modules, mod)
			app.modcontainers[sing] = injector.NewContainer()
			if m, ok := sing.(withImports); ok {
				for _, child := range m.Imports() {
					loadModule(child, newVisited(sing))
				}
			}
			if m, ok := mod.(withServer); ok {
				app.eng = m.Server(app.eng)
			}
		}
	}
	loadModule(main, make(map[interface{}]interface{}))
	main = getSingleton(app.modules, main)

	var loadProviders func(mod interface{})
	loadProviders = func(mod interface{}) {
		sing := getSingleton(app.modules, mod)
		if m, ok := sing.(withImports); ok {
			for _, child := range m.Imports() {
				loadProviders(child)
				childsing := getSingleton(app.modules, child)
				var loadExports func(c injector.Container, exp interface{})
				loadExports = func(c injector.Container, exp interface{}) {
					if goutils.Try(func() {
						expsing := getSingleton(app.modules, exp)
						if expmod, ok := expsing.(withExports); ok {
							for _, expexp := range expmod.Exports() {
								loadExports(app.modcontainers[expsing], expexp)
							}
						}
					}) != nil {
						expsing := getSingleton(c, exp)
						app.modcontainers[sing].Bind(expsing)
					}
				}
				if m, ok := childsing.(withExports); ok {
					for _, exp := range m.Exports() {
						loadExports(app.modcontainers[childsing], exp)
					}
				}
			}
		}
		if m, ok := sing.(withProviders); ok {
			sorted := make([]interface{}, 0)
			var lastsorted int
			var sort func()
			sort = func() {
				for _, p := range m.Providers() {
					t := goutils.UnwrapType(reflect.TypeOf(p))
					if t.Kind() == reflect.Func {
						resolvable := true
						for i := 0; i < t.NumIn(); i++ {
							in := goutils.UnwrapType(t.In(i))
							insing := reflect.New(in).Interface()
							if goutils.Try(func() { app.modules.Resolve(insing) }) != nil {
								resolvable = false
								break
							}
						}
						if resolvable {
							sorted = append(sorted, p)
						}
					} else {
						sorted = append(sorted, p)
					}
				}
				if lastsorted == len(sorted) {
					panic(fmt.Errorf("providers in module %v have circular dependency", goutils.UnwrapType(reflect.TypeOf(sing)).Name()))
				}
				lastsorted = len(sorted)
				if len(sorted) != len(m.Providers()) {
					sort()
				}
			}
			sort()
			for _, p := range sorted {
				app.modcontainers[sing].Bind(p)
				loadJob(getSingleton(app.modcontainers[sing], p))
			}
		}
	}

	loadProviders(main)

	var loadControllers func(mod interface{})
	loadControllers = func(mod interface{}) {
		sing := getSingleton(app.modules, mod)
		if m, ok := sing.(withImports); ok {
			for _, child := range m.Imports() {
				loadControllers(child)
			}
		}
		if m, ok := mod.(withControllers); ok {
			for _, c := range m.Controllers() {
				app.modcontainers[sing].Bind(c)
				csing := getSingleton(app.modcontainers[sing], c)
				path := ""
				if p, ok := csing.(pathed); ok {
					path = p.Path()
				}
				grp := app.eng.Group(path)
				getHandler := func(ctl func(*gin.Context) interface{}) gin.HandlerFunc {
					return func(c *gin.Context) {
						var res interface{}
						err := goutils.Try(func() { res = ctl(c) })
						var status int
						var contentType string
						var response []byte
						var body interface{}
						populateRes := func(res interface{}, defStatus int, defBody interface{}) {
							if s, ok := res.(withStatus); ok {
								status = s.Status()
							} else {
								status = defStatus
							}
							if r, ok := res.(withBody); ok {
								body = r.Body()
							} else {
								body = defBody
							}
						}
						if err == nil {
							populateRes(res, 200, res)
						} else {
							populateRes(err, 400, err.Error())
						}
						if b, ok := body.([]byte); ok {
							response = b
							contentType = "text/plain"
						} else if r, ok := body.(string); ok {
							response = []byte(r)
							contentType = "text/plain"
						} else {
							if response, err = json.Marshal(body); err != nil {
								status = 500
								response = []byte("failed to serialize response body")
								contentType = "text/plain"
								fmt.Println(err)
							} else {
								contentType = "application/json"
							}
						}
						c.Data(status, contentType, response)
					}
				}
				if h, ok := csing.(getRouted); ok {
					grp.GET(path, getHandler(h.Get))
				}
				if h, ok := csing.(postRouted); ok {
					grp.POST(path, getHandler(h.Post))
				}
				if h, ok := csing.(putRouted); ok {
					grp.PUT(path, getHandler(h.Put))
				}
				if h, ok := csing.(deleteRouted); ok {
					grp.DELETE(path, getHandler(h.Delete))
				}
			}
		}
	}

	loadControllers(main)

	return &app
}
