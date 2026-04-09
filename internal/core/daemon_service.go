package core

import "fmt"

// DaemonService runs all scheduled jobs inside a long-lived process.
// It is driven by signals from the input adapter (SIGHUP to reload, SIGTERM to stop).
type DaemonService struct {
	configRepo   ConfigRepository
	jobScheduler JobScheduler
	archiveSvc   *ArchiveService
}

func NewDaemonService(r ConfigRepository, js JobScheduler, a *ArchiveService) *DaemonService {
	return &DaemonService{configRepo: r, jobScheduler: js, archiveSvc: a}
}

// Run loads the config, starts the scheduler, then blocks until stopCh is closed.
// When reloadCh receives, the config is reloaded and the scheduler is updated.
func (s *DaemonService) Run(configPath string, reloadCh <-chan struct{}, stopCh <-chan struct{}) error {
	cfg, err := s.configRepo.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if err := s.jobScheduler.Start(cfg.Jobs, s.archiveSvc.RunJob); err != nil {
		return fmt.Errorf("starting scheduler: %w", err)
	}
	defer s.jobScheduler.Stop()

	for {
		select {
		case <-stopCh:
			return nil
		case <-reloadCh:
			cfg, err = s.configRepo.Load(configPath)
			if err != nil {
				// Log and continue with the previous config.
				continue
			}
			if err := s.jobScheduler.Reload(cfg.Jobs, s.archiveSvc.RunJob); err != nil {
				continue
			}
		}
	}
}
