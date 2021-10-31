package gim

import (
	"github.com/robfig/cron/v3"
)

type job interface {
	Run()
}

type cronjob interface {
	Cron() string
}

type immediatejob interface {
	Blocking() bool
}

func loadJob(instance interface{}) {
	if j, ok := instance.(job); ok {
		if c, ok := j.(cronjob); ok {
			cron.New().AddFunc(c.Cron(), j.Run)
		}
		if im, ok := j.(immediatejob); ok {
			if im.Blocking() {
				j.Run()
			} else {
				go j.Run()
			}
		}
	}
}
