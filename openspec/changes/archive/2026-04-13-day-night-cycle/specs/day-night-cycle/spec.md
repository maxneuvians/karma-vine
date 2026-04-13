## ADDED Requirements

### Requirement: Time-of-day advances continuously each tick
The system SHALL maintain a `timeOfDay float64` field in `Model` in the range `[0, 1)`, where `0` is midnight, `0.25` is 6 AM, `0.5` is noon, and `0.75` is 6 PM. On each `TickMsg`, `timeOfDay` SHALL advance by `timeScale / (ticksPerSecond * secondsPerDay)`. One full cycle SHALL take 30 real seconds at `timeScale == 1`. `timeOfDay` SHALL wrap around at 1.0.

#### Scenario: Time advances on each tick at 1× speed
- **WHEN** sixty `TickMsg` events are dispatched at `timeScale == 1`
- **THEN** `timeOfDay` has advanced by approximately `1.0` (one full day)

#### Scenario: Time wraps at midnight
- **WHEN** `timeOfDay` is `0.99` and a tick advances it past `1.0`
- **THEN** `timeOfDay` wraps to a value less than `1.0` (not ≥ 1.0)

### Requirement: Time speed is adjustable with [ and ] keys
The system SHALL support a `timeScale int` field with discrete values `1, 2, 5, 10`. Pressing `]` SHALL advance to the next higher scale (clamped at 10). Pressing `[` SHALL retreat to the next lower scale (clamped at 1).

#### Scenario: ] key increases time scale
- **WHEN** `timeScale == 1` and the player presses `]`
- **THEN** `timeScale == 2`

#### Scenario: [ key decreases time scale
- **WHEN** `timeScale == 5` and the player presses `[`
- **THEN** `timeScale == 2`

#### Scenario: Time scale is clamped at maximum
- **WHEN** `timeScale == 10` and the player presses `]`
- **THEN** `timeScale` remains `10`

#### Scenario: Time scale is clamped at minimum
- **WHEN** `timeScale == 1` and the player presses `[`
- **THEN** `timeScale` remains `1`

### Requirement: Fire cells are generated on local maps and illuminate a radius
`Ground.HasFire bool` SHALL indicate that a cell contains a fire. `GenerateLocalMap` SHALL place fire cells deterministically based on biome content tables. At render time, each fire cell SHALL illuminate all cells within Manhattan distance 4 of it, applying `dimFactor = 1.0` to those cells regardless of the global time-of-day dim.

#### Scenario: Fire cell is rendered with a fire glyph
- **WHEN** a local map cell has `HasFire == true` and is rendered
- **THEN** that cell displays a fire glyph (e.g., `♨` or `*`) regardless of objects or animals

#### Scenario: Cells within radius 4 of a fire are fully lit at night
- **WHEN** `timeOfDay == 0.0` (midnight) and a fire cell is at `{10, 10}`
- **THEN** the cell at `{13, 10}` (Manhattan distance 3) receives `effectiveDim == 1.0`

#### Scenario: Cells beyond radius 4 of a fire use global dim at night
- **WHEN** `timeOfDay == 0.0` (midnight) and a fire cell is at `{10, 10}`
- **THEN** the cell at `{15, 10}` (Manhattan distance 5) receives `effectiveDim < 0.5`

### Requirement: HUD displays the current in-game time and time scale
The HUD SHALL include a 24-hour formatted clock derived from `timeOfDay` (e.g., `06:30`) and the current `timeScale` as a multiplier suffix (e.g., `2×`).

#### Scenario: Clock shows noon at timeOfDay 0.5
- **WHEN** `timeOfDay == 0.5`
- **THEN** the HUD contains the text `12:00`

#### Scenario: Clock shows midnight at timeOfDay 0.0
- **WHEN** `timeOfDay == 0.0`
- **THEN** the HUD contains the text `00:00`

#### Scenario: HUD shows time scale
- **WHEN** `timeScale == 5`
- **THEN** the HUD contains the text `5×`
