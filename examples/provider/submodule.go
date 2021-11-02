package main

type SubModule struct{}

func (m *SubModule) Providers() []interface{} { return []interface{}{newSubService} }
func (m *SubModule) Exports() []interface{}   { return []interface{}{newSubService} }

type subService struct {
	greetingIndex      int
	greetingcandidates []string
}

func newSubService() *subService {
	var svc subService
	svc.greetingcandidates = []string{"hello", "world"}
	return &svc
}
func (svc *subService) getGreeting() string {
	res := svc.greetingcandidates[svc.greetingIndex]
	svc.greetingIndex++
	if svc.greetingIndex >= len(svc.greetingcandidates) {
		svc.greetingIndex = 0
	}
	return res
}
