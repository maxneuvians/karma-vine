## Context

Karma Vine is a terminal roguelike written in Go using Bubble Tea v2. The game has several polish gaps that undermine player experience: loot items have no stat effects, there is no HP recovery, defeat terminates the process immediately, enemy portraits default to a generic blob for all lowercase chars (every dungeon enemy), and early dungeon floors spawn a single enemy. These are all narrow, well-scoped changes with no shared dependencies.

Relevant files:
- `internal/game/model.go` — starting equipment definitions, model state
- `internal/game/enemy.go` — dungeon loot tables
- `internal/game/dungeon.go` — enemy count per floor (`enemyCount := depth`)
- `internal/game/combat_portraits.go` — `enemyPortrait(char rune)` selection function
- `internal/game/render.go` — `renderEnemyPanel` constructs `enemyChar` from enemy data; calls `enemyPortrait`
- `internal/game/input.go` — key handlers, `tea.Quit` on defeat, `ScreenCombat` dismissal
- `internal/game/types.go` — `ScreenMode` enum, model field declarations
- `internal/game/local.go` — `Ground.HasFire` field on each cell

## Goals / Non-Goals

**Goals:**
- Item stat values: Rusty Dagger → +1 DMG, Short Sword → +2 DMG; Cloth Tunic / Cloth Pants → +1 ARM each; Leather Boots → +1 ARM.
- Campfire resting: press `r` on a local-map fire cell to restore 5 HP (up to MaxHP), with a cooldown of 60 game-ticks to prevent abuse.
- Death screen: on combat defeat show `ScreenDeath` with killer name; `r` restarts (new `NewModel()`), `q` quits.
- Portrait fix: replace char-based `enemyPortrait(rune)` with a name-based lookup `enemyPortraitByName(string)` used by `renderEnemyPanel`; keep `enemyPortrait(rune)` for the existing test API or update tests.
- Minimum enemy floor count: clamp `enemyCount` to `max(3, depth)`.

**Non-Goals:**
- HP regeneration over time, potions, or any other healing beyond campfire resting.
- Persistent world state, save/load.
- Changing combat resolution logic, initiative, or any other stat system.
- Adding new biomes, tiles, or enemy types.

## Decisions

### 1. Item stat changes — direct literal edits

The item definitions live as inline struct literals in `model.go` (starting equipment) and `enemy.go` (loot tables). The cleanest approach is to add `ArmourBonus` and `DamageBonus` fields directly to those literals. No new data structures or item registries are warranted for five items.

*Alternative considered*: a central item registry / constants file. Rejected — premature abstraction for five items; the codebase has no such pattern yet.

### 2. Campfire resting — add `restCooldown int` to `Model`

`Ground.HasFire` is already set per-cell in the local map. A new `r` key branch in `input.go` checks:
1. `m.mode == ModeLocal`
2. Player is standing on a cell where `m.currentLocal.Ground[px][py].HasFire == true`
3. `m.restCooldown == 0`

On success: `m.playerHP = min(m.playerHP+5, MaxPlayerHP)`, `m.restCooldown = 60`. The existing 500 ms `TickMsg` loop decrements `restCooldown` each tick (add one line to `Update`'s tick branch). The HUD already shows HP — no new UI element is strictly necessary, though a short status message ("Rested." flash) can be added as a stretch.

*Alternative considered*: event/message-based cooldown. Rejected — the model's tick already counts time; a simple int counter is idiomatic in this codebase.

### 3. Death screen — new `ScreenDeath ScreenMode` constant

Add `ScreenDeath` to the `ScreenMode` enum in `types.go`. Add a `deathKiller string` field to `Model`. On defeat (currently `return m, tea.Quit`), instead set `m.screenMode = ScreenDeath` and `m.deathKiller = cs.Enemy.Name`. Add a `renderDeathScreen` function in `render.go` (analogous to `renderCombatScreen`). Add key handling in `input.go` for `ScreenDeath`: `r` → reset model to `NewModel()`, `q` / `ctrl+c` → `tea.Quit`.

*Alternative considered*: a dedicated `ModeGameOver` in the `Mode` enum (world/local/dungeon tier). Rejected — `ScreenMode` is the correct abstraction (it already handles fullscreen overlays like combat and inventory, which are orthogonal to map tier).

### 4. Portrait fix — name-based lookup, backward-compatible

Add `enemyPortraitByName(name string) portrait` alongside the existing `enemyPortrait(rune)`. Map known enemy names:
- Humanoid: "Goblin", "Bandit", "Jungle Troll", "Frost Giant", "Stone Golem"
- Beast: "Cave Crustacean", "Cave Rat"
- Undead: "Sand Wraith", "Ice Wraith"
- Default: `fallbackPortrait`

Update `renderEnemyPanel` in `render.go` to call `enemyPortraitByName(cs.Enemy.Name)` instead of `enemyPortrait(enemyChar)`. Keep `enemyPortrait(rune)` in place for now (tests reference it) — update the tests for the new function in the same PR.

*Alternative considered*: modify `enemyPortrait` signature to accept both char and name. Rejected — overloading is unidiomatic in Go; two functions with clear names is cleaner.

### 5. Minimum enemy count — one-line clamp

Change `enemyCount := depth` to `enemyCount := max(3, depth)` in `dungeon.go`. The existing guard `if enemyCount > len(enemyPositions)` already caps it against available room positions, so no additional bounds logic is needed.

## Risks / Trade-offs

- **Campfire cooldown tick-counting**: The tick fires every 500 ms; 60 ticks ≈ 30 seconds real time. If the user changes tick frequency in future, the cooldown duration drifts. → Acceptable for now; document the coupling.
- **Death screen restart**: `NewModel()` reinitialises all state. If any global mutable state is added in future, the restart path may leave stale data. → Not a current concern; no global mutable state exists today.
- **Portrait name lookup**: Any new enemy whose name doesn't appear in the lookup silently falls back to the generic portrait. → This is the same behavior as today; it is explicitly the fallback intent.
- **Minimum 3 enemies on floor 1**: Players who descend immediately at low HP may find floor 1 suddenly harder. → This is the intended design improvement; early floors feeling empty was the bug.
- **Clothing armour values**: Tests in `game_test.go` check that clothing items have no bonuses (lines ~220). These tests must be updated. → Small test update, no risk.

## Migration Plan

No data migration — the game has no persistence. All changes take effect on the next binary build. No rollback strategy needed beyond reverting the commit.

## Open Questions

- Should campfire resting show a brief status flash ("You rest by the fire. +5 HP") in the HUD? — Suggested yes, but not blocking; can be added as a follow-on.
- Should the death screen show total enemies killed / floors reached? — Out of scope for this change (no kill counter exists); deferred.
