## ADDED Requirements

### Requirement: DungeonEnemy glyphs are rendered on visible cells
The dungeon render loop SHALL draw each `DungeonEnemy` in `m.currentDungeon.Enemies` whose cell `(X, Y)` is in the visibility set. The enemy SHALL be rendered using `enemy.Template.Char` in `enemy.Template.Color`. The player glyph (`@`) takes priority over all other layers; enemies take priority over floor/object glyphs on the same cell but are overridden by the player.

#### Scenario: Enemy on visible floor cell is rendered
- **WHEN** a `DungeonEnemy` is at a cell within `playerViewRadius`
- **THEN** the cell renders `enemy.Template.Char` in `enemy.Template.Color`

#### Scenario: Enemy on non-visible cell is not rendered
- **WHEN** a `DungeonEnemy` is at a cell outside the visibility set
- **THEN** the cell renders as a blank space (fog of war)

#### Scenario: Player glyph overrides enemy glyph on same cell
- **WHEN** the player and an enemy occupy the same cell (mid-combat transition)
- **THEN** `'@'` is rendered, not the enemy glyph
