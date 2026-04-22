## 1. Tray menu update

- [ ] 1.1 In `internal/adapters/tray/tray.go` `onReady()`, replace the flat job loop with a "Backed-up folders" parent item and `AddSubMenuItem` calls for each job
- [ ] 1.2 Verify the parent item is NOT disabled (so the sub-menu is accessible on macOS)
- [ ] 1.3 Build and launch the app; confirm "Backed-up folders" appears as a sub-menu in the menu bar
