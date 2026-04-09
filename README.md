# A-rc

A macOS CLI archiver.
It archives configured folders on a cron schedule, and uploads the result to Google Drive. Runs as a persistent
background daemon managed by launchd.

## Install

```bash
git clone <repo>
cd a-rc
go build -o a-rc
mv a-rc /usr/local/bin/a-rc
```

## Configuration

Default config location: `~/.config/a-rc/config.yaml`

```yaml
log_dir: ~/Library/Logs/a-rc

jobs:
  - path: ~/Projects/a-rc
    schedule: "0 */6 * * *"  # every 6 hours
```

`schedule` uses standard 5-field crontab syntax (minute hour day month weekday).

### Google Drive

Add a `gdrive` section to `config.yaml`:

```yaml
gdrive:
  service_account_file: ~/.config/a-rc/gdrive-service-account.json
  folder: Backups   # Drive folder name or ID shared with the service account
```

#### Setting up a service account

1. Go to [console.cloud.google.com](https://console.cloud.google.com) and create a new project (e.g. `a-rc`).
2. Enable the **Google Drive API**: APIs & Services - Enable APIs - search "Google Drive API".
3. Create a service account: IAM & Admin - Service Accounts - Create. No special roles needed.
4. Generate a key: open the service account - Keys - Add Key - JSON. Save the downloaded file as your `service_account_file` path.
5. In Google Drive, create a folder (e.g. `Backups`) and share it with the service account email (e.g. `a-rc@your-project.iam.gserviceaccount.com`) as **Editor**.

No browser login or token refresh - the service account key file is the only credential needed.

## Usage

```bash
# Register the daemon with launchd and start it
a-rc schedule

# After editing config.yaml, reload the daemon
a-rc schedule

# Show configured jobs and daemon status
a-rc list

# Run a single job immediately
a-rc run ~/Projects/a-rc

# Stop and remove the daemon
a-rc unschedule
```

Use `--config` to point at a non-default config file:

```bash
a-rc --config /path/to/config.yaml schedule
```

## How it works

`a-rc schedule` installs a single launchd `LaunchAgent` (`com.a-rc.daemon`) that keeps the daemon process alive with
`KeepAlive=true`. All job schedules are driven inside the daemon by [robfig/cron](https://github.com/robfig/cron) — no
per-job plists.

Re-running `a-rc schedule` on a running daemon sends `SIGHUP` to reload the config without restarting. `SIGTERM`
triggers a clean stop.

Archives are written to a temporary directory, uploaded to the configured GDrive folder, then deleted locally.

Logs are written to `log_dir` as defined in `config.yaml`:

- `a-rc.log` - stdout
- `a-rc.err` - stderr

## Notes

- The plist hardcodes the binary path at `a-rc schedule` time. If you move or rebuild the binary to a different
  location, re-run `a-rc schedule`.
- Google Drive upload is not yet implemented (Phase 2).
