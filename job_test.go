package gim_test

import (
	"testing"

	"github.com/onichandame/gim"
	"github.com/stretchr/testify/assert"
)

type JobModule struct{}

func (mod *JobModule) Providers() []interface{} {
	return []interface{}{newImmediateJobService, &storage}
}

type ImmediateJobService struct {
	storage *JobStorageService
}

func newImmediateJobService(storage *JobStorageService) *ImmediateJobService {
	return &ImmediateJobService{storage: storage}
}
func (svc *ImmediateJobService) Run() {
	svc.storage.Int++
}
func (svc *ImmediateJobService) Blocking() bool {
	return true
}

type JobStorageService struct {
	Int int
}

var storage JobStorageService

func TestJob(t *testing.T) {
	assert.NotPanics(t, func() { gim.Bootstrap(&JobModule{}) })
	assert.Equal(t, 1, storage.Int)
}
