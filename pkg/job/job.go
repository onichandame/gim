package job

import "github.com/robfig/cron/v3"

type JobModule struct{}

func (m *JobModule) Providers() []interface{} { return []interface{}{newJobService} }
func (m *JobModule) Exports() []interface{}   { return []interface{}{&JobService{}} }

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
