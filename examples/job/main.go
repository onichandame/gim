package main

import (
	"fmt"

	"github.com/onichandame/gim"
	"github.com/onichandame/gim/pkg/job"
)

var MainModule = gim.Module{
	Imports:   []*gim.Module{&job.JobModule},
	Providers: []interface{}{newMainService},
}

type MainService struct{}

func newMainService(jobs *job.JobService) *MainService {
	var svc MainService
	jobs.Cron("@every 1s", svc.print)
	return &svc
}

func (svc *MainService) print() {
	fmt.Println("hello")
}

func main() {
	MainModule.Bootstrap()
	done := make(chan int)
	<-done
}
