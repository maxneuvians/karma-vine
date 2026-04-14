## Requirements

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

### Requirement: Enter or > descends from world map to local map
The system SHALL handle `enter` and `>` key messages. In `ModeWorld`, pressing either key SHALL: call `LocalMapFor` for the current `worldPos`, set `model.localMap` to the result, set `model.mode = ModeLocal`, and set `playerPos` to the first non-blocking cell scanning outward from `{21, 9}`.

#### Scenario: Descend loads and displays local map
- **WHEN** the player presses `enter` in `ModeWorld`
- **THEN** `model.mode == ModeLocal` and `model.localMap` is non-nil after the update

#### Scenario: Player is placed at or near map centre on descend
- **WHEN** the player descends and cell `{21, 9}` is not blocking
- **THEN** `playerPos == {21, 9}`

### Requirement: Escape or < ascends from local map to world map
The system SHALL handle `esc` and `<` key messages. In `ModeLocal`, pressing either key SHALL set `model.mode = ModeWorld`. `model.localMap` SHALL NOT be set to nil. This handler SHALL only fire when `m.mode == ModeLocal`; it SHALL NOT fire when `m.mode == ModeDungeon`.

#### Scenario: Ascend returns to world map
- **WHEN** the player presses `esc` in `ModeLocal`
- **THEN** `model.mode == ModeWorld` after the update

#### Scenario: Local map is preserved in cache after ascend
- **WHEN** the player ascends and then descends to the same world tile
- **THEN** no new `GenerateLocalMap` call is made (cache hit)

#### Scenario: esc in ModeDungeon does not go to world map
- **WHEN** the player presses `esc` in `ModeDungeon`
- **THEN** `model.mode != ModeWorld`

### Requirement: Movement keys work in ModeDungeon
The system SHALL route `up`/`w`, `down`/`s`, `left`/`a`, `right`/`d` through `applyDelta` when `m.mode == ModeDungeon`, applying dungeon-specific bounds and wall collision rules.

#### Scenario: WASD moves player in ModeDungeon
- **WHEN** the player presses `d` in `ModeDungeon` and the cell to the right is `CellFloor`
- **THEN** `playerPos.X` increases by 1

### Requirement: Stair transitions handled in ModeDungeon
The system SHALL handle `enter`/`>` and `esc`/`<` in `ModeDungeon` as dungeon-level transitions (not world/local transitions). The existing `esc` handler for `ModeLocal → ModeWorld` SHALL NOT trigger when `m.mode == ModeDungeon`.

#### Scenario: esc in ModeDungeon does not ascend to world map
- **WHEN** `m.mode == ModeDungeon` and the player presses `esc`
- **THEN** `m.mode` is either `ModeDungeon` (if depth > 1) or `ModeLocal` (if depth == 1) — never `ModeWorld`

#### Scenario: enter in ModeLocal on dungeon entrance enters dungeon
- **WHEN** `m.mode == ModeLocal` and `playerPos` is at a cell with `Object.Char == '>'`
- **THEN** `m.mode == ModeDungeon` after the key event

### Requirement: q and ctrl+c quit the program
The system SHALL handle `q` and `ctrl+c` at all times by returning `tea.Quit`.

#### Scenario: q quits from world mode
- **WHEN** the player presses `q` in `ModeWorld`
- **THEN** the program exits with code 0

#### Scenario: ctrl+c quits from local mode
- **WHEN** the player presses `ctrl+c` in `ModeLocal`
- **THEN** the program exits with code 0

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

### Requirement: [ and ] keys adjust time speed
The system SHALL handle `[` and `]` key messages at all times. `]` SHALL advance `timeScale` to the next value in `{1, 2, 5, 10}`, clamped at `10`. `[` SHALL retreat to the previous value, clamped at `1`.

#### Scenario: ] increases time scale
- **WHEN** `timeScale == 1` and the player presses `]`
- **THEN** `timeScale == 2`

#### Scenario: [ decreases time scale
- **WHEN** `timeScale == 5` and the player presses `[`
- **THEN** `timeScale == 2`

#### Scenario: ] clamps at maximum
- **WHEN** `timeScale == 10` and the player presses `]`
- **THEN** `timeScale` remains `10`
