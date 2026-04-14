## Context

The game currently has two active modes (`ModeWorld`, `ModeLocal`) and a rendering pipeline in `render.go` that dispatches based on `m.mode`. Local maps are 160×48 grids of `Ground`/`Object`/`Animal` cells stored in `localCache map[WorldCoord]*LocalMap`. Dungeon levels need a parallel structure: a grid of cells that persists across visits, keyed by `(worldX, worldY, depth)`.

The player descends from a local map into a dungeon via a staircase-down object already on the local map. Dungeon levels are self-contained: they have their own grid, their own object placements (torches/braziers), their own staircase objects, and their own fog-of-war visibility state.

## Goals / Non-Goals

**Goals:**
- New `ModeDungeon` tier that plugs cleanly into the existing mode enum and `buildView` dispatch
- Procedural BSP generation per level, seeded deterministically from `(globalSeed, worldX, worldY, depth)` so re-entering gives the same layout
- Fog-of-war: only cells within a radius of the player (or within radius of lit torches/braziers) are visible
- Wall objects (torches) and floor items (braziers) never occupy the same cell as the player; player can walk over floor items
- Dungeon entrance (`>`) injected into local map generation; max depth (5–10) randomised per entrance and stored with the cache entry
- Up-staircase (`<`) at dungeon level start position; down-staircase (`>`) in a generated room (absent on final level)
- `esc`/`<` ascends; `enter`/`>` descends
- ≥90% test coverage on new generation and rendering code

**Non-Goals:**
- Combat, inventory, or item pickup mechanics
- Lighting propagation through walls (simple radius-based lighting only)
- Persistent explored/revealed map (fog-of-war resets each visit is acceptable)
- Multiple dungeon entrances per local map (one per tile)
- Dungeon biome theming or branching paths

## Decisions

### Decision 1: New `DungeonLevel` type, not reusing `LocalMap`

`LocalMap` carries `Ground/Objects/Animals/LitMap` arrays sized to `LocalMapW×LocalMapH` (160×48). Dungeon levels need different dimensions (e.g., 80×24), have no animals, and need a `Visited` mask for fog-of-war. Reusing `LocalMap` would require nullable fields and confusing size assumptions.

**Alternative considered**: Tag `LocalMap` with a `isDungeon bool`. Rejected — coupling unrelated concerns, harder to test independently.

### Decision 2: BSP room generation

Binary Space Partitioning splits the level area recursively into sub-rectangles, then places a room inside each leaf. Rooms are connected by L-shaped corridors. This produces natural-looking room layouts without loops, avoids overlapping rooms without explicit collision checks, and is deterministic given a fixed seed.

**Alternative considered**: Drunkard's walk (cave carving). Produces organic caves rather than discrete rooms, harder to guarantee stair placement in accessible areas. Could be added later as a biome variant.

### Decision 3: Dungeon cache keyed by `(worldX, worldY, depth)`

Type `dungeonKey struct { wx, wy, depth int }` maps to `*DungeonLevel`. Stored in `Model.dungeonCache`. Levels are generated on first entry and never regenerated, matching the `localCache` pattern.

### Decision 4: Max depth stored in `DungeonMeta`, not `DungeonLevel`

A separate `DungeonMeta` (or `dungeonMeta map[WorldCoord]DungeonMeta`) stores the randomised max depth per entrance. This separates "how deep does this dungeon go?" from any individual level.

**Alternative considered**: Storing max depth on the deepest `DungeonLevel` and looking it up. Requires generating all intermediate levels first; wasteful.

### Decision 5: Fog-of-war as a computed visibility set, not persisted `LitMap`

On each render frame, compute the set of visible cells: cells within `playerViewRadius` (e.g., 6) of the player, plus cells within `torchRadius` (e.g., 4) of each lit torch/brazier on the current level. No persistent `VisitedMap` — keep it simple; the level layout is always re-renderable once generated.

**Alternative considered**: Persist a `Revealed [W][H]bool` array per level so explored areas stay visible even after moving away. This is a common roguelike pattern and should be added as a future enhancement; for now, simple radius visibility suffices.

### Decision 6: Dungeon dimensions 80×24

Smaller than local maps (160×48) to keep generation fast and ensure rooms are discoverable. Still fits within a typical terminal at normal zoom. Tunable constants `DungeonW = 80`, `DungeonH = 24`.

### Decision 7: Wall/floor item placement rules

Torches are placed on `CellWall` cells adjacent to open floor; they block passage (like `Object.Blocking = true`). Braziers are placed on `CellFloor` cells inside rooms; they do not block passage (`Blocking = false`). Player spawn point (`findDungeonSpawn`) scans outward from the up-staircase to find the nearest non-blocking floor cell.

## Risks / Trade-offs

[Determinism] Seeding `rand` with `(globalSeed ^ worldX*31 ^ worldY*97 ^ depth*7)` must produce the same level on re-entry. → Use `rand.New(rand.NewSource(int64(seed)))` in the generator, never the global `rand`.

[Terminal width] 80-column dungeon may overflow narrow terminals. → `renderDungeonMap` clips to `viewportW` like the local map renderer.

[Test coverage] Procedural generation is hard to unit test deterministically. → Test `GenerateDungeonLevel` with a fixed seed and assert: stair count, room count ≥ 1, all floor cells reachable from up-stair.

[Deep stack risk] BSP recursion depth ≤ log2(min(W,H)) ≈ 4–5 levels; no stack concern.

## Migration Plan

1. Add new types to `types.go` (no behaviour change).
2. Add `dungeonCache`, `dungeonDepth`, `dungeonMaxDepth`, `currentDungeon` fields to `Model` (zero values safe).
3. Add `dungeon.go` with generation logic (pure functions, fully testable in isolation).
4. Inject dungeon entrance into `local.go`'s `LocalMapFor`.
5. Extend `input.go` with `ModeDungeon` cases.
6. Add `renderDungeonMap` to `render.go`; extend `buildView` dispatch.
7. Run full test suite; add/fix tests to maintain ≥90% coverage.

Rollback: all changes are additive until step 4 (local map injection). Step 4 is a small, isolated change to `LocalMapFor` with test coverage — straightforward to revert.

## Open Questions

- Should explored cells remain visible after the player moves away (persistent fog)? → Deferred; start with radius-only visibility.
- Should torches be extinguishable? → Out of scope for this change.
- One dungeon entrance per local tile, or could there be cave mouths and separate dungeon entrances? → One per tile for now.
