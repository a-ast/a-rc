package tray

import (
	_ "embed"
	"fmt"

	"a-rc/internal/app"

	"github.com/getlantern/systray"
)

//go:embed bow.png
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
func (t *TrayApp) Run() {
	systray.Run(t.onReady, t.onExit)
}

func (t *TrayApp) onReady() {
	systray.SetTemplateIcon(iconPNG, iconPNG)
	systray.SetTooltip("a-rc archiver")

	cfg, err := t.configRepo.Load(*t.configPath)
	if err != nil {
		item := systray.AddMenuItem(fmt.Sprintf("Error: %s", err), "")
		item.Disable()
	} else if len(cfg.Jobs) == 0 {
		item := systray.AddMenuItem("No jobs configured", "")
		item.Disable()
	} else {
		for _, j := range cfg.Jobs {
			item := systray.AddMenuItem(fmt.Sprintf("%s", j.Name), "")
			item.Disable()
		}
		if err := t.scheduler.Start(cfg.Jobs, t.archiveSvc.RunJob); err != nil {
			item := systray.AddMenuItem(fmt.Sprintf("Scheduler error: %s", err), "")
			item.Disable()
		}
	}

	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "Stop a-rc")

	go func() {
		<-quit.ClickedCh
		systray.Quit()
	}()
}

func (t *TrayApp) onExit() {
	t.scheduler.Stop()
}
