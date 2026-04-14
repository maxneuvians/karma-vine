## 1. Types

- [x] 1.1 Add `LootEntry struct { Item Item; Weight int }` to `types.go`
- [x] 1.2 Add `EnemyTemplate struct { Name string; Char rune; Color string; BaseHP, MaxHP, BaseArmour, MaxArmour, BaseMinDamage, MaxMinDamage, BaseMaxDamage, MaxMaxDamage, BaseInitiative, MaxInitiative int; LootTable []LootEntry }` to `types.go`
- [x] 1.3 Add `DungeonEnemy struct { X, Y int; Template *EnemyTemplate; HP, MaxHP, Armour, MinDamage, MaxDamage, Initiative int }` to `types.go`
- [x] 1.4 Add `Enemies []*DungeonEnemy` field to `DungeonLevel` in `types.go`
- [x] 1.5 Add `Biome Biome` field to `DungeonMeta` in `types.go`

## 2. Model

- [x] 2.1 Add `combatDungeonEnemy *DungeonEnemy` field to `Model` in `model.go`

## 3. Enemy Roster and Spawn

- [x] 3.1 Create `internal/game/enemy.go` with `dungeonEnemyRoster map[Biome]*EnemyTemplate` covering all 14 biomes per spec (Forest/DenseForest → Goblin, Plains → Bandit, Desert/AridSteppe → Sand Wraith, Jungle/Savanna → Jungle Troll, Tundra → Frost Giant, Taiga/Snow → Ice Wraith, Mountains → Stone Golem, DeepOcean/ShallowWater/Beach → Cave Crustacean); include a fallback entry
- [x] 3.2 Define loot tables for each template: each table SHALL include at least one no-drop entry (`Item{}`, `Weight: 3`) and 1–3 themed item drops (e.g. Goblin → `{"Goblin Ear", Weight:2}`, `{"Rusty Dagger", Weight:1}`)
- [x] 3.3 Implement `spawnEnemy(tmpl *EnemyTemplate, x, y, depth, maxDepth int) *DungeonEnemy` using linear interpolation: `fraction = float64(depth-1)/float64(max(maxDepth-1,1))`; compute each stat as `base + int(fraction*float64(maxStat-base))`; set `HP = MaxHP`

## 4. Dungeon Generation — Biome Parameter and Enemy Spawning

- [x] 4.1 Update `GenerateDungeonLevel` signature to `GenerateDungeonLevel(globalSeed, wx, wy, depth, maxDepth int, biome Biome) *DungeonLevel`
- [x] 4.2 After all room/corridor/object placement, look up `dungeonEnemyRoster[biome]` (with fallback), collect all `CellFloor` positions not occupied by objects, shuffle them using the level RNG, pick the first `depth` positions, call `spawnEnemy` for each, store in `level.Enemies`
- [x] 4.3 Update `DungeonLevelFor` to pass `biome` from `m.dungeonMeta[key].Biome` to `GenerateDungeonLevel`
- [x] 4.4 Update the dungeon entry handler (`ModeLocal` `enter`/`>`) to record `Biome: TileAt(m.worldPos, &m).Biome` into `DungeonMeta` on first entry

## 5. Enemy Pathfinding (per-tick)

- [x] 5.1 Implement `moveEnemies(m Model) Model` in `enemy.go`: iterate `m.currentDungeon.Enemies`; for each, if Chebyshev distance to `m.playerPos` > 8, skip; otherwise run BFS from enemy cell over `CellFloor` cells (skip cells occupied by other enemies), find next step toward player, move enemy; if after moving an enemy's position equals `m.playerPos`, initiate combat (build combatants using `buildDungeonEnemyCombatant`, call `resolveCombat`, set `m.combatState`, `m.combatDungeonEnemy`, `m.screenMode = ScreenCombat`, `m.paused = true`) and break from the loop
- [x] 5.2 In `Update`'s `TickMsg` handler, call `moveEnemies(m)` when `m.mode == ModeDungeon && !m.paused`

## 6. Player Movement — Combat on Enemy Cell

- [x] 6.1 In `handleKey` dungeon movement branches (`up`/`down`/`left`/`right` in `ModeDungeon`), before applying the move, check if the target cell contains a `DungeonEnemy`; if so, initiate combat with that enemy (`buildDungeonEnemyCombatant`, `resolveCombat`, set `ScreenCombat`) and return without moving the player

## 7. Combat — DungeonEnemy Combatant Builder

- [x] 7.1 Implement `buildDungeonEnemyCombatant(e *DungeonEnemy) Combatant` in `combat.go`: copy all stat fields from the enemy; `Name = e.Template.Name`

