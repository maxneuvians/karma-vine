## Context

BubbleTea delivers key events as `tea.KeyMsg`. Matching is done on `msg.String()` (returns the key name) or `msg.Type`. The brief specifies arrow keys and WASD, Enter/`>` to descend, Escape/`<` to ascend, and quit. Collision in local mode is against `Object.Blocking` on the target cell.

## Goals / Non-Goals

**Goals:**
- A single `handleKey(msg tea.KeyMsg, m Model) (Model, tea.Cmd)` function covers all bindings to keep `Update` readable
- World movement is unbounded; `worldPos` can be any integer coordinate
- Local movement is bounded to `[0, 41] × [0, 17]` and blocked by `Object.Blocking == true`
- Descend loads the local map via `LocalMapFor` and positions player at `{21, 9}` (centre of 42×18)
- Ascend preserves `localMap` pointer (stays in `localCache`) and sets `mode = ModeWorld`

**Non-Goals:**
- Diagonal movement (brief does not specify it for player input)
- Sprint / run mode
- Menu or pause screen

## Decisions

**Normalise WASD and arrows to a `(dx, dy)` delta** — A single `applyDelta(dx, dy int, m Model) Model` function handles both `ModeWorld` and `ModeLocal` moves cleanly. Alternative: separate branches for each key — leads to duplicated bounds/collision logic.

**Descend places player at map centre** — The brief does not specify a spawn point, so the centre `{21, 9}` of the 42×18 grid is the most sensible default. Alternative: top-left `{0, 0}` — risks spawning inside a blocking object.

**Ascend does not nil-out `localMap`** — The pointer stays valid; it's the same pointer stored in `localCache`. This means the rendering change will always find a `localMap` until the player explicitly descends elsewhere. Alternative: nil-out on ascend — safe but requires an extra nil-guard in the renderer.

## Risks / Trade-offs

- **World position integer overflow** → At go's `int` bounds (±2^62 on 64-bit) this would require moving ~9 quintillion tiles. Not a practical concern.
- **Player spawns on a blocking object at centre** → Rare but possible. Mitigation: scan outwards from centre for the first non-blocking cell at descend time.
