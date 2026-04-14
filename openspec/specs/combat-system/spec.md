## ADDED Requirements

### Requirement: Combatant struct holds all combat-relevant stats
The system SHALL define a `Combatant` struct in `types.go` with the fields: `Name string`, `HP int`, `MaxHP int`, `Armour int`, `MinDamage int`, `MaxDamage int`, `Initiative int`. It SHALL also define `type RoundHook func(self, opponent *Combatant)` and `CombatState` struct containing `Player Combatant`, `Enemy Combatant`, `Hooks []RoundHook`, `Log []string`, and `Round int`.

#### Scenario: Combatant fields are accessible
- **WHEN** a `Combatant` is constructed with all fields set
- **THEN** each field can be read and written directly

### Requirement: Combat is initiated when the player presses g on an animal cell
In `ModeLocal`, when the player presses `g` and the player's current cell contains an `Animal`, the system SHALL construct two `Combatant` values (player from base stats + equipment, enemy from the animal's stat table), call `resolveCombat`, store the resulting `CombatState` on the model, set `m.screenMode = ScreenCombat`, and set `m.paused = true`. If the cell contains no animal, the existing item-pickup behaviour applies unchanged.

#### Scenario: g on animal cell initiates combat
- **WHEN** the player is in `ModeLocal`, presses `g`, and their current cell has an `Animal`
- **THEN** `m.screenMode == ScreenCombat` and `m.combatState` is non-nil

#### Scenario: g on empty cell still picks up items
- **WHEN** the player is in `ModeLocal`, presses `g`, and their cell has no `Animal`
- **THEN** `m.screenMode` remains `ScreenNormal` and item pickup logic runs as before

### Requirement: resolveCombat resolves all rounds and returns a complete CombatState
`resolveCombat(player, enemy Combatant, hooks []RoundHook, rng *rand.Rand) CombatState` SHALL iterate rounds until one combatant's HP reaches ≤ 0. Each round SHALL: (1) run all hooks with `(attacker, defender)` in initiative order, (2) have the combatant with higher `Initiative` deal `rng.Intn(MaxDamage-MinDamage+1) + MinDamage` damage reduced by the defender's `Armour` (floor 0), (3) have the other combatant do the same. Ties in initiative are broken in the player's favour. The function SHALL append a human-readable line to `CombatState.Log` for each event (hook effects, attacks, damage). It SHALL NOT modify the original `Combatant` values passed in.

#### Scenario: Higher-initiative combatant attacks first
- **WHEN** player has `Initiative 10` and enemy has `Initiative 5`
- **THEN** the first log entry describes the player's attack

#### Scenario: Armour reduces damage to floor zero
- **WHEN** attacker deals `MinDamage == MaxDamage == 3` and defender `Armour == 5`
- **THEN** the defender's HP decreases by 0 (not negative)

#### Scenario: Combat ends when HP reaches zero
- **WHEN** one combatant's HP is reduced to ≤ 0
- **THEN** no further rounds are resolved and the log contains a "defeated" entry

#### Scenario: Original Combatant values are not mutated
- **WHEN** `resolveCombat` returns
- **THEN** the `player` and `enemy` values passed in are unchanged

### Requirement: RoundHooks are called before each round's attacks
Before each round's attack exchange, the system SHALL call every `RoundHook` in `CombatState.Hooks` with `(attacker *Combatant, defender *Combatant)` where `attacker` is the higher-initiative combatant and `defender` is the other. Hooks may mutate the combatant structs in place. If a hook changes HP to ≤ 0, combat ends without executing attacks for that round.

#### Scenario: Hook fires before attacks each round
- **WHEN** a hook is registered that increments `attacker.MinDamage` by 1
- **THEN** every round's damage roll uses the incremented value

#### Scenario: Hook that kills defender ends combat before attacks
- **WHEN** a hook sets `defender.HP = 0`
- **THEN** no attack log entries appear for that round

### Requirement: CombatResult is encoded in the final CombatState
After `resolveCombat` returns, `CombatState` SHALL expose which side won via a `PlayerWon bool` field (true if enemy HP ≤ 0 at end). The caller is responsible for acting on this result.

#### Scenario: PlayerWon is true when enemy is defeated
- **WHEN** the enemy's HP reaches 0 before the player's
- **THEN** `combatState.PlayerWon == true`

#### Scenario: PlayerWon is false when player is defeated
- **WHEN** the player's HP reaches 0 before the enemy's
- **THEN** `combatState.PlayerWon == false`

### Requirement: Game ends when player HP reaches zero
After `resolveCombat` returns with `PlayerWon == false`, the model Update SHALL call `tea.Quit`. The game exits; no further input is accepted.

#### Scenario: Game quits on player defeat
- **WHEN** combat resolves with `PlayerWon == false`
- **THEN** the returned `tea.Cmd` is `tea.Quit`

### Requirement: Exploration resumes when enemy is defeated
After `resolveCombat` returns with `PlayerWon == true`, the model SHALL set `m.screenMode = ScreenNormal`, `m.paused = false`, and remove the defeated animal from the local map.

#### Scenario: Screen returns to normal on victory
- **WHEN** combat resolves with `PlayerWon == true` and the player dismisses the result
- **THEN** `m.screenMode == ScreenNormal`

#### Scenario: Defeated animal is removed from local map
- **WHEN** combat resolves with `PlayerWon == true`
- **THEN** the animal that was fought is absent from `m.localMap.Animals`

### Requirement: Player combat stats are derived from base values and equipment
The player `Combatant` SHALL be constructed with: `HP = MaxHP = 20 + sum of equipped item HPBonus`, `Armour = sum of equipped item ArmourBonus`, `MinDamage = 1 + sum of equipped item MinDamageBonus`, `MaxDamage = 3 + sum of equipped item MaxDamageBonus`, `Initiative = 5 + sum of equipped item InitiativeBonus`. For this iteration all item bonus fields are 0 (no items grant bonuses yet); the structure exists for future use.

#### Scenario: Player with no bonuses has base stats
- **WHEN** the player has no equipped items with bonuses
- **THEN** player combatant has `HP=20`, `Armour=0`, `MinDamage=1`, `MaxDamage=3`, `Initiative=5`
