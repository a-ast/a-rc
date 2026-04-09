package yaml

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"a-rc/internal/core"
	goyaml "gopkg.in/yaml.v3"
)

// Loader implements core.ConfigRepository using a YAML file.
type Loader struct{}

func New() *Loader { return &Loader{} }

func (l *Loader) Load(path string) (*core.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	var cfg core.Config
	if err := goyaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	expandTilde(&cfg)
	return &cfg, nil
}

// expandTilde replaces leading ~ in all path fields with the user's home directory.
func expandTilde(cfg *core.Config) {
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
	cfg.GDrive.ServiceAccountFile = expand(cfg.GDrive.ServiceAccountFile)
	for i := range cfg.Jobs {
		cfg.Jobs[i].Path = expand(cfg.Jobs[i].Path)
	}
}
