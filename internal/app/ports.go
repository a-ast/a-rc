package app

import "a-rc/internal/domain"

// Output ports — defined by the application, implemented by adapters.

// ConfigRepository loads and parses application configuration.
type ConfigRepository interface {
	Load(path string) (*domain.Config, error)
}

// Archiver creates an archive from a job's source path and returns the archive file path.
type Archiver interface {
	Archive(job domain.Job) (archivePath string, err error)
}

// Uploader transfers a local archive file to the configured GDrive folder.
type Uploader interface {
	Upload(localPath string) error
}

// JobScheduler runs jobs on their configured cron schedules.
type JobScheduler interface {
	Start(jobs []domain.Job, runner func(domain.Job) error) error
	Stop()
	Reload(jobs []domain.Job, runner func(domain.Job) error) error
}
