## ADDED Requirements

### Requirement: Player movement in dungeon mode
The system SHALL handle `up`/`w`, `down`/`s`, `left`/`a`, `right`/`d` key messages in `ModeDungeon`. Movement SHALL be blocked by `CellWall` cells and by `Object` cells where `Blocking == true`. Movement SHALL be clamped to `[0, DungeonW-1] × [0, DungeonH-1]`.

#### Scenario: Player moves onto floor cell
- **WHEN** the player presses `right` in `ModeDungeon` and the cell to the right is `CellFloor` with no blocking object
- **THEN** `playerPos.X` increases by 1

#### Scenario: Player cannot move into wall
- **WHEN** the player presses `up` and the cell above is `CellWall`
- **THEN** `playerPos` remains unchanged

#### Scenario: Player cannot move outside dungeon bounds
- **WHEN** `playerPos.X == 0` and the player presses `left`
- **THEN** `playerPos.X` remains 0

#### Scenario: Player can walk over floor item (brazier)
- **WHEN** the player presses `right` and the target cell is `CellFloor` with a non-blocking brazier object
- **THEN** `playerPos.X` increases by 1 (movement succeeds)

### Requirement: Descending into dungeon from local map
The system SHALL handle `enter`/`>` in `ModeLocal` when `playerPos` is on a cell containing a dungeon entrance object (`Char == '>'`). On this event the system SHALL:
1. Look up or create `DungeonMeta` for current `worldPos` (randomise `MaxDepth` in `[5,10]` on first entry)
2. Call `DungeonLevelFor(wx, wy, 1, m)` to obtain level 1
3. Set `m.currentDungeon` to the level, `m.dungeonDepth = 1`, `m.mode = ModeDungeon`
4. Place `playerPos` at the cell adjacent to `level.UpStair` (or on `UpStair` if no adjacent floor cell)

#### Scenario: Descend from local map activates dungeon mode
- **WHEN** the player presses `enter` while standing on a dungeon entrance (`>`) in `ModeLocal`
- **THEN** `m.mode == ModeDungeon` and `m.dungeonDepth == 1`

#### Scenario: Enter key on non-staircase cell in ModeLocal descends to world map
- **WHEN** the player presses `enter` in `ModeLocal` and `playerPos` is NOT on a dungeon entrance
- **THEN** the existing world-map descent behaviour is NOT triggered (ModeLocal has no world-map descent via enter; this tests no regression)

### Requirement: Descending to next dungeon level
The system SHALL handle `enter`/`>` in `ModeDungeon` when `m.currentDungeon.HasDownStair == true` and `playerPos == m.currentDungeon.DownStair`. On this event the system SHALL:
1. Call `DungeonLevelFor(wx, wy, m.dungeonDepth+1, m)`
2. Set `m.currentDungeon` to the new level, increment `m.dungeonDepth`
3. Place `playerPos` at `level.UpStair`

#### Scenario: Pressing > on down-staircase descends one level
- **WHEN** the player is standing on the down-staircase and presses `>`
- **THEN** `m.dungeonDepth` increases by 1 and `m.currentDungeon` is the new level

#### Scenario: Pressing > away from down-staircase does nothing
- **WHEN** the player is NOT standing on the down-staircase and presses `>`
- **THEN** `m.dungeonDepth` and `m.currentDungeon` are unchanged

#### Scenario: No descent possible on final level
- **WHEN** the player is on the final dungeon level (`HasDownStair == false`) and presses `>`
- **THEN** `m.dungeonDepth` and `m.currentDungeon` are unchanged

### Requirement: Ascending dungeon levels
The system SHALL handle `esc`/`<` in `ModeDungeon`. On this event:
- If `m.dungeonDepth > 1`: call `DungeonLevelFor(wx, wy, m.dungeonDepth-1, m)`, set `m.currentDungeon` to that level, decrement `m.dungeonDepth`, place `playerPos` at `level.DownStair`
- If `m.dungeonDepth == 1`: set `m.mode = ModeLocal`, restore `playerPos` to the dungeon entrance cell on the local map

#### Scenario: Ascending from depth > 1 goes to previous level
- **WHEN** `m.dungeonDepth == 3` and the player presses `esc`
- **THEN** `m.dungeonDepth == 2` and `m.currentDungeon` is level 2

#### Scenario: Ascending from depth 1 returns to local map
- **WHEN** `m.dungeonDepth == 1` and the player presses `esc`
- **THEN** `m.mode == ModeLocal`

#### Scenario: Local map position restored near dungeon entrance
- **WHEN** the player ascends from depth 1
- **THEN** `playerPos` is at the dungeon entrance cell (or the nearest non-blocking adjacent cell)
