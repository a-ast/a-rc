package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run <name>",
		Short: "Run a single archive job immediately (identified by its name)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			cfg, err := svc.Config.Load(ConfigPath)
			if err != nil {
				return err
			}
			for _, job := range cfg.Jobs {
				if job.Name == target {
					fmt.Printf("running job %q\n", job.Name)
					return svc.Archive.RunJob(job)
				}
			}
			return fmt.Errorf("no job configured with name %q", target)
		},
	}
}
