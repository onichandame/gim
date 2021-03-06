package gim

import (
	"fmt"
	"reflect"
	"time"

	"github.com/onichandame/gim/pkg/injector"
	goutils "github.com/onichandame/go-utils"
	"github.com/sirupsen/logrus"
)

// Module defines a self-contained module with all its internal and external dependencies.
//
// Every module should be uniquely defined by one instance of this type.
//
// The dependency graph of all modules in an application must be a directed acyclic graph(DAG), in
// other words, cyclic dependency is not allowed. If the rule is broken, Bootstrap will fail.
type Module struct {
	Name          string        // defines the displaying name for this module used for logging purposes.
	Imports       []*Module     // defines the dependency on external modules.
	Exports       []interface{} // defines the providers in the module which can be used by other modules.
	Providers     []interface{} // defines all the providers in this module.
	modcontainers map[*Module]injector.Container
}

// Get returns the singleton of the specified type in a bootstrapped module. The calling module must
// be the root(where Bootstrap is called). It returns nil if not found
//
// The singleton need not be exported. It is not recommended to rely on this method when many private
// singletons exist in many modules, as it's behaviour is not defined in such circumstance.
func (a *Module) Get(prov interface{}) interface{} {
	for _, c := range a.modcontainers {
		t := getType(prov)
		sing := c[t]
		if sing != nil {
			return sing
		}
	}
	return nil
}

// Bootstrap injects the dependencies as declared into every module in the tree.
// It panics on any error occurred
func (main *Module) Bootstrap() {
	logger := logrus.New().WithFields(logrus.Fields{
		"pkg": "Gim",
	})
	logger.Logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
	main.modcontainers = make(map[*Module]injector.Container)
	var loadModule func(mod *Module, visited map[interface{}]interface{})
	loadModule = func(mod *Module, visited map[interface{}]interface{}) {
		getNewVisited := func() map[interface{}]interface{} {
			res := make(map[interface{}]interface{})
			for k, v := range visited {
				res[k] = v
			}
			res[mod] = nil
			return res
		}
		if _, ok := visited[mod]; ok {
			panic(fmt.Errorf("circular module dependency detected for module %v", mod.Name))
		}
		if _, ok := main.modcontainers[mod]; !ok {
			main.modcontainers[mod] = injector.NewContainer()
			for _, child := range mod.Imports {
				loadModule(child, getNewVisited())
				var loadExports func(c injector.Container, exp interface{})
				loadExports = func(c injector.Container, exp interface{}) {
					if expmod, ok := exp.(*Module); ok {
						for _, expexp := range expmod.Exports {
							loadExports(main.modcontainers[expmod], expexp)
						}
					} else {
						exptype, expsing := getTypeAndSingleton(c, exp)
						main.modcontainers[mod][exptype] = expsing
					}
				}
				for _, childexp := range child.Exports {
					loadExports(main.modcontainers[child], childexp)
				}
			}
			startTime := time.Now()
			logger := logger.WithField("module", mod.Name)
			unsorted := make([]interface{}, len(mod.Providers))
			for i, v := range mod.Providers {
				unsorted[i] = v
			}
			lastunsorted := len(unsorted)
			sortedIndicis := make([]int, 0)
			var sort func()
			sort = func() {
				for i, p := range unsorted {
					t := goutils.UnwrapType(reflect.TypeOf(p))
					if t.Kind() == reflect.Func {
						resolvable := true
						for i := 0; i < t.NumIn(); i++ {
							in := goutils.UnwrapType(t.In(i))
							insing := main.modcontainers[mod].Resolve(reflect.New(in).Interface())
							if insing == nil {
								resolvable = false
								break
							}
						}
						if resolvable {
							main.modcontainers[mod].Bind(p)
							sortedIndicis = append(sortedIndicis, i)
						}
					} else {
						main.modcontainers[mod].Bind(p)
						sortedIndicis = append(sortedIndicis, i)
					}
				}
				removeElement := func(ind int) {
					unsorted = append(unsorted[:ind], unsorted[ind+1:]...)
				}
				for i := len(sortedIndicis) - 1; i >= 0; i-- {
					ind := sortedIndicis[i]
					removeElement(ind)
				}
				sortedIndicis = make([]int, 0)
				if len(unsorted) != 0 && lastunsorted == len(unsorted) {
					panic(fmt.Errorf("providers in module %v have dependencies unresolvable. it can be a circular dependency or a missing dependency", mod.Name))
				}
				lastunsorted = len(unsorted)
				if len(unsorted) != 0 {
					sort()
				}
			}
			sort()
			logger.WithField("duration", time.Since(startTime).String()).Infoln("module loaded")
		}
	}
	loadModule(main, make(map[interface{}]interface{}))
}
