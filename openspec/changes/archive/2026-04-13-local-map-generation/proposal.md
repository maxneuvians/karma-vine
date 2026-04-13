## Why

When the player descends into a world map tile the explorer must show a detailed 42×18 local map. This map must be deterministic (same world tile always produces the same local map) and biome-appropriate (a forest tile produces a forest floor, a desert tile produces sand and cacti). Without this, the "drill-down" half of the two-tier architecture is non-functional.

## What Changes

- Define the `Ground`, `Object`, and `Animal` data types that make up local map cells
- Implement the `hash(x, y, seed int) float64` function from the brief for deterministic local seeding
- Implement `GenerateLocalMap(worldX, worldY, globalSeed int, biome Biome) *LocalMap` using two noise passes at scale `0.12`
- Populate ground tiles, objects (placed where noise exceeds biome thresholds), and initial animal placements per biome using the tables in the brief
- Cache generated local maps in `Model.localCache`; revisiting a world tile re-uses the cached map

## Capabilities

### New Capabilities
- `local-map-generation`: deterministic local map generation — `Ground`/`Object`/`Animal` types, hash-based seeding, biome content tables, noise-driven object placement, `LocalMap` struct, and `LocalMapFor` accessor

### Modified Capabilities
<!-- none — project-scaffold defines a LocalMap placeholder; this change replaces it -->

## Impact

- Replaces the placeholder `LocalMap` struct from `project-scaffold` with the full implementation
- Depends on `Biome` constants introduced by `world-map-generation`
- Animal tick / movement logic is intentionally out of scope here (handled by `animal-system`)
