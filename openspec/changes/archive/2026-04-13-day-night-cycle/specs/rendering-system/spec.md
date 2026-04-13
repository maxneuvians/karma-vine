## MODIFIED Requirements

### Requirement: Night mode dims all tile colours
When `Model.nightMode` is `true`, the system SHALL multiply the R, G, B components of every tile colour by `0.35` (rounded to nearest integer) before applying the Lipgloss style. The player glyph colour `#f0f6fc` SHALL NOT be dimmed.
**Reason**: Replaced by time-of-day automatic dimming.
**Migration**: Remove `nightMode bool` from `Model`. Replace all callers of `applyColor(hex, nightMode)` with `applyColor(hex, dimFactor)`, where `dimFactor` is replaced by new `dimFactor float64` derived from `timeOfDay`.

## REMOVED Requirements

### Requirement: Night mode dims all tile colours
**Reason**: The boolean flag `nightMode` is superseded by the continuous `dimFactor float64` derived from `timeOfDay`. Manual toggle is removed in favour of automatic simulation.
**Migration**: Remove `Model.nightMode bool`. Replace `applyColor(hex, bool)` with `applyColor(hex, dimFactor float64)`. The `dimFactor` is computed via a cosine curve: `clamp(0.5*(1+cos(2π*timeOfDay))*0.85 + 0.15, 0.15, 1.0)`. At noon (`timeOfDay=0.5`) `dimFactor≈1.0`; at midnight (`timeOfDay=0.0`) `dimFactor≈0.15`.

## ADDED Requirements

### Requirement: Tile colors are dimmed by a continuous time-of-day factor
The system SHALL compute `dimFactor float64` from `Model.timeOfDay` using the formula `dimFactor = clamp(0.5*(1+cos(2π*timeOfDay))*0.85 + 0.15, 0.15, 1.0)`. All tile colours on both world and local maps SHALL be multiplied channel-wise by `dimFactor` before applying Lipgloss styles. The player glyph colour SHALL NOT be dimmed.

#### Scenario: Full brightness at noon
- **WHEN** `timeOfDay == 0.5`
- **THEN** `dimFactor` is approximately `1.0` and tile colours are unmodified

#### Scenario: Dim at midnight
- **WHEN** `timeOfDay == 0.0`
- **THEN** `dimFactor` is approximately `0.15` and tile colours are multiplied by `~0.15`

#### Scenario: Player glyph unaffected by dim
- **WHEN** `timeOfDay == 0.0`
- **THEN** the player `@` is rendered in full `#f0f6fc`

### Requirement: Local map fire cells override dim within illumination radius
On the local map, for each cell within Manhattan distance 4 of any `HasFire == true` ground cell, the effective dim factor SHALL be `1.0` (fully lit), overriding the global `dimFactor`.

#### Scenario: Fire radius overrides midnight dim
- **WHEN** `timeOfDay == 0.0` and a fire cell is at `{20, 20}`
- **THEN** the cell at `{20, 23}` (distance 3) is rendered at full brightness
