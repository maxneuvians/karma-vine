## MODIFIED Requirements

### Requirement: Tile type carries biome metadata
The system SHALL define a `Tile` struct with fields `Biome Biome`, `Char rune`, `Color string`, `Elevation float64`, `Moisture float64`, `Temperature float64`. `Biome` SHALL be an `int` type with named constants for all fourteen biomes: `DeepOcean`, `ShallowWater`, `Beach`, `Forest`, `Plains`, `DenseForest`, `Desert`, `Mountains`, `Snow`, `Jungle`, `Savanna`, `AridSteppe`, `Tundra`, `Taiga`.

#### Scenario: Tile fields are populated after generation
- **WHEN** `TileAt` is called for any world coordinate
- **THEN** the returned `Tile` has a non-zero `Char`, a non-empty `Color`, an `Elevation` in the range `[0, 1]`, and a `Temperature` in the range `[0, 1]`

### Requirement: Biome thresholds use elevation, moisture, and temperature
The system SHALL assign biomes via `classifyBiome(e, m, temperature float64) (Biome, rune, string)` using a three-axis classification as specified in the temperature-zones capability spec. The nine original biomes (`DeepOcean`, `ShallowWater`, `Beach`, `Forest`, `Plains`, `DenseForest`, `Desert`, `Mountains`, `Snow`) are retained and the five new biomes (`Jungle`, `Savanna`, `AridSteppe`, `Tundra`, `Taiga`) are added.

#### Scenario: Tile with elevation 0.20 is classified as DeepOcean regardless of temperature
- **WHEN** `classifyBiome(0.20, 0.50, 0.80)` is called
- **THEN** the biome is `DeepOcean`, `Char` is `≋`, and `Color` is `#1a6fa8`

#### Scenario: Temperate tile with elevation 0.48 and moisture 0.60 is classified as Forest
- **WHEN** `classifyBiome(0.48, 0.60, 0.50)` is called
- **THEN** the biome is `Forest`, `Char` is `♣`, and `Color` is `#2d7a1f`
