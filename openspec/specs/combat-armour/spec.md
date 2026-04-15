## ADDED Requirements

### Requirement: Combatant has a CurrentArmour field
The system SHALL add `CurrentArmour int` to the `Combatant` struct. `resolveCombat` SHALL initialise `CurrentArmour = Armour` on both combatants at the start of the function before the round loop begins.

#### Scenario: CurrentArmour initialised to Armour at combat start
- **WHEN** `resolveCombat` is called with a combatant whose `Armour == 3`
- **THEN** that combatant's `CurrentArmour == 3` before the first round

#### Scenario: Combatant with zero Armour has zero CurrentArmour
- **WHEN** `resolveCombat` is called with a combatant whose `Armour == 0`
- **THEN** `CurrentArmour == 0` before the first round

### Requirement: applyDamage drains CurrentArmour before HP
The system SHALL replace `rollDamage` with `applyDamage(raw int, defender *Combatant) (armourDrain, hpDrain int)` that applies incoming `raw` damage as follows:
1. If `defender.CurrentArmour > 0`: subtract `raw` from `CurrentArmour`. If `CurrentArmour` goes negative, the absolute value is `hpDrain`; clamp `CurrentArmour` to 0 and subtract `hpDrain` from `defender.HP`. Return `(original CurrentArmour value drained, hpDrain)`.
2. If `defender.CurrentArmour == 0`: subtract `raw` from `defender.HP`. Return `(0, raw)`.
3. `raw == 0` is a no-op; return `(0, 0)`.

#### Scenario: Hit fully absorbed by armour
- **WHEN** `defender.CurrentArmour == 3` and `raw == 2`
- **THEN** `defender.CurrentArmour == 1`, `defender.HP` is unchanged, `armourDrain == 2`, `hpDrain == 0`

#### Scenario: Hit depletes armour and overflows to HP
- **WHEN** `defender.CurrentArmour == 2` and `raw == 5`
- **THEN** `defender.CurrentArmour == 0`, `defender.HP` decreases by 3, `armourDrain == 2`, `hpDrain == 3`

#### Scenario: Hit applied directly to HP when armour is zero
- **WHEN** `defender.CurrentArmour == 0` and `raw == 4`
- **THEN** `defender.HP` decreases by 4, `armourDrain == 0`, `hpDrain == 4`

#### Scenario: Zero raw damage is a no-op
- **WHEN** `raw == 0`
- **THEN** `defender.CurrentArmour` and `defender.HP` are unchanged

### Requirement: resolveCombat uses applyDamage and emits structured log lines
`resolveCombat` SHALL call `applyDamage` in place of `rollDamage` and emit log lines using the following format rules based on the returned `(armourDrain, hpDrain)`:
- Armour-only hit (`hpDrain == 0`, `armourDrain > 0`): `"Round %d: %s attacks %s — absorbs %d (Armour: %d)"`
- Armour-broken hit (`armourDrain > 0`, `hpDrain > 0`): `"Round %d: %s attacks %s — armour broken, %d HP damage (%d HP, 0 Armour)"`
- Pure HP hit (`armourDrain == 0`, `hpDrain > 0`): `"Round %d: %s attacks %s for %d damage (%d HP left)"`
- No damage (`armourDrain == 0`, `hpDrain == 0`): `"Round %d: %s attacks %s but deals no damage"`

#### Scenario: Armour-only hit produces correct log line
- **WHEN** `applyDamage` returns `(armourDrain=2, hpDrain=0)` for round 1, attacker "Player", defender "Wolf"
- **THEN** log contains `"Round 1: Player attacks Wolf — absorbs 2 (Armour: %d)"` where `%d` is remaining armour

#### Scenario: Armour-broken hit produces correct log line
- **WHEN** `applyDamage` returns `(armourDrain=2, hpDrain=3)` for round 1
- **THEN** log contains `"Round 1: %s attacks %s — armour broken, 3 HP damage (%d HP, 0 Armour)"`

#### Scenario: Pure HP hit produces correct log line
- **WHEN** `defender.CurrentArmour == 0` and `applyDamage` returns `(0, 4)` for round 2
- **THEN** log contains `"Round 2: %s attacks %s for 4 damage (%d HP left)"`

#### Scenario: No-damage hit produces correct log line
- **WHEN** `raw == 0`
- **THEN** log contains `"Round %d: %s attacks %s but deals no damage"`

### Requirement: armourAtRound computes current armour from log
The system SHALL implement `armourAtRound(log []string, roundIndex int, combatantName string, startArmour int) int` that scans log lines up to `roundIndex` (using the same boundary logic as `combatLogLinesUpTo`) and reconstructs the combatant's current armour value by parsing armour drain events. Returns `startArmour` when `roundIndex == 0`.

#### Scenario: armourAtRound returns startArmour before any rounds
- **WHEN** `armourAtRound` is called with `roundIndex == 0`
- **THEN** returns `startArmour`

#### Scenario: armourAtRound decrements on armour-only hits
- **WHEN** log contains `"Round 1: Enemy attacks Player — absorbs 1 (Armour: 2)"` and `roundIndex == 1`
- **THEN** `armourAtRound(..., "Player", 3)` returns `2`

#### Scenario: armourAtRound returns 0 after armour-broken hit
- **WHEN** log contains `"Round 1: Enemy attacks Player — armour broken, 3 HP damage (7 HP, 0 Armour)"` and `roundIndex == 1`
- **THEN** `armourAtRound(..., "Player", 2)` returns `0`
