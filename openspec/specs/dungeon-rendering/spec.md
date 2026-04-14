## Requirements

### Requirement: Dungeon map is rendered with fog-of-war
The system SHALL implement `renderDungeonMap(m Model) string` that renders the dungeon grid. Cells outside the visible set SHALL be rendered as a space (` `) with no styling. Cells inside the visible set SHALL be rendered based on their content. The rendered output SHALL be clipped to `m.viewportW × (m.viewportH - 1)`.

#### Scenario: Cells far from player are hidden
- **WHEN** the dungeon is rendered and a cell is more than `playerViewRadius` steps from the player and not within `torchRadius` of any torch/brazier
- **THEN** that cell renders as a blank space

#### Scenario: Cells within player radius are visible
- **WHEN** the dungeon is rendered and a cell is within `playerViewRadius` (6) of the player
- **THEN** that cell renders its content (wall, floor, object, or player)

### Requirement: Dungeon visibility is computed by radius
The system SHALL compute a visibility set on each render frame. A cell is visible if:
- Its Chebyshev distance to `m.playerPos` is ≤ `playerViewRadius` (constant 6), OR
- Its Chebyshev distance to any torch or brazier object on the current level is ≤ `torchRadius` (constant 4)

#### Scenario: Torch illuminates surrounding cells
- **WHEN** a torch is at `(10, 10)` and `playerViewRadius` would not reach `(12, 12)`
- **THEN** `(12, 12)` is still visible because its Chebyshev distance to the torch is 2 (≤ `torchRadius` 4)

#### Scenario: Player always visible to themselves
- **WHEN** the dungeon is rendered
- **THEN** the player's own cell is always in the visible set

### Requirement: Dungeon cell glyphs and colors
The system SHALL render dungeon cells as follows:
- `CellWall`: `'#'` in `"#666666"`
- `CellFloor`: `'.'` in `"#444444"`
- Wall torch object on a `CellWall` cell: `'†'` in `"#e8c96a"`
- Brazier object on a `CellFloor` cell: `'Ω'` in `"#e07030"`
- Up-staircase object: `'<'` in `"#e8c96a"`
- Down-staircase object: `'>'` in `"#e8c96a"`
- Player at `playerPos`: `'@'` in `"#f0f6fc"` bold, overrides all other layers

#### Scenario: Player glyph overrides floor glyph
- **WHEN** `playerPos` is at a `CellFloor` cell
- **THEN** `'@'` is rendered at that position, not `'.'`

#### Scenario: Player can stand on brazier cell
- **WHEN** `playerPos` is at a cell containing a brazier object
- **THEN** `'@'` is rendered (player glyph, not brazier), and the brazier is not blocking

### Requirement: Dungeon HUD shows depth and mode
The system SHALL extend `renderHUD` so that when `m.mode == ModeDungeon`, the status bar displays:
- The text `Dungeon` as the mode label
- Current depth as `Depth: N` where N is `m.dungeonDepth`
- World coordinates `(wx, wy)`

#### Scenario: HUD shows depth in dungeon mode
- **WHEN** `m.mode == ModeDungeon` and `m.dungeonDepth == 3`
- **THEN** the status bar contains `Depth: 3`

#### Scenario: HUD shows dungeon label
- **WHEN** `m.mode == ModeDungeon`
- **THEN** the status bar contains the text `Dungeon`

### Requirement: Key bar shows dungeon hints
The system SHALL update `renderKeyBar` so that when `m.mode == ModeDungeon`, the hint string includes `< up`, `> down`, and `esc exit`.

#### Scenario: Key bar shows dungeon hints in dungeon mode
- **WHEN** `m.mode == ModeDungeon`
- **THEN** the key bar contains `< up` and `> down`

#### Scenario: Key bar does not show dungeon hints in other modes
- **WHEN** `m.mode == ModeWorld`
- **THEN** the key bar does not contain `< up`
