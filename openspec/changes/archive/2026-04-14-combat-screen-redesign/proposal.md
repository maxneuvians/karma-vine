## Why

The current combat screen resolves all rounds instantly and displays the result in a flat list. There is no sense of drama or pacing — the player can't follow what happened. A round-by-round animated layout with distinct character portraits and a live combat log makes fights engaging and readable.

## What Changes

- Combat screen layout replaced with a **three-panel design**: left panel (hero portrait + stats), right panel (enemy portrait + stats), bottom panel (combat log)
- Combat resolution stays pre-computed (full `CombatState.Log` is already stored), but **playback is stepped**: one round's log lines are revealed every N seconds
- **Three playback speeds**: Slow (3 s/round), Normal (1 s/round), Fast (0.2 s/round); toggled with `[` and `]` keys during combat
- **Hero portrait**: multi-line ASCII art of the player character displayed in the left panel (reusing/extending the existing ragdoll art)
- **Enemy portrait**: single large glyph centred in a box in the right panel (uses `enemy.Template.Char` and `Color`)
- **HP progress bars** for each combatant update visually as the log advances round by round
- When playback reaches the final round, the result banner (Victory / Defeated) and dismiss hint are shown
- **BREAKING** (display only): the previous `renderCombatScreen` layout is replaced in full

## Capabilities

### New Capabilities
- `combat-playback`: Round-by-round stepped playback of pre-computed `CombatState.Log`; speed control (`[`/`]`); `CombatTickMsg` drives advancement; model fields `combatLogIndex int`, `combatSpeed int`

### Modified Capabilities
- `combat-rendering`: Three-panel layout replacing the current flat layout; hero and enemy portraits; per-panel HP bars; log reveals one round at a time; result banner shown after final round
- `input-navigation`: `[` and `]` keys change `combatSpeed` during `ScreenCombat`; all other non-dismiss keys remain suppressed
- `rendering-system`: No structural change; `renderCombatScreen` is updated in place

## Impact

- `internal/game/types.go` — `CombatTickMsg` type; `CombatSpeed` constants
- `internal/game/model.go` — `combatLogIndex int`, `combatSpeed int` fields; `CombatTickMsg` scheduling
- `internal/game/render.go` — `renderCombatScreen` fully rewritten; hero/enemy portrait functions
- `internal/game/input.go` — `[` / `]` speed controls in `ScreenCombat` block
- No new dependencies
