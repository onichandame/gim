package gim

import (
	"context"
)

type ProviderArgs struct {
}

type Provider struct {
	Provide interface{}
	Factory func(app context.Context) interface{}
	Token     interface{}
}

func(p*Provider)getToken()interface{}{
	if p.Token==nil{return p}else{return p.Token}
}

func (p *Provider) bootstrap(app context.Context) context.Context {
	key := p.getToken()
	val := p.Provide
	if p.Factory != nil {
		val = p.Factory(app)
	}
	return context.WithValue(app, key, val)
}
