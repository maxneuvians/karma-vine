## Context

The world map renders each tile using `tile.Char` and `tile.Color`, driven by the biome classification in `renderWorldMap`. Tiles already carry `Elevation float64`, `Temperature float64`, and `Moisture float64`, but nothing exposes these values visually. The model has a `showSidebar bool` toggled by `?` as the only existing overlay pattern. There is no sub-navigation concept yet — all key events go directly to movement or global commands.

## Goals / Non-Goals

**Goals:**
- Add four map modes: `Default` (existing biome view), `Temperature`, `Elevation`, `Political` (contour lines).
- Introduce a mode-picker overlay that opens/closes with `m`, navigated with `up`/`down`, confirmed with `enter`, cancelled with `esc`.
- Establish a reusable sub-navigation pattern: while any picker overlay is open, arrow keys and enter are consumed by the overlay rather than the player movement handler.
- Keep tile rendering changes isolated to the world map path (local map is unaffected).

**Non-Goals:**
- Animated transitions between modes.
- Per-mode legend or color-scale HUD widgets.
- Applying map modes to the local map view.
- Persisting the selected mode across sessions.

## Decisions

### D1: `MapMode` as a named int type on the Model

Add `MapMode int` to `types.go` with constants `MapModeDefault`, `MapModeTemperature`, `MapModeElevation`, `MapModePolitical`. Store `mapMode MapMode` on `Model`. This mirrors the existing `Mode` / `worldZoom` pattern — a single scalar field drives a branching render path — keeping the model flat and easy to test.

**Alternative considered:** a `[]bool` feature-flag slice. Rejected: modes are mutually exclusive, a scalar is cleaner.

### D2: Per-mode tile override function

Add `tileVisual(t Tile, mode MapMode) (ch rune, color string)` in `render.go`. The world-map render loop calls this instead of reading `t.Char`/`t.Color` directly when `mode != MapModeDefault`. Each mode maps a float value (elevation, temperature) to a color gradient using linear interpolation between two hex colors, or derives a contour character for `Political` mode.

- **Temperature**: cold (`#4488ff`) → hot (`#ff4422`), use `·` as char.
- **Elevation**: low (`#1a6fa8`) → high (`#f0f6fc`), use `·` as char.
- **Political**: draw `+` at contour boundaries (integer step of `elevation*10`), else `·`.

**Alternative considered:** lookup tables (discrete buckets). Rejected: a linear lerp on existing float values requires no new tile fields and is trivially testable.

### D3: Sub-navigation via `showMapPicker bool` + `mapPickerCursor int`

Add `showMapPicker bool` and `mapPickerCursor int` to `Model`. In `handleKey`, check `showMapPicker` first and route `up`/`down`/`enter`/`esc`/`m` to picker logic; other keys fall through normally. This avoids a separate `Mode` variant (which would require updating every mode switch) and keeps the overlay as a thin boolean guard.

**Alternative considered:** a new `ModeMapPicker` variant of the `Mode` enum. Rejected: it would shadow `ModeWorld`/`ModeLocal` and break the ascend/descend logic; an overlay bool is less invasive.

### D4: Picker renders as a narrow floating panel, right-aligned

`renderMapPicker(m Model, height int) string` draws a bordered list of mode names with the cursor highlighted, positioned at the right edge of the viewport (similar width to the sidebar). `buildView` composites it over the map by splitting the viewport width if `showMapPicker` is true, exactly as it does for `showSidebar`. This reuses the existing sidebar composition path without duplicating layout logic.

**Alternative considered:** an inline HUD row at the bottom. Rejected: the bottom row is already the key bar; vertical space is limited.

## Risks / Trade-offs

- **Simultaneous sidebar + picker open**: Both `showSidebar` and `showMapPicker` could be true at once, eating too much horizontal space on narrow terminals. Mitigation: close the sidebar when opening the picker and vice versa (mutual exclusion in the key handler).
- **Political contour legibility at high zoom**: At `worldZoom > 1`, the contour char `+` may be spaced too coarsely to read as a line. Mitigation: acceptable for now; a future zoom-aware contour algorithm is a non-goal.
- **Color lerp on non-color-safe terminals**: Hex interpolation works on true-color terminals; 256-color terminals will quantize. Mitigation: Lipgloss handles degradation gracefully; no action needed.
- **Test coverage**: `tileVisual` and `renderMapPicker` must hit ≥90% branch coverage. Mitigation: table-driven tests for each mode boundary (min, mid, max float values).

## Migration Plan

No data migration or deployment steps. The feature is purely additive: `mapMode` defaults to `MapModeDefault` (zero value) and `showMapPicker` defaults to `false`, so existing behavior is unchanged unless `m` is pressed.

## Open Questions

- Should the `m` key binding work in `ModeLocal` too (as a no-op, or disabled)? Currently proposed as world-only; the key bar hint will omit it in local mode.
- Should there be a zoom-aware contour algorithm for `Political` mode, or is the current simple approach sufficient for the first version?
