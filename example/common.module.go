package main

import "github.com/onichandame/gim"

var CommonModule = gim.Module{
	Name:      `CommonModule`,
	Providers: []interface{}{newCommonService},
	Exports:   []interface{}{newCommonService},
}

type CommonService struct {
	counter int
}

func newCommonService() *CommonService {
	var svc CommonService
	return &svc
}

func (svc *CommonService) Add(i int) int {
	svc.counter += i
	return svc.counter
}
