## 1. Setup

- [x] 1.1 Add `nightMode bool` field to `Model` struct in `model.go`
- [x] 1.2 Create `render.go` to house all rendering helpers; move/replace the stub `View()` to call the new functions

## 2. Night Mode Helper

- [x] 2.1 Implement `dimColor(hex string) string` that parses a `#rrggbb` string, multiplies each channel by `0.35`, and returns a new `#rrggbb` string
- [x] 2.2 Implement `applyColor(hex string, night bool) string` that calls `dimColor` when `night` is true, otherwise returns `hex` unchanged
- [x] 2.3 Write unit tests for `dimColor` covering the Forest green example `#2d7a1f` → `#0f2a0a`

## 3. World Map Renderer

- [x] 3.1 Implement `renderWorldMap(m Model) string` that iterates `viewportW × (viewportH-1)` screen cells, computes world coordinates with the brief's viewport math, calls `TileAt`, and styles each `Char` with Lipgloss foreground
- [x] 3.2 Apply `applyColor` to each tile colour
- [x] 3.3 Assemble rows with `strings.Builder` and join with `"\n"`

## 4. Local Map Renderer

- [x] 4.1 Implement `renderLocalMap(m Model) string` that iterates the 42×18 `LocalMap`
- [x] 4.2 Apply layering logic: animal → object → ground, with player `@` overriding at `playerPos`
- [x] 4.3 Style player `@` with bold `#f0f6fc` regardless of `nightMode`
- [x] 4.4 Style all other glyphs with `applyColor` applied to their colour field

## 5. HUD Status Bar

- [x] 5.1 Implement `renderHUD(m Model) string` using `lipgloss.NewStyle()` to create a single-row bar
- [x] 5.2 Include biome name string (convert `Biome` constant to human-readable via a `biomeName(b Biome) string` helper), elevation formatted to 2 d.p., world coords `(x, y)`, chunk coords `chunk (cx, cy)`
- [x] 5.3 Compose map view and HUD with `lipgloss.JoinVertical(lipgloss.Left, mapView, hud)`

## 6. View and Resize Wiring

- [x] 6.1 Update `View()` on `Model` to call `renderWorldMap` or `renderLocalMap` based on `m.mode`, then compose with `renderHUD`
- [x] 6.2 Handle `tea.WindowSizeMsg` in `Update()`: set `m.viewportW = msg.Width`, `m.viewportH = msg.Height`

## 7. Verification

- [x] 7.1 Run the binary and confirm the world map renders with coloured biome tiles
- [x] 7.2 Descend into a local tile and confirm the local map renders without blank rows
- [x] 7.3 Resize the terminal window and confirm the display reflows without panic
- [x] 7.4 Toggle night mode (add temporary `n` key binding for testing) and confirm visible dimming
