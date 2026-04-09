package cmd

import (
	"os"
	"path/filepath"

	"a-rc/internal/core"
	"github.com/spf13/cobra"
)

// ConfigPath is exported so adapters constructed in main.go can hold a pointer to it
// and read the correct value after cobra parses the --config flag.
var ConfigPath string

// Services is the container passed from main to all subcommands.
type Services struct {
	Archive  *core.ArchiveService
	Schedule *core.ScheduleService
	Daemon   *core.DaemonService
	Config   core.ConfigRepository
}

var svc *Services

func NewRootCmd(services *Services) *cobra.Command {
	svc = services

	root := &cobra.Command{
		Use:   "a-rc",
		Short: "A command line archiver",
	}

	home, _ := os.UserHomeDir()
	defaultConfig := filepath.Join(home, ".config", "a-rc", "config.yaml")
	root.PersistentFlags().StringVar(&ConfigPath, "config", defaultConfig, "config file path")

	root.AddCommand(
		newScheduleCmd(),
		newUnscheduleCmd(),
		newDaemonCmd(),
		newRunCmd(),
		newListCmd(),
	)
	return root
}
