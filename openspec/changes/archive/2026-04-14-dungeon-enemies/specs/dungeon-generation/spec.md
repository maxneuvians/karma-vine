## MODIFIED Requirements

### Requirement: GenerateDungeonLevel produces a valid level
The system SHALL implement `GenerateDungeonLevel(globalSeed, wx, wy, depth, maxDepth int, biome Biome) *DungeonLevel` that:
1. Seeds a local `rand.Rand` with `int64(globalSeed ^ wx*31337 ^ wy*1619 ^ depth*7919)` for determinism
2. Fills all cells as `CellWall`, then uses BSP to split the area into at least 4 leaf rectangles
3. Carves a room inside each leaf (room is at least 3×3 cells, walls kept as `CellWall`, interior set to `CellFloor`)
4. Connects adjacent rooms with L-shaped corridors of `CellFloor` cells
5. Places an up-staircase object (`Char: '<'`, `Color: "#e8c96a"`, `Blocking: false`, `Name: "Staircase Up"`) on a floor cell near the top-left room; stores position in `DungeonLevel.UpStair`
6. If `depth < maxDepth`, places a down-staircase object (`Char: '>'`, `Color: "#e8c96a"`, `Blocking: false`, `Name: "Staircase Down"`) on a floor cell in a different room; stores position in `DungeonLevel.DownStair` and sets `HasDownStair = true`
7. Seeds torches (`Char: '†'`, `Color: "#e8c96a"`, `Blocking: true`, `Name: "Torch"`) on `CellWall` cells adjacent to `CellFloor` cells; approximately 1 torch per 5 rooms
8. Seeds braziers (`Char: 'Ω'`, `Color: "#e07030"`, `Blocking: false`, `Name: "Brazier"`) on `CellFloor` cells inside rooms; approximately 1 brazier per 6 rooms
9. **After placing objects**, looks up `biome` in `dungeonEnemyRoster` to get the enemy template, spawns `depth` enemies on random `CellFloor` cells that are not occupied by objects or staircases, calls `spawnEnemy` for each, and stores them in `DungeonLevel.Enemies`

#### Scenario: Same inputs produce the same level
- **WHEN** `GenerateDungeonLevel(42, 3, 7, 1, 5, Jungle)` is called twice
- **THEN** both returned `DungeonLevel` values have identical `Cells` arrays and identical enemy positions

#### Scenario: Generated level has at least one floor cell
- **WHEN** `GenerateDungeonLevel` is called with any valid inputs
- **THEN** at least one cell has `Kind == CellFloor`

#### Scenario: Up-staircase is placed on a floor cell
- **WHEN** `GenerateDungeonLevel` is called
- **THEN** `level.Cells[level.UpStair.X][level.UpStair.Y].Kind == CellFloor`

#### Scenario: Down-staircase present on non-final level
- **WHEN** `depth < maxDepth`
- **THEN** `level.HasDownStair == true` and `level.Cells[level.DownStair.X][level.DownStair.Y].Kind == CellFloor`

#### Scenario: No down-staircase on final level
- **WHEN** `depth == maxDepth`
- **THEN** `level.HasDownStair == false`

#### Scenario: Torch placed on wall cell
- **WHEN** a torch object is placed
- **THEN** the cell it occupies has `Kind == CellWall`

#### Scenario: No dungeon object has an empty Name
- **WHEN** `GenerateDungeonLevel` is called
- **THEN** every non-nil `Object` in any `DungeonCell` has a non-empty `Name`

#### Scenario: Enemy count equals depth
- **WHEN** `GenerateDungeonLevel` is called with `depth=3`
- **THEN** `len(level.Enemies) == 3`

#### Scenario: Enemies are placed on floor cells
- **WHEN** `GenerateDungeonLevel` is called
- **THEN** every enemy in `level.Enemies` occupies a `CellFloor` cell

### Requirement: DungeonMeta stores per-entrance max depth and biome
The system SHALL define `DungeonMeta struct { MaxDepth int; Biome Biome }` and store one per world coordinate in `Model.dungeonMeta map[WorldCoord]DungeonMeta`. When first generating a dungeon for a world tile, the system SHALL randomise `MaxDepth` in `[5, 10]` using the global seed and record the tile's `Biome`.

#### Scenario: MaxDepth is in range [5, 10]
- **WHEN** a `DungeonMeta` is created for a new world tile
- **THEN** `meta.MaxDepth >= 5 && meta.MaxDepth <= 10`

#### Scenario: MaxDepth is stable across re-entries
- **WHEN** the player enters and re-enters the same dungeon entrance
- **THEN** `DungeonMeta.MaxDepth` is the same value on both entries

#### Scenario: Biome is recorded on first entry
- **WHEN** a `DungeonMeta` is created for a tile with biome Tundra
- **THEN** `meta.Biome == Tundra`
