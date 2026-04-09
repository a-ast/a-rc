package cron

import (
	"fmt"
	"sync"

	"a-rc/internal/core"
	robfigcron "github.com/robfig/cron/v3"
)

// Adapter implements core.JobScheduler using robfig/cron.
type Adapter struct {
	mu   sync.Mutex
	cron *robfigcron.Cron
}

func New() *Adapter {
	return &Adapter{}
}

func (a *Adapter) Start(jobs []core.Job, runner func(core.Job) error) error {
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

func (a *Adapter) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cron != nil {
		a.cron.Stop()
		a.cron = nil
	}
}

func (a *Adapter) Reload(jobs []core.Job, runner func(core.Job) error) error {
	a.Stop()
	return a.Start(jobs, runner)
}
