## ADDED Requirements

### Requirement: Axe object is injected into forested local maps
When generating a local map for a cell whose biome is `BiomeForest`, `BiomeDenseForest`, `BiomeJungle`, or `BiomeTaiga`, the system SHALL have a chance to place an `Axe` object at a random floor cell (not occupied by another object or blocking terrain). The axe SHALL be defined as:

```go
&Object{Char: '⚒', Color: "#a0a0a0", Name: "Axe", Blocking: false, Pickupable: true}
```

The spawn probability SHALL be configurable as a constant (e.g., `axeSpawnChance = 0.3`); exactly one axe per map at most.

#### Scenario: Forest map may contain an Axe
- **WHEN** a local map is generated for a Forest biome cell
- **THEN** at most one `Object` with `Name == "Axe"` and `Pickupable == true` exists on the map

#### Scenario: Non-forest maps do not contain an Axe
- **WHEN** a local map is generated for a `BiomeDesert` or `BiomeTundra` biome cell
- **THEN** no `Object` with `Name == "Axe"` exists on the map

#### Scenario: Axe is placed on a passable cell
- **WHEN** an axe is spawned
- **THEN** its cell is not occupied by any other object and is not a blocking terrain tile
