## Requirements

### Requirement: Arrow keys and WASD move the player
The system SHALL handle `up`/`w`, `down`/`s`, `left`/`a`, `right`/`d` key messages. In `ModeWorld`, each shall increment/decrement `worldPos.X` or `worldPos.Y` by 1. In `ModeLocal`, each shall move `playerPos` by 1 in the corresponding axis, subject to bounds and collision rules.

#### Scenario: Arrow key moves player in ModeWorld
- **WHEN** the player presses the right arrow key in `ModeWorld`
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

### Requirement: Enter or > descends from world map to local map
The system SHALL handle `enter` and `>` key messages. In `ModeWorld`, pressing either key SHALL: call `LocalMapFor` for the current `worldPos`, set `model.localMap` to the result, set `model.mode = ModeLocal`, and set `playerPos` to the first non-blocking cell scanning outward from `{21, 9}`.

#### Scenario: Descend loads and displays local map
- **WHEN** the player presses `enter` in `ModeWorld`
- **THEN** `model.mode == ModeLocal` and `model.localMap` is non-nil after the update

#### Scenario: Player is placed at or near map centre on descend
- **WHEN** the player descends and cell `{21, 9}` is not blocking
- **THEN** `playerPos == {21, 9}`

### Requirement: Escape or < ascends from local map to world map
The system SHALL handle `esc` and `<` key messages. In `ModeLocal`, pressing either key SHALL set `model.mode = ModeWorld`. `model.localMap` SHALL NOT be set to nil.

#### Scenario: Ascend returns to world map
- **WHEN** the player presses `esc` in `ModeLocal`
- **THEN** `model.mode == ModeWorld` after the update

#### Scenario: Local map is preserved in cache after ascend
- **WHEN** the player ascends and then descends to the same world tile
- **THEN** no new `GenerateLocalMap` call is made (cache hit)

### Requirement: q and ctrl+c quit the program
The system SHALL handle `q` and `ctrl+c` at all times by returning `tea.Quit`.

#### Scenario: q quits from world mode
- **WHEN** the player presses `q` in `ModeWorld`
- **THEN** the program exits with code 0

#### Scenario: ctrl+c quits from local mode
- **WHEN** the player presses `ctrl+c` in `ModeLocal`
- **THEN** the program exits with code 0

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
