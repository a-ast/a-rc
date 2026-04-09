# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

A-rc is a macOS CLI archiver. It archives configured folders and uploads them to Google Drive on a cron schedule, managed as a persistent background daemon via launchd.

## Commands

```bash
go run main.go        # Run
go build -o a-rc      # Build
go test ./...         # Test
go vet ./...          # Vet
go fmt ./...          # Format
```

CLI usage:
```bash
a-rc schedule         # Install and start the launchd daemon
a-rc unschedule       # Stop and remove the launchd daemon
a-rc run <path>       # Run one job immediately (identified by its configured path)
a-rc list             # Show configured jobs and daemon status
# daemon subcommand is hidden — invoked by launchd only
```

## Architecture

Hexagonal architecture: core has no external imports; adapters implement ports; `main.go` wires everything.

```
main.go                        # DI wiring only
cmd/                           # Input adapter (cobra CLI)
internal/core/
  domain.go                    # Config, Job types
  ports.go                     # Interfaces: ConfigRepository, Archiver, Uploader,
                               #   ProcessManager, JobScheduler
  archive_service.go           # RunJob: archive + upload
  schedule_service.go          # Install/Uninstall/Status via ProcessManager
  daemon_service.go            # Blocking cron loop; driven by reloadCh/stopCh channels
internal/adapters/
  yaml/      → ConfigRepository
  launchd/   → ProcessManager  (plist + launchctl bootstrap/bootout)
  cron/      → JobScheduler    (robfig/cron)
  archiver/  → Archiver        (stub — Phase 2)
  gdrive/    → Uploader        (stub — Phase 2)
```

## Config

Default location: `~/.config/a-rc/config.yaml`

```yaml
log_dir: ~/Library/Logs/a-rc

jobs:
  - path: ~/Documents
    schedule: "0 2 * * *"   # 5-field cron expression
```

## Google Drive (Phase 2)

Configured in `config.yaml` under the `gdrive` key:

```yaml
gdrive:
  service_account_file: ~/.config/a-rc/gdrive-service-account.json
  folder: Backups
```

The `gdrive.Adapter` holds a pointer to `cmd.ConfigPath` and loads the config at upload time, so it always reflects the current `--config` flag value.

## Scheduling design

A single launchd `LaunchAgent` (`com.a-rc.daemon`) keeps `a-rc daemon` alive with `KeepAlive=true`. All cron schedules are driven inside the daemon process by `robfig/cron` — no per-job plists, no cron-expression-to-launchd conversion. `SIGHUP` triggers a config reload; `SIGTERM` triggers a clean stop.

Re-running `a-rc schedule` after editing the config sends `SIGHUP` to the running daemon instead of reinstalling.

## Archiving design (Phase 2)

Archives are written to `os.MkdirTemp` (no persistent local copy), uploaded to the configured GDrive folder, then the temp file is deleted by `ArchiveService`.
