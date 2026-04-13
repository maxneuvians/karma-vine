## Context

BubbleTea calls `View()` on every model update and expects a single string containing the full terminal frame. Lipgloss styles are applied inline to each character. Performance matters: the render loop runs at interactive frame rates, so it must avoid allocations in the hot path. The viewport dimensions come from `tea.WindowSizeMsg` and are stored on the model so `View()` is a pure function of model state.

## Goals / Non-Goals

**Goals:**
- World-map render uses the viewport math from the brief: `worldX = playerWorldX + (screenX - viewportWidth/2)`
- Local-map render layers ground → object → animal → player; only the top-most non-nil glyph is drawn per cell
- HUD is built with `lipgloss.JoinVertical` and shows: biome name, elevation (2 d.p.), world coords `(x, y)`, chunk `(cx, cy)`
- Night mode applies a `0.35` multiplier to R, G, B channels parsed from each `#rrggbb` colour string
- Player `@` is always `#f0f6fc` bold with no night-mode dimming
- Window resize is handled cleanly with no layout glitch

**Non-Goals:**
- Wish SSH multiplayer rendering
- Day/night cycle animation (night mode is a static toggle for v1)
- Particle effects or animation beyond animal movement glyphs

## Decisions

**`strings.Builder` for the frame** — Building each row into a `strings.Builder` then joining rows avoids quadratic string concatenation. Alternative: `[]string` rows joined with `"\n"` — acceptable but Builder benchmarks slightly better.

**Night mode as a per-render pass** — Rather than storing pre-dimmed colours, colours are dimmed at render time from the canonical `#rrggbb` value. This keeps tile data clean and allows the toggle to take effect immediately. Trade-off: parsing 6 hex digits per tile per frame (~756 tiles) — negligible at 60 fps.

**`lipgloss.Color` per-tile** — Each tile creates a `lipgloss.NewStyle().Foreground(lipgloss.Color(hex))` inline. Lipgloss styles are lightweight value types so this is safe. Alternative: global style map keyed by colour string — premature optimisation.

**HUD below map, fixed height 1 row** — Simple and predictable. The map occupies `viewportH - 1` rows; the last row is the status bar. Alternative: overlay with a border box — harder to align and wastes rows.

## Risks / Trade-offs

- **Night mode hex parsing is fragile for non-`#rrggbb` strings** → All colours in the biome tables are strictly `#rrggbb`, so this is safe. Mitigation: add a `sanitiseColor` helper that returns the original string unchanged if it doesn't match.
- **World render calls `TileAt` per visible cell every frame** → `TileAt` uses the chunk cache, so it is O(1) after warm-up. Cold frames (new chunks) may stutter briefly — acceptable.
