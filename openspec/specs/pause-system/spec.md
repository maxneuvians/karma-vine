## ADDED Requirements

### Requirement: Model has a paused field
`Model` SHALL have a `paused bool` field. `NewModel()` SHALL leave `paused` at its zero value (`false`), meaning the game starts unpaused.

#### Scenario: New model starts unpaused
- **WHEN** `NewModel()` is called
- **THEN** `m.paused == false`

### Requirement: Space key toggles pause
The system SHALL handle the space key (`" "`) in `handleKey` at all times (all modes, all screen modes). Pressing space SHALL toggle `m.paused`: `false → true` and `true → false`.

#### Scenario: Space pauses when unpaused
- **WHEN** `m.paused == false` and the player presses space
- **THEN** `m.paused == true`

#### Scenario: Space unpauses when paused
- **WHEN** `m.paused == true` and the player presses space
- **THEN** `m.paused == false`

#### Scenario: Space toggles pause in ModeWorld
- **WHEN** `m.mode == ModeWorld` and the player presses space
- **THEN** `m.paused` is toggled

#### Scenario: Space toggles pause in ScreenInventory
- **WHEN** `m.screenMode == ScreenInventory` and the player presses space
- **THEN** `m.paused` is toggled

### Requirement: TickMsg is a no-op when paused
When `m.paused == true` and a `TickMsg` is received, the system SHALL skip advancing `timeOfDay` and skip calling `moveAnimals`. The system SHALL still reschedule the next tick by returning `tickCmd()`.

#### Scenario: Time does not advance while paused
- **WHEN** `m.paused == true` and `Update` receives a `TickMsg`
- **THEN** `m.timeOfDay` is unchanged after the update

#### Scenario: Animals do not move while paused
- **WHEN** `m.paused == true`, `m.mode == ModeLocal`, and `Update` receives a `TickMsg`
- **THEN** animal positions in `m.localMap` are unchanged after the update

#### Scenario: Tick loop continues while paused
- **WHEN** `m.paused == true` and `Update` receives a `TickMsg`
- **THEN** the returned `tea.Cmd` is non-nil (tick is rescheduled)
