## MODIFIED Requirements

### Requirement: Items may register RoundHooks for combat side effects
The system SHALL support an optional `CombatHooks(self *Combatant, opponent *Combatant) []RoundHook` pattern where items can produce hooks at combat-start time. In this iteration no items implement this; the mechanism is documented so future item interactions can extend it without restructuring the combat loop.

When `buildPlayerCombatant` is called, it SHALL iterate over all equipped items and, for any item that has associated hooks (future), append those hooks to `CombatState.Hooks`. Currently all items return an empty hook slice.

`buildPlayerCombatant` SHALL also accumulate `ArmourBonus` and `DamageBonus` from all equipped items and apply them to the returned `Combatant`'s `Armour`, `MinDamage`, and `MaxDamage` fields.

#### Scenario: No items produce hooks in base implementation
- **WHEN** `buildPlayerCombatant` is called with only the default outfit equipped (Cloth Tunic, Cloth Pants, Leather Boots, Wooden Sword, Wooden Shield)
- **THEN** `len(combatState.Hooks) == 0`

#### Scenario: Equipped items contribute DamageBonus to player combatant
- **WHEN** `buildPlayerCombatant` is called with the Wooden Sword equipped (`DamageBonus: 1`) and base `MinDamage=1, MaxDamage=3`
- **THEN** the returned `Combatant` has `MinDamage == 2` and `MaxDamage == 4`

#### Scenario: Equipped items contribute ArmourBonus to player combatant
- **WHEN** `buildPlayerCombatant` is called with the Wooden Shield equipped (`ArmourBonus: 1`) and base `Armour=0`
- **THEN** the returned `Combatant` has `Armour == 1`

#### Scenario: Multiple bonus items stack additively
- **WHEN** two equipped items each have `DamageBonus: 1` and base `MaxDamage=3`
- **THEN** the returned `Combatant` has `MaxDamage == 5`

#### Scenario: Empty equipment slots contribute no bonus
- **WHEN** all equipment slots are empty
- **THEN** bonuses are 0 and base stats are unchanged
