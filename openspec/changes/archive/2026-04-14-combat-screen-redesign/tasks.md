## 1. Types and Constants

- [x] 1.1 Add `type CombatTickMsg struct{}` to `types.go`
- [x] 1.2 Add constants `CombatSpeedSlow = 0`, `CombatSpeedNormal = 1`, `CombatSpeedFast = 2` to `types.go`
- [x] 1.3 Implement `combatSpeedDuration(speed int) time.Duration` in `combat.go` (or `types.go`): returns `3*time.Second`, `1*time.Second`, `200*time.Millisecond` for 0/1/2; default `1*time.Second`

## 2. Model

- [x] 2.1 Add `combatLogIndex int` field to `Model` in `model.go`
- [x] 2.2 Add `combatSpeed int` field to `Model` in `model.go`
- [x] 2.3 Initialise `combatSpeed = CombatSpeedNormal` in `NewModel()`
- [x] 2.4 Add `CombatTickMsg` case to `Update`'s message switch: if `m.screenMode != ScreenCombat` return no-op; else if `m.combatLogIndex < m.combatState.Round` increment `m.combatLogIndex` and return `tea.Tick(combatSpeedDuration(m.combatSpeed), func(t time.Time) tea.Msg { return CombatTickMsg{} })`; else return nil cmd

## 3. Combat Initiation — Schedule First Tick and Reset Index

- [x] 3.1 In all combat initiation sites (surface `g`-on-animal, dungeon movement into enemy, enemy moves onto player cell): set `m.combatLogIndex = 0` immediately after `m.screenMode = ScreenCombat`
- [x] 3.2 Return `tea.Tick(combatSpeedDuration(m.combatSpeed), func(t time.Time) tea.Msg { return CombatTickMsg{} })` as the `tea.Cmd` when entering `ScreenCombat` (replace the current `nil` return)

## 4. Input — Speed Controls

- [x] 4.1 In the `ScreenCombat` key handler block in `handleKey`, add `case "]":` that sets `m.combatSpeed = min(CombatSpeedFast, m.combatSpeed+1)` and returns `m, nil`
- [x] 4.2 Add `case "[":` that sets `m.combatSpeed = max(CombatSpeedSlow, m.combatSpeed-1)` and returns `m, nil`
- [x] 4.3 Ensure `[` and `]` are handled before the `ScreenCombat` early-return so they do not fall through to the world time-scale handler

## 5. Combat Log Grouping Helper

- [x] 5.1 Implement `combatLogLinesUpTo(log []string, roundIndex int) []string` in `combat.go`: scan `log`; collect all lines up to (but not including) the `(roundIndex+1)`-th line that starts with `\"Round \"` prefix; return the collected lines
- [x] 5.2 Implement `hpAtRound(startHP int, log []string, roundIndex int, combatantName string) int` in `combat.go` (or `render.go`): scan `combatLogLinesUpTo(log, roundIndex)` for lines matching `\"<combatantName> takes N damage\"`; subtract each N from `startHP`; clamp to 0; return result

## 6. Combat Screen Renderer — Three-Panel Layout

- [x] 6.1 Rewrite `renderCombatScreen(m Model) string` in `render.go`:
  - Compute `logRows = m.viewportH / 3`, `topH = m.viewportH - logRows`
  - Compute `leftW = max(20, m.viewportW*40/100)`, `rightW = max(20, m.viewportW*40/100)`
- [x] 6.2 Implement `renderHeroPanel(m Model, width, height int) string`: centre the `ragdoll` art vertically in the top portion; below it render `"Hero"` or player name, HP progress bar (`renderProgressBar` from hud-player-stats), `HP x/y`, `ARM:N`, `DMG: min-max`, `Initiative: N`; use `combatLogIndex` to derive current HP via `hpAtRound`; pad to `height` rows; constrain to `width`
- [x] 6.3 Implement `renderEnemyPanel(m Model, width, height int) string`: render enemy `Template.Char` (or `combatState.Enemy.Name[0]` for surface animals) centred in a lipgloss bordered box; below render enemy name, HP bar, stats using same pattern as hero panel; pad to `height` rows; constrain to `width`
- [x] 6.4 Implement `renderCombatLog(m Model, width, height int) string`: header line with speed label (`Slow`/`Normal`/`Fast`) and `[ ] speed` hint; compute visible lines via `combatLogLinesUpTo`; show the last `height-2` lines that fit; if `combatLogIndex >= combatState.Round`, append Victory/Defeated banner and dismiss hint at the bottom
- [x] 6.5 Assemble: `topRow = lipgloss.JoinHorizontal(lipgloss.Top, heroPanel, enemyPanel)`; final = `lipgloss.JoinVertical(lipgloss.Left, topRow, logPanel)`; return final

