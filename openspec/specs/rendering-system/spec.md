## Requirements

### Requirement: World map is rendered using viewport math
In `ModeWorld`, the system SHALL render a grid of `viewportW × (viewportH - 1)` cells. For each screen cell `(screenX, screenY)`, the world coordinate SHALL be `worldX = playerWorldX + (screenX - viewportW/2)`, `worldY = playerWorldY + (screenY - (viewportH-1)/2)`. Each cell SHALL display the char and color returned by `tileVisual(tile, m.mapMode)` as the Lipgloss foreground, rather than `tile.Char` and `tile.Color` directly. When `m.mapMode != MapModeDefault`, the dim-factor color scaling SHALL still be applied to the color returned by `tileVisual`.

#### Scenario: Player tile appears at screen centre
- **WHEN** the world map is rendered
- **THEN** the tile at screen position `(viewportW/2, (viewportH-1)/2)` corresponds to `worldPos.X, worldPos.Y`

#### Scenario: All visible cells are drawn without blank gaps
- **WHEN** the world map is rendered with a 40×20 viewport
- **THEN** the rendered string contains exactly 39 newline characters (one per row, last row has no trailing newline)

#### Scenario: Temperature mode renders with temperature-derived color
- **WHEN** `m.mapMode == MapModeTemperature` and the world map is rendered
- **THEN** each cell's color is derived from `tileVisual(tile, MapModeTemperature)`, not from `tile.Color`

### Requirement: View composition includes map picker overlay when active
The system SHALL, in `buildView`, check `m.showMapPicker`. When true, it SHALL render `renderMapPicker(m, mapH)` on the right side of the viewport (using the same composition approach as `showSidebar`), reducing the available map width by the picker panel width. When both `showMapPicker` and `showSidebar` are false, layout is unchanged from the current implementation.

#### Scenario: Map picker reduces map width when open
- **WHEN** `showMapPicker == true` and the viewport is 80×24
- **THEN** the rendered map occupies fewer than 80 columns (picker panel takes the remainder)

#### Scenario: No layout change when picker is closed
- **WHEN** `showMapPicker == false` and `showSidebar == false`
- **THEN** the map renders at full viewport width (minus HUD rows), identical to current behavior

### Requirement: Local map is rendered with layered glyphs
In `ModeLocal`, the system SHALL render the 42×18 `LocalMap` grid. For each cell, it SHALL display: the animal's `Char` if an animal occupies that cell, else the object's `Char` if `Objects[x][y]` is non-nil, else the ground's `Char`. If the cell matches `playerPos`, the player glyph `@` SHALL be drawn in `#f0f6fc` bold, overriding all other layers.

#### Scenario: Player glyph overrides ground/object/animal
- **WHEN** the local map is rendered and `playerPos` is `{5, 3}`
- **THEN** the character at screen position `(5, 3)` is `@` in bold `#f0f6fc`

#### Scenario: Object appears above ground
- **WHEN** a cell has a non-nil `Object` and no animal or player
- **THEN** `Object.Char` is rendered at that cell position

### Requirement: HUD status bar displays game state
The system SHALL render a status bar as the bottom row of the frame. It SHALL display the current biome name, elevation formatted to 2 decimal places, world coordinates as `(x, y)`, and chunk coordinates as `chunk (cx, cy)`.

#### Scenario: HUD shows correct biome name
- **WHEN** the player is on a `Forest` tile
- **THEN** the status bar contains the text `Forest`

#### Scenario: HUD shows world coordinates
- **WHEN** the player's world position is `{10, -5}`
- **THEN** the status bar contains `(10, -5)`

### Requirement: Tile colors are dimmed by a continuous time-of-day factor
The system SHALL compute `dimFactor(timeOfDay float64) float64` using a cosine curve: `0.5*(1+cos(2π*(tod-0.5)))*0.85+0.15`, clamped to `[0.15, 1.0]`. This yields `1.0` at noon (`tod=0.5`) and `0.15` at midnight (`tod=0.0`). Each tile's R, G, B colour components SHALL be scaled by this factor before applying the Lipgloss style. The player glyph colour `#f0f6fc` SHALL NOT be dimmed.

#### Scenario: dimFactor is 1.0 at noon
- **WHEN** `timeOfDay == 0.5`
- **THEN** `dimFactor(0.5) == 1.0`

#### Scenario: dimFactor is 0.15 at midnight
- **WHEN** `timeOfDay == 0.0`
- **THEN** `dimFactor(0.0) == 0.15`

#### Scenario: Player colour is unaffected by dim factor
- **WHEN** `timeOfDay == 0.0` (midnight)
- **THEN** the player glyph `@` is rendered in `#f0f6fc`

### Requirement: Local map fire cells override dim within illumination radius
When rendering a local map cell, the effective dim factor for that cell SHALL be `max(globalDim, LitMap[x][y])`, where `globalDim = dimFactor(timeOfDay)` and `LitMap[x][y]` is the precomputed fire illumination intensity. This means fire cells and nearby cells appear brighter than the global dim level at night.

#### Scenario: Fire cell at midnight is brighter than ambient dim
- **WHEN** `timeOfDay == 0.0` (midnight) and a cell has `LitMap[x][y] == 1.0`
- **THEN** the effective dim factor for that cell is `1.0` (not `0.15`)

#### Scenario: Distant cell unaffected by fire at midnight
- **WHEN** `timeOfDay == 0.0` and a cell has `LitMap[x][y] == 0`
- **THEN** the effective dim factor for that cell is `0.15` (ambient only)

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

### Requirement: Viewport dimensions update on window resize
The system SHALL handle `tea.WindowSizeMsg` in `Update()`, storing `msg.Width` in `viewportW` and `msg.Height` in `viewportH`. The next `View()` call SHALL use the updated dimensions.

#### Scenario: Resize message updates stored dimensions
- **WHEN** a `tea.WindowSizeMsg{Width: 120, Height: 40}` is dispatched
- **THEN** `model.viewportW == 120` and `model.viewportH == 40` after the update

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

### Requirement: HUD shows inventory item count
The HUD row SHALL include the current inventory count in the format `Items: N/8` (where 8 = `InventoryMaxSlots`). This count SHALL be visible in all three modes when the HUD is displayed.

#### Scenario: HUD shows item count
- **WHEN** the player has 2 items and the HUD is rendered
- **THEN** the HUD contains "Items: 2/8" (or equivalent display)

### Requirement: Key bar shows inventory and use hints
The key bar SHALL include the following additional hints:
- `i inv` — toggle inventory
- `g pick` — pick up (shown in ModeLocal and ModeDungeon only)
- `d drop` — drop item (shown in ModeLocal and ModeDungeon only)
- `u use` — use item (shown in ModeLocal and ModeDungeon only)

#### Scenario: Key bar in ModeLocal includes inventory keys
- **WHEN** `m.mode == ModeLocal`
- **THEN** the key bar contains "i", "g", "d", and "u" hints
