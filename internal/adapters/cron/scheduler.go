package cron

import (
	"fmt"
	"sync"

	"a-rc/internal/domain"

	robfigcron "github.com/robfig/cron/v3"
)

// Scheduler implements domain.JobScheduler using robfig/cron.
type Scheduler struct {
	mu   sync.Mutex
	cron *robfigcron.Cron
}

func New() *Scheduler {
	return &Scheduler{}
}

func (a *Scheduler) Start(jobs []domain.Job, runner func(domain.Job) error) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	c := robfigcron.New()
	for _, job := range jobs {
		j := job // capture loop variable
		if _, err := c.AddFunc(j.Schedule, func() {
			_ = runner(j) // errors are logged inside runner
		}); err != nil {
			return fmt.Errorf("invalid schedule for path %q: %w", j.Path, err)
		}
	}
	c.Start()
	a.cron = c
	return nil
}

func (a *Scheduler) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cron != nil {
		a.cron.Stop()
		a.cron = nil
	}
}

func (a *Scheduler) Reload(jobs []domain.Job, runner func(domain.Job) error) error {
	a.Stop()
	return a.Start(jobs, runner)
}
