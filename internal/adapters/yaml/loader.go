package yaml

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"a-rc/internal/domain"

	robfigcron "github.com/robfig/cron/v3"
	goyaml "gopkg.in/yaml.v3"
)

// Loader implements domain.ConfigRepository using a YAML file.
type Loader struct{}

func New() *Loader { return &Loader{} }

// rawConfig mirrors domain.Config but uses a map for jobs so names come from YAML keys.
type rawConfig struct {
	LogDir string                   `yaml:"log_dir"`
	GDrive domain.GoogleDriveConfig `yaml:"gdrive"`
	Jobs   map[string]rawJob        `yaml:"jobs"`
}

type rawJob struct {
	Path     string `yaml:"path"`
	Schedule string `yaml:"schedule"`
}

func (l *Loader) Load(path string) (*domain.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer func() { _ = f.Close() }()

	var raw rawConfig
	if err := goyaml.NewDecoder(f).Decode(&raw); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	cfg := &domain.Config{
		LogDir: raw.LogDir,
		GDrive: raw.GDrive,
	}
	for name, j := range raw.Jobs {
		cfg.Jobs = append(cfg.Jobs, domain.Job{Name: name, Path: j.Path, Schedule: j.Schedule})
	}

	expandTilde(cfg)
	if err := validate(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func validate(cfg *domain.Config) error {
	parser := robfigcron.NewParser(robfigcron.Minute | robfigcron.Hour | robfigcron.Dom | robfigcron.Month | robfigcron.Dow)
	for _, j := range cfg.Jobs {
		if j.Name == "" {
			return fmt.Errorf("job is missing a name")
		}
		if j.Path == "" {
			return fmt.Errorf("job %q: path is required", j.Name)
		}
		if j.Schedule == "" {
			return fmt.Errorf("job %q: schedule is required", j.Name)
		}
		if _, err := os.Stat(j.Path); err != nil {
			return fmt.Errorf("job %q: path %q does not exist", j.Name, j.Path)
		}
		if _, err := parser.Parse(j.Schedule); err != nil {
			return fmt.Errorf("job %q: invalid schedule %q: %w", j.Name, j.Schedule, err)
		}
	}
	return nil
}

// expandTilde replaces leading ~ in all path fields with the user's home directory.
func expandTilde(cfg *domain.Config) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	expand := func(p string) string {
		if strings.HasPrefix(p, "~/") {
			return filepath.Join(home, p[2:])
		}
		return p
	}
	cfg.LogDir = expand(cfg.LogDir)
	cfg.GDrive.CredentialsFile = expand(cfg.GDrive.CredentialsFile)
	cfg.GDrive.TokenFile = expand(cfg.GDrive.TokenFile)
	for i := range cfg.Jobs {
		cfg.Jobs[i].Path = expand(cfg.Jobs[i].Path)
	}
}
