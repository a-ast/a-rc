package gdrive

import (
	"errors"
	"fmt"

	"a-rc/internal/core"
)

// Adapter implements core.Uploader.
// It reads GDrive config from the config file at upload time via ConfigRepository.
// Phase 2: real Google Drive upload using credentials_file + token_file.
type Adapter struct {
	configRepo core.ConfigRepository
	configPath *string // pointer to cmd.ConfigPath — valid after cobra flag parsing
}

func New(repo core.ConfigRepository, configPath *string) *Adapter {
	return &Adapter{configRepo: repo, configPath: configPath}
}

func (a *Adapter) Upload(localPath string) error {
	cfg, err := a.configRepo.Load(*a.configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	g := cfg.GDrive
	if g.ServiceAccountFile == "" {
		return errors.New("gdrive.service_account_file is not set in config")
	}
	if g.Folder == "" {
		return errors.New("gdrive.folder is not set in config")
	}
	// Phase 2: authenticate with service account key g.ServiceAccountFile and upload localPath to g.Folder.
	return fmt.Errorf("gdrive: not implemented (Phase 2) — would upload %q to folder %q", localPath, g.Folder)
}
