## 1. Model & Types

- [x] 1.1 Add `combatPaused bool` field to `Model` struct in `types.go`
- [x] 1.2 Update `NewModel()` in `model.go` to initialise `combatPaused = false` (default zero value is correct, but verify no explicit init needed)
- [x] 1.3 Update the enter-combat path in `model.go` to set `combatPaused = true` and `combatLogIndex = 0` when `ScreenCombat` is entered — remove the auto-`CombatTickMsg` from the `combatLogIndex == 0` branch

## 2. Input Handling

- [x] 2.1 In `input.go`, add a handler for `ScreenCombat` + `combatPaused == true`: on Space or Enter key press, set `m.combatPaused = false` and return the first `tea.Tick(combatSpeedDuration(m.combatSpeed), CombatTickMsg{})` command
- [x] 2.2 Ensure Space key during active playback (`combatPaused == false`) is a no-op in combat (does not double-schedule ticks)

## 3. Portrait Data

- [x] 3.1 Create `internal/game/combat_portraits.go` and define the `portraitCell` struct with `r rune` and `color lipgloss.Color` fields
- [x] 3.2 Define `playerPortrait [][]portraitCell` as a 20×40 literal — humanoid heroic silhouette using `█`, `▓`, `▒`, `░`, `▄`, `▀`, `▌`, `▐` with appropriate skin/armour colours
- [x] 3.3 Define at least three archetype enemy portraits (humanoid, beast, undead) each as 20×40 `[][]portraitCell` literals
- [x] 3.4 Define a generic `fallbackPortrait [][]portraitCell` (20×40) for unrecognised enemy types
- [x] 3.5 Implement `enemyPortrait(char rune) [][]portraitCell` with a `switch` on `char` mapping known runes to archetypes and returning `fallbackPortrait` for unknowns

## 4. Portrait Rendering

- [x] 4.1 Implement `renderPortrait(p [][]portraitCell, panelWidth int) string` in `combat_portraits.go` — renders each cell with lipgloss foreground colour, clips rows to `panelWidth`, joins rows with `\n`
- [x] 4.2 Update `renderHeroPanel` in `render.go` to call `renderPortrait(playerPortrait, panelWidth)` in place of the ragdoll ASCII art loop
- [x] 4.3 Update the enemy panel rendering in `render.go` to call `renderPortrait(enemyPortrait(cs.Enemy.Template.Char), panelWidth)` in place of the single-glyph centred box

## 5. Paused UI

- [x] 5.1 Update `renderCombatLog` in `render.go`: when `m.combatPaused == true`, render only a centred `"[ Space ] Begin Combat"` prompt line styled with the accent colour; suppress all log lines and the speed hint
- [x] 5.2 Verify the bottom log panel height calculation is unchanged regardless of pause state

## 6. Tests

- [x] 6.1 Write unit tests for `portraitCell`, `playerPortrait` dimensions (20 rows × 40 cols), and `enemyPortrait` fallback behaviour in `combat_portraits_test.go`
- [x] 6.2 Write unit tests for `renderPortrait`: full-width output has 20 lines, clipping to `panelWidth < 40` truncates rows correctly
- [x] 6.3 Write unit tests for pause flow in `model_test.go` or `input_test.go`: combat entry sets `combatPaused = true` and no tick cmd; Space press sets `combatPaused = false` and returns non-nil cmd
- [x] 6.4 Update `renderCombatScreen` tests in `render_test.go`: assert portrait block characters present in output, assert pause prompt shown when `combatPaused = true`, assert speed hint shown when `combatPaused = false`
- [x] 6.5 Run `go test ./internal/game/... -cover` and confirm coverage remains ≥ 90%
