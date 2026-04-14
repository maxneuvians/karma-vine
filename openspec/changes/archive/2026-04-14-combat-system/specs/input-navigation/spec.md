## ADDED Requirements

### Requirement: ScreenCombat suppresses all normal input
When `m.screenMode == ScreenCombat`, `handleKey` SHALL ignore all keys except `enter`, `space` (dismiss result / quit), and `q`/`ctrl+c` (quit game). Movement keys, inventory keys, and all other bindings SHALL be no-ops.

#### Scenario: Movement is suppressed in ScreenCombat
- **WHEN** `m.screenMode == ScreenCombat` and the player presses `up`
- **THEN** `m.playerPos` is unchanged

#### Scenario: Inventory key is suppressed in ScreenCombat
- **WHEN** `m.screenMode == ScreenCombat` and the player presses `i`
- **THEN** `m.screenMode` remains `ScreenCombat`

### Requirement: enter/space dismisses the combat result screen
When `m.screenMode == ScreenCombat` and the player presses `enter` or `space`: if `m.combatState.PlayerWon == true`, the system SHALL set `m.screenMode = ScreenNormal`, `m.paused = false`, and remove the defeated animal; if `m.combatState.PlayerWon == false`, the system SHALL return `tea.Quit`.

#### Scenario: enter on victory returns to exploration
- **WHEN** `m.screenMode == ScreenCombat`, `m.combatState.PlayerWon == true`, and the player presses `enter`
- **THEN** `m.screenMode == ScreenNormal` and `m.paused == false`

#### Scenario: enter on defeat quits the game
- **WHEN** `m.screenMode == ScreenCombat`, `m.combatState.PlayerWon == false`, and the player presses `enter`
- **THEN** the returned `tea.Cmd` is `tea.Quit`
