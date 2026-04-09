package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newScheduleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schedule",
		Short: "Install and start the background archiver daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			binary, err := os.Executable()
			if err != nil {
				return fmt.Errorf("resolving binary path: %w", err)
			}
			if err := svc.Schedule.Install(ConfigPath, binary); err != nil {
				return err
			}
			fmt.Println("daemon scheduled")
			return nil
		},
	}
}
