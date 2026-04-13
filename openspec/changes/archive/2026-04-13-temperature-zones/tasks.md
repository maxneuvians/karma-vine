## 1. Types — Tile and Biome constants

- [x] 1.1 Add `Temperature float64` field to the `Tile` struct in `types.go`
- [x] 1.2 Add five new `Biome` constants to the iota block in `types.go`: `Jungle`, `Savanna`, `AridSteppe`, `Tundra`, `Taiga`

## 2. World generation — temperature derivation

- [x] 2.1 In `world.go`, add a `computeTemperature(worldY int, noise opensimplex.Noise) float64` helper that computes a cosine latitudinal gradient with half-period 800 tiles blended with noise perturbation amplitude 0.12, clamped to `[0, 1]`
- [x] 2.2 In `generateChunk`, instantiate a temperature noise object seeded with `globalSeed + 7`
- [x] 2.3 Call `computeTemperature` for each tile and store the result in `tile.Temperature`

## 3. Biome classification — three-axis table

- [x] 3.1 Update `classifyBiome` signature to `classifyBiome(e, m, temperature float64) (Biome, rune, string)`
- [x] 3.2 Implement the water/coast cases first (DeepOcean, ShallowWater, Beach) — unchanged regardless of temperature
- [x] 3.3 Implement the hot band (temperature ≥ 0.65): Jungle, Savanna, AridSteppe, Desert, Mountains, Snow
- [x] 3.4 Implement the temperate band (0.35 ≤ temperature < 0.65): Forest, Plains, DenseForest, Desert, Mountains, Snow
- [x] 3.5 Implement the cold band (temperature < 0.35): Taiga, Tundra, Mountains, Snow
- [x] 3.6 Update the `generateChunk` call site to pass `temperature` as third argument to `classifyBiome`

## 4. Rendering — new biome names and colours

- [x] 4.1 Add cases for all five new biomes to `biomeName()` in `render.go`
- [x] 4.2 Verify `renderWorldMap` and `renderLocalMap` require no changes (they use `tile.Char` and `tile.Color` which are already set by `classifyBiome`)

## 5. Tests

- [x] 5.1 Update `TestClassifyBiome_DeepOcean` to pass a temperature argument (e.g., 0.5)
- [x] 5.2 Update `TestClassifyBiome_Forest` to pass a temperate temperature argument (e.g., 0.50)
- [x] 5.3 Update `TestClassifyBiome_Snow` to pass a temperature argument (e.g., 0.5)
- [x] 5.4 Add `TestClassifyBiome_Jungle`: `classifyBiome(0.45, 0.70, 0.80)` → `Jungle`, `♣`, `#1a7a2e`
- [x] 5.5 Add `TestClassifyBiome_Savanna`: `classifyBiome(0.45, 0.45, 0.80)` → `Savanna`, `ˬ`, `#b5a04a`
- [x] 5.6 Add `TestClassifyBiome_AridSteppe`: `classifyBiome(0.45, 0.20, 0.80)` → `AridSteppe`, `·`, `#c9a97a`
- [x] 5.7 Add `TestClassifyBiome_Tundra`: `classifyBiome(0.45, 0.20, 0.20)` → `Tundra`, `∙`, `#8ab08a`
- [x] 5.8 Add `TestClassifyBiome_Taiga`: `classifyBiome(0.45, 0.70, 0.20)` → `Taiga`, `♠`, `#3a6b52`
- [x] 5.9 Update `TestTileAt_Smoke` to assert `tile.Temperature ∈ [0, 1]`
- [x] 5.10 Add `TestComputeTemperature_Equator`: Y=0 → temperature > 0.80
- [x] 5.11 Add `TestComputeTemperature_NorthPole`: Y=400 → temperature < 0.40

## 6. Verification

- [x] 6.1 Run `go build ./...` — no compile errors
- [x] 6.2 Run `go test ./internal/game/` — all tests pass
- [x] 6.3 Run the game and visually confirm jungle/savanna near Y=0 and tundra/taiga at higher Y values
