## ADDED Requirements

### Requirement: Animals have combat stats defined in a stat table
The system SHALL define a `AnimalCombatStats` struct with fields `HP int`, `Armour int`, `MinDamage int`, `MaxDamage int`, `Initiative int`. A package-level map `animalStatTable map[string]AnimalCombatStats` SHALL provide default stats keyed on `Animal.Name`. Any animal name absent from the table SHALL fall back to a default entry (`HP=5`, `Armour=0`, `MinDamage=1`, `MaxDamage=2`, `Initiative=3`).

#### Scenario: Known animal name returns its stats
- **WHEN** `animalStatTable["Wolf"]` is looked up
- **THEN** the returned struct has non-zero `HP` and `Initiative`

#### Scenario: Unknown animal name returns fallback stats
- **WHEN** `animalStatTable["UnknownBeast"]` is looked up (key absent)
- **THEN** the returned struct has `HP=5`, `Armour=0`, `MinDamage=1`, `MaxDamage=2`, `Initiative=3`

### Requirement: buildEnemyCombatant constructs a Combatant from an Animal
The system SHALL provide `buildEnemyCombatant(a Animal) Combatant` that looks up `a.Name` in `animalStatTable` (with fallback) and returns a fully populated `Combatant` with `Name = a.Name` and stats from the table.

#### Scenario: Combatant built from animal has correct name
- **WHEN** `buildEnemyCombatant(Animal{Name: "Deer"})` is called
- **THEN** the returned `Combatant.Name == "Deer"`
