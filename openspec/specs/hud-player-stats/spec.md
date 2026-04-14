## ADDED Requirements

### Requirement: Model tracks live player HP and MaxHP
The system SHALL add `playerHP int` and `playerMaxHP int` fields to `Model`. `NewModel()` SHALL initialise both to `20`. These fields represent the player's current and maximum hit points, updated by combat outcomes.

#### Scenario: NewModel initialises HP to 20
- **WHEN** `NewModel()` is called
- **THEN** `m.playerHP == 20` and `m.playerMaxHP == 20`

### Requirement: HUD displays a HP progress bar
`renderHUD` SHALL render a HP progress bar of the form `[████░░░░░░]` where filled cells (`█`) represent current HP and empty cells (`░`) represent missing HP. The bar width SHALL be proportional to the viewport: the bar MUST NOT overflow `m.viewportW`. The bar SHALL be followed by a numeric `HP current/max` label.

#### Scenario: Full HP renders all filled cells
- **WHEN** `m.playerHP == m.playerMaxHP`
- **THEN** the HUD contains no `░` characters in the HP bar segment

#### Scenario: Zero HP renders all empty cells
- **WHEN** `m.playerHP == 0`
- **THEN** the HUD contains no `█` characters in the HP bar segment

#### Scenario: HP bar does not exceed viewport width
- **WHEN** `m.viewportW == 40`
- **THEN** the total rendered HUD line is ≤ 40 characters wide (ANSI escape sequences excluded)

### Requirement: HUD displays armour value
`renderHUD` SHALL include an armour badge derived from the sum of `ArmourBonus` fields across all equipped items. The badge format SHALL be `ARM:N` where N is the total armour value (0 when nothing is equipped that grants armour).

#### Scenario: No armour-granting equipment shows ARM:0
- **WHEN** the player has the default outfit with no armour bonuses
- **THEN** the HUD contains `ARM:0`

### Requirement: HUD retains tile info, clock, speed, pause indicator, and help hint
The HUD row SHALL continue to show the mode-contextual tile/dungeon info, clock, time scale, and `[PAUSED]` when `m.paused == true`. It SHALL end with `  ? help` as the sole key binding hint. The separate key-bar row SHALL be removed.

#### Scenario: HUD contains help hint
- **WHEN** `renderHUD` is called in any mode
- **THEN** the output contains `? help`

#### Scenario: HUD contains PAUSED when paused
- **WHEN** `m.paused == true`
- **THEN** the output contains `[PAUSED]`

#### Scenario: HUD does not contain old key-bar binding strings
- **WHEN** `buildView` is called in `ScreenNormal`
- **THEN** the output does NOT contain `↑↓←→/wasd move`

### Requirement: renderProgressBar is a pure width-aware function
The system SHALL implement `renderProgressBar(current, max, width int, fillColor, emptyColor string) string` returning a lipgloss-styled bar of exactly `width` characters (`█` filled, `░` empty), scaled proportionally. If `max <= 0` the bar SHALL be all empty. If `current > max` the bar SHALL be all filled.

#### Scenario: Half HP renders half-filled bar
- **WHEN** `renderProgressBar(5, 10, 10, ...)` is called
- **THEN** the returned string contains exactly 5 `█` and 5 `░` characters

#### Scenario: Width of 0 returns empty string
- **WHEN** `renderProgressBar(5, 10, 0, ...)` is called
- **THEN** the returned string is `""`
