package gim

import (
	"github.com/robfig/cron/v3"
)

type ImmediateJobConfig struct {
	Blocking bool
}

type Job struct {
	Cron      string
	Immediate *ImmediateJobConfig
	Run       func()
}

func (j *Job) bootstrap() {
	run := func() { j.Run() }
	if j.Cron != "" {
		cron := cron.New()
		cron.AddFunc(j.Cron, run)
		cron.Start()
	}
	if j.Immediate != nil {
		if j.Immediate.Blocking {
			run()
		} else {
			go run()
		}
	}
}
