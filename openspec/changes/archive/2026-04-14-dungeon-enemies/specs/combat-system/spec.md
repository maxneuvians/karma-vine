## ADDED Requirements

### Requirement: buildDungeonEnemyCombatant constructs a Combatant from a DungeonEnemy
The system SHALL implement `buildDungeonEnemyCombatant(e *DungeonEnemy) Combatant` that returns a `Combatant` with `Name = e.Template.Name`, `HP = e.HP`, `MaxHP = e.MaxHP`, `Armour = e.Armour`, `MinDamage = e.MinDamage`, `MaxDamage = e.MaxDamage`, `Initiative = e.Initiative`.

#### Scenario: Combatant built from DungeonEnemy reflects live stats
- **WHEN** `buildDungeonEnemyCombatant` is called with an enemy that has `HP=15`, `MaxHP=20`
- **THEN** the returned `Combatant` has `HP=15` and `MaxHP=20`

#### Scenario: Combatant name matches template name
- **WHEN** `buildDungeonEnemyCombatant` is called with an enemy whose `Template.Name == "Frost Giant"`
- **THEN** the returned `Combatant.Name == "Frost Giant"`

### Requirement: combatEnemy field supports both Animal and DungeonEnemy references
The model SHALL store the combat target as an interface or via two separate nullable fields so the victory dismiss handler can distinguish between an `*Animal` (surface) and a `*DungeonEnemy` (dungeon) and apply the correct post-combat cleanup. The existing `combatEnemy *Animal` field SHALL be supplemented by `combatDungeonEnemy *DungeonEnemy`. Exactly one of these SHALL be non-nil during active combat.

#### Scenario: Surface combat sets combatEnemy non-nil and combatDungeonEnemy nil
- **WHEN** the player initiates combat with a surface animal
- **THEN** `m.combatEnemy != nil` and `m.combatDungeonEnemy == nil`

#### Scenario: Dungeon combat sets combatDungeonEnemy non-nil and combatEnemy nil
- **WHEN** the player initiates combat with a DungeonEnemy
- **THEN** `m.combatDungeonEnemy != nil` and `m.combatEnemy == nil`
