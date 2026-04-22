## Why

The tray menu lists backed-up folder names as flat top-level items, mixing informational content with actionable items (Quit). Moving them into a dedicated sub-menu reduces visual clutter and makes the menu easier to scan.

## What Changes

- The flat list of job names in the tray menu is replaced by a single top-level item "Backed-up folders" that reveals a sub-menu listing each job name on hover.

## Capabilities

### New Capabilities

- `tray-folders-submenu`: The tray menu shows backed-up folder names grouped under a collapsible "Backed-up folders" sub-menu item rather than as flat top-level entries.

### Modified Capabilities

<!-- No existing specs to modify. -->

## Impact

- `internal/adapters/tray/tray.go`: `onReady()` — replace flat `systray.AddMenuItem` loop with a parent item and `AddSubMenuItem` calls.
- No API, config, or dependency changes.
