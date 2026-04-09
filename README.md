# A-rc

A macOS archiver that lives in the menu bar.
It archives configured folders on a cron schedule and uploads them to Google Drive.

## Install

```bash
git clone https://github.com/a-ast/a-rc
cd a-rc
go build -o a-rc
mv a-rc /usr/local/bin/a-rc
```

## Configuration

Default config location: `~/.config/a-rc/config.yaml`

```yaml
log_dir: ~/Library/Logs/a-rc

gdrive:
  credentials_file: ~/.config/a-rc/gdrive-credentials.json
  token_file: ~/.config/a-rc/gdrive-token.json
  folder: Backups

jobs:
  - path: ~/Documents
    schedule: "0 2 * * *"    # daily at 2am
  - path: ~/Projects/a-rc
    schedule: "0 */6 * * *"  # every 6 hours
```

`schedule` uses standard 5-field cron syntax (minute hour day month weekday).

### Google Drive

#### Getting credentials

1. Go to [console.cloud.google.com](https://console.cloud.google.com) and create a new project (e.g. `a-rc`).
2. Enable the **Google Drive API**: APIs & Services - Enable APIs - search "Google Drive API".
3. Create credentials: APIs & Services - Credentials - Create Credentials - **OAuth client ID** - Application type: **Desktop app**. Download the JSON and save it as `credentials_file`.
4. Configure the consent screen: OAuth consent screen - External - add your Google account as a test user.

#### Authorization (first run)

`token_file` is created automatically — you do not need to obtain it manually.

On the first upload, a-rc opens a browser for authorization:

```bash
a-rc run ~/some/path
# Browser opens → click Allow → paste the authorization code into the terminal
```

a-rc exchanges the code for a token and saves it to `token_file`. Every subsequent run uses that file silently and refreshes it automatically when it expires.

#### Upload

Each archive is named `<folder>.zip` (e.g. `Documents.zip`).
If a file with the same name already exists in the Drive folder it is overwritten in place, no duplicates.

## Usage

```bash
# Launch the menu bar tray app
a-rc

# Run a single job immediately
a-rc run ~/Documents

# Show configured jobs
a-rc list
```

Use `--config` to point at a non-default config file:

```bash
a-rc --config /path/to/config.yaml
```

## How it works

`a-rc` runs as a macOS menu bar app. The icon in the top-right bar shows a bow-and-arrow. Clicking it opens a menu listing each configured job (path + schedule) and a Quit item.

All job schedules are driven inside the app process by [robfig/cron](https://github.com/robfig/cron). When a job fires, a-rc zips the folder into a temporary directory, uploads the archive to the configured Drive folder (overwriting any previous version), then deletes the local copy.
