## Context

The world map must feel infinite and consistent. The player can wander arbitrarily far from the origin and always get the same world back when they return. OpenSimplex noise at a fixed scale provides smooth biome transitions without tile-boundary seams. Chunks are the unit of lazy generation: only the chunks visible in the current viewport are ever computed, keeping memory bounded during a typical session.

## Goals / Non-Goals

**Goals:**
- Two independent `opensimplex.New(seed)` instances per world — one for elevation, one for moisture — so the two noise fields are uncorrelated
- Chunk generation is a pure function of `(chunkX, chunkY, globalSeed)` — no randomness beyond the seed
- `TileAt` is the single point of access; callers never touch `Model.chunks` directly
- All nine biome variants from the brief are covered with correct thresholds and display characters

**Non-Goals:**
- Rendering (separate change)
- Local map content inside biomes (separate change)
- Points of interest / noise-peak landmarks (v1 out of scope)

## Decisions

**Two separate noise objects (elevation + moisture)** — Using a single noise object with an offset would introduce subtle correlations at large distances. Separate objects with seeds derived from `globalSeed` and `globalSeed+1` give true independence. Alternative: single noise object with 2D offset — rejected due to potential periodicity artifacts.

**Chunk size 32×32** — Matches the brief exactly. 1 024 tiles per chunk is cheap to generate and fits comfortably in a single allocation. Alternative: variable chunk size — unnecessary complexity.

**`map[ChunkCoord]*Chunk` with pointer values** — Pointer semantics mean the map never copies 1 024-tile arrays. Alternative: value semantics — would copy the entire array on every map lookup.

**Biome assigned at generation time, not at render time** — Storing `Biome` on the tile avoids re-running threshold logic during every frame render. Trade-off: slightly more memory per tile, but the render loop stays O(viewport) with no branching on noise values.

## Risks / Trade-offs

- **Global seed = 0 produces valid but fixed world** → Acceptable; callers should seed from `time.Now().UnixNano()` in `main.go`.
- **Unbounded chunk cache** → Memory grows proportionally to explored area. Mitigation: document that a cache eviction strategy is a v2 concern; the brief makes no requirement for it.
- **Noise scale 0.07 is hardcoded** → Changing it would alter the world. Acceptable as a named constant `WorldNoiseScale`.
