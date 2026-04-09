package core

import "fmt"

// ScheduleService manages the lifecycle of the background daemon via the OS process manager.
type ScheduleService struct {
	configRepo ConfigRepository
	procMgr    ProcessManager
}

func NewScheduleService(r ConfigRepository, p ProcessManager) *ScheduleService {
	return &ScheduleService{configRepo: r, procMgr: p}
}

// Install validates the config, then installs and starts the daemon.
// If the daemon is already running it signals it to reload instead.
func (s *ScheduleService) Install(configPath, binaryPath string) error {
	cfg, err := s.configRepo.Load(configPath)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	return s.procMgr.Install(binaryPath, configPath, cfg.LogDir)
}

// Uninstall stops and removes the daemon.
func (s *ScheduleService) Uninstall() error {
	return s.procMgr.Uninstall()
}

// Status returns whether the daemon is running and its PID.
func (s *ScheduleService) Status() (running bool, pid int, err error) {
	return s.procMgr.Status()
}
