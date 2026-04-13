## Why

The current biome classification uses only elevation and moisture, producing uniform biome bands regardless of latitude. This means the world has no sense of climate gradient — jungles and tundra cannot coexist on the same map, and the tropics look identical to the poles. Adding a temperature axis derived from latitude (Y coordinate) gives the world genuine climate zones that feel geographically believable.

## What Changes

- Add a `temperature float64` field to `Tile` computed from world Y coordinate and elevation
- Introduce five new biomes: `Jungle`, `Savanna`, `AridSteppe`, `Tundra`, `Taiga`
- Modify `classifyBiome` to accept temperature as a third axis alongside elevation and moisture
- Revise the biome thresholds to use a three-axis (elevation, moisture, temperature) classification table
- `Tile` struct gains a `Temperature float64` field in `[0, 1]` (0 = cold pole, 1 = hot equator)
- `generateChunk` computes temperature from a latitudinal gradient blended with a minor noise offset for variation

## Capabilities

### New Capabilities

- `temperature-zones`: Temperature axis derived from latitude (Y coordinate) and a small noise perturbation; drives climate-zone biome classification with hot equatorial, temperate mid-latitude, and cold polar bands

### Modified Capabilities

- `world-map-generation`: Biome classification gains a temperature parameter; nine existing biomes are remapped and five new ones added; `Tile` struct gains `Temperature float64`

## Impact

- `internal/game/types.go` — `Tile` struct, `Biome` constants
- `internal/game/world.go` — `generateChunk`, `classifyBiome` signature and logic
- `internal/game/render.go` — new biome name/colour/char entries
- `internal/game/world_test.go` — `TestClassifyBiome_*` tests updated; new biome tests added
- No changes to local map generation, rendering pipeline, or input system
