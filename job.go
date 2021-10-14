package gim

import (
	"context"

	"github.com/robfig/cron/v3"
)

type ImmediateJobConfig struct {
	Blocking bool
}

type Job struct {
	Cron      string
	Immediate *ImmediateJobConfig
	Run       func(app context.Context)
}

func (j *Job) bootstrap(app context.Context) {
	run := func() { j.Run(app) }
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
