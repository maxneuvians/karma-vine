## ADDED Requirements

### Requirement: combatPaused field tracks pause state
The system SHALL add `combatPaused bool` to `Model`. `combatPaused` SHALL be set to `true` each time `ScreenCombat` is entered. It SHALL be set to `false` when the player unpauses. `combatLogIndex` SHALL also be reset to `0` on each combat entry (existing behaviour preserved).

#### Scenario: combatPaused is true on combat start
- **WHEN** the player initiates combat and `m.screenMode` is set to `ScreenCombat`
- **THEN** `m.combatPaused == true`

#### Scenario: combatPaused is false after unpause
- **WHEN** the player presses Space or Enter during the paused state
- **THEN** `m.combatPaused == false`

## MODIFIED Requirements

### Requirement: Combat is entered with a CombatTickMsg scheduled
When `m.screenMode` is set to `ScreenCombat`, the system SHALL set `m.combatPaused = true` and shall NOT schedule a `CombatTickMsg`. The first tick SHALL only be scheduled when the player presses Space or Enter to unpause.

#### Scenario: Entering combat does NOT schedule a CombatTickMsg
- **WHEN** combat is initiated and `m.screenMode` is set to `ScreenCombat`
- **THEN** the returned `tea.Cmd` does NOT include a `CombatTickMsg` tick (combatPaused is true)

#### Scenario: Pressing Space while paused schedules the first CombatTickMsg
- **WHEN** `m.combatPaused == true` and the Space key is pressed
- **THEN** `m.combatPaused == false` and the returned `tea.Cmd` is non-nil (schedules a tick)

#### Scenario: Pressing Enter while paused schedules the first CombatTickMsg
- **WHEN** `m.combatPaused == true` and the Enter key is pressed
- **THEN** `m.combatPaused == false` and the returned `tea.Cmd` is non-nil (schedules a tick)

#### Scenario: Space key during active playback is a no-op
- **WHEN** `m.combatPaused == false` and Space is pressed during ongoing playback
- **THEN** playback continues unchanged and no additional tick is scheduled
