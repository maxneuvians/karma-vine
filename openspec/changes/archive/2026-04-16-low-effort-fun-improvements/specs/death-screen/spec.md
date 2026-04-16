## ADDED Requirements

### Requirement: ScreenDeath is a valid ScreenMode
The system SHALL add `ScreenDeath` to the `ScreenMode` enum in `types.go`.

#### Scenario: ScreenDeath constant is defined
- **WHEN** `ScreenDeath` is referenced in code
- **THEN** it compiles without error as a valid `ScreenMode` value distinct from `ScreenNormal`, `ScreenInventory`, and `ScreenCombat`

### Requirement: Model stores the name of the last killer
`Model` SHALL include a `deathKiller string` field. `NewModel()` SHALL initialise it to `""`.

#### Scenario: New model has empty deathKiller
- **WHEN** `NewModel()` is called
- **THEN** `m.deathKiller == ""`

### Requirement: Combat defeat transitions to death screen instead of quitting
When the combat screen is dismissed (Space or Enter pressed) and `m.combatState.PlayerWon == false`, the system SHALL:
1. Set `m.deathKiller = m.combatState.Enemy.Name`.
2. Set `m.screenMode = ScreenDeath`.
3. NOT call `tea.Quit`.

#### Scenario: Defeat sets ScreenDeath and records killer
- **WHEN** `m.combatState.PlayerWon == false` and the player dismisses the combat screen
- **THEN** `m.screenMode == ScreenDeath` and `m.deathKiller` equals the enemy's name

#### Scenario: Defeat does not quit the program
- **WHEN** `m.combatState.PlayerWon == false` and the player dismisses the combat screen
- **THEN** no `tea.Quit` command is returned

### Requirement: Death screen is rendered as a fullscreen overlay
When `m.screenMode == ScreenDeath`, the system SHALL render a fullscreen death screen. The screen SHALL display:
- A prominent "YOU DIED" heading.
- The text "Killed by: <deathKiller>".
- Instructions: "Press R to restart  |  Press Q to quit".

#### Scenario: Death screen contains YOU DIED text
- **WHEN** `View()` is called with `m.screenMode == ScreenDeath`
- **THEN** the output contains the text "YOU DIED"

#### Scenario: Death screen shows killer name
- **WHEN** `m.deathKiller == "Goblin"` and `m.screenMode == ScreenDeath`
- **THEN** the output contains "Goblin"

#### Scenario: Death screen shows restart instruction
- **WHEN** `m.screenMode == ScreenDeath`
- **THEN** the output contains "R" and "restart" (case-insensitive)

### Requirement: Death screen key handling
When `m.screenMode == ScreenDeath`:
- Pressing `r` (or `R`) SHALL reset the model to `NewModel()` and return no command.
- Pressing `q` or `ctrl+c` SHALL return `tea.Quit`.
- All other keys SHALL be no-ops.

#### Scenario: r on death screen restarts the game
- **WHEN** `m.screenMode == ScreenDeath` and the player presses `r`
- **THEN** the returned model equals `NewModel()` in all key fields (playerHP, mode, screenMode)

#### Scenario: q on death screen quits
- **WHEN** `m.screenMode == ScreenDeath` and the player presses `q`
- **THEN** the returned command is `tea.Quit`
