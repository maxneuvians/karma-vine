## MODIFIED Requirements

### Requirement: Arrow keys and WASD move the player
The system SHALL handle `up`/`w`, `down`/`s`, `left`/`a`, `right`/`d` key messages. In `ModeWorld`, each shall increment/decrement `worldPos.X` or `worldPos.Y` by 1, **unless `showMapPicker == true`**, in which case `up` and `down` SHALL instead move `mapPickerCursor` by -1 or +1 respectively (clamped to `[0, 3]`), and `left`/`right`/`a`/`d`/`w`/`s` SHALL be no-ops. In `ModeLocal`, each shall move `playerPos` by 1 in the corresponding axis, subject to bounds and collision rules. **When `m.screenMode == ScreenInventory`**, `up`/`w` and `down`/`s` SHALL move `inventoryCursor` instead of the player or picker, and `left`/`right`/`a`/`d` SHALL be no-ops. **When `m.paused == true` and `m.screenMode == ScreenNormal`**, all directional keys SHALL be no-ops (player and world-pos movement are suppressed; picker and inventory cursor movement are unaffected by pause).

#### Scenario: Arrow key moves player in ModeWorld
- **WHEN** the player presses the right arrow key in `ModeWorld` and `showMapPicker == false` and `screenMode == ScreenNormal`
- **THEN** `worldPos.X` increases by 1

#### Scenario: WASD moves player in ModeLocal
- **WHEN** the player presses `s` (down) in `ModeLocal`, `screenMode == ScreenNormal`, and the cell below is passable
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

#### Scenario: Up/down navigate inventory cursor when inventory is open
- **WHEN** `screenMode == ScreenInventory` and the player presses `down`
- **THEN** `inventoryCursor` increments (clamped at `len(inventory.Items)-1`) and `playerPos` is unchanged

#### Scenario: Movement is blocked when paused in ScreenNormal
- **WHEN** `m.paused == true`, `screenMode == ScreenNormal`, `mode == ModeLocal`, and the player presses `right`
- **THEN** `playerPos` is unchanged

#### Scenario: Inventory cursor still moves when paused in ScreenInventory
- **WHEN** `m.paused == true`, `screenMode == ScreenInventory`, and the player presses `down`
- **THEN** `inventoryCursor` increments normally

## ADDED Requirements

### Requirement: Left-click movement is blocked when paused
When `m.paused == true` and `m.screenMode == ScreenNormal`, `handleMouseClick` SHALL treat left-click movement as a no-op (same as when sidebar or inventory is open).

#### Scenario: Click movement is a no-op while paused
- **WHEN** `m.paused == true`, `screenMode == ScreenNormal`, `mode == ModeLocal`, no panels open, and the player left-clicks a cell to the right
- **THEN** `playerPos` is unchanged
