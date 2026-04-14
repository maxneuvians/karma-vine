## ADDED Requirements

### Requirement: MapMode type defines the active world map overlay
The system SHALL define a `MapMode int` type in `types.go` with four constants: `MapModeDefault` (0), `MapModeTemperature` (1), `MapModeElevation` (2), `MapModePolitical` (3). The `Model` struct SHALL include a `mapMode MapMode` field (zero value = `MapModeDefault`) and a `showMapPicker bool` field (default `false`) and a `mapPickerCursor int` field (default 0).

#### Scenario: Default model has MapModeDefault
- **WHEN** `NewModel()` is called
- **THEN** `model.mapMode == MapModeDefault` and `model.showMapPicker == false`

#### Scenario: MapMode constants are distinct
- **WHEN** the four constants are compared pairwise
- **THEN** no two constants have the same integer value

### Requirement: tileVisual returns mode-specific char and color for world map rendering
The system SHALL provide `tileVisual(t Tile, mode MapMode) (ch rune, color string)` in `render.go`. For `MapModeDefault` it SHALL return `t.Char` and `t.Color` unchanged. For `MapModeTemperature` it SHALL return `·` as the char and a hex color linearly interpolated between `#4488ff` (cold, temperature=0) and `#ff4422` (hot, temperature=1). For `MapModeElevation` it SHALL return `·` and a hex color linearly interpolated between `#1a6fa8` (low, elevation=0) and `#f0f6fc` (high, elevation=1). For `MapModePolitical` it SHALL return `+` and `#aabbcc` when `int(t.Elevation*10) != int((t.Elevation+0.05)*10)` (contour boundary), else `·` and `#334455`.

#### Scenario: Default mode returns tile values unchanged
- **WHEN** `tileVisual(tile, MapModeDefault)` is called
- **THEN** the returned char equals `tile.Char` and the returned color equals `tile.Color`

#### Scenario: Temperature mode midpoint returns blended color
- **WHEN** `tileVisual(Tile{Temperature: 0.5}, MapModeTemperature)` is called
- **THEN** the returned char is `·` and the color is approximately `#aa6688` (midpoint blend of `#4488ff` and `#ff4422`)

#### Scenario: Elevation mode max returns high-elevation color
- **WHEN** `tileVisual(Tile{Elevation: 1.0}, MapModeElevation)` is called
- **THEN** the returned char is `·` and the color is `#f0f6fc`

#### Scenario: Elevation mode min returns low-elevation color
- **WHEN** `tileVisual(Tile{Elevation: 0.0}, MapModeElevation)` is called
- **THEN** the returned char is `·` and the color is `#1a6fa8`

### Requirement: Map mode picker panel renders as a right-side overlay
The system SHALL provide `renderMapPicker(m Model, height int) string` that renders a bordered list of the four mode names (`Default`, `Temperature`, `Elevation`, `Political`) with `m.mapPickerCursor` highlighted using a distinct style. The panel SHALL be `22` characters wide (including border). Each row SHALL display a cursor indicator (`>` for selected, ` ` otherwise) followed by the mode name.

#### Scenario: Picker highlights the cursor row
- **WHEN** `renderMapPicker` is called with `mapPickerCursor == 2`
- **THEN** the rendered string contains `> Elevation` and does not contain `> Default`

#### Scenario: Picker contains all four mode names
- **WHEN** `renderMapPicker` is called
- **THEN** the rendered string contains `Default`, `Temperature`, `Elevation`, and `Political`

### Requirement: Opening the map picker closes the sidebar, and vice versa
When `showMapPicker` is set to `true`, `showSidebar` SHALL be set to `false`. When `showSidebar` is set to `true`, `showMapPicker` SHALL be set to `false`.

#### Scenario: Opening picker closes sidebar
- **WHEN** `showSidebar == true` and the `m` key is pressed
- **THEN** `showMapPicker == true` and `showSidebar == false`

#### Scenario: Opening sidebar closes picker
- **WHEN** `showMapPicker == true` and the `?` key is pressed
- **THEN** `showSidebar == true` and `showMapPicker == false`
