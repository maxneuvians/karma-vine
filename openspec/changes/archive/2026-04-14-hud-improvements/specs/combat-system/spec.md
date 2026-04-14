## MODIFIED Requirements

### Requirement: Player combat stats are derived from base values and equipment
The player `Combatant` SHALL be constructed with: `HP = m.playerHP`, `MaxHP = m.playerMaxHP`, `Armour = sum of equipped item ArmourBonus`, `MinDamage = 1 + sum of equipped item MinDamageBonus`, `MaxDamage = 3 + sum of equipped item MaxDamageBonus`, `Initiative = 5 + sum of equipped item InitiativeBonus`. For this iteration all item bonus fields are 0 (no items grant bonuses yet); the structure exists for future use. After a victorious combat, `m.playerHP` SHALL be updated to the surviving player combatant's HP value so carry-over HP is reflected in the HUD.

#### Scenario: Player combatant HP seeds from model field
- **WHEN** `m.playerHP == 15` and `m.playerMaxHP == 20`
- **THEN** the constructed player `Combatant` has `HP=15` and `MaxHP=20`

#### Scenario: Player with no bonuses retains correct damage and initiative
- **WHEN** the player has no equipped items with bonuses and `m.playerHP == 20`
- **THEN** player combatant has `Armour=0`, `MinDamage=1`, `MaxDamage=3`, `Initiative=5`

#### Scenario: playerHP updated after victory
- **WHEN** combat resolves with `PlayerWon == true` and the player's surviving HP is `12`
- **THEN** `m.playerHP == 12` after dismissing the combat screen
