## MODIFIED Requirements

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
