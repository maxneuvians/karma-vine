## 1. Model

- [x] 1.1 Add `playerHP int` and `playerMaxHP int` fields to `Model` in `model.go`
- [x] 1.2 Add `showHelpPanel bool` field to `Model` in `model.go`
- [x] 1.3 Initialise `playerHP = 20` and `playerMaxHP = 20` in `NewModel()`
- [x] 1.4 Confirm `showHelpPanel` zero-value (`false`) requires no explicit init

## 2. Progress Bar Utility

- [x] 2.1 Implement `renderProgressBar(current, max, width int, fillColor, emptyColor string) string` in `render.go`: compute `filled = max(0, min(width, current*width/max))`, build string of `filled` × `█` + `(width-filled)` × `░`, apply `fillColor` to filled segment and `emptyColor` to empty segment via lipgloss, return concatenated string
- [x] 2.2 Guard edge cases: `width <= 0` returns `""`, `max <= 0` returns all-empty bar, `current >= max` returns all-filled bar

## 3. HUD Rewrite

- [x] 3.1 Rewrite `renderHUD` in `render.go`: compute armour sum from `m.inventory.Equipped` item `ArmourBonus` fields; call `renderProgressBar` for HP bar with width proportional to viewport (e.g. 10–20 chars, capped so full HUD line ≤ `m.viewportW`)
- [x] 3.2 Format HUD as: `[HP bar] HP current/max  ARM:N  <tile/mode info>  <clock>  <speed>  [PAUSED]  ? help`
- [x] 3.3 Remove the `[PAUSED]` duplicate path (keep it in the rewritten HUD) and remove old `text +=` append
- [x] 3.4 Keep mode-contextual tile info (biome, temp, coords, dungeon depth) in the rewritten HUD

## 4. Remove Key-Bar Row

- [x] 4.1 Delete the `renderKeyBar` function from `render.go`
- [x] 4.2 In `buildView`, change `lipgloss.JoinVertical(lipgloss.Left, mapView, renderHUD(m), renderKeyBar(m))` to `lipgloss.JoinVertical(lipgloss.Left, mapView, renderHUD(m))`
- [x] 4.3 Update map height calculation: change `mapH := m.viewportH - 2` to `mapH := m.viewportH - 1` (only one chrome row now)

## 5. Help Panel Renderer

- [x] 5.1 Implement `renderHelpPanel(m Model) string` in `render.go`: build a `[]string` of lines starting with a `" Key Bindings"` header and separator
- [x] 5.2 Add mode-specific binding sections: universal bindings (q quit, space pause, i inventory, \ sidebar), then ModeWorld-specific, ModeLocal-specific, ModeDungeon-specific blocks
- [x] 5.3 Clamp output: `lines = lines[:min(len(lines), m.viewportH)]` before `strings.Join(lines, "\n")`
- [x] 5.4 Ensure each line is rendered with `lipgloss.NewStyle().MaxWidth(m.viewportW)` to prevent horizontal overflow

## 6. Input — Help Panel Toggle

- [x] 6.1 In `handleKey`, find the `case "?"` branch that currently toggles `m.showSidebar` and replace it with `m.showHelpPanel = !m.showHelpPanel`
- [x] 6.2 Add `case "\\"` branch to toggle `m.showSidebar` (sidebar moves to backslash)
- [x] 6.3 Ensure the `ScreenCombat` and `ScreenInventory` early-return blocks appear before the `?` handling so `?` is suppressed in those modes (already the case structurally; verify no fallthrough)

## 7. buildView Dispatch

- [x] 7.1 In `buildView`, add a `showHelpPanel` early-return after the `ScreenCombat` check and before the `ScreenInventory` check: `if m.showHelpPanel { return renderHelpPanel(m) }`

## 8. Combat — HP Carry-Over

- [x] 8.1 In `buildPlayerCombatant`, change `HP: 20, MaxHP: 20` to `HP: m.playerHP, MaxHP: m.playerMaxHP`
- [x] 8.2 In the victory dismiss handler (`handleKey`, `ScreenCombat` + `enter`/`space` + `PlayerWon == true`), after clearing `combatState`, set `m.playerHP = m.combatState.Player.HP` (capture before clearing) so carry-over HP is written back to the model

## 9. Tests

- [x] 9.1 `TestRenderProgressBar_HalfFull` in `render_test.go`: `renderProgressBar(5, 10, 10, ...)` → exactly 5 `█` and 5 `░`
- [x] 9.2 `TestRenderProgressBar_Full` in `render_test.go`: `current == max` → no `░`
- [x] 9.3 `TestRenderProgressBar_Zero` in `render_test.go`: `current == 0` → no `█`
- [x] 9.4 `TestRenderProgressBar_ZeroWidth` in `render_test.go`: width 0 → empty string
- [x] 9.5 `TestRenderProgressBar_ZeroMax` in `render_test.go`: max 0 → all empty
- [x] 9.6 `TestRenderHUD_ContainsHelpHint` in `render_test.go`: `renderHUD` output contains `? help`
- [x] 9.7 `TestRenderHUD_ContainsPaused` in `render_test.go`: `m.paused == true` → output contains `[PAUSED]`
- [x] 9.8 `TestRenderHUD_ContainsArmour` in `render_test.go`: default model → output contains `ARM:0`
- [x] 9.9 `TestRenderHUD_HPBarPresent` in `render_test.go`: output contains `█` or `░`
- [x] 9.10 `TestBuildView_NoKeyBar` in `render_test.go`: `buildView` output does NOT contain `↑↓←→/wasd move`
- [x] 9.11 `TestBuildView_HelpPanelShown` in `render_test.go`: `m.showHelpPanel == true` → `buildView` equals `renderHelpPanel(m)`
- [x] 9.12 `TestRenderHelpPanel_HeightClamped` in `render_test.go`: `m.viewportH == 5` → output has ≤ 5 lines
- [x] 9.13 `TestRenderHelpPanel_LocalBindings` in `render_test.go`: `m.mode == ModeLocal` → output contains `g` and `i`
- [x] 9.14 `TestHandleKey_QuestionMarkTogglesHelpPanel` in `input_test.go`: `?` in `ScreenNormal` → `showHelpPanel` toggled
- [x] 9.15 `TestHandleKey_QuestionMarkSuppressedInInventory` in `input_test.go`: `?` in `ScreenInventory` → `showHelpPanel` unchanged
- [x] 9.16 `TestHandleKey_BackslashTogglesSidebar` in `input_test.go`: `\` in `ScreenNormal` → `showSidebar` toggled
- [x] 9.17 `TestNewModel_PlayerHP` in `game_test.go`: `NewModel()` → `playerHP == 20`, `playerMaxHP == 20`, `showHelpPanel == false`
- [x] 9.18 `TestBuildPlayerCombatant_UsesModelHP` in `combat_test.go`: `m.playerHP == 15` → combatant `HP == 15`
- [x] 9.19 Update any existing tests that assert on `renderKeyBar` output or check that `buildView` contains key-bar text — remove or rewrite those assertions
