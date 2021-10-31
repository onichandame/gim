package gim_test

import (
	"testing"

	"github.com/onichandame/gim"
	"github.com/stretchr/testify/assert"
)

type MainModule struct{}

func TestModule(t *testing.T) {
	assert.NotPanics(t, func() { gim.Bootstrap(&MainModule{}) })
}
