package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run <path>",
		Short: "Run a single archive job immediately (identified by its configured path)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := filepath.Clean(args[0])
			cfg, err := svc.Config.Load(ConfigPath)
			if err != nil {
				return err
			}
			for _, job := range cfg.Jobs {
				if filepath.Clean(job.Path) == target {
					fmt.Printf("running job for path %q\n", job.Path)
					return svc.Archive.RunJob(job)
				}
			}
			return fmt.Errorf("no job configured for path %q", target)
		},
	}
}
