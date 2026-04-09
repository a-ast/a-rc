package core

// Config holds the full application configuration loaded from the config file.
type Config struct {
	LogDir string       `yaml:"log_dir"`
	GDrive GDriveConfig `yaml:"gdrive"`
	Jobs   []Job        `yaml:"jobs"`
}

// GDriveConfig holds Google Drive settings.
type GDriveConfig struct {
	ServiceAccountFile string `yaml:"service_account_file"`
	Folder             string `yaml:"folder"`
}

// Job describes a single archive task.
type Job struct {
	Path     string `yaml:"path"`
	Schedule string `yaml:"schedule"` // 5-field cron expression
}
