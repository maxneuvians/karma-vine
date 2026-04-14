## ADDED Requirements

### Requirement: ScreenMode type controls overlay routing
The system SHALL define `ScreenMode int` in `types.go` with constants `ScreenNormal ScreenMode = 0` and `ScreenInventory ScreenMode = 1`. `Model` SHALL have a `screenMode ScreenMode` field. `NewModel()` SHALL leave `screenMode` at its zero value (`ScreenNormal`).

#### Scenario: New model starts in ScreenNormal
- **WHEN** `NewModel()` is called
- **THEN** `m.screenMode == ScreenNormal`

### Requirement: showInventory field is replaced by screenMode
The `showInventory bool` field SHALL be removed from `Model`. All code that previously checked `m.showInventory` or set `m.showInventory = true/false` SHALL be updated to use `m.screenMode == ScreenInventory` or `m.screenMode = ScreenInventory / ScreenNormal` respectively.

#### Scenario: i key sets screenMode to ScreenInventory
- **WHEN** `m.screenMode == ScreenNormal` and the player presses `i`
- **THEN** `m.screenMode == ScreenInventory`

#### Scenario: Second i key returns to ScreenNormal
- **WHEN** `m.screenMode == ScreenInventory` and the player presses `i`
- **THEN** `m.screenMode == ScreenNormal`

### Requirement: renderFullscreenInventory renders a full-viewport inventory
The system SHALL implement `renderFullscreenInventory(m Model) string` that fills the entire viewport (`m.viewportW × m.viewportH`). The layout SHALL be two columns:
- **Left column** (60% of viewport width): item list with header "Inventory", separator, item rows (`[glyph] Name  xN`), cursor-highlighted row, "Empty" placeholder when no items, and hint row at bottom.
- **Right column** (remaining width): ASCII ragdoll body outline centred vertically, with named slot labels (`Head`, `Chest`, `Left Hand`, `Right Hand`, `Legs`, `Feet`) rendered as `SlotName : [ Empty ]`.

The panel SHALL be rendered regardless of the current `m.mode` (world/local/dungeon are all valid).

#### Scenario: Fullscreen panel fills viewport
- **WHEN** `renderFullscreenInventory` is called with `viewportW=120, viewportH=40`
- **THEN** the returned string contains exactly 39 newline characters (40 rows, last row has no trailing newline)

#### Scenario: Item list shows glyph and count
- **WHEN** inventory has one item "Axe" with Count=2
- **THEN** the rendered output contains "Axe" and "x2"

#### Scenario: Empty inventory shows placeholder
- **WHEN** `len(m.inventory.Items) == 0`
- **THEN** the rendered output contains "Empty"

#### Scenario: Ragdoll column shows all slot labels
- **WHEN** `renderFullscreenInventory` is called
- **THEN** the output contains "Head", "Chest", "Left Hand", "Right Hand", "Legs", and "Feet"

#### Scenario: Each slot shows [Empty] placeholder
- **WHEN** no items are equipped (all placeholders)
- **THEN** every slot label is followed by "[ Empty ]"

#### Scenario: Cursor row is visually distinct
- **WHEN** `inventoryCursor == 1` and there are 3 items
- **THEN** the second item row is rendered in a distinct highlight style compared to the others

### Requirement: buildView dispatches to fullscreen inventory when ScreenInventory
When `m.screenMode == ScreenInventory`, `buildView` SHALL return `renderFullscreenInventory(m)` directly, bypassing the normal map/HUD/keybar composition. The HUD and key bar SHALL NOT be visible when the inventory is open.

#### Scenario: buildView returns fullscreen panel when ScreenInventory
- **WHEN** `m.screenMode == ScreenInventory`
- **THEN** `buildView` output contains "Inventory" and does not contain a map glyph

#### Scenario: buildView returns normal map when ScreenNormal
- **WHEN** `m.screenMode == ScreenNormal` and `m.mode == ModeWorld`
- **THEN** `buildView` output does not contain "Inventory"

### Requirement: esc closes the fullscreen inventory
When `m.screenMode == ScreenInventory` and the player presses `esc`, `screenMode` SHALL be set back to `ScreenNormal`. This is in addition to the existing `i` key toggle.

#### Scenario: esc closes inventory
- **WHEN** `m.screenMode == ScreenInventory` and the player presses `esc`
- **THEN** `m.screenMode == ScreenNormal`
