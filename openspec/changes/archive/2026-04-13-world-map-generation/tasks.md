## 1. Type Definitions

- [x] 1.1 Define `Biome` type as `int` in `types.go` with nine named constants matching the brief
- [x] 1.2 Define `Tile` struct in `types.go` with fields `Biome`, `Char rune`, `Color string`, `Elevation float64`, `Moisture float64`
- [x] 1.3 Replace the placeholder `Chunk` struct (from `project-scaffold`) with `Chunk{ Tiles [32][32]Tile }` in `world.go`

## 2. Noise Initialisation

- [x] 2.1 In `world.go`, create a package-level or per-generation function that instantiates two `opensimplex.New(int64)` noise objects — one seeded with `globalSeed`, one with `globalSeed+1`
- [x] 2.2 Define the constant `WorldNoiseScale = 0.07`

## 3. Biome Classification

- [x] 3.1 Implement `classifyBiome(e, m float64) (Biome, rune, string)` in `world.go` using the nine threshold conditions from the brief in the correct evaluation order
- [x] 3.2 Write a unit test in `world_test.go` covering at minimum: `DeepOcean` (e=0.20), `Forest` (e=0.48, m=0.60), `Snow` (e=0.85)

## 4. Chunk Generation

- [x] 4.1 Implement `generateChunk(cx, cy, globalSeed int) *Chunk` that iterates all 32×32 positions, computes world-space coordinates, samples both noise functions, calls `classifyBiome`, and populates each `Tile`
- [x] 4.2 Ensure noise is sampled at `float64(worldX) * WorldNoiseScale`, `float64(worldY) * WorldNoiseScale`

## 5. TileAt Accessor

- [x] 5.1 Implement `TileAt(worldX, worldY int, m *Model) Tile` that computes `cx = worldX / 32`, `cy = worldY / 32` (handle negative coordinates with `math.Floor` division)
- [x] 5.2 Check `m.chunks[ChunkCoord{cx, cy}]`; if absent, call `generateChunk` and store the result
- [x] 5.3 Return `m.chunks[ChunkCoord{cx, cy}].Tiles[localX][localY]` where `localX = worldX - cx*32`

## 6. Verification

- [x] 6.1 Run `go test ./...` and confirm all tests pass
- [x] 6.2 Add a quick smoke test in `world_test.go` that calls `TileAt` 1 000 times over a large coordinate range and asserts no panics and all returned biomes are valid constants
