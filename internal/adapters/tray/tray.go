package tray

import (
	_ "embed"
	"fmt"

	"a-rc/internal/app"

	"github.com/getlantern/systray"
)

//go:embed icon.png
var iconPNG []byte

// TrayApp implements a macOS menu bar tray app.
type TrayApp struct {
	configRepo app.ConfigRepository
	configPath *string
	scheduler  app.JobScheduler
	archiveSvc *app.ArchiveService
}

func New(configRepo app.ConfigRepository, configPath *string, scheduler app.JobScheduler, archiveSvc *app.ArchiveService) *TrayApp {
	return &TrayApp{
		configRepo: configRepo,
		configPath: configPath,
		scheduler:  scheduler,
		archiveSvc: archiveSvc,
	}
}

// Run starts the tray app. Blocks until the user quits.
func (a *TrayApp) Run() {
	systray.Run(a.onReady, a.onExit)
}

func (a *TrayApp) onReady() {
	systray.SetTemplateIcon(iconPNG, iconPNG)
	systray.SetTooltip("a-rc archiver")

	cfg, err := a.configRepo.Load(*a.configPath)
	if err != nil {
		item := systray.AddMenuItem(fmt.Sprintf("Error: %s", err), "")
		item.Disable()
	} else if len(cfg.Jobs) == 0 {
		item := systray.AddMenuItem("No jobs configured", "")
		item.Disable()
	} else {
		for _, j := range cfg.Jobs {
			item := systray.AddMenuItem(fmt.Sprintf("%s  %s  %s", j.Name, j.Schedule, j.Path), "")
			item.Disable()
		}
		_ = a.scheduler.Start(cfg.Jobs, a.archiveSvc.RunJob)
	}

	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "Stop a-rc")

	go func() {
		<-quit.ClickedCh
		systray.Quit()
	}()
}

func (a *TrayApp) onExit() {
	a.scheduler.Stop()
}
