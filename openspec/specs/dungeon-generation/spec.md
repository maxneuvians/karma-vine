## Requirements

### Requirement: Dungeon types are defined
The system SHALL define the following types in `internal/game/types.go`:
- `CellKind int` with constants `CellWall`, `CellFloor`
- `DungeonCell struct { Kind CellKind; Object *Object }` — a single grid cell
- `DungeonLevel struct { Cells [DungeonW][DungeonH]DungeonCell; UpStair LocalCoord; DownStair LocalCoord; HasDownStair bool }` — one dungeon floor
- `dungeonKey struct { wx, wy, depth int }` — cache key
- Constants `DungeonW = 80`, `DungeonH = 24`

#### Scenario: DungeonLevel zero-value is safe
- **WHEN** a `DungeonLevel` is allocated with `&DungeonLevel{}`
- **THEN** all `Cells` default to `CellWall` (zero value) without panicking

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

### Requirement: DungeonMeta stores per-entrance max depth
The system SHALL define `DungeonMeta struct { MaxDepth int }` and store one per world coordinate in `Model.dungeonMeta map[WorldCoord]DungeonMeta`. When first generating a dungeon for a world tile, the system SHALL randomise `MaxDepth` in `[5, 10]` using the global seed.

#### Scenario: MaxDepth is in range [5, 10]
- **WHEN** a `DungeonMeta` is created for a new world tile
- **THEN** `meta.MaxDepth >= 5 && meta.MaxDepth <= 10`

#### Scenario: MaxDepth is stable across re-entries
- **WHEN** the player enters a dungeon, ascends, and enters again
- **THEN** the same `MaxDepth` is used (looked up from `dungeonMeta`, not re-randomised)

### Requirement: DungeonLevelFor accessor caches levels
The system SHALL provide `DungeonLevelFor(wx, wy, depth int, m *Model) *DungeonLevel` that looks up `m.dungeonCache`, calls `GenerateDungeonLevel` only on a cache miss, stores the result, and returns it.

#### Scenario: Cache miss triggers generation
- **WHEN** `DungeonLevelFor` is called for a key not in `m.dungeonCache`
- **THEN** a new `DungeonLevel` is generated, stored, and returned

#### Scenario: Cache hit returns same instance
- **WHEN** `DungeonLevelFor` is called twice for the same `(wx, wy, depth)`
- **THEN** both calls return a pointer to the same `DungeonLevel`

### Requirement: Unlit torches and braziers are pickupable
Torches and braziers that are generated in the unlit state (`Lit == false`) SHALL have `Pickupable == true` set at generation time. Torches and braziers that start lit (e.g., those adjacent to stairs) SHALL have `Pickupable == false` (or the default) to keep them as fixed light sources.

> **Note:** In the current implementation all torches and braziers start unlit. This requirement simply ensures `Pickupable: true` is set on those objects so the inventory system can pick them up.

#### Scenario: Unlit torch is pickupable
- **WHEN** a dungeon level is generated
- **THEN** every `Object` with `Name == "Torch"` and `Lit == false` has `Pickupable == true`

#### Scenario: Brazier is pickupable when unlit
- **WHEN** a dungeon level is generated
- **THEN** every `Object` with `Name == "Brazier"` and `Lit == false` has `Pickupable == true`

#### Scenario: Lit torch is not pickupable
- **WHEN** a torch has `Lit == true` (manually lit by player)
- **THEN** `Pickupable == false` (or unchanged from generation; lighting a torch removes its pickupability)
