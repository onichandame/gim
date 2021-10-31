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
type withControllers interface {
	Controllers() []interface{}
}
type withMiddlewares interface {
	Middlewares() []interface{}
}
type withProviders interface {
	Providers() []interface{}
}
type ProviderConfig struct {
	Global bool
	Public bool
	// either an instance or a factory function
	Provide interface{}
}

func bootstrap(module interface{}, parentContainer *injector.Container) {
	rootContainer := parentContainer.GetRoot()
	mod := reflect.New(goutils.UnwrapType(reflect.TypeOf(module))).Interface()
	// if module alread loaded, do not load again
	err := goutils.Try(func() { rootContainer.Resolve(mod) })
	if err == nil {
		return
	}
	// if module not loaded, load it to root
	rootContainer.Bind(module, true)
	container := injector.NewContainer()
	container.SetParent(parentContainer)
	// load sub-modules
	if mod, ok := module.(withImports); ok {
		for _, child := range mod.Imports() {
			bootstrap(child, container)
		}
	}
	// load providers
	if mod, ok := module.(withProviders); ok {
		providers := make([]ProviderConfig, 0)
		for _, prov := range mod.Providers() {
			v := goutils.UnwrapValue(reflect.ValueOf(prov))
			if v.Type() == reflect.TypeOf(ProviderConfig{}) {
				providers = append(providers, v.Interface().(ProviderConfig))
			} else {
				providers = append(providers, ProviderConfig{Provide: prov})
			}
		}
		bindProvider := func(prov ProviderConfig) {
			if prov.Global {
				rootContainer.Bind(prov.Provide, true)
			} else {
				container.Bind(prov.Provide, prov.Public)
			}
			// get singleton for job loading
			p := newEntity(prov.Provide)
			container.Resolve(p)
			// init jobs defined in singleton
			loadJob(p)
		}
		loadedInd := make([]int, 0)
		sort := func() {
			for _, ind := range loadedInd {
				providers = append(providers[:ind], providers[ind+1:]...)
			}
			loadedInd = make([]int, 0)
		}
		sortProviders := func() {
			// static providers
			for ind, prov := range providers {
				t := goutils.UnwrapType(reflect.TypeOf(prov.Provide))
				if t.Kind() != reflect.Func {
					bindProvider(prov)
					loadedInd = append(loadedInd, ind)
				}
			}
			sort()
			for ind, prov := range providers {
				t := reflect.TypeOf(prov.Provide)
				resolvable := true
				for i := 0; i < t.NumIn(); i++ {
					in := goutils.UnwrapType(t.In(i))
					if goutils.Try(func() { container.Resolve(reflect.New(in).Interface()) }) != nil {
						resolvable = false
						break
					}
				}
				if resolvable {
					bindProvider(prov)
					loadedInd = append(loadedInd, ind)
				}
			}
			sort()
		}
		lastlength := len(providers)
		for {
			sortProviders()
			if len(providers) == lastlength {
				panic(fmt.Errorf("providers in module %v have circular dependency", module))
			} else {
				lastlength = len(providers)
			}
			if len(providers) == 0 {
				break
			}
		}
	}
	// load middlewares
	if mod, ok := module.(withMiddlewares); ok {
		var eng gin.Engine
		container.Resolve(&eng)
		for _, mw := range mod.Middlewares() {
			container.Bind(mw, false)
			m := newEntity(mw)
			container.Resolve(m)
			if v, ok := m.(withMiddleware); ok {
				eng.Use(v.Use())
			}
		}
	}
	// load controllers
	if mod, ok := module.(withControllers); ok {
		for _, ctlr := range mod.Controllers() {
			container.Bind(ctlr, false)
			var eng gin.Engine
			container.Resolve(&eng)
			instance := newEntity(ctlr)
			container.Resolve(instance)
			path := ""
			if p, ok := instance.(pathed); ok {
				path = p.Path()
			}
			grp := eng.Group(path)
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
			if h, ok := instance.(getRouted); ok {
				grp.GET("", getHandler(h.Get))
			}
			if h, ok := instance.(postRouted); ok {
				grp.POST("", getHandler(h.Post))
			}
			if h, ok := instance.(putRouted); ok {
				grp.PUT("", getHandler(h.Put))
			}
			if h, ok := instance.(deleteRouted); ok {
				grp.DELETE("", getHandler(h.Delete))
			}
		}
	}
}

type rootModule struct {
	mainMod interface{}
}

func (r rootModule) Imports() []interface{} { return []interface{}{&HTTPModule{}, r.mainMod} }

func Bootstrap(main interface{}) *injector.Container {
	var root rootModule
	root.mainMod = main
	rootContainer := injector.NewContainer()
	bootstrap(&root, rootContainer)
	return rootContainer
}
