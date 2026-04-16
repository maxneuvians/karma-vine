## 1. Item Stat Fixes

- [x] 1.1 In `internal/game/model.go`, add `ArmourBonus: 1` to the Cloth Tunic, Cloth Pants, and Leather Boots item literals in `NewModel()`
- [x] 1.2 In `internal/game/enemy.go`, add `DamageBonus: 1` to all Rusty Dagger item literals in loot tables
- [x] 1.3 In `internal/game/enemy.go`, add `DamageBonus: 2` to all Short Sword item literals in loot tables
- [x] 1.4 Update `game_test.go` tests that assert clothing items have zero bonuses to assert the new values (ArmourBonus 1 for tunic/pants/boots)
- [x] 1.5 Update `combat_test.go` if any test hardcodes the player's starting armour pool (now 4 instead of 1)
- [x] 1.6 Run `go test ./...` and confirm all tests pass

## 2. Campfire Resting

- [x] 2.1 Add `restCooldown int` field to `Model` struct in `internal/game/types.go` or `model.go`
- [x] 2.2 In `NewModel()`, ensure `restCooldown` initialises to `0` (zero value, no explicit set needed)
- [x] 2.3 In `Update()` tick branch (`internal/game/model.go`), decrement `m.restCooldown` by 1 if it is greater than 0
- [x] 2.4 In `internal/game/input.go`, add a `case "r":` branch in the main key handler that checks `m.mode == ModeLocal`, player cell `HasFire`, and `m.restCooldown == 0`; on success add 5 HP (capped at MaxPlayerHP) and set `restCooldown = 60`
- [x] 2.5 Write unit tests in `input_test.go` covering: rest on fire cell heals 5 HP, rest caps at MaxPlayerHP, rest on non-fire cell is no-op, rest during cooldown is no-op, rest in dungeon mode is no-op
- [x] 2.6 Run `go test ./...` and confirm all tests pass

## 3. Death Screen

- [x] 3.1 Add `ScreenDeath ScreenMode` constant to the `ScreenMode` enum in `internal/game/types.go`
- [x] 3.2 Add `deathKiller string` field to `Model` struct
- [x] 3.3 In `internal/game/input.go`, find the defeat branch (`m.combatState.PlayerWon == false` → currently `return m, tea.Quit`) and replace it with: set `m.deathKiller = m.combatState.Enemy.Name`, set `m.screenMode = ScreenDeath`, return `m, nil`
- [x] 3.4 In `internal/game/input.go`, add a `ScreenDeath` case in the screen-mode dispatch: `r` → `return NewModel(), nil`; `q` / `ctrl+c` → `return m, tea.Quit`; all others → no-op
- [x] 3.5 In `internal/game/render.go`, add a `renderDeathScreen(m Model) string` function that produces a centred fullscreen view with "YOU DIED", "Killed by: <deathKiller>", and "Press R to restart  |  Press Q to quit"
- [x] 3.6 Wire `renderDeathScreen` into `View()` — add a `case ScreenDeath:` branch that returns `renderDeathScreen(m)`
- [x] 3.7 Write unit tests covering: defeat sets `ScreenDeath` and `deathKiller`, defeat does not return `tea.Quit`, `r` on death screen returns a fresh model, `q` on death screen returns `tea.Quit`, `View()` with `ScreenDeath` contains "YOU DIED" and the killer name
- [x] 3.8 Run `go test ./...` and confirm all tests pass

## 4. Portrait Name-Based Lookup

- [x] 4.1 In `internal/game/combat_portraits.go`, add `enemyPortraitByName(name string) portrait` with a `switch name` mapping all 9 dungeon enemy names to their archetypes (humanoid/beast/undead/fallback)
- [x] 4.2 In `internal/game/render.go`, update `renderEnemyPanel` to call `enemyPortraitByName(cs.Enemy.Name)` instead of `enemyPortrait(enemyChar)` (the `enemyChar` derivation block can be removed)
- [x] 4.3 Update or add tests in `combat_portraits_test.go` to cover `enemyPortraitByName` for at least one humanoid name, one beast name, one undead name, and one unknown name
- [x] 4.4 Run `go test ./...` and confirm all tests pass

## 5. Minimum Enemy Count per Floor

- [x] 5.1 In `internal/game/dungeon.go`, change `enemyCount := depth` to `enemyCount := max(3, depth)`
- [x] 5.2 Update the dungeon generation test in `dungeon_test.go` that asserts enemy count equals depth to assert it equals `max(3, depth)`
- [x] 5.3 Run `go test ./...` and confirm all tests pass
