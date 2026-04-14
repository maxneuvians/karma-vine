## 1. Model & Types

- [x] 1.1 Add `MapMode int` type and `MapModeDefault`, `MapModeTemperature`, `MapModeElevation`, `MapModePolitical` constants to `internal/game/types.go`
- [x] 1.2 Add `mapMode MapMode`, `showMapPicker bool`, and `mapPickerCursor int` fields to `Model` in `internal/game/model.go`
- [x] 1.3 Verify `NewModel()` zero-values these fields correctly (no explicit init needed; zero value = `MapModeDefault` and `false`)

## 2. Tile Visual Override

- [x] 2.1 Add hex-lerp helper `lerpHex(a, b string, t float64) string` in `internal/game/render.go` that interpolates R, G, B components
- [x] 2.2 Implement `tileVisual(t Tile, mode MapMode) (ch rune, color string)` in `internal/game/render.go` with branches for all four modes
- [x] 2.3 Write table-driven unit tests for `tileVisual` covering: default pass-through, temperature min/mid/max, elevation min/max, political contour boundary vs. non-boundary
- [x] 2.4 Write unit test for `lerpHex` at t=0, t=1, and t=0.5

## 3. World Map Render Integration

- [x] 3.1 Update the world-map render loop in `internal/game/render.go` to call `tileVisual(tile, m.mapMode)` instead of reading `tile.Char`/`tile.Color` directly
- [x] 3.2 Ensure the dim-factor color scaling is applied to the color returned by `tileVisual` (not bypassed for non-default modes)
- [x] 3.3 Add/update render test asserting that `MapModeTemperature` changes the output color vs. `MapModeDefault` for the same tile

## 4. Map Picker Panel

- [x] 4.1 Add `renderMapPicker(m Model, height int) string` in `internal/game/render.go` that renders a 22-char-wide bordered list of the four mode names with cursor indicator
- [x] 4.2 Integrate `renderMapPicker` into `buildView`: when `showMapPicker == true`, compose it on the right side (reducing map width), mirroring the `showSidebar` path
- [x] 4.3 Write unit test for `renderMapPicker` asserting: all four names present, cursor row shows `>` prefix, non-cursor rows show ` ` prefix

## 5. Input Handling

- [x] 5.1 In `internal/game/input.go`, add early-exit block at the top of `handleKey`: when `showMapPicker == true`, route `up`/`down` to `mapPickerCursor` adjustment (clamped 0–3), `enter` to apply selection + close, `esc`/`m` to close without applying, and `return` early (skip movement/other handlers)
- [x] 5.2 Add `m` key case to `handleKey` for `ModeWorld` (when picker is closed): set `showMapPicker = true`, `showSidebar = false`, `mapPickerCursor = int(m.mapMode)`
- [x] 5.3 Update `?` key handler to set `showMapPicker = false` when opening the sidebar (mutual exclusion)
- [x] 5.4 Verify `enter` key still descends in `ModeWorld` when `showMapPicker == false` (no regression)
- [x] 5.5 Write input tests: `m` opens picker, `m` again closes picker, `enter` applies cursor, `esc` cancels, `up`/`down` move cursor, cursor clamps, movement keys are blocked while picker open, `?` closes picker

## 6. Key Bar Update

- [x] 6.1 Update `renderKeyBar` in `internal/game/render.go` to include `m map` in the world-mode hints string
- [x] 6.2 Verify key bar hint does not appear in `ModeLocal` hints string