## 8. Combat — Loot Resolution

- [x] 8.1 Implement `resolveEnemyLoot(table []LootEntry, rng *rand.Rand) Item` in `enemy.go` (or `combat.go`): compute total weight, pick random value in `[0, totalWeight)`, walk entries to find winner; return zero `Item` if table empty or total weight 0
- [x] 8.2 In the `ScreenCombat` dismiss handler, after `PlayerWon == true` cleanup: if `m.combatDungeonEnemy != nil`, call `resolveEnemyLoot` with a fresh RNG; if result `Name != ""` and inventory not full, add to inventory (stack if same Name exists, else append `Count=1`); then remove enemy from `m.currentDungeon.Enemies`; clear `m.combatDungeonEnemy`
- [x] 8.3 Ensure the existing `*Animal` victory path (remove from `localMap.Animals`) is unchanged and still runs when `m.combatEnemy != nil`

## 9. Rendering

- [x] 9.1 In `renderDungeonMap`, after drawing the floor/object/staircase layers but before drawing the player, iterate `m.currentDungeon.Enemies`; for each enemy whose cell is in the visibility set and whose `(X, Y)` does not equal `m.playerPos`, draw `enemy.Template.Char` in `enemy.Template.Color`

## 10. Tests

- [x] 10.1 `TestSpawnEnemy_BaseStats` in `enemy_test.go`: depth 1, maxDepth 5 → stats equal `Base` values
- [x] 10.2 `TestSpawnEnemy_MaxStats` in `enemy_test.go`: depth 5, maxDepth 5 → stats equal `Max` values
- [x] 10.3 `TestSpawnEnemy_MidStats` in `enemy_test.go`: depth 3, maxDepth 5 → stats between Base and Max
- [x] 10.4 `TestDungeonEnemyRoster_JungleTroll` in `enemy_test.go`: `dungeonEnemyRoster[Jungle].Name == "Jungle Troll"`
- [x] 10.5 `TestDungeonEnemyRoster_FrostGiant` in `enemy_test.go`: `dungeonEnemyRoster[Tundra].Name == "Frost Giant"`
- [x] 10.6 `TestResolveEnemyLoot_SingleEntry` in `enemy_test.go`: single entry always returns that item
- [x] 10.7 `TestResolveEnemyLoot_ZeroWeight` in `enemy_test.go`: zero-weight table returns empty item
- [x] 10.8 `TestResolveEnemyLoot_NoDrop` in `enemy_test.go`: no-drop entry (empty Name) can be returned
- [x] 10.9 `TestGenerateDungeonLevel_EnemyCount` in `dungeon_test.go` (or `enemy_test.go`): `depth=3` → `len(level.Enemies) == 3`
- [x] 10.10 `TestGenerateDungeonLevel_EnemiesOnFloor` in `dungeon_test.go`: every enemy cell is `CellFloor`
- [x] 10.11 `TestGenerateDungeonLevel_Deterministic_WithBiome` in `dungeon_test.go`: same args produce same enemy positions
- [x] 10.12 `TestMoveEnemies_ApproachesPlayer` in `enemy_test.go`: enemy within sight radius moves closer each tick
- [x] 10.13 `TestMoveEnemies_IdleOutsideRadius` in `enemy_test.go`: enemy > 8 cells away does not move
- [x] 10.14 `TestMoveEnemies_PausedNoMove` in `enemy_test.go`: `m.paused == true` → enemy positions unchanged
- [x] 10.15 `TestMoveEnemies_TriggersScreenCombat` in `enemy_test.go`: enemy reaches player cell → `screenMode == ScreenCombat`
- [x] 10.16 `TestHandleKey_MovingIntoEnemyTriggersCombat` in `input_test.go`: directional move into enemy cell → `screenMode == ScreenCombat`, player pos unchanged
- [x] 10.17 `TestBuildDungeonEnemyCombatant` in `combat_test.go`: stat fields match enemy; Name matches template
- [x] 10.18 `TestLootAddedToInventoryAfterVictory` in `input_test.go` or `game_test.go`: victory with non-empty loot and space → item appears in inventory
- [x] 10.19 `TestLootDiscardedWhenInventoryFull` in same file: victory with full inventory → inventory unchanged
- [x] 10.20 `TestDungeonMeta_BiomeRecorded` in `input_test.go`: first dungeon descent records correct biome in `DungeonMeta`
- [x] 10.21 Update existing `TestGenerateDungeonLevel_*` tests to pass `biome` argument (e.g. `Plains`) to the updated signature
