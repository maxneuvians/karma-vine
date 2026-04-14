## ADDED Requirements

### Requirement: buildView dispatches to dungeon render path
The system SHALL extend `buildView` to dispatch to `renderDungeonMap` when `m.mode == ModeDungeon`. The dungeon render path SHALL compose the dungeon map with the HUD and optional key bar, matching the structure of the local map render path.

#### Scenario: Dungeon map renders when mode is ModeDungeon
- **WHEN** `buildView` is called with `m.mode == ModeDungeon`
- **THEN** the returned string contains dungeon cell characters (`#`, `.`, `@`) and does not contain world-map tile characters

#### Scenario: HUD is present in dungeon view
- **WHEN** `buildView` is called with `m.mode == ModeDungeon`
- **THEN** the returned string contains the HUD status bar with dungeon depth information
