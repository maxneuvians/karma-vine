## MODIFIED Requirements

### Requirement: Descending into dungeon from local map
The system SHALL handle `enter`/`>` in `ModeLocal` when `playerPos` is on a cell containing a dungeon entrance object (`Char == '>'`). On this event the system SHALL:
1. Look up or create `DungeonMeta` for current `worldPos` (randomise `MaxDepth` in `[5,10]` on first entry; record `Biome` from `TileAt(m.worldPos).Biome` on first entry)
2. Call `DungeonLevelFor(wx, wy, 1, m)` to obtain level 1 (which now uses `meta.Biome` internally)
3. Set `m.currentDungeon` to the level, `m.dungeonDepth = 1`, `m.mode = ModeDungeon`
4. Place `playerPos` at the cell adjacent to `level.UpStair` (or on `UpStair` if no adjacent floor cell)

#### Scenario: Descend from local map activates dungeon mode
- **WHEN** the player presses `enter` while standing on a dungeon entrance (`>`) in `ModeLocal`
- **THEN** `m.mode == ModeDungeon` and `m.dungeonDepth == 1`

#### Scenario: Biome is recorded in DungeonMeta on first descent
- **WHEN** the player enters a dungeon for the first time at a Tundra world tile
- **THEN** `m.dungeonMeta[m.worldPos].Biome == Tundra`

#### Scenario: Enter key on non-staircase cell in ModeLocal descends to world map
- **WHEN** the player presses `enter` in `ModeLocal` and `playerPos` is NOT on a dungeon entrance
- **THEN** the existing world-map descent behaviour is NOT triggered

### Requirement: Player moving onto enemy cell triggers combat
When the player presses a directional key in `ModeDungeon` and the target cell is a `CellFloor` cell occupied by a `DungeonEnemy`, the system SHALL initiate combat with that enemy (build combatants, call `resolveCombat`, set `ScreenCombat`) instead of moving the player. The player's position SHALL remain unchanged until after combat is dismissed.

#### Scenario: Moving into enemy cell triggers ScreenCombat
- **WHEN** the player presses `right` in `ModeDungeon` and the cell to the right contains a `DungeonEnemy`
- **THEN** `m.screenMode == ScreenCombat` and `m.playerPos.X` is unchanged

#### Scenario: Moving into empty floor cell proceeds normally
- **WHEN** the player presses `right` in `ModeDungeon` and the cell to the right is an empty `CellFloor`
- **THEN** `m.playerPos.X` increases by 1 and `m.screenMode == ScreenNormal`
