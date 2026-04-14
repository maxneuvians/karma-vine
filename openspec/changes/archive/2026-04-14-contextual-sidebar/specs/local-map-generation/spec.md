## MODIFIED Requirements

### Requirement: Local map cell types are defined
The system SHALL define: `Ground{Char rune, Color string, Passable bool}`, `Object{Char rune, Color string, Blocking bool, Name string}`, and `Animal{X, Y int, Char rune, Color string, Flee bool, Name string}`. `LocalMap` SHALL have fields `Ground [LocalMapW][LocalMapH]Ground`, `Objects [LocalMapW][LocalMapH]*Object`, `Animals []*Animal`.

#### Scenario: LocalMap zero-value is safe to use
- **WHEN** a `LocalMap` is allocated with `&LocalMap{}`
- **THEN** all `Ground` cells default to zero-value without requiring explicit initialisation, and `Objects` cells default to `nil`

#### Scenario: Object Name field is accessible
- **WHEN** an `Object` is constructed with `Name: "Tree"`
- **THEN** `obj.Name == "Tree"`

### Requirement: GenerateLocalMap populates Name on all objects and animals
The system SHALL set the `Name` field on every `Object` and `Animal` created during `GenerateLocalMap`. No object or animal SHALL be created with an empty `Name`. The dungeon entrance object SHALL have `Name: "Dungeon Entrance"`.

#### Scenario: Dungeon entrance object has Name "Dungeon Entrance"
- **WHEN** `GenerateLocalMap` is called
- **THEN** the object at the dungeon entrance cell has `Name == "Dungeon Entrance"`

#### Scenario: No object has an empty Name
- **WHEN** `GenerateLocalMap` is called for any biome
- **THEN** every non-nil `Objects[x][y]` has a non-empty `Name`

#### Scenario: No animal has an empty Name
- **WHEN** `GenerateLocalMap` is called for any biome
- **THEN** every `Animal` in `Animals` has a non-empty `Name`