## 7. Tests

- [x] 7.1 `TestCombatSpeedDuration_Slow` in `combat_test.go`: `combatSpeedDuration(0)` == `3*time.Second`
- [x] 7.2 `TestCombatSpeedDuration_Fast` in `combat_test.go`: `combatSpeedDuration(2)` == `200*time.Millisecond`
- [x] 7.3 `TestCombatSpeedDuration_OutOfRange` in `combat_test.go`: `combatSpeedDuration(99)` == `1*time.Second`
- [x] 7.4 `TestNewModel_CombatSpeed` in `game_test.go`: `NewModel()` → `combatSpeed == CombatSpeedNormal`
- [x] 7.5 `TestCombatTickMsg_AdvancesIndex` in `game_test.go`: model in `ScreenCombat` with `combatLogIndex=2`, `combatState.Round=5` receives `CombatTickMsg` → `combatLogIndex==3`, cmd non-nil
- [x] 7.6 `TestCombatTickMsg_NoAdvancePastFinal` in `game_test.go`: `combatLogIndex == combatState.Round` receives `CombatTickMsg` → index unchanged, cmd nil
- [x] 7.7 `TestCombatTickMsg_NoopOutsideCombat` in `game_test.go`: `screenMode == ScreenNormal` receives `CombatTickMsg` → model unchanged
- [x] 7.8 `TestHandleKey_BracketIncreasesSpeed` in `input_test.go`: `]` in `ScreenCombat` with `combatSpeed=CombatSpeedSlow` → `combatSpeed==CombatSpeedNormal`
- [x] 7.9 `TestHandleKey_BracketClampsFast` in `input_test.go`: `]` with `combatSpeed=CombatSpeedFast` → unchanged
- [x] 7.10 `TestHandleKey_LeftBracketDecreasesSpeed` in `input_test.go`: `[` with `combatSpeed=CombatSpeedFast` → `CombatSpeedNormal`
- [x] 7.11 `TestHandleKey_LeftBracketClampsSlow` in `input_test.go`: `[` with `combatSpeed=CombatSpeedSlow` → unchanged
- [x] 7.12 `TestCombatLogLinesUpTo_ReturnsRound1Only` in `combat_test.go`: log with "Round 1:" and "Round 2:" lines; `roundIndex=1` → only round 1 lines returned
- [x] 7.13 `TestCombatLogLinesUpTo_ZeroIndex` in `combat_test.go`: `roundIndex=0` → empty slice
- [x] 7.14 `TestHpAtRound_NoDamage` in `combat_test.go`: no matching damage lines → `startHP` returned unchanged
- [x] 7.15 `TestHpAtRound_TakesOneDamage` in `combat_test.go`: log contains `"Hero takes 5 damage"` → `startHP-5` returned
- [x] 7.16 `TestRenderCombatScreen_ContainsRagdoll` in `render_test.go`: output contains `~O~` (ragdoll line)
- [x] 7.17 `TestRenderCombatScreen_ContainsEnemyGlyph` in `render_test.go`: enemy with `Name="Wolf"` → output contains `W` or `"Wolf"`
- [x] 7.18 `TestRenderCombatScreen_NoBannerBeforePlaybackComplete` in `render_test.go`: `combatLogIndex < combatState.Round` → output does NOT contain `"Victory!"` or `"Defeated!"`
- [x] 7.19 `TestRenderCombatScreen_VictoryBannerAfterPlayback` in `render_test.go`: `combatLogIndex >= combatState.Round`, `PlayerWon=true` → output contains `"Victory!"`
- [x] 7.20 `TestRenderCombatScreen_SpeedLabelSlow` in `render_test.go`: `combatSpeed=CombatSpeedSlow` → output contains `"Slow"`
- [x] 7.21 `TestRenderCombatScreen_SpeedLabelFast` in `render_test.go`: `combatSpeed=CombatSpeedFast` → output contains `"Fast"`
- [x] 7.22 Update any existing render tests that assert on the old flat combat layout (single-panel stat lines) to expect the new three-panel structure
