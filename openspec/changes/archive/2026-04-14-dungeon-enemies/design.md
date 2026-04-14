## Context

The dungeon system generates multi-level rooms connected by staircases. A `DungeonLevel` holds `Cells`, `UpStair`, `DownStair`, and `HasDownStair`. The combat system resolves fights via `resolveCombat(player, enemy Combatant, hooks, rng)` and already handles animal enemies via `buildEnemyCombatant(Animal)`. The inventory is capped at `InventoryMaxSlots = 8` items. The world tile's `Biome` is available at the dungeon entrance via `m.worldPos`.

The key design challenges are:
1. Biome is known at dungeon entry time but not stored on `DungeonLevel` — needs to be stored so levels can be regenerated deterministically with the right enemy roster.
2. Pathfinding must be simple enough to run every tick for multiple enemies without frame drops in a terminal.
3. Loot resolution must plug into the existing inventory without restructuring it.

## Goals / Non-Goals

**Goals:**
- Biome-keyed enemy roster (14 biomes → distinct named enemy types)
- Depth-scaled stats: enemies on level 1 are weak; enemies on the final level are dangerous
- BFS pathfinding limited to a sight radius (8 cells) so enemies only hunt when the player is nearby
- Auto-combat when player and enemy share a cell
- Weighted loot table per enemy type; at most one item dropped per fight; added to inventory if not full
- Enemies rendered on the dungeon map with a distinct glyph and colour
- `DungeonLevel` stores `Biome` so regenerated levels stay consistent

**Non-Goals:**
- Enemy ranged attacks or special abilities (future)
- Multiple enemies attacking simultaneously (fights are still 1v1)
- Persistent enemy state across level transitions (enemies respawn on re-entry — same deterministic seed)
- Enemy AI beyond line-of-sight BFS chase

## Decisions

### 1. DungeonEnemy is a runtime struct; enemies are not baked into the level seed

**Decision:** `DungeonEnemy` holds `X, Y int`, `Name string`, `Char rune`, `Color string`, plus a `*EnemyTemplate` pointer (stats, loot table). On first entry `GenerateDungeonLevel` populates `Enemies`; on re-entry (level already cached) the existing slice is used. If a player defeats an enemy and re-enters the level, enemies have respawned (same seed, same positions) — acceptable for this iteration.

**Rationale:** Keeping enemies out of the deterministic seed means we don't need to serialise kill state. Simpler, and respawn is a common roguelike convention.

**Alternative considered:** Storing defeated-enemy positions in a set and filtering at render time. Adds per-level persistent state; deferred.

### 2. Biome stored on DungeonMeta, passed to GenerateDungeonLevel

**Decision:** `DungeonMeta` gains a `Biome Biome` field set on first dungeon entry from `TileAt(m.worldPos).Biome`. `GenerateDungeonLevel` signature becomes `GenerateDungeonLevel(globalSeed, wx, wy, depth, maxDepth int, biome Biome) *DungeonLevel`. All call sites updated.

**Rationale:** The biome is fixed per entrance. Storing it in `DungeonMeta` (which is already keyed by `WorldCoord`) ensures all depths of the same dungeon share the same biome.

### 3. Stat scaling formula: linear interpolation by depth fraction

**Decision:** For a level at `depth d` out of `maxDepth D`:
```
fraction = float64(d-1) / float64(max(D-1, 1))   // 0.0 at depth 1, 1.0 at max depth
HP       = baseHP + int(fraction * (maxHP - baseHP))
Damage   = similarly scaled
```
`baseHP` and `maxHP` are defined per enemy template. This gives a smooth ramp from easy (depth 1) to hard (max depth) using values from the template rather than a global formula.

**Rationale:** Simple, testable, no magic constants at call sites. Each template author controls the difficulty range.

### 4. BFS pathfinding with sight radius cap

**Decision:** Each tick, for each enemy, if the Chebyshev distance to the player is ≤ 8, run a BFS from the enemy's cell over `CellFloor` cells (walls block, other enemies block) to find the shortest path to the player; move one step. If distance > 8, the enemy stays put (idle). BFS is bounded to the sight radius so worst-case nodes visited = π × 8² ≈ 200 per enemy.

**Rationale:** BFS guarantees optimal paths through winding corridors. The sight-radius bound keeps per-tick cost predictable regardless of map size. A\* would be faster in open areas but adds implementation complexity with no perceptible benefit at terminal frame rates.

**Alternative considered:** Dijkstra's from player outward (precomputed, amortised). More optimal but requires storing a distance map on the model and invalidating it on every player move — more state to manage.

### 5. Auto-combat triggers on cell overlap, not on adjacent cells

**Decision:** After each enemy move, if the enemy's `(X, Y)` equals `m.playerPos`, combat initiates immediately (same flow as `g`-on-animal: build combatants, call `resolveCombat`, set `ScreenCombat`). Player moving into an enemy cell also triggers combat. Multiple enemies can be adjacent without triggering.

**Rationale:** Consistent with the existing "occupying the same cell" combat model. Avoids needing to reason about facing or range.

### 6. Loot: at most one item per fight, resolved after victory in combat dismiss handler

**Decision:** `EnemyTemplate.LootTable []LootEntry` where `LootEntry { Item Item; Weight int }`. After victory, roll a weighted random choice from the table; if the result is non-empty and inventory is not full, add to inventory. Roll happens in the `ScreenCombat` dismiss handler using `rand.New(rand.NewSource(rand.Int63()))`.

**Rationale:** One item per fight keeps inventory pressure manageable. The dismiss handler already runs post-combat cleanup; loot resolution fits naturally there. Using a fresh RNG (not the combat RNG) means loot isn't correlated with fight outcomes.

## Risks / Trade-offs

- **[Multiple enemies per room slow the tick]** → Capped BFS + sight radius means ≤ 8–10 enemies active at once on a level with no per-tick cost for idle enemies.
- **[Respawn on re-entry may feel cheap]** → Accepted for v1; persistent kill state is a future change.
- **[GenerateDungeonLevel signature change breaks existing call sites]** → All call sites are internal; updated in task list. Tests updated in parallel.
- **[Inventory full means loot is silently dropped]** → Acceptable; a HUD flash or log message can be added in a polish pass.

## Open Questions

- Should enemies be visible through walls in the render (they're on floor cells, so always visible if lit)? Answer: yes — same lighting rules as objects.
- Should enemies block player movement (force combat on move-into), or can the player pass through? Answer: move-into triggers combat; enemy cell is effectively blocking for the player.
