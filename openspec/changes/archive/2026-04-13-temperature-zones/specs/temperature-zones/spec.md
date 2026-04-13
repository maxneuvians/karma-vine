## ADDED Requirements

### Requirement: Each tile carries a temperature value derived from latitude
The system SHALL compute a `temperature float64 ∈ [0, 1]` for every tile, where values near 1.0 represent hot equatorial conditions and values near 0.0 represent cold polar conditions. Temperature SHALL be derived from a cosine latitudinal gradient centred on Y=0 with a half-period of 800 world tiles (so climate bands repeat every 1600 tiles), blended with a small noise perturbation of amplitude 0.12 to break horizontal symmetry. The `Tile` struct SHALL include a `Temperature float64` field populated by `generateChunk`.

#### Scenario: Tile at Y=0 has high temperature
- **WHEN** a tile is generated at world coordinate (0, 0)
- **THEN** `tile.Temperature` is greater than `0.80`

#### Scenario: Tile far north has low temperature
- **WHEN** a tile is generated at world coordinate (0, 400)
- **THEN** `tile.Temperature` is less than `0.40`

#### Scenario: Temperature field is populated for all tiles
- **WHEN** `TileAt` is called for any world coordinate
- **THEN** the returned `Tile.Temperature` is in the range `[0, 1]`

### Requirement: Biome classification uses temperature as a third axis
The `classifyBiome` function SHALL accept a third parameter `temperature float64` and use it alongside elevation and moisture to assign biomes. The classification SHALL produce distinct biome zones for hot (temperature ≥ 0.65), temperate (0.35–0.65), and cold (< 0.35) bands. The full biome table SHALL be:

**Water/coast (all temperature bands):**
1. `e < 0.28` → `DeepOcean`
2. `e < 0.36` → `ShallowWater`
3. `e < 0.40` → `Beach`

**Hot band (temperature ≥ 0.65), land only:**
4. `e < 0.50` and `m > 0.55` → `Jungle`
5. `e < 0.50` and `m > 0.30` → `Savanna`
6. `e < 0.50` → `AridSteppe`
7. `e < 0.62` and `m > 0.45` → `Jungle`
8. `e < 0.62` → `Desert`
9. `e < 0.78` → `Mountains`
10. else → `Snow`

**Temperate band (0.35 ≤ temperature < 0.65), land only:**
11. `e < 0.50` and `m > 0.55` → `Forest`
12. `e < 0.50` → `Plains`
13. `e < 0.62` and `m > 0.45` → `DenseForest`
14. `e < 0.62` and `m < 0.35` → `Desert`
15. `e < 0.78` → `Mountains`
16. else → `Snow`

**Cold band (temperature < 0.35), land only:**
17. `e < 0.50` and `m > 0.50` → `Taiga`
18. `e < 0.50` → `Tundra`
19. `e < 0.62` → `Taiga`
20. `e < 0.78` → `Mountains`
21. else → `Snow`

#### Scenario: Hot wet lowland is Jungle
- **WHEN** `classifyBiome(0.45, 0.70, 0.80)` is called
- **THEN** the returned biome is `Jungle`

#### Scenario: Hot dry lowland is AridSteppe
- **WHEN** `classifyBiome(0.45, 0.20, 0.80)` is called
- **THEN** the returned biome is `AridSteppe`

#### Scenario: Hot mid-moisture lowland is Savanna
- **WHEN** `classifyBiome(0.45, 0.45, 0.80)` is called
- **THEN** the returned biome is `Savanna`

#### Scenario: Cold wet lowland is Taiga
- **WHEN** `classifyBiome(0.45, 0.70, 0.20)` is called
- **THEN** the returned biome is `Taiga`

#### Scenario: Cold dry lowland is Tundra
- **WHEN** `classifyBiome(0.45, 0.20, 0.20)` is called
- **THEN** the returned biome is `Tundra`

#### Scenario: Temperate wet lowland is Forest
- **WHEN** `classifyBiome(0.45, 0.70, 0.50)` is called
- **THEN** the returned biome is `Forest`

#### Scenario: Deep ocean unaffected by temperature
- **WHEN** `classifyBiome(0.20, 0.50, 0.80)` is called
- **THEN** the returned biome is `DeepOcean`

### Requirement: New biomes have distinct glyphs and colours
The system SHALL define five new `Biome` constants and register their display properties in the rendering system:

| Biome | Char | Colour |
|-------|------|--------|
| `Jungle` | `♣` | `#1a7a2e` |
| `Savanna` | `ˬ` | `#b5a04a` |
| `AridSteppe` | `·` | `#c9a97a` |
| `Tundra` | `∙` | `#8ab08a` |
| `Taiga` | `♠` | `#3a6b52` |

#### Scenario: Jungle tile renders with correct glyph and colour
- **WHEN** `classifyBiome` returns `Jungle`
- **THEN** the returned `rune` is `♣` and the colour is `#1a7a2e`

#### Scenario: Tundra tile renders with correct glyph and colour
- **WHEN** `classifyBiome` returns `Tundra`
- **THEN** the returned `rune` is `∙` and the colour is `#8ab08a`
