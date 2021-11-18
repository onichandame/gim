package job

import (
	"github.com/onichandame/gim"
	"github.com/robfig/cron/v3"
)

var JobModule = gim.Module{
	Name:      "JobModule",
	Providers: []interface{}{newJobService},
	Exports:   []interface{}{newJobService},
}

type JobService struct {
	cron *cron.Cron
}

func newJobService() *JobService {
	var svc JobService
	svc.cron = cron.New()
	svc.cron.Start()
	return &svc
}

func (svc *JobService) Cron(c string, fn func()) {
	svc.cron.AddFunc(c, fn)
}
