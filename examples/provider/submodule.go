package main

import "fmt"

type SubModule struct{}

func (m *SubModule) Providers() []interface{} { return []interface{}{newSubService} }
func (m *SubModule) Exports() []interface{}   { return []interface{}{newSubService} }

type SubService struct {
	greetingIndex      int
	greetingcandidates []string
}

var subsvc *SubService

func newSubService() *SubService {
	var svc SubService
	svc.greetingcandidates = []string{"hello", "world"}
	fmt.Println("svc")
	fmt.Printf("%p\n", &svc)
	return &svc
}
func (svc *SubService) getGreeting() string {
	res := svc.greetingcandidates[svc.greetingIndex]
	svc.greetingIndex++
	if svc.greetingIndex >= len(svc.greetingcandidates) {
		svc.greetingIndex = 0
	}
	return res
}
, 