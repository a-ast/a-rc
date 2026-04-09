package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func newDaemonCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "daemon",
		Short:  "Run the archiver daemon (managed by launchd)",
		Hidden: true, // invoked by launchd, not directly by users
		RunE: func(cmd *cobra.Command, args []string) error {
			stopCh := make(chan struct{})
			reloadCh := make(chan struct{}, 1)

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGHUP)
			go func() {
				for sig := range sigCh {
					switch sig {
					case syscall.SIGHUP:
						fmt.Println("a-rc daemon: reloading config")
						select {
						case reloadCh <- struct{}{}:
						default:
						}
					case syscall.SIGTERM:
						fmt.Println("a-rc daemon: stopping")
						close(stopCh)
						return
					}
				}
			}()

			return svc.Daemon.Run(ConfigPath, reloadCh, stopCh)
		},
	}
}
