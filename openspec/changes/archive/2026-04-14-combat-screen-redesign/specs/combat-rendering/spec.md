## MODIFIED Requirements

### Requirement: ScreenCombat renders a fullscreen combat view
When `m.screenMode == ScreenCombat`, `buildView` SHALL delegate entirely to `renderCombatScreen(m)`. The combat screen SHALL fill the full viewport (`m.viewportW × m.viewportH`) and suppress the normal map, HUD, and key bar.

#### Scenario: buildView returns combat screen when ScreenCombat is active
- **WHEN** `m.screenMode == ScreenCombat`
- **THEN** the output of `buildView` equals the output of `renderCombatScreen(m)` and contains no map tiles

### Requirement: Combat screen uses a three-panel layout
`renderCombatScreen` SHALL render a three-panel layout:
- **Top-left panel** (width ≈ 40% of viewport, height = viewport height minus log rows): hero portrait (ragdoll ASCII art centred) followed by player name, HP progress bar (`HP current/max`), Armour, Damage range, Initiative
- **Top-right panel** (width ≈ 40% of viewport, same height): enemy portrait (enemy `Char` centred in a bordered box), followed by enemy name, HP progress bar, Armour, Damage range, Initiative
- **Bottom log panel** (full viewport width, height = `viewportH / 3`): combat log lines revealed so far; a speed indicator showing current playback speed (`Slow` / `Normal` / `Fast`) and `[ ] speed` hint

The two top panels SHALL be joined horizontally. The combined top and bottom SHALL be joined vertically to fill `m.viewportH` rows.

#### Scenario: Left panel contains hero portrait lines
- **WHEN** `renderCombatScreen` is called
- **THEN** the output contains lines from the ragdoll art (e.g. `~O~`)

#### Scenario: Right panel contains enemy glyph
- **WHEN** the enemy has `Template.Char == 'G'`
- **THEN** the right panel contains `G`

#### Scenario: Both panels show stat labels
- **WHEN** `renderCombatScreen` is called with a player `HP=15, MaxHP=20` and enemy `HP=8, MaxHP=12`
- **THEN** the output contains `15/20` and `8/12`

#### Scenario: Top panels do not overflow viewport height
- **WHEN** `m.viewportH == 20`
- **THEN** the top panel section occupies exactly `viewportH - viewportH/3` rows

### Requirement: Combat log reveals lines up to the current round index
The log panel SHALL display only the log lines associated with rounds 1 through `m.combatLogIndex`. Lines from later rounds SHALL NOT appear. The most recent lines that fit the log panel height SHALL be shown (truncated from the top if needed).

#### Scenario: combatLogIndex 0 shows no log lines
- **WHEN** `m.combatLogIndex == 0`
- **THEN** the log panel contains no combat narration lines (only the header/speed hint)

#### Scenario: combatLogIndex 2 shows rounds 1 and 2 lines
- **WHEN** `m.combatLogIndex == 2` and the log has lines for rounds 1–5
- **THEN** only lines belonging to rounds 1 and 2 are visible

### Requirement: HP progress bars reflect HP at the current playback round
The HP bars in both panels SHALL show HP values updated to reflect damage dealt through round `m.combatLogIndex`. The renderer SHALL scan visible log lines for damage patterns and compute current HP from `MaxHP` minus cumulative damage. Full HP is shown when `combatLogIndex == 0`.

#### Scenario: HP bar shows full HP before any rounds are revealed
- **WHEN** `m.combatLogIndex == 0`
- **THEN** player HP bar shows `MaxHP/MaxHP` (all filled)

#### Scenario: HP bar decreases as rounds are revealed
- **WHEN** `m.combatLogIndex` advances and the log contains damage events
- **THEN** the affected combatant's HP bar shows reduced HP

### Requirement: Result banner shown only after playback completes
The Victory/Defeated banner and dismiss hint SHALL appear only when `m.combatLogIndex >= m.combatState.Round` (all rounds revealed). Before that, no banner is shown.

#### Scenario: No banner before playback is complete
- **WHEN** `m.combatLogIndex < m.combatState.Round`
- **THEN** the output does NOT contain `"Victory!"` or `"Defeated!"`

#### Scenario: Victory banner shown after all rounds revealed
- **WHEN** `m.combatLogIndex >= m.combatState.Round` and `m.combatState.PlayerWon == true`
- **THEN** the output contains `"Victory!"`

#### Scenario: Defeat banner shown after all rounds revealed
- **WHEN** `m.combatLogIndex >= m.combatState.Round` and `m.combatState.PlayerWon == false`
- **THEN** the output contains `"Defeated!"`

### Requirement: Log panel shows current playback speed and speed hint
The log panel header SHALL display the current speed label (`Slow`, `Normal`, or `Fast`) and the key hint `[ ] speed`. This is always visible regardless of playback progress.

#### Scenario: Speed label shows Slow when combatSpeed is CombatSpeedSlow
- **WHEN** `m.combatSpeed == CombatSpeedSlow`
- **THEN** the output contains `Slow`

#### Scenario: Speed label shows Fast when combatSpeed is CombatSpeedFast
- **WHEN** `m.combatSpeed == CombatSpeedFast`
- **THEN** the output contains `Fast`
