## ADDED Requirements

### Requirement: CombatTickMsg type and CombatSpeed constants are defined
The system SHALL define `type CombatTickMsg struct{}` in `types.go`. It SHALL also define integer constants `CombatSpeedSlow = 0`, `CombatSpeedNormal = 1`, `CombatSpeedFast = 2` and a helper `combatSpeedDuration(speed int) time.Duration` returning `3000ms`, `1000ms`, `200ms` respectively. Any out-of-range speed SHALL return `1000ms`.

#### Scenario: CombatSpeedSlow maps to 3000ms
- **WHEN** `combatSpeedDuration(CombatSpeedSlow)` is called
- **THEN** the result is `3 * time.Second`

#### Scenario: CombatSpeedFast maps to 200ms
- **WHEN** `combatSpeedDuration(CombatSpeedFast)` is called
- **THEN** the result is `200 * time.Millisecond`

#### Scenario: Out-of-range speed returns 1000ms
- **WHEN** `combatSpeedDuration(99)` is called
- **THEN** the result is `1 * time.Second`

### Requirement: Model tracks combat playback state
The system SHALL add `combatLogIndex int` and `combatSpeed int` to `Model`. `NewModel()` SHALL initialise `combatSpeed = CombatSpeedNormal`. `combatLogIndex` is reset to `0` each time `ScreenCombat` is entered. `combatSpeed` persists across combats.

#### Scenario: NewModel initialises combatSpeed to Normal
- **WHEN** `NewModel()` is called
- **THEN** `m.combatSpeed == CombatSpeedNormal`

#### Scenario: combatLogIndex resets on combat start
- **WHEN** the player initiates combat and `m.screenMode` is set to `ScreenCombat`
- **THEN** `m.combatLogIndex == 0`

### Requirement: Combat is entered with a CombatTickMsg scheduled
When `m.screenMode` is set to `ScreenCombat`, the system SHALL return a `tea.Tick(combatSpeedDuration(m.combatSpeed), ...)` command that fires a `CombatTickMsg`. This schedules the first round reveal.

#### Scenario: Entering combat schedules a CombatTickMsg
- **WHEN** combat is initiated
- **THEN** the returned `tea.Cmd` is non-nil (schedules a tick)

### Requirement: CombatTickMsg advances combatLogIndex by one round block
When `Update` receives a `CombatTickMsg` and `m.screenMode == ScreenCombat`:
- If `m.combatLogIndex < m.combatState.Round`, increment `m.combatLogIndex` by 1 and return a new `tea.Tick(combatSpeedDuration(m.combatSpeed), ...)` command to schedule the next advance.
- If `m.combatLogIndex >= m.combatState.Round`, do not increment and return no tick command (playback is complete).
- If `m.screenMode != ScreenCombat`, treat the tick as a no-op and return nil.

#### Scenario: CombatTickMsg advances index when playback is incomplete
- **WHEN** `combatLogIndex == 2` and `combatState.Round == 5` and a `CombatTickMsg` is received
- **THEN** `combatLogIndex == 3` and a non-nil cmd is returned

#### Scenario: CombatTickMsg does not advance past final round
- **WHEN** `combatLogIndex == combatState.Round` and a `CombatTickMsg` is received
- **THEN** `combatLogIndex` is unchanged and the returned cmd is nil

#### Scenario: Stale CombatTickMsg outside ScreenCombat is a no-op
- **WHEN** `m.screenMode == ScreenNormal` and a `CombatTickMsg` is received
- **THEN** model is unchanged and cmd is nil

### Requirement: Speed change takes effect on the next scheduled tick
When `m.combatSpeed` is changed via `[` or `]` during `ScreenCombat`, the new speed applies to the next `CombatTickMsg` reschedule — the current in-flight tick fires at the old speed but the next one uses the new duration.

#### Scenario: Speed change is reflected in combatSpeed field
- **WHEN** `]` is pressed during ScreenCombat with `combatSpeed == CombatSpeedSlow`
- **THEN** `m.combatSpeed == CombatSpeedNormal`
