## Why

Dungeons are currently empty rooms with no threat — descending feels consequence-free. Adding biome-themed enemies that actively hunt the player and drop loot creates a core dungeon-crawl loop: risk scales with depth, rewards scale with difficulty, and the world's surface biome gives each cave its own identity.

## What Changes

- New `DungeonEnemy` type spawned on dungeon levels at generation time; population and stats scale with `depth`
- Enemy roster is keyed on the world biome above the dungeon entrance — e.g. Jungle → Jungle Troll, Tundra → Frost Giant, Plains → Bandit, Desert → Sand Wraith
- Each enemy type has a `LootTable`: a weighted list of `Item` drops evaluated on defeat
- Enemy pathfinding: each `TickMsg` in `ModeDungeon` moves visible enemies one step toward the player using simple BFS on `CellFloor` cells
- Combat initiates automatically when an enemy steps onto the player's cell (or player onto enemy's cell)
- `DungeonLevel` gains an `Enemies []*DungeonEnemy` slice; generation populates it based on depth and biome
- Loot is added to `m.inventory` (up to `InventoryMaxSlots`) after a victorious combat
- Defeated enemy is removed from `DungeonLevel.Enemies`
- Dungeon rendering updated to draw enemy glyphs on floor cells

## Capabilities

### New Capabilities
- `dungeon-enemy-system`: `DungeonEnemy` type, biome-keyed roster, depth-scaled stat generation, loot tables, per-tick BFS pathfinding, auto-combat trigger
- `dungeon-enemy-loot`: Loot table definition, weighted random drop resolution, inventory integration after combat victory

### Modified Capabilities
- `dungeon-generation`: `GenerateDungeonLevel` gains a `biome Biome` parameter and populates `DungeonLevel.Enemies` after room placement
- `dungeon-rendering`: Render loop draws `DungeonEnemy` glyphs on their current cells
- `dungeon-navigation`: Descend/ascend handler passes biome to `GenerateDungeonLevel`; enemy-onto-player and player-onto-enemy cells trigger combat
- `combat-system`: `buildEnemyCombatant` gains an overload/path for `DungeonEnemy` (in addition to the existing `Animal` path); loot is resolved and granted after victory
- `animal-system`: No change to existing surface animals

## Impact

- `internal/game/types.go` — `DungeonEnemy` struct, `LootEntry`, `LootTable` type, `DungeonLevel.Enemies` field
- `internal/game/dungeon.go` (or `dungeon_gen.go`) — `GenerateDungeonLevel` signature change, enemy spawning logic, biome roster table
- `internal/game/enemy.go` — new file: BFS pathfinding, per-tick enemy movement, auto-combat trigger
- `internal/game/combat.go` — `buildDungeonEnemyCombatant`, loot resolution
- `internal/game/render.go` — dungeon render loop draws enemy glyphs
- `internal/game/input.go` — movement into enemy cell triggers combat; biome passed at dungeon entry
- No new dependencies
