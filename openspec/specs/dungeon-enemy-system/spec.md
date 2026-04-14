## ADDED Requirements

### Requirement: DungeonEnemy and EnemyTemplate types are defined
The system SHALL define in `types.go`:
- `LootEntry struct { Item Item; Weight int }` — one entry in a weighted loot table
- `EnemyTemplate struct { Name string; Char rune; Color string; BaseHP int; MaxHP int; BaseArmour int; MaxArmour int; BaseMinDamage int; MaxMinDamage int; BaseMaxDamage int; MaxMaxDamage int; BaseInitiative int; MaxInitiative int; LootTable []LootEntry }` — archetype for one enemy kind
- `DungeonEnemy struct { X, Y int; Template *EnemyTemplate; HP int; MaxHP int; Armour int; MinDamage int; MaxDamage int; Initiative int }` — a live enemy instance on a level
- `DungeonLevel` SHALL gain an `Enemies []*DungeonEnemy` field

#### Scenario: DungeonEnemy zero-value is safe
- **WHEN** a `DungeonEnemy` is allocated with `&DungeonEnemy{}`
- **THEN** it does not panic

#### Scenario: DungeonLevel.Enemies is non-nil after generation
- **WHEN** `GenerateDungeonLevel` returns a level
- **THEN** `level.Enemies` is non-nil (may be empty slice)

### Requirement: Biome-keyed enemy roster maps each biome to an EnemyTemplate
The system SHALL define `dungeonEnemyRoster map[Biome]*EnemyTemplate` covering all 14 biomes. Each entry SHALL have distinct `Name`, `Char`, and `Color` values. A fallback entry SHALL be used for any biome not explicitly listed. The roster SHALL include at minimum:
- `Forest`, `DenseForest` → `Goblin` (`g`, `#4a9e4a`)
- `Plains` → `Bandit` (`b`, `#c8a060`)
- `Desert`, `AridSteppe` → `Sand Wraith` (`w`, `#d4b050`)
- `Jungle`, `Savanna` → `Jungle Troll` (`T`, `#2a7a2a`)
- `Tundra` → `Frost Giant` (`G`, `#80c0e0`)
- `Taiga`, `Snow` → `Ice Wraith` (`W`, `#b0d8f0`)
- `Mountains` → `Stone Golem` (`O`, `#888888`)
- `DeepOcean`, `ShallowWater`, `Beach` → `Cave Crustacean` (`c`, `#c06030`)

#### Scenario: Jungle biome resolves to Jungle Troll
- **WHEN** `dungeonEnemyRoster[Jungle]` is looked up
- **THEN** the returned template has `Name == "Jungle Troll"`

#### Scenario: Tundra biome resolves to Frost Giant
- **WHEN** `dungeonEnemyRoster[Tundra]` is looked up
- **THEN** the returned template has `Name == "Frost Giant"`

### Requirement: spawnEnemy creates a DungeonEnemy with depth-scaled stats
`spawnEnemy(tmpl *EnemyTemplate, x, y, depth, maxDepth int) *DungeonEnemy` SHALL compute stats by linearly interpolating between `tmpl.BaseX` and `tmpl.MaxX` using `fraction = float64(depth-1)/float64(max(maxDepth-1,1))`. Stats SHALL be floored at their `Base` value. The resulting `DungeonEnemy` SHALL have `X=x`, `Y=y`, `Template=tmpl`, and all combat fields set to the interpolated values. `HP` SHALL equal `MaxHP` (full health at spawn).

#### Scenario: depth 1 yields base stats
- **WHEN** `spawnEnemy(tmpl, 0, 0, 1, 5)` is called with `tmpl.BaseHP=10`, `tmpl.MaxHP=30`
- **THEN** the returned enemy has `HP=10` and `MaxHP=10`

#### Scenario: max depth yields max stats
- **WHEN** `spawnEnemy(tmpl, 0, 0, 5, 5)` is called with `tmpl.BaseHP=10`, `tmpl.MaxHP=30`
- **THEN** the returned enemy has `HP=30` and `MaxHP=30`

#### Scenario: mid depth interpolates stats
- **WHEN** `spawnEnemy(tmpl, 0, 0, 3, 5)` is called with `tmpl.BaseHP=10`, `tmpl.MaxHP=30`
- **THEN** the returned enemy has `HP` between 10 and 30

### Requirement: Enemies move one step per tick toward visible players via BFS
On each `TickMsg` in `ModeDungeon`, for each enemy in `m.currentDungeon.Enemies`, the system SHALL compute Chebyshev distance to `m.playerPos`. If distance ≤ 8 (sight radius), the system SHALL run BFS from the enemy's cell over `CellFloor` cells (walls block; other enemy cells block) to find the next step toward the player and move the enemy to that cell. If distance > 8 the enemy SHALL remain stationary. If BFS finds no path the enemy SHALL remain stationary.

#### Scenario: Enemy within sight radius moves closer each tick
- **WHEN** an enemy is 3 cells away from the player and has a clear floor path
- **THEN** after one `TickMsg` the enemy's Chebyshev distance to the player is 2

#### Scenario: Enemy outside sight radius does not move
- **WHEN** an enemy is 10 cells away from the player
- **THEN** after one `TickMsg` the enemy's position is unchanged

#### Scenario: Enemy blocked by wall does not move through it
- **WHEN** an enemy is adjacent to the player but separated by a `CellWall`
- **THEN** after one `TickMsg` the enemy has not moved through the wall

#### Scenario: Enemy movement is suppressed when game is paused
- **WHEN** `m.paused == true` and a `TickMsg` is received
- **THEN** all enemy positions remain unchanged

### Requirement: Auto-combat triggers when enemy and player share a cell
After each enemy move step, if the enemy's `(X, Y)` equals `m.playerPos`, the system SHALL immediately initiate combat: build a `Combatant` from the enemy, build the player combatant, call `resolveCombat`, store result in `m.combatState`, store the enemy in `m.combatEnemy`, set `m.screenMode = ScreenCombat`, and `m.paused = true`. Only the first overlapping enemy SHALL trigger combat per tick (subsequent enemies skip their move).

#### Scenario: Enemy reaching player cell triggers ScreenCombat
- **WHEN** an enemy moves onto `m.playerPos`
- **THEN** `m.screenMode == ScreenCombat` and `m.combatState` is non-nil

#### Scenario: Player moving onto enemy cell triggers combat
- **WHEN** the player presses a movement key and the target cell contains a `DungeonEnemy`
- **THEN** `m.screenMode == ScreenCombat`

### Requirement: Defeated dungeon enemy is removed from level
When combat resolves with `PlayerWon == true` and `m.combatEnemy` references a `*DungeonEnemy`, the system SHALL remove that enemy from `m.currentDungeon.Enemies`. The `*Animal` path for surface combat is unaffected.

#### Scenario: Defeated enemy no longer appears on level
- **WHEN** the player defeats a `DungeonEnemy`
- **THEN** the enemy is absent from `m.currentDungeon.Enemies` after dismissing the result
