package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUnscheduleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unschedule",
		Short: "Stop and remove the background archiver daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Schedule.Uninstall(); err != nil {
				return err
			}
			fmt.Println("daemon unscheduled")
			return nil
		},
	}
}
