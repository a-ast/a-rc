# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

A-rc is a macOS archiver that runs as a menu bar tray app. It archives configured folders on a cron schedule and uploads them to Google Drive.

## Commands

```bash
go run main.go        # Run
go build -o a-rc      # Build
go test ./...         # Test
go vet ./...          # Vet
go fmt ./...          # Format
go run ./tools/genicon  # Regenerate internal/adapters/tray/icon.png
```

CLI usage:
```bash
a-rc                  # Launch the menu bar tray app (default)
a-rc run <path>       # Run one job immediately (identified by its configured path)
a-rc list             # Show configured jobs
```

## Architecture

Hexagonal architecture: core has no external imports; adapters implement ports; `main.go` wires everything.

```
main.go                        # DI wiring only
cmd/                           # Input adapter (cobra CLI)
internal/domain/
  domain.go                    # Config, GDriveConfig, Job — pure business types, no imports
internal/app/
  ports.go                     # Output port interfaces: ConfigRepository, Archiver, Uploader, JobScheduler
  archive_service.go           # ArchiveService: RunJob use case (archive + upload + cleanup)
internal/adapters/
  yaml/loader.go    → ConfigRepository
  cron/scheduler.go → JobScheduler    (robfig/cron)
  archiver/zip.go   → Archiver        (zip via WalkDir, staged in os.MkdirTemp)
  gdrive/drive.go   → Uploader        (OAuth2 Desktop app flow)
  tray/tray.go      → TrayApp         (systray menu bar app)
tools/genicon/        # Icon generator for tray/icon.png
```

## Config

Default location: `~/.config/a-rc/config.yaml`

```yaml
log_dir: ~/Library/Logs/a-rc

gdrive:
  credentials_file: ~/.config/a-rc/gdrive-credentials.json
  token_file: ~/.config/a-rc/gdrive-token.json
  folder: Backups

jobs:
  - path: ~/Documents
    schedule: "0 2 * * *"   # 5-field cron expression
```

## Google Drive

Uses OAuth2 (Desktop app flow). 
On first upload the browser opens for authorization. 
The token is saved to `token_file` and refreshed automatically. 
The `gdrive.Adapter` holds a pointer to `cmd.ConfigPath` and loads 
the config at upload time.

## Tray app design

`a-rc` (no subcommand) launches a macOS menu bar tray app. 
On startup it loads the config, displays each job (`schedule  path`) as a disabled menu item, and starts the cron scheduler (`robfig/cron`). A Quit item stops the scheduler and exits.

The menu bar icon (`tray/icon.png`) is a 22x22 template PNG (black on transparent) embedded via `//go:embed`. macOS inverts it automatically for dark mode.

## Archiving design

Archives are written to `os.MkdirTemp` (no persistent local copy), uploaded to the configured GDrive folder (overwriting any existing file with the same name), then the temp dir is deleted by `ArchiveService`.
