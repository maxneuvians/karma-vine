## Why

The world map is the top-level navigation layer of the explorer. Without it there is nothing to see or move through. The infinite, chunk-based approach with deterministic biome assignment is the core technical bet of the project: it lets the world feel vast without storing any state beyond the global seed.

## What Changes

- Implement `ChunkCoord` → `*Chunk` generation using two OpenSimplex noise passes (elevation `e`, moisture `m`) at scale `0.07`
- Each `Chunk` is 32×32 `Tile` cells; chunks generated lazily on demand and cached in `Model.chunks`
- Assign biome to each tile using the nine-way threshold table from the brief (`DeepOcean`, `ShallowWater`, `Beach`, `Forest`, `Plains`, `DenseForest`, `Desert`, `Mountains`, `Snow`)
- Each `Tile` stores: `Biome`, display `Char rune`, `Color string`, `Elevation float64`, `Moisture float64`
- Expose a `TileAt(worldX, worldY int, m *Model) Tile` helper that computes the chunk coordinate, generates the chunk if absent, and returns the tile

## Capabilities

### New Capabilities
- `world-map-generation`: chunk-based infinite world map — noise sampling, chunk caching, biome classification, tile data, and the `TileAt` accessor

### Modified Capabilities
<!-- none — project-scaffold defines the Chunk placeholder; this change replaces it with the real implementation -->

## Impact

- Replaces the placeholder `Chunk` struct from `project-scaffold` with a fully populated type
- No external API changes; purely internal to the Go package
- Adds `opensimplex-go` usage (already in `go.mod` from `project-scaffold`)
