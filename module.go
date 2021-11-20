package gim

import (
	"fmt"
	"reflect"
	"time"

	"github.com/onichandame/gim/pkg/injector"
	goutils "github.com/onichandame/go-utils"
	"github.com/sirupsen/logrus"
)

type Module struct {
	Name          string
	Imports       []*Module
	Exports       []interface{}
	Providers     []interface{}
	modcontainers map[*Module]injector.Container
}

func (a *Module) Get(prov interface{}) interface{} {
	for _, c := range a.modcontainers {
		ent := newEntity(prov)
		sing := c.Resolve(ent)
		if sing != nil {
			return sing
		}
	}
	return nil
}

func (main *Module) Bootstrap() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		PadLevelText:     true,
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
						expsing := getSingleton(c, exp)
						main.modcontainers[mod].Bind(expsing)
					}
				}
				for _, childexp := range child.Exports {
					loadExports(main.modcontainers[child], childexp)
				}
			}
			startTime := time.Now()
			logger.Infoln(fmt.Sprintf("[%v] loading\n", mod.Name))
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
					panic(fmt.Errorf("providers in module %v have dependencies unresolvable. it can be a circular dependency or a missing dependenckjy", mod.Name))
				}
				lastunsorted = len(unsorted)
				if len(unsorted) != 0 {
					sort()
				}
			}
			sort()
			logger.Infoln(fmt.Sprintf("[%v] loaded in %v \n", mod.Name, time.Since(startTime).String()))
		}
	}
	loadModule(main, make(map[interface{}]interface{}))

}
