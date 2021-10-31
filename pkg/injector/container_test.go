package injector_test

import (
	"testing"

	"github.com/onichandame/gim/pkg/injector"
	"github.com/stretchr/testify/assert"
)

func TestContainer(t *testing.T) {
	t.Run("bind singleton", func(t *testing.T) {
		container := injector.NewContainer()
		type Entity struct{}
		ent := new(Entity)
		assert.NotPanics(t, func() {
			container.Bind(ent, false)
		})
	})
	t.Run("resolve singleton", func(t *testing.T) {
		container := injector.NewContainer()
		type Entity struct{}
		ent := new(Entity)
		container.Bind(ent, false)
		var resolved Entity
		assert.NotPanics(t, func() {
			container.Resolve(&resolved)
		})
		assert.Equal(t, ent, &resolved)
		assert.Panics(t, func() {
			type Entity struct{}
			var ent Entity
			container.Resolve(&ent)
		})
	})
	t.Run("dynamic binding", func(t *testing.T) {
		container := injector.NewContainer()
		type Entity struct{}
		type Box struct {
			Entity *Entity
		}
		var ent Entity
		container.Bind(&ent, false)
		assert.NotPanics(t, func() {
			container.Bind(func(ent *Entity) *Box { return &Box{Entity: ent} }, false)
		})
		var box Box
		assert.NotPanics(t, func() {
			container.Resolve(&box)
		})
		assert.Equal(t, &ent, box.Entity)
	})
	t.Run("resolve children", func(t *testing.T) {
		t.Run("private", func(t *testing.T) {
			parent := injector.NewContainer()
			child := injector.NewContainer()
			child.SetParent(parent)
			type Entity struct{}
			var ent Entity
			child.Bind(&ent, false)
			assert.Panics(t, func() {
				parent.Resolve(new(Entity))
			})
		})
		t.Run("public", func(t *testing.T) {
			parent := injector.NewContainer()
			child := injector.NewContainer()
			child.SetParent(parent)
			type Entity struct{}
			var ent Entity
			child.Bind(&ent, true)
			assert.NotPanics(t, func() {
				var resolved Entity
				parent.Resolve(&resolved)
				assert.Equal(t, &ent, &resolved)
			})
		})
		t.Run("grandchildren", func(t *testing.T) {
			parent := injector.NewContainer()
			child := injector.NewContainer()
			grandchild := injector.NewContainer()
			child.SetParent(parent)
			grandchild.SetParent(child)
			type Entity struct{}
			var ent Entity
			grandchild.Bind(&ent, true)
			assert.Panics(t, func() {
				parent.Resolve(new(Entity))
			})
		})
	})
	t.Run("resolve root", func(t *testing.T) {
		root := injector.NewContainer()
		leaf := injector.NewContainer()
		leaf.SetParent(root)
		type Entity struct{}
		var ent Entity
		root.Bind(&ent, true)
		assert.NotPanics(t, func() {
			var resolved Entity
			leaf.Resolve(&resolved)
			assert.Equal(t, &ent, &resolved)
		})
	})
}
