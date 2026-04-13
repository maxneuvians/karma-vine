## ADDED Requirements

### Requirement: World map is rendered using viewport math
In `ModeWorld`, the system SHALL render a grid of `viewportW × (viewportH - 1)` cells. For each screen cell `(screenX, screenY)`, the world coordinate SHALL be `worldX = playerWorldX + (screenX - viewportW/2)`, `worldY = playerWorldY + (screenY - (viewportH-1)/2)`. Each cell SHALL display `tile.Char` styled with `tile.Color` as the Lipgloss foreground.

#### Scenario: Player tile appears at screen centre
- **WHEN** the world map is rendered
- **THEN** the tile at screen position `(viewportW/2, (viewportH-1)/2)` corresponds to `worldPos.X, worldPos.Y`

#### Scenario: All visible cells are drawn without blank gaps
- **WHEN** the world map is rendered with a 40×20 viewport
- **THEN** the rendered string contains exactly 39 newline characters (one per row, last row has no trailing newline)

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

### Requirement: Night mode dims all tile colours
When `Model.nightMode` is `true`, the system SHALL multiply the R, G, B components of every tile colour by `0.35` (rounded to nearest integer) before applying the Lipgloss style. The player glyph colour `#f0f6fc` SHALL NOT be dimmed.

#### Scenario: Night mode reduces brightness of biome colours
- **WHEN** `nightMode` is `true` and a tile has `Color = "#2d7a1f"` (Forest green)
- **THEN** the applied colour is approximately `#0f2a0a` (each channel × 0.35)

#### Scenario: Player colour is unaffected by night mode
- **WHEN** `nightMode` is `true`
- **THEN** the player glyph `@` is rendered in `#f0f6fc`

### Requirement: Viewport dimensions update on window resize
The system SHALL handle `tea.WindowSizeMsg` in `Update()`, storing `msg.Width` in `viewportW` and `msg.Height` in `viewportH`. The next `View()` call SHALL use the updated dimensions.

#### Scenario: Resize message updates stored dimensions
- **WHEN** a `tea.WindowSizeMsg{Width: 120, Height: 40}` is dispatched
- **THEN** `model.viewportW == 120` and `model.viewportH == 40` after the update
