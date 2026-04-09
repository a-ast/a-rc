package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := svc.Config.Load(ConfigPath)
			if err != nil {
				return err
			}

			if len(cfg.Jobs) == 0 {
				fmt.Println("no jobs configured")
				return nil
			}
			fmt.Printf("%-30s  %-20s  %s\n", "NAME", "SCHEDULE", "PATH")
			for _, j := range cfg.Jobs {
				fmt.Printf("%-30s  %-20s  %s\n", j.Name, j.Schedule, j.Path)
			}
			return nil
		},
	}
}
