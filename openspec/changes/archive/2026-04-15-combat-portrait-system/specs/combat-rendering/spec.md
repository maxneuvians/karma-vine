## MODIFIED Requirements

### Requirement: Combat screen uses a three-panel layout
`renderCombatScreen` SHALL render a three-panel layout:
- **Top-left panel** (width ≈ 40% of viewport, height = viewport height minus log rows): hero portrait (40×20 unicode block-character portrait centred) followed by player name, HP progress bar (`HP current/max`), Armour, Damage range, Initiative
- **Top-right panel** (width ≈ 40% of viewport, same height): enemy portrait (40×20 unicode block-character portrait centred), followed by enemy name, HP progress bar, Armour, Damage range, Initiative
- **Bottom log panel** (full viewport width, height = `viewportH / 3`): when `m.combatPaused == true`, shows a centred `"[ Space ] Begin Combat"` prompt and no log lines; when unpaused, shows combat log lines revealed so far, a speed indicator showing current playback speed (`Slow` / `Normal` / `Fast`) and `[ ] speed` hint

The two top panels SHALL be joined horizontally. The combined top and bottom SHALL be joined vertically to fill `m.viewportH` rows.

#### Scenario: Left panel contains hero portrait block characters
- **WHEN** `renderCombatScreen` is called
- **THEN** the output contains at least one unicode block character from the player portrait (rune ≥ U+2580)

#### Scenario: Right panel contains enemy portrait block characters
- **WHEN** `renderCombatScreen` is called with any enemy
- **THEN** the output contains at least one unicode block character from the enemy portrait

#### Scenario: Both panels show stat labels
- **WHEN** `renderCombatScreen` is called with a player `HP=15, MaxHP=20` and enemy `HP=8, MaxHP=12`
- **THEN** the output contains `15/20` and `8/12`

#### Scenario: Top panels do not overflow viewport height
- **WHEN** `m.viewportH == 20`
- **THEN** the top panel section occupies exactly `viewportH - viewportH/3` rows

#### Scenario: Log panel shows pause prompt when combat is paused
- **WHEN** `m.combatPaused == true`
- **THEN** the bottom panel contains `"Begin Combat"` and does NOT contain `Slow`, `Normal`, or `Fast`

#### Scenario: Log panel shows speed hint when combat is unpaused
- **WHEN** `m.combatPaused == false`
- **THEN** the bottom panel contains the speed label (`Slow`, `Normal`, or `Fast`) and the `[ ] speed` hint
