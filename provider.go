package gim

import (
	"context"
)

type ProviderArgs struct {
}

type Provider struct {
	Provide interface{}
	Factory func(app context.Context) interface{}
	Key     interface{}
}

func (p *Provider) bootstrap(app context.Context) context.Context {
	key := p.Key
	if key == nil {
		key = p
	}
	val := p.Provide
	if p.Factory != nil {
		val = p.Factory(app)
	}
	return context.WithValue(app, key, val)
}
