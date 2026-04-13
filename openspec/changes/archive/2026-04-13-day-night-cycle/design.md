## Context

The game engine runs a BubbleTea tick loop that fires every 500 ms. Color rendering is handled in `render.go` via `applyColor(hex, nightMode bool)`, which currently applies a fixed 0.35 multiplier. A manual `nightMode` toggle (`n` key) controls this. There is no time model; day and night are a static user choice rather than a living world state.

The local map is 160×48 cells. Each cell has a `Ground` struct (char, color, passable) and an optional `*Object`. The rendering pipeline layers ground → object → animal → player per cell.

## Goals / Non-Goals

**Goals:**
- A continuous `timeOfDay float64` in `[0, 1)` advances each tick; 0 = midnight, 0.25 = 6 AM, 0.5 = noon, 0.75 = 6 PM
- One full cycle (0 → 1) takes 30 real seconds at 1× speed; `timeScale` multiplies tick advancement (max 10×, so fastest cycle = 3 s)
- A smooth `dimFactor float64` in `[0, 1]` is derived from `timeOfDay`; 1.0 at noon, ~0.15 at midnight
- `dimFactor` replaces the boolean `nightMode` in all color calculations
- On the local map, cells where `Ground.HasFire == true` illuminate a configurable radius (default 4 cells); cells inside the radius receive `dimFactor = 1.0` regardless of global dim
- Fire cells are generated deterministically in biome tables during `GenerateLocalMap`
- `[` / `]` keys cycle through time speeds: 1×, 2×, 5×, 10×
- HUD shows a formatted 24-hour clock derived from `timeOfDay`
- Remove `nightMode bool` and the `n` key binding

**Non-Goals:**
- Weather or seasonal effects
- Player-placeable fires
- Smooth blending between illuminated and dark cells (hard cutoff at radius)
- Per-biome sunrise/sunset times

## Decisions

**Single float `timeOfDay` in `[0, 1)` advanced per tick** — The tick fires every 500 ms. At 1× speed, one full day = 30 s = 60 ticks. Each tick advances `timeOfDay` by `1/(60 * (1/timeScale))`. At 10× speed, 6 ticks per second × 10 = advancement of `10/60` per second, so a cycle in 3 s. Alternative: integer millisecond clock — unnecessarily complex for a game without subsecond time precision needs.

**Cosine dim curve** — `dimFactor = clamp(0.5*(1 + cos(2π * timeOfDay)) * 0.85 + 0.15, 0.15, 1.0)` produces a smooth sunrise/sunset and a gentle floor at midnight (0.15) so the world is never pitch black. Alternative: linear ramp — harsher transitions, less atmospheric.

**Per-cell illumination on local map only** — The world map is abstract (infinite tile grid); fire illumination is a close-up feature. On the world map, only the global `dimFactor` applies. On the local map, a second pass computes per-cell `effectiveDim = max(globalDim, fireInfluence(x, y))` where `fireInfluence` is 1.0 within radius 4 of any fire cell, 0 otherwise. Alternative: distance-weighted falloff — more realistic but adds a per-cell distance loop; hard cutoff is cheaper and sufficient for the aesthetic.

**`HasFire bool` on `Ground` rather than a new `FireObject`** — Fires are always floor-level and don't block movement. Bolting onto `Ground` avoids a new type and a third render layer. The fire glyph is rendered as part of the ground layer when `HasFire` is true.

**Discrete `timeScale` steps: 1, 2, 5, 10** — Avoids floating-point UX. `[` decreases, `]` increases, wrapping at limits. Shown in HUD as `1×`, `2×`, etc.

## Risks / Trade-offs

- **Local map illumination pass cost** — Each frame iterates all visible cells (mapW × mapH ≈ 2,000–6,000) and for each checks fire cells. A naive O(visible × fires) loop could be slow if there are many fires. Mitigation: precompute a `lit [LocalMapW][LocalMapH]bool` array once per tick (not per frame) whenever `timeOfDay` or animal positions change — since fires are static, recompute only when `timeOfDay` crosses a threshold (e.g., every 0.01 change). Actually fires are static so the `lit` grid only needs building once after GenerateLocalMap and doesn't change. Precompute at descend time.
- **Removing `nightMode`** → existing tests that reference `nightMode` will fail until updated.
- **Cosine curve passes through twilight quickly** — dawn/dusk each last ~3 s at 1× speed, which may feel rushed. Mitigation: the 10× speed cap is mainly for testing; the default 1× feel can be tuned by adjusting the curve exponent later.
