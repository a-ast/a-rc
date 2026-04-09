package main

import (
	"fmt"
	"os"

	"a-rc/cmd"
	"a-rc/internal/adapters/archiver"
	adptcron "a-rc/internal/adapters/cron"
	"a-rc/internal/adapters/gdrive"
	"a-rc/internal/adapters/tray"
	adptyaml "a-rc/internal/adapters/yaml"
	"a-rc/internal/app"
)

func main() {
	// Output adapters.
	configRepo := adptyaml.New()
	jobScheduler := adptcron.New()
	arc := archiver.New()
	upl := gdrive.New(configRepo, &cmd.ConfigPath)

	// Core services.
	archiveSvc := app.NewArchiveService(arc, upl)

	// Tray adapter.
	trayApp := tray.New(configRepo, &cmd.ConfigPath, jobScheduler, archiveSvc)

	// Input adapter (CLI).
	root := cmd.NewRootCmd(&cmd.Services{
		Archive: archiveSvc,
		Config:  configRepo,
		Tray:    trayApp,
		GDrive:  upl,
	})

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
