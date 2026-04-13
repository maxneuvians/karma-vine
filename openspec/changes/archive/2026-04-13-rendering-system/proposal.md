## Why

Without a rendering layer the game is invisible. The `View()` method on `Model` must translate the in-memory world and local maps into a styled terminal string every frame. Lipgloss provides the colour and layout primitives; this change wires them together into a playable display with a HUD status bar and optional night-mode dimming.

## What Changes

- Implement world-map rendering: iterate the visible chunk range around `worldPos`, draw each tile's `Char` with its Lipgloss foreground colour
- Implement local-map rendering: iterate the full 42×18 `LocalMap`, layer ground → object → animal → player `@`
- Implement a HUD status bar composed via `lipgloss.JoinVertical`: shows biome name, elevation value, world coords, and chunk coord
- Implement night mode: multiply each tile's RGB channels by `0.35` when `Model.nightMode` is `true`
- Player character `@` is always rendered in `#f0f6fc` bold regardless of night mode
- Replace the stub `View()` with the full implementation
- Use `tea.WindowSizeMsg` to update `viewportW` / `viewportH` dynamically

## Capabilities

### New Capabilities
- `rendering-system`: Lipgloss-based terminal rendering — world-map view, local-map view, HUD status bar, night-mode RGB dimming, viewport math, and player glyph

### Modified Capabilities
<!-- none — the stub View from project-scaffold is replaced, not spec-extended -->

## Impact

- Replaces stub `View()` in `model.go` / `render.go`
- Adds `nightMode bool` field to `Model` struct
- No new dependencies (Lipgloss already in `go.mod`)
