package gim

import (
	"context"
)

type ProviderArgs struct {
}

type Provider struct {
	Provide interface{}
	Key     interface{}
}

func (p *Provider) bootstrap(app context.Context) context.Context {
	key := p.Key
	if key == nil {
		key = p
	}
	return context.WithValue(app, key, p.Provide)
}
