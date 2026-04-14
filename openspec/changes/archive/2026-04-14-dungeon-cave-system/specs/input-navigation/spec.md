## ADDED Requirements

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

## MODIFIED Requirements

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
