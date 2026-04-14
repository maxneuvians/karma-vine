## ADDED Requirements

### Requirement: buildView dispatches to renderHelpPanel when showHelpPanel is true
When `m.showHelpPanel == true` and `m.screenMode == ScreenNormal`, `buildView` SHALL return `renderHelpPanel(m)` immediately after the `ScreenCombat` and `ScreenInventory` early-return checks.

#### Scenario: buildView returns help panel for showHelpPanel
- **WHEN** `m.showHelpPanel == true` and `m.screenMode == ScreenNormal`
- **THEN** `buildView(m)` equals `renderHelpPanel(m)`

### Requirement: buildView uses only one chrome row (HUD; no key-bar)
`buildView` in `ScreenNormal` SHALL compose the view as `mapView + "\n" + renderHUD(m)` (one chrome row). The `renderKeyBar` function SHALL be removed. Map height SHALL be `m.viewportH - 1`.

#### Scenario: buildView output does not contain old key-bar text
- **WHEN** `buildView` is called in `ScreenNormal` with `m.mode == ModeLocal`
- **THEN** the output does NOT contain the string `↑↓←→/wasd move`

#### Scenario: buildView still contains HUD content
- **WHEN** `buildView` is called in `ScreenNormal`
- **THEN** the output contains `? help`
