## Context

The tray app is built on `github.com/getlantern/systray` v1.2.2. Menu items are created in `onReady()` in `internal/adapters/tray/tray.go`. Currently, each job is added as a top-level disabled item via `systray.AddMenuItem`. The library supports nested items via `(*MenuItem).AddSubMenuItem(title, tooltip string) *MenuItem`.

## Goals / Non-Goals

**Goals:**
- Group job names under a single "Backed-up folders" parent item in the tray menu.

**Non-Goals:**
- Making job items interactive (run-on-click, etc.)
- Changing what information is displayed per job
- Any changes outside `tray.go`

## Decisions

**Use `AddSubMenuItem` on a parent `*MenuItem`**

The systray library's `(*MenuItem).AddSubMenuItem` creates a nested sub-menu item. The parent item must remain enabled (not `.Disable()`'d) — on macOS, disabling a parent hides the sub-menu. Children are still disabled to keep them informational.

No alternatives considered; `AddSubMenuItem` is the only sub-menu API in the library.

## Risks / Trade-offs

- **macOS menu bar rendering** — sub-menu arrow and hover behavior are handled by the OS; no custom rendering risk.
- Enabling the parent item means clicking "Backed-up folders" directly does nothing visible (no `ClickedCh` listener). This is acceptable since the item is purely a container.
