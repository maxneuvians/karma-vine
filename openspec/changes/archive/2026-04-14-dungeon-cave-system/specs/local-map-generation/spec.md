## ADDED Requirements

### Requirement: Local maps contain a dungeon entrance object
The system SHALL place exactly one dungeon entrance object on each generated local map. The entrance object SHALL have `Char: '>'`, `Color: "#e8c96a"`, `Blocking: false`. The entrance SHALL be placed on a passable ground cell that is not already occupied by another object, chosen deterministically from the local seed.

#### Scenario: Generated local map contains exactly one dungeon entrance
- **WHEN** `GenerateLocalMap` is called for any biome
- **THEN** exactly one cell in `Objects` has `Char == '>'`

#### Scenario: Dungeon entrance is not blocking
- **WHEN** a dungeon entrance is placed on the local map
- **THEN** its `Blocking` field is `false` (player can stand on it)

#### Scenario: Dungeon entrance does not overlap another object
- **WHEN** the dungeon entrance cell is determined
- **THEN** no other object occupies the same cell
