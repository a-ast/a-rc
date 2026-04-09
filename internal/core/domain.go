package core

// Config holds the full application configuration loaded from the config file.
type Config struct {
	LogDir string       `yaml:"log_dir"`
	GDrive GDriveConfig `yaml:"gdrive"`
	Jobs   []Job        `yaml:"jobs"`
}

// GDriveConfig holds Google Drive settings.
type GDriveConfig struct {
	CredentialsFile string `yaml:"credentials_file"` // OAuth2 client secret JSON from Google Console
	TokenFile       string `yaml:"token_file"`       // cached OAuth2 token (written on first auth)
	Folder          string `yaml:"folder"`
}

// Job describes a single archive task.
type Job struct {
	Path     string `yaml:"path"`
	Schedule string `yaml:"schedule"` // 5-field cron expression
}
