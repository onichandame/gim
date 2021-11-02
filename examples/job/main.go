package main

import (
	"fmt"

	"github.com/onichandame/gim"
	"github.com/onichandame/gim/pkg/job"
)

type MainModule struct{}

func (m *MainModule) Imports() []interface{}   { return []interface{}{&job.JobModule{}} }
func (m *MainModule) Providers() []interface{} { return []interface{}{newMainService} }

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
	app := gim.Bootstrap(&MainModule{})
	app.Server().Run("0.0.0.0:80")
}
