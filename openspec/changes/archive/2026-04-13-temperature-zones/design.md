## Context

The world generator currently classifies biomes from two noise fields: elevation (`e`) and moisture (`m`). The biome table is a simple nested switch with nine biomes. This produces plausible terrain but no sense of latitude — a tile at Y=0 and Y=10000 is equally likely to be jungle or tundra. The new temperature axis introduces a latitudinal gradient so that equatorial regions are hot and polar regions are cold, with all biome variety shifting accordingly.

## Goals / Non-Goals

**Goals:**
- Derive a `temperature float64 ∈ [0, 1]` for each tile from its Y coordinate and a small noise perturbation
- Add five new biomes: `Jungle`, `Savanna`, `AridSteppe`, `Tundra`, `Taiga`
- Replace the 2-axis biome table with a 3-axis (elevation, moisture, temperature) classification
- Keep existing ocean/water/beach/snow biomes mapped intuitively within the new scheme
- Add `Temperature float64` to the `Tile` struct

**Non-Goals:**
- Hemisphere mirroring (north/south symmetry) — the world is infinite and gradient is one-directional
- Per-biome local map content changes — that is a separate change
- Seasonal temperature variation — temperature is static per tile

## Decisions

**Temperature derivation from Y coordinate**

A cosine-based latitudinal band `lat = cos(π × |Y| / halfPeriod)` maps Y=0 to temperature 1.0 (equator) and ramps down as |Y| increases, cycling back periodically to simulate multiple climate belts across the infinite map. A half-period of 800 tiles gives ~1600-tile wide climate bands, which feels continental at the current noise scale.

A small noise perturbation (amplitude 0.12) is added to break perfect horizontal symmetry and produce natural-looking climate boundaries rather than straight lines.

Alternative considered: simple linear falloff per |Y|. Rejected because an infinite world with linear falloff is permanently frozen beyond a certain Y value, leaving large unusable polar wastelands.

**Three-axis biome table layout**

The classification splits first by temperature, then by elevation, then by moisture:

| Temperature | Low elevation | Mid elevation (moisture split) | High elevation |
|-------------|--------------|-------------------------------|----------------|
| Hot (≥0.65) | → ocean/beach | dry→Savanna/AridSteppe, wet→Jungle | → Mountains/VolcanicPeak (Snow reused) |
| Temperate (0.35–0.65) | → ocean/beach | dry→Plains/Desert, wet→Forest/DenseForest | → Mountains/Snow |
| Cold (<0.35) | → ocean/beach | dry→Tundra, wet→Taiga | → Mountains/Snow |

This gives a smooth gradient from equatorial jungle through temperate forest to arctic tundra as Y increases, matching real-world climate patterns.

Alternative considered: a lookup table of (temp_band, elev_band, moist_band) enum indices. Rejected in favour of readable nested conditions that are easier to tune visually.

**New biome glyphs and colours**

- `Jungle` — `♣` (reused from Forest), `#1a7a2e` (deep green)
- `Savanna` — `ˬ`, `#b5a04a` (yellow-brown grass)
- `AridSteppe` — `·`, `#c9a97a` (pale tan — distinct from Desert `~`)
- `Tundra` — `∙`, `#8ab08a` (muted grey-green)
- `Taiga` — `♠` (reused from DenseForest), `#3a6b52` (cold dark green)

**No change to local map generation in this change**

Local map content (objects, animals, fire thresholds) is keyed on `Biome`. New biomes will have empty/default content until a follow-up change adds their local content. This is intentional — world appearance ships first.

## Risks / Trade-offs

- **Existing biome constants are renumbered** — The `Biome` iota list gains 5 entries. Any persisted/serialised biome values would break, but this game has no save files, so there is no migration concern.
- **`classifyBiome` tests require full rewrite** — All `TestClassifyBiome_*` test inputs must supply a temperature argument; old test cases will fail to compile. Mitigation: rewrite tests as part of this change.
- **Visual regression on existing seeds** — Adding a temperature axis changes which biome is assigned to many existing tiles. The world will look different after this change. This is expected and desired.
- **Cosine climate bands repeat** — At multiples of 1600 tiles north/south, the climate cycles back to hot. This is intentional (multiple continent belts) but may feel repetitive at extreme Y values. Mitigation: acceptable for now; a future change could add hemispheric noise to break symmetry further.
