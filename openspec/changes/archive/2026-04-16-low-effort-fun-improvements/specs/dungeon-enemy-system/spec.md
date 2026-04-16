## MODIFIED Requirements

### Requirement: GenerateDungeonLevel spawns a minimum of 3 enemies per floor
`GenerateDungeonLevel` SHALL compute enemy count as `max(3, depth)`, capped by the number of available enemy spawn positions. The previous behaviour of `enemyCount = depth` (yielding 1 enemy on floor 1) is replaced.

#### Scenario: Floor 1 spawns at least 3 enemies
- **WHEN** `GenerateDungeonLevel` is called with `depth == 1` and sufficient room positions exist
- **THEN** `len(level.Enemies) == 3`

#### Scenario: Floor 5 spawns 5 enemies
- **WHEN** `GenerateDungeonLevel` is called with `depth == 5` and sufficient room positions exist
- **THEN** `len(level.Enemies) == 5`

#### Scenario: Enemy count is still capped by available positions
- **WHEN** `GenerateDungeonLevel` is called and only 2 room positions are available
- **THEN** `len(level.Enemies) == 2` (capped by positions, not floored)
