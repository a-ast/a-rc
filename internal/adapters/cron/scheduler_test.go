package cron_test

import (
	"sync/atomic"
	"testing"
	"time"

	"a-rc/internal/adapters/cron"
	"a-rc/internal/domain"
)

var validJob = domain.Job{Name: "job", Path: "/tmp", Schedule: "@every 1h"}

func TestStart_InvalidSchedule(t *testing.T) {
	s := cron.New()
	jobs := []domain.Job{{Name: "bad", Path: "/tmp", Schedule: "not-a-cron"}}
	if err := s.Start(jobs, func(domain.Job) error { return nil }); err == nil {
		t.Error("expected error for invalid schedule, got nil")
	}
}

func TestStart_ValidSchedule(t *testing.T) {
	s := cron.New()
	if err := s.Start([]domain.Job{validJob}, func(domain.Job) error { return nil }); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s.Stop()
}

func TestStop_WithoutStart(t *testing.T) {
	cron.New().Stop() // must not panic
}

// robfig/cron rounds @every intervals under 1s up to 1s, so the minimum
// meaningful interval for execution tests is 1s.
const tick = "@every 1s"
const tickWait = 1500 * time.Millisecond

func TestStart_RunnerIsCalled(t *testing.T) {
	s := cron.New()
	var count atomic.Int32
	jobs := []domain.Job{{Name: "job", Path: "/tmp", Schedule: tick}}
	if err := s.Start(jobs, func(domain.Job) error { count.Add(1); return nil }); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	time.Sleep(tickWait)
	s.Stop()
	if count.Load() == 0 {
		t.Error("runner was never called")
	}
}

func TestReload_ReplacesJobs(t *testing.T) {
	s := cron.New()
	var oldCount, newCount atomic.Int32

	oldJobs := []domain.Job{{Name: "old", Path: "/tmp", Schedule: tick}}
	if err := s.Start(oldJobs, func(domain.Job) error { oldCount.Add(1); return nil }); err != nil {
		t.Fatalf("start: %v", err)
	}

	newJobs := []domain.Job{{Name: "new", Path: "/tmp", Schedule: tick}}
	if err := s.Reload(newJobs, func(domain.Job) error { newCount.Add(1); return nil }); err != nil {
		t.Fatalf("reload: %v", err)
	}

	time.Sleep(tickWait)
	s.Stop()

	if newCount.Load() == 0 {
		t.Error("new runner was never called after reload")
	}
}
