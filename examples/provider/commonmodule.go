package main

import "github.com/onichandame/gim"

var CommonModule = gim.Module{
	Exports:   []interface{}{new(CommonService)},
	Providers: []interface{}{new(CommonService)},
}

type CommonService struct {
	tick int
}

func (svc *CommonService) GetNextTick() int {
	svc.tick = 1 - svc.tick
	return svc.tick
}
