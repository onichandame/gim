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
			t.Run("existing", func(t *testing.T) {
				container := injector.NewContainer()
				ent := new(Entity)
				container.Bind(ent)
				assert.Equal(t, ent, container.Resolve(new(Entity)))
			})
			t.Run("non-existing", func(t *testing.T) {
				container := injector.NewContainer()
				assert.Nil(t, container.Resolve(new(Entity)))
			})
		})
		t.Run("constructor", func(t *testing.T) {
			container := injector.NewContainer()
			var ent Entity
			container.Bind(func() *Entity { return &ent })
			var resolved Entity
			assert.Equal(t, &ent, container.Resolve(&resolved))
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
			assert.Equal(t, &ent, container.Resolve(new(Dependent)).(*Dependent).ent)
		})
		t.Run("rebind to another container", func(t *testing.T) {
			c1 := injector.NewContainer()
			c2 := injector.NewContainer()
			var ent Entity
			c1.Bind(&ent)
			res1 := c1.Resolve(new(Entity)).(*Entity)
			c2.Bind(&res1)
			res2 := c2.Resolve(new(Entity)).(*Entity)
			assert.Equal(t, &ent, res2)
			assert.Equal(t, res1, res2)
		})
		t.Run("resolve or panic", func(t *testing.T) {
			c := injector.NewContainer()
			assert.Panics(t, func() { c.ResolveOrPanic("") })
			c.Bind(new(Entity))
			assert.NotPanics(t, func() { c.ResolveOrPanic(new(Entity)) })
		})
	})
}
