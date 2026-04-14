## MODIFIED Requirements

### Requirement: View() returns tea.View
The `View()` method on `Model` SHALL return `tea.View` instead of `string`. The content SHALL be set via `tea.NewView(buildView(m))`. The returned `tea.View` SHALL have `AltScreen = true` and `MouseMode = tea.MouseModeCellMotion`. `buildView` itself continues to return a `string` and is unchanged in its signature.

#### Scenario: View() wraps buildView output in tea.View
- **WHEN** `View()` is called
- **THEN** the returned `tea.View`'s content equals what `buildView` returns

### Requirement: buildView dispatches to dungeon render path
The system SHALL extend `buildView` to dispatch on **both** `m.screenMode` and `m.mode`. When `m.screenMode == ScreenInventory`, `buildView` SHALL return `renderFullscreenInventory(m)` immediately, bypassing all other render paths. When `m.screenMode == ScreenNormal`, all existing dispatch logic (sidebar, map picker, local/dungeon/world map selection) is unchanged.

#### Scenario: Dungeon map renders when mode is ModeDungeon and screenMode is ScreenNormal
- **WHEN** `buildView` is called with `m.mode == ModeDungeon` and `m.screenMode == ScreenNormal`
- **THEN** the returned string contains dungeon cell characters and does not contain world-map tile characters

#### Scenario: Fullscreen inventory renders when screenMode is ScreenInventory
- **WHEN** `buildView` is called with `m.screenMode == ScreenInventory`
- **THEN** the returned string contains "Inventory" and does not contain HUD depth/mode information

#### Scenario: HUD is present in dungeon view when ScreenNormal
- **WHEN** `buildView` is called with `m.mode == ModeDungeon` and `m.screenMode == ScreenNormal`
- **THEN** the returned string contains the HUD status bar with dungeon depth information

### Requirement: renderSidebar dispatches on active mode
The system SHALL update `renderSidebar` to switch on `m.mode`, rendering world content for `ModeWorld`, local content for `ModeLocal`, and dungeon content for `ModeDungeon`. The `localCharNames` lookup map SHALL be removed; object and animal names SHALL be read directly from `obj.Name` and `a.Name`.

#### Scenario: renderSidebar called in ModeDungeon returns dungeon content
- **WHEN** `renderSidebar` is called with `m.mode == ModeDungeon`
- **THEN** the returned string contains `Dungeon` and `Depth:`

#### Scenario: renderSidebar called in ModeLocal uses Name field
- **WHEN** `renderSidebar` is called with `m.mode == ModeLocal` and an object has `Name: "Dungeon Entrance"`
- **THEN** the returned string contains `Dungeon Entrance`

## ADDED Requirements

### Requirement: renderFullscreenInventory produces a full-viewport two-column layout
The system SHALL implement `renderFullscreenInventory(m Model) string` in `render.go`. The function SHALL use `lipgloss.JoinHorizontal` to compose a left inventory column and a right ragdoll column, with the combined width equal to `m.viewportW` and height equal to `m.viewportH`.

#### Scenario: Output fills viewport dimensions
- **WHEN** `renderFullscreenInventory` is called with `viewportW=100, viewportH=30`
- **THEN** the rendered string has exactly 29 newline characters

#### Scenario: Left column contains inventory item list
- **WHEN** inventory has items
- **THEN** the left column lists item names and counts

#### Scenario: Right column contains ragdoll slots
- **WHEN** `renderFullscreenInventory` is called
- **THEN** the right column contains "Head", "Chest", "Left Hand", "Right Hand", "Legs", and "Feet"

### Requirement: Inventory panel side-column path removed
The old `inventoryPanelW` side-column append in `buildView` (which reduced map width by `inventoryPanelW`) SHALL be removed. The inventory is now exclusively rendered via `renderFullscreenInventory`. The `inventoryPanelW` constant SHALL be removed.

#### Scenario: Map renders at full width when ScreenNormal and no sidebar
- **WHEN** `m.screenMode == ScreenNormal`, `showSidebar == false`, and `showMapPicker == false`
- **THEN** the map is rendered at `m.viewportW` columns (not `m.viewportW - inventoryPanelW`)
