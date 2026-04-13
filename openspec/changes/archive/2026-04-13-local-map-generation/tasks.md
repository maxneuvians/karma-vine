## 1. Type Definitions

- [x] 1.1 Define `Ground`, `Object`, and `Animal` structs in `types.go`
- [x] 1.2 Replace the placeholder `LocalMap` struct (from `project-scaffold`) with the full struct: `Ground [42][18]Ground`, `Objects [42][18]*Object`, `Animals []*Animal`

## 2. Hash Function

- [x] 2.1 Implement `hash(x, y, seed int) float64` in `local.go` using the exact formula from the brief
- [x] 2.2 Write a unit test in `local_test.go` verifying `hash(5, 10, 42)` is deterministic across two calls
- [x] 2.3 Write a unit test verifying the output is in the range `[0, 1)`

## 3. Biome Content Tables

- [x] 3.1 Define per-biome ground variant tables (character + colour slices) for all six content biomes: Forest, Desert, Plains, Mountains, Snow, Beach
- [x] 3.2 Define per-biome object tables (character + colour + blocking) for each biome
- [x] 3.3 Define per-biome animal tables (character + colour + flee) for each biome
- [x] 3.4 Define per-biome `objectThreshold float64` and `maxAnimals int` constants

## 4. GenerateLocalMap

- [x] 4.1 Implement `GenerateLocalMap(worldX, worldY, globalSeed int, biome Biome) *LocalMap` in `local.go`
- [x] 4.2 Derive `localSeed` using `hash(worldX, worldY, globalSeed)` scaled to int64
- [x] 4.3 Instantiate two `opensimplex.New` objects (terrain + object noise) from derived seeds
- [x] 4.4 Iterate all 42×18 cells; sample terrain noise at scale `0.12`; pick ground variant from biome table using noise value quantised to table length
- [x] 4.5 For each cell where object-noise > `objectThreshold`, pick a random object from the biome table and set `Objects[x][y]`
- [x] 4.6 Place animals at up to `maxAnimals` random passable positions using biome animal table
- [x] 4.7 Write unit tests verifying Forest map contains a tree, Desert map contains a cactus, and same inputs produce identical output

## 5. LocalMapFor Accessor

- [x] 5.1 Implement `LocalMapFor(worldX, worldY int, m *Model) *LocalMap` in `local.go`
- [x] 5.2 On cache miss, call `TileAt` to get biome, call `GenerateLocalMap`, store in `m.localCache`
- [x] 5.3 On cache hit, return the existing pointer
- [x] 5.4 Write a unit test that calls `LocalMapFor` twice and asserts pointer equality

## 6. Verification

- [x] 6.1 Run `go test ./...` and confirm all tests pass with no data-race warnings (`go test -race ./...`)
