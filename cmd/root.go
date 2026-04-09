package cmd

import (
	"os"
	"path/filepath"

	"a-rc/internal/adapters/gdrive"
	"a-rc/internal/adapters/tray"
	"a-rc/internal/app"

	"github.com/spf13/cobra"
)

// ConfigPath is exported so adapters constructed in main.go can hold a pointer to it
// and read the correct value after cobra parses the --config flag.
var ConfigPath string

// Services is the container passed from main to all subcommands.
type Services struct {
	Archive *app.ArchiveService
	Config  app.ConfigRepository
	Tray    *tray.TrayApp
	GDrive  *gdrive.Drive
}

var svc *Services

func NewRootCmd(services *Services) *cobra.Command {
	svc = services

	root := &cobra.Command{
		Use:   "a-rc",
		Short: "A macOS archiver with menu bar control",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc.Tray.Run()
			return nil
		},
	}

	defaultConfig := ""
	if home, err := os.UserHomeDir(); err == nil {
		defaultConfig = filepath.Join(home, ".config", "a-rc", "config.yaml")
	}
	root.PersistentFlags().StringVar(&ConfigPath, "config", defaultConfig, "config file path")

	root.AddCommand(
		newRunCmd(),
		newListCmd(),
		newAuthorizeGDriveCmd(),
	)
	return root
}
