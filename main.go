package main

import (
	"fmt"
	"os"

	"a-rc/cmd"
	"a-rc/internal/adapters/archiver"
	adptcron "a-rc/internal/adapters/cron"
	"a-rc/internal/adapters/gdrive"
	"a-rc/internal/adapters/launchd"
	adptyaml "a-rc/internal/adapters/yaml"
	"a-rc/internal/core"
)

func main() {
	// Output adapters.
	configRepo := adptyaml.New()
	procMgr := launchd.New()
	jobScheduler := adptcron.New()
	arc := archiver.New()
	upl := gdrive.New(configRepo, &cmd.ConfigPath)

	// Core services.
	archiveSvc := core.NewArchiveService(arc, upl)
	scheduleSvc := core.NewScheduleService(configRepo, procMgr)
	daemonSvc := core.NewDaemonService(configRepo, jobScheduler, archiveSvc)

	// Input adapter (CLI).
	root := cmd.NewRootCmd(&cmd.Services{
		Archive:  archiveSvc,
		Schedule: scheduleSvc,
		Daemon:   daemonSvc,
		Config:   configRepo,
	})

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
