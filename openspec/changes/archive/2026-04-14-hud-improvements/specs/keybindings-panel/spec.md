## ADDED Requirements

### Requirement: Model tracks help panel visibility
The system SHALL add `showHelpPanel bool` to `Model`. `NewModel()` SHALL initialise it to `false`.

#### Scenario: NewModel initialises showHelpPanel to false
- **WHEN** `NewModel()` is called
- **THEN** `m.showHelpPanel == false`

### Requirement: ? key toggles the help panel in ScreenNormal
When `m.screenMode == ScreenNormal`, pressing `?` SHALL toggle `m.showHelpPanel`. When `m.screenMode == ScreenInventory`, `ScreenCombat`, or any non-normal screen, `?` SHALL be a no-op.

#### Scenario: ? opens help panel when closed
- **WHEN** `m.screenMode == ScreenNormal`, `m.showHelpPanel == false`, and `?` is pressed
- **THEN** `m.showHelpPanel == true`

#### Scenario: ? closes help panel when open
- **WHEN** `m.screenMode == ScreenNormal`, `m.showHelpPanel == true`, and `?` is pressed
- **THEN** `m.showHelpPanel == false`

#### Scenario: ? is no-op in ScreenInventory
- **WHEN** `m.screenMode == ScreenInventory` and `?` is pressed
- **THEN** `m.showHelpPanel` is unchanged

#### Scenario: ? is no-op in ScreenCombat
- **WHEN** `m.screenMode == ScreenCombat` and `?` is pressed
- **THEN** `m.showHelpPanel` is unchanged

### Requirement: renderHelpPanel fills the viewport with contextual key bindings
`renderHelpPanel(m Model) string` SHALL render a fullscreen help overlay. It SHALL show a header `" Key Bindings"` and a separator, then list key bindings relevant to `m.mode` and `m.screenMode`. The rendered output SHALL be clamped: the number of lines SHALL NOT exceed `m.viewportH`, and each line SHALL NOT exceed `m.viewportW` visible characters.

#### Scenario: Help panel output height does not exceed viewportH
- **WHEN** `renderHelpPanel` is called with `m.viewportH == 10`
- **THEN** the output contains at most 10 newline-separated lines

#### Scenario: Help panel shows mode-specific bindings for ModeLocal
- **WHEN** `m.mode == ModeLocal` and the help panel is rendered
- **THEN** the output contains `g` (pick up / fight) and `i` (inventory)

#### Scenario: Help panel shows mode-specific bindings for ModeWorld
- **WHEN** `m.mode == ModeWorld` and the help panel is rendered
- **THEN** the output contains `m` (map mode) and `i` (inventory)

#### Scenario: Help panel shows mode-specific bindings for ModeDungeon
- **WHEN** `m.mode == ModeDungeon` and the help panel is rendered
- **THEN** the output contains `f` (torch) and `<`/`>` (up/down)

### Requirement: buildView returns help panel when showHelpPanel is true
When `m.showHelpPanel == true`, `buildView` SHALL return `renderHelpPanel(m)` immediately, before map/HUD rendering. This takes priority over the sidebar and map picker but is suppressed by `ScreenInventory` and `ScreenCombat` (those early-returns come first).

#### Scenario: buildView returns help panel when showHelpPanel is true
- **WHEN** `m.showHelpPanel == true` and `m.screenMode == ScreenNormal`
- **THEN** `buildView(m)` returns the same string as `renderHelpPanel(m)`

#### Scenario: buildView does not return help panel in ScreenInventory
- **WHEN** `m.showHelpPanel == true` and `m.screenMode == ScreenInventory`
- **THEN** `buildView(m)` does NOT return `renderHelpPanel(m)`
