package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured jobs and daemon status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := svc.Config.Load(ConfigPath)
			if err != nil {
				return err
			}

			running, pid, _ := svc.Schedule.Status()
			status := "stopped"
			if running && pid > 0 {
				status = fmt.Sprintf("running (pid %d)", pid)
			} else if running {
				status = "registered (not running)"
			}
			fmt.Printf("daemon: %s\n\n", status)

			if len(cfg.Jobs) == 0 {
				fmt.Println("no jobs configured")
				return nil
			}
			fmt.Printf("%-20s  %s\n", "SCHEDULE", "PATH")
			for _, j := range cfg.Jobs {
				fmt.Printf("%-20s  %s\n", j.Schedule, j.Path)
			}
			return nil
		},
	}
}
