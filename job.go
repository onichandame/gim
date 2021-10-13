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
	if j.Cron != "" {
		cron := cron.New()
		cron.AddFunc(j.Cron, j.Run)
		cron.Start()
	}
	if j.Immediate != nil {
		if j.Immediate.Blocking {
			j.Run()
		} else {
			go j.Run()
		}
	}
}
