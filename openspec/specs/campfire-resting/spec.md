## ADDED Requirements

### Requirement: Model tracks campfire rest cooldown
`Model` SHALL include a `restCooldown int` field. `NewModel()` SHALL initialise it to `0`. On each `TickMsg` where `m.restCooldown > 0`, the system SHALL decrement `m.restCooldown` by 1.

#### Scenario: New model has zero rest cooldown
- **WHEN** `NewModel()` is called
- **THEN** `m.restCooldown == 0`

#### Scenario: Rest cooldown decrements on each tick
- **WHEN** `m.restCooldown == 5` and a `TickMsg` is received
- **THEN** `m.restCooldown == 4`

#### Scenario: Rest cooldown does not go below zero
- **WHEN** `m.restCooldown == 0` and a `TickMsg` is received
- **THEN** `m.restCooldown == 0`

### Requirement: Player can rest at a campfire cell to restore HP
When `m.mode == ModeLocal` and the player presses `r`:
1. If the current player cell has `m.currentLocal.Ground[px][py].HasFire == false`, the key SHALL be a no-op.
2. If `m.restCooldown > 0`, the key SHALL be a no-op.
3. Otherwise, `m.playerHP` SHALL increase by 5, capped at `MaxPlayerHP`. `m.restCooldown` SHALL be set to `60`.

#### Scenario: r on a fire cell restores 5 HP
- **WHEN** `m.mode == ModeLocal`, player is on a `HasFire == true` cell, `m.restCooldown == 0`, and `m.playerHP == 10`
- **THEN** `m.playerHP == 15` and `m.restCooldown == 60`

#### Scenario: r on a fire cell does not exceed MaxPlayerHP
- **WHEN** player is on a fire cell, `m.restCooldown == 0`, and `m.playerHP == MaxPlayerHP - 2`
- **THEN** `m.playerHP == MaxPlayerHP`

#### Scenario: r on a non-fire cell is a no-op
- **WHEN** `m.mode == ModeLocal`, player is on a cell where `HasFire == false`
- **THEN** `m.playerHP` is unchanged and `m.restCooldown` is unchanged

#### Scenario: r during cooldown is a no-op
- **WHEN** player is on a fire cell and `m.restCooldown > 0`
- **THEN** `m.playerHP` is unchanged

#### Scenario: r in world or dungeon mode is a no-op
- **WHEN** `m.mode == ModeWorld` or `m.mode == ModeDungeon` and `r` is pressed
- **THEN** `m.playerHP` and `m.restCooldown` are unchanged
