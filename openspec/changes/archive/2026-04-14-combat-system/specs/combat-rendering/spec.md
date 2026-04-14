## ADDED Requirements

### Requirement: ScreenCombat renders a fullscreen combat view
When `m.screenMode == ScreenCombat`, `buildView` SHALL delegate entirely to `renderCombatScreen(m)`. The combat screen SHALL fill the full viewport (`m.viewportW × m.viewportH`) and suppress the normal map, HUD, and key bar.

#### Scenario: buildView returns combat screen when ScreenCombat is active
- **WHEN** `m.screenMode == ScreenCombat`
- **THEN** the output of `buildView` equals the output of `renderCombatScreen(m)` and contains no map tiles

### Requirement: Combat screen shows two combatant stat panels
`renderCombatScreen` SHALL render a header line showing `"⚔ Combat"` and two stat panels side by side: the left panel for the player and the right panel for the enemy. Each panel SHALL show: `Name`, current `HP / MaxHP`, `Armour`, `Damage` range (`MinDamage–MaxDamage`), and `Initiative`.

#### Scenario: Player panel shows correct stats
- **WHEN** `m.combatState.Player` has `HP=15`, `MaxHP=20`, `Armour=2`, `MinDamage=1`, `MaxDamage=4`, `Initiative=5`
- **THEN** the rendered output contains `15/20`, `Armour: 2`, `1-4`, and `Initiative: 5`

#### Scenario: Enemy panel shows correct stats
- **WHEN** `m.combatState.Enemy` has `Name="Wolf"`, `HP=8`, `MaxHP=12`
- **THEN** the rendered output contains `"Wolf"` and `8/12`

### Requirement: Combat screen shows a combat log
Below the stat panels, `renderCombatScreen` SHALL render the lines from `m.combatState.Log`, one per row, newest at the bottom. If there are more log lines than available rows, only the most recent lines that fit SHALL be shown (no scroll).

#### Scenario: Log lines appear in order
- **WHEN** `m.combatState.Log` contains `["Round 1", "Player attacks", "Enemy attacks"]`
- **THEN** the rendered output contains all three lines in that order

#### Scenario: Excess log lines are truncated from the top
- **WHEN** the log has more lines than available rows
- **THEN** only the most recent lines that fit in the viewport are shown

### Requirement: Combat screen shows a result banner and dismiss hint when combat is over
When `m.combatState` is set and `HP ≤ 0` for either combatant, `renderCombatScreen` SHALL display a prominent result banner: `"Victory!"` if `m.combatState.PlayerWon`, or `"Defeated!"` otherwise. Below the banner it SHALL show a hint: `"  press enter to continue"` (victory) or `"  press enter to quit"` (defeat).

#### Scenario: Victory banner appears when player won
- **WHEN** `m.combatState.PlayerWon == true`
- **THEN** the rendered output contains `"Victory!"`

#### Scenario: Defeat banner appears when player lost
- **WHEN** `m.combatState.PlayerWon == false`
- **THEN** the rendered output contains `"Defeated!"`
