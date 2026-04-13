## ADDED Requirements

### Requirement: Local map cell types are defined
The system SHALL define: `Ground{Char rune, Color string, Passable bool}`, `Object{Char rune, Color string, Blocking bool}`, and `Animal{X, Y int, Char rune, Color string, Flee bool}`. `LocalMap` SHALL have fields `Ground [42][18]Ground`, `Objects [42][18]*Object`, `Animals []*Animal`.

#### Scenario: LocalMap zero-value is safe to use
- **WHEN** a `LocalMap` is allocated with `&LocalMap{}`
- **THEN** all `Ground` cells default to zero-value without requiring explicit initialisation, and `Objects` cells default to `nil`

### Requirement: Hash function matches the brief specification
The system SHALL implement `hash(x, y, seed int) float64` using this exact formula:
```
h := x*1619 + y*31337 + seed*6971
h = (h ^ (h >> 16)) * 0x45d9f3b
h = (h ^ (h >> 16)) * 0x45d9f3b
h = h ^ (h >> 16)
return float64(uint32(h)) / float64(0xffffffff)
```
The return value SHALL be in `[0, 1)`.

#### Scenario: Hash is deterministic
- **WHEN** `hash(5, 10, 42)` is called twice
- **THEN** both calls return the identical `float64` value

#### Scenario: Hash distributes across world coordinates
- **WHEN** `hash` is called for 100 distinct `(x, y)` pairs with the same seed
- **THEN** at least 90 distinct values are returned (no degenerate clustering)

### Requirement: GenerateLocalMap produces a biome-appropriate 42×18 map
The system SHALL implement `GenerateLocalMap(worldX, worldY, globalSeed int, biome Biome) *LocalMap` that:
1. Derives `localSeed` via `hash(worldX, worldY, globalSeed)`
2. Creates two `opensimplex.New` instances seeded from `localSeed`
3. Iterates all 42×18 cells, sets ground character and colour from biome ground table at noise scale `0.12`
4. Places objects where the object-noise value exceeds the biome's `objectThreshold`
5. Places up to a biome-defined maximum number of animals at random passable positions

#### Scenario: Same world coordinate always produces the same map
- **WHEN** `GenerateLocalMap(3, 7, 12345, Forest)` is called twice
- **THEN** both returned `LocalMap` values have identical `Ground` and `Objects` arrays

#### Scenario: Forest biome map contains at least one tree object
- **WHEN** `GenerateLocalMap` is called for the `Forest` biome
- **THEN** at least one cell in `Objects` is non-nil with `Char` equal to `♣` or `♠`

#### Scenario: Desert biome map contains at least one cactus object
- **WHEN** `GenerateLocalMap` is called for the `Desert` biome
- **THEN** at least one cell in `Objects` is non-nil with `Char` equal to `ψ`

### Requirement: LocalMapFor accessor caches local maps
The system SHALL provide `LocalMapFor(worldX, worldY int, m *Model) *LocalMap` that looks up `m.localCache`, calls `GenerateLocalMap` only on a cache miss (using `TileAt` to obtain the biome), stores the result, and returns it.

#### Scenario: Cache miss triggers generation and storage
- **WHEN** `LocalMapFor` is called for a coordinate not in `m.localCache`
- **THEN** the result is added to `m.localCache[WorldCoord{worldX, worldY}]`

#### Scenario: Cache hit skips generation
- **WHEN** `LocalMapFor` is called twice for the same coordinate
- **THEN** both calls return a pointer to the same `LocalMap` instance
