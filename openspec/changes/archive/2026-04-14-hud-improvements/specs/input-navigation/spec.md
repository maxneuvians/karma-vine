## ADDED Requirements

### Requirement: ? key toggles showHelpPanel in ScreenNormal
When `m.screenMode == ScreenNormal`, pressing `?` SHALL toggle `m.showHelpPanel` between `true` and `false`. The existing binding of `?` to `showSidebar` SHALL be removed and `?` SHALL exclusively control the help panel. The sidebar SHALL be toggled by `\` (backslash) instead.

#### Scenario: ? opens help panel
- **WHEN** `m.screenMode == ScreenNormal`, `m.showHelpPanel == false`, and `?` is pressed
- **THEN** `m.showHelpPanel == true` and `m.showSidebar` is unchanged

#### Scenario: \ toggles sidebar
- **WHEN** `m.screenMode == ScreenNormal` and `\` is pressed
- **THEN** `m.showSidebar` is toggled

#### Scenario: ? is suppressed in ScreenInventory and ScreenCombat
- **WHEN** `m.screenMode == ScreenInventory` or `m.screenMode == ScreenCombat` and `?` is pressed
- **THEN** `m.showHelpPanel` is unchanged
