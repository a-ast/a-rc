### Requirement: Backed-up folders sub-menu
The tray menu SHALL display a top-level item labelled "Backed-up folders" that contains each configured job name as a non-interactive sub-menu item.

#### Scenario: Jobs are configured
- **WHEN** one or more jobs are configured
- **THEN** the tray menu shows a "Backed-up folders" item with a sub-menu listing each job name

#### Scenario: Sub-menu items are non-interactive
- **WHEN** the user opens the "Backed-up folders" sub-menu
- **THEN** each job name item is disabled and cannot be clicked

#### Scenario: No jobs configured
- **WHEN** no jobs are configured
- **THEN** "Backed-up folders" sub-menu is not shown; the existing "No jobs configured" message is shown instead
