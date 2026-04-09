package cron

import (
	"fmt"
	"log"
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

func (s *Scheduler) Start(jobs []domain.Job, runner func(domain.Job) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := robfigcron.New()
	for _, job := range jobs {
		j := job // capture loop variable
		if _, err := c.AddFunc(j.Schedule, func() {
			if err := runner(j); err != nil {
				log.Printf("job %q failed: %v", j.Name, err)
			}
		}); err != nil {
			return fmt.Errorf("invalid schedule for path %q: %w", j.Path, err)
		}
	}
	c.Start()
	s.cron = c
	return nil
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cron != nil {
		s.cron.Stop()
		s.cron = nil
	}
}

func (s *Scheduler) Reload(jobs []domain.Job, runner func(domain.Job) error) error {
	s.Stop()
	return s.Start(jobs, runner)
}
