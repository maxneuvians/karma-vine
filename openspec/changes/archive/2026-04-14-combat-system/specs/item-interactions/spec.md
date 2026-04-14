## ADDED Requirements

### Requirement: Items may register RoundHooks for combat side effects
The system SHALL support an optional `CombatHooks(self *Combatant, opponent *Combatant) []RoundHook` pattern where items can produce hooks at combat-start time. In this iteration no items implement this; the mechanism is documented so future item interactions can extend it without restructuring the combat loop.

When `buildPlayerCombatant` is called, it SHALL iterate over all equipped items and, for any item that has associated hooks (future), append those hooks to `CombatState.Hooks`. Currently all items return an empty hook slice.

#### Scenario: No items produce hooks in base implementation
- **WHEN** `buildPlayerCombatant` is called with the default outfit (Cloth Tunic, Cloth Pants, Leather Boots)
- **THEN** `len(combatState.Hooks) == 0`
