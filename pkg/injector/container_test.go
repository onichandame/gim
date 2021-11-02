package injector_test

import (
	"testing"

	"github.com/onichandame/gim/pkg/injector"
	"github.com/stretchr/testify/assert"
)

func TestContainer(t *testing.T) {
	t.Run("bind singleton", func(t *testing.T) {
		type Entity struct{}
		t.Run("concret", func(t *testing.T) {
			assert.NotPanics(t, func() {
				container := injector.NewContainer()
				container.Bind(new(Entity))
			})
			assert.Panics(t, func() {
				container := injector.NewContainer()
				container.Bind(new(Entity))
				container.Bind(new(Entity))
			})
		})
		t.Run("constructor", func(t *testing.T) {
			assert.NotPanics(t, func() {
				container := injector.NewContainer()
				container.Bind(func() *Entity { return new(Entity) })
			})
			assert.Panics(t, func() {
				container := injector.NewContainer()
				container.Bind(func() *Entity { return new(Entity) })
				container.Bind(new(Entity))
			})
		})
		t.Run("dependency", func(t *testing.T) {
			assert.NotPanics(t, func() {
				container := injector.NewContainer()
				type Dependent struct{ ent *Entity }
				ent := new(Entity)
				container.Bind(ent)
				container.Bind(func(ent *Entity) *Dependent {
					var d Dependent
					d.ent = ent
					return &d
				})
			})
			assert.Panics(t, func() {
				container := injector.NewContainer()
				type Dependent struct{ ent *Entity }
				ent := new(Entity)
				container.Bind(ent)
				container.Bind(func(ent *Entity) *Dependent {
					var d Dependent
					d.ent = ent
					return &d
				})
				container.Bind(new(Dependent))
			})
		})
	})
	t.Run("resolve singleton", func(t *testing.T) {
		type Entity struct{}
		t.Run("concret", func(t *testing.T) {
			assert.NotPanics(t, func() {
				container := injector.NewContainer()
				ent := new(Entity)
				container.Bind(ent)
				var resolved Entity
				container.Resolve(&resolved)
				assert.Equal(t, ent, &resolved)
			})
			assert.Panics(t, func() {
				container := injector.NewContainer()
				container.Resolve(new(Entity))
			})
		})
		t.Run("constructor", func(t *testing.T) {
			container := injector.NewContainer()
			var ent Entity
			container.Bind(func() *Entity { return &ent })
			var resolved Entity
			container.Resolve(&resolved)
			assert.Equal(t, &ent, &resolved)
		})
		t.Run("dependency", func(t *testing.T) {
			container := injector.NewContainer()
			type Dependent struct{ ent *Entity }
			var ent Entity
			container.Bind(&ent)
			container.Bind(func(ent *Entity) *Dependent {
				var d Dependent
				d.ent = ent
				return &d
			})
			var d Dependent
			container.Resolve(&d)
			assert.Equal(t, &ent, d.ent)
		})
	})
}
