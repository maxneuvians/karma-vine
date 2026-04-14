## MODIFIED Requirements

### Requirement: GenerateDungeonLevel produces a valid level
The system SHALL implement `GenerateDungeonLevel(globalSeed, wx, wy, depth, maxDepth int) *DungeonLevel` that:
1. Seeds a local `rand.Rand` with `int64(globalSeed ^ wx*31337 ^ wy*1619 ^ depth*7919)` for determinism
2. Fills all cells as `CellWall`, then uses BSP to split the area into at least 4 leaf rectangles
3. Carves a room inside each leaf (room is at least 3×3 cells, walls kept as `CellWall`, interior set to `CellFloor`)
4. Connects adjacent rooms with L-shaped corridors of `CellFloor` cells
5. Places an up-staircase object (`Char: '<'`, `Color: "#e8c96a"`, `Blocking: false`, `Name: "Staircase Up"`) on a floor cell near the top-left room; stores position in `DungeonLevel.UpStair`
6. If `depth < maxDepth`, places a down-staircase object (`Char: '>'`, `Color: "#e8c96a"`, `Blocking: false`, `Name: "Staircase Down"`) on a floor cell in a different room; stores position in `DungeonLevel.DownStair` and sets `HasDownStair = true`
7. Seeds torches (`Char: '†'`, `Color: "#e8c96a"`, `Blocking: true`, `Name: "Torch"`) on `CellWall` cells adjacent to `CellFloor` cells; approximately 1 torch per 5 rooms
8. Seeds braziers (`Char: 'Ω'`, `Color: "#e07030"`, `Blocking: false`, `Name: "Brazier"`) on `CellFloor` cells inside rooms; approximately 1 brazier per 6 rooms

#### Scenario: Same inputs produce the same level
- **WHEN** `GenerateDungeonLevel(42, 3, 7, 1, 5)` is called twice
- **THEN** both returned `DungeonLevel` values have identical `Cells` arrays

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
- **THEN** the cell it occupies has `Kind == CellWall` (torches are wall-mounted)

#### Scenario: No dungeon object has an empty Name
- **WHEN** `GenerateDungeonLevel` is called
- **THEN** every non-nil `Object` in any `DungeonCell` has a non-empty `Name`
