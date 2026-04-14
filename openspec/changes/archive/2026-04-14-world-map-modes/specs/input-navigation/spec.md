## MODIFIED Requirements

### Requirement: Arrow keys and WASD move the player
The system SHALL handle `up`/`w`, `down`/`s`, `left`/`a`, `right`/`d` key messages. In `ModeWorld`, each shall increment/decrement `worldPos.X` or `worldPos.Y` by 1, **unless `showMapPicker == true`**, in which case `up` and `down` SHALL instead move `mapPickerCursor` by -1 or +1 respectively (clamped to `[0, 3]`), and `left`/`right`/`a`/`d`/`w`/`s` SHALL be no-ops. In `ModeLocal`, each shall move `playerPos` by 1 in the corresponding axis, subject to bounds and collision rules (unchanged; map picker is world-only).

#### Scenario: Arrow key moves player in ModeWorld
- **WHEN** the player presses the right arrow key in `ModeWorld` and `showMapPicker == false`
- **THEN** `worldPos.X` increases by 1

#### Scenario: WASD moves player in ModeLocal
- **WHEN** the player presses `s` (down) in `ModeLocal` and the cell below is passable
- **THEN** `playerPos.Y` increases by 1

#### Scenario: Player cannot move outside local map bounds
- **WHEN** `playerPos` is `{0, 0}` and the player presses `up`
- **THEN** `playerPos` remains `{0, 0}`

#### Scenario: Player cannot move into a blocking object cell
- **WHEN** `playerPos` is `{5, 5}` and `Objects[5][4].Blocking == true` and the player presses `up`
- **THEN** `playerPos` remains `{5, 5}`

#### Scenario: Up/down navigate picker cursor when picker is open
- **WHEN** `showMapPicker == true`, `mapPickerCursor == 1`, and the player presses `down`
- **THEN** `mapPickerCursor == 2` and `worldPos` is unchanged

#### Scenario: Picker cursor clamps at bottom
- **WHEN** `showMapPicker == true`, `mapPickerCursor == 3`, and the player presses `down`
- **THEN** `mapPickerCursor` remains `3`

#### Scenario: Picker cursor clamps at top
- **WHEN** `showMapPicker == true`, `mapPickerCursor == 0`, and the player presses `up`
- **THEN** `mapPickerCursor` remains `0`

## ADDED Requirements

### Requirement: m key toggles the map mode picker in ModeWorld
The system SHALL handle the `m` key in `ModeWorld`. Pressing `m` when `showMapPicker == false` SHALL set `showMapPicker = true`, `showSidebar = false`, and set `mapPickerCursor` to the index of the current `mapMode`. Pressing `m` (or `esc`) when `showMapPicker == true` SHALL close the picker without changing `mapMode`. The `m` key SHALL be a no-op in `ModeLocal`.

#### Scenario: m opens picker in ModeWorld
- **WHEN** `showMapPicker == false` and the player presses `m` in `ModeWorld`
- **THEN** `showMapPicker == true` and `mapPickerCursor == int(m.mapMode)`

#### Scenario: m closes picker without changing mode
- **WHEN** `showMapPicker == true` and the player presses `m`
- **THEN** `showMapPicker == false` and `mapMode` is unchanged

#### Scenario: m is a no-op in ModeLocal
- **WHEN** the player presses `m` in `ModeLocal`
- **THEN** `showMapPicker` remains `false` and `mode` remains `ModeLocal`

### Requirement: Enter confirms map mode selection from picker
The system SHALL handle `enter` when `showMapPicker == true` in `ModeWorld`. Pressing `enter` SHALL set `mapMode` to the `MapMode` constant corresponding to `mapPickerCursor` and then close the picker (`showMapPicker = false`).

#### Scenario: Enter applies cursor selection and closes picker
- **WHEN** `showMapPicker == true`, `mapPickerCursor == 2`, and the player presses `enter`
- **THEN** `mapMode == MapModeElevation` and `showMapPicker == false`

#### Scenario: Enter in ModeWorld without picker open descends to local map
- **WHEN** `showMapPicker == false` and the player presses `enter` in `ModeWorld`
- **THEN** `model.mode == ModeLocal` (existing descend behavior is unchanged)
