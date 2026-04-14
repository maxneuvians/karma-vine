## Requirements

### Requirement: Sidebar content is determined by the active mode
The system SHALL render the sidebar with content specific to `m.mode`. When `m.mode == ModeWorld`, it SHALL show the biome legend or map-mode overlay hint. When `m.mode == ModeLocal`, it SHALL show the local legend with named objects and wildlife. When `m.mode == ModeDungeon`, it SHALL show dungeon-specific content. Content from a different mode SHALL NOT appear in the sidebar.

#### Scenario: Local sidebar does not show dungeon depth
- **WHEN** `m.mode == ModeLocal` and the sidebar is rendered
- **THEN** the sidebar does not contain text `Depth:`

#### Scenario: Dungeon sidebar does not show biome legend
- **WHEN** `m.mode == ModeDungeon` and the sidebar is rendered
- **THEN** the sidebar does not contain biome names such as `Deep Ocean` or `Forest`

#### Scenario: World sidebar does not show local legend
- **WHEN** `m.mode == ModeWorld` and the sidebar is rendered
- **THEN** the sidebar does not contain the `Legend` header

### Requirement: World sidebar shows map-mode overlay hint for non-default modes
When `m.mapMode != MapModeDefault`, the sidebar SHALL replace the biome colour legend with a header naming the active overlay and a one-line description of how to interpret the colours. When `m.mapMode == MapModeDefault`, the sidebar SHALL show the biome colour legend as before.

#### Scenario: Temperature mode shows overlay hint
- **WHEN** `m.mode == ModeWorld`, `m.mapMode == MapModeTemperature`, and the sidebar is rendered
- **THEN** the sidebar contains `Temperature` and does not contain `Deep Ocean`

#### Scenario: Default mode shows biome legend
- **WHEN** `m.mode == ModeWorld`, `m.mapMode == MapModeDefault`, and the sidebar is rendered
- **THEN** the sidebar contains `Deep Ocean`

#### Scenario: Elevation mode shows overlay hint
- **WHEN** `m.mode == ModeWorld`, `m.mapMode == MapModeElevation`, and the sidebar is rendered
- **THEN** the sidebar contains `Elevation` and does not contain `Beach`

### Requirement: Local sidebar shows fully named objects and wildlife
In `ModeLocal`, the sidebar SHALL list each unique object and animal visible on the current local map by its `Name` field. No object or animal SHALL be listed with an empty name or a raw glyph string as its name.

#### Scenario: Dungeon entrance shows correct label
- **WHEN** `m.mode == ModeLocal` and the local map contains a dungeon entrance object
- **THEN** the sidebar contains `Dungeon Entrance` (not `>`)

#### Scenario: All animals are named
- **WHEN** `m.mode == ModeLocal` and the local map contains animals
- **THEN** every animal listed in the sidebar has a non-empty name that is not a single character glyph

#### Scenario: All objects are named
- **WHEN** `m.mode == ModeLocal` and the local map contains objects
- **THEN** every object listed in the sidebar has a non-empty name

### Requirement: Dungeon sidebar shows depth and level contents
In `ModeDungeon`, the sidebar SHALL display:
- A `Dungeon` header
- Current depth as `Depth: N`
- A `Contents` sub-section listing each unique named object type present on the current level (torches, braziers, up/down staircases)

#### Scenario: Dungeon sidebar shows depth
- **WHEN** `m.mode == ModeDungeon` and `m.dungeonDepth == 2`
- **THEN** the sidebar contains `Depth: 2`

#### Scenario: Dungeon sidebar lists torch when present
- **WHEN** `m.mode == ModeDungeon` and the current level contains at least one torch
- **THEN** the sidebar contains `Torch`

#### Scenario: Dungeon sidebar lists down staircase only when present
- **WHEN** `m.mode == ModeDungeon` and `m.currentDungeon.HasDownStair == false`
- **THEN** the sidebar does not contain `Staircase Down`

#### Scenario: Dungeon sidebar lists up staircase
- **WHEN** `m.mode == ModeDungeon`
- **THEN** the sidebar contains `Staircase Up`
