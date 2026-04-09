package core

// Output ports — defined by the core, implemented by adapters.

// ConfigRepository loads and parses application configuration.
type ConfigRepository interface {
	Load(path string) (*Config, error)
}

// Archiver creates an archive from a job's source path and returns the archive file path.
type Archiver interface {
	Archive(job Job) (archivePath string, err error)
}

// Uploader transfers a local archive file to the configured GDrive folder.
// The destination folder is resolved by the adapter from environment variables.
type Uploader interface {
	Upload(localPath string) error
}

// ProcessManager manages the lifecycle of the background daemon process via the OS scheduler.
type ProcessManager interface {
	Install(binaryPath, configPath, logDir string) error
	Uninstall() error
	Status() (running bool, pid int, err error)
}

// JobScheduler runs jobs on their configured cron schedules inside the daemon process.
type JobScheduler interface {
	Start(jobs []Job, runner func(Job) error) error
	Stop()
	Reload(jobs []Job, runner func(Job) error) error
}
