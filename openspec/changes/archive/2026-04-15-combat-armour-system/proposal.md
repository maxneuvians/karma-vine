## Why

Combat currently uses armour as a static stat displayed in the UI but not modelled as an absorbing resource — all damage goes directly to HP. Making armour a depletable combat buffer that resets between encounters adds tactical depth and makes equipment choices meaningful, starting with two new starter items (Wooden Shield, Wooden Sword) that immediately demonstrate those choices.

## What Changes

- Armour is now a depletable combat resource: incoming damage is applied to the attacker's current armour first; only damage exceeding current armour carries over to HP
- When current armour reaches 0, all remaining damage from that hit and all subsequent hits applies directly to HP
- At the start of each new enemy encounter, the player's current armour resets to `MaxArmour` (the full armour value derived from stats + equipment); HP does not reset
- Enemies' armour also behaves as a depletable resource within a single combat (resets only at combat start, not between rounds)
- `Item` gains two new optional fields: `ArmourBonus int` and `DamageBonus int`; these are summed across all equipped items and applied when building the player `Combatant`
- `NewModel()` pre-equips the player with a **Wooden Sword** (`SlotRightHand`, `DamageBonus: +1`) and a **Wooden Shield** (`SlotLeftHand`, `ArmourBonus: +1`)
- The combat log and HP/armour bars in the rendering layer reflect the live armour value as it changes during playback

## Capabilities

### New Capabilities
- `combat-armour`: Armour-as-absorber mechanic in the combat engine — armour absorbs incoming damage before HP, depletes with each hit, resets to `MaxArmour` at the start of each encounter

### Modified Capabilities
- `equipment-system`: `NewModel()` pre-equips Wooden Sword and Wooden Shield; `Item` struct gains `ArmourBonus` and `DamageBonus` fields
- `item-interactions`: `buildPlayerCombatant` sums `ArmourBonus` and `DamageBonus` from all equipped items and applies them to the player `Combatant`'s `Armour` and damage range fields

## Impact

- `internal/game/types.go`: `Item` struct gains `ArmourBonus int` and `DamageBonus int`; `CombatState` or `Combatant` may need a `CurrentArmour int` field separate from the static `Armour` (max) field
- `internal/game/combat.go`: damage resolution updated to drain `CurrentArmour` before HP; armour reset injected at encounter start
- `internal/game/model.go`: `buildPlayerCombatant` (or equivalent) sums item bonuses
- `internal/game/model.go` / `NewModel()`: Wooden Sword and Wooden Shield added to `Equipped`
- `internal/game/render.go`: armour bar playback uses `CurrentArmour` values scraped from log (analogous to existing HP playback); armour values in log lines updated
- No new external dependencies
