## 1. Model Fields

- [x] 1.1 Add `timeOfDay float64` to `Model` in `model.go` (initial value `0.25` — 6 AM)
- [x] 1.2 Add `timeScale int` to `Model` in `model.go` (initial value `1`)
- [x] 1.3 Remove `nightMode bool` from `Model`
- [x] 1.4 Update `NewModel()` to initialise `timeOfDay: 0.25` and `timeScale: 1`

## 2. Ground HasFire + Lit Grid

- [x] 2.1 Add `HasFire bool` field to `Ground` struct in `types.go`
- [x] 2.2 Add `hasFires bool` and `fireThreshold float64` fields to `biomeContent` in `local.go`
- [x] 2.3 Add campfire entries to biome tables: Plains (`♨`, `#ff8800`, threshold `0.90`), Beach (`♨`, `#ff8800`, threshold `0.88`), Mountains (`♨`, `#ff6600`, threshold `0.92`)
- [x] 2.4 In `GenerateLocalMap`, after placing objects, place fire cells using a third noise sample: if `fireNoise > content.fireThreshold`, set `lm.Ground[x][y].HasFire = true`
- [x] 2.5 Add `LitMap [LocalMapW][LocalMapH]bool` field to `LocalMap` in `types.go`
- [x] 2.6 Implement `buildLitMap(lm *LocalMap)` in `local.go` that sets `LitMap[x][y] = true` for all cells within Manhattan distance 4 of any fire cell
- [x] 2.7 Call `buildLitMap(lm)` at the end of `GenerateLocalMap`

## 3. Time Advancement

- [x] 3.1 In `model.go` `Update()`, in the `TickMsg` branch, advance `m.timeOfDay` by `float64(m.timeScale) / 60.0` (60 ticks = 30 s at 1×, giving a full cycle per 30 s)
- [x] 3.2 Wrap `timeOfDay` using `math.Mod(m.timeOfDay+delta, 1.0)` to keep it in `[0, 1)`

## 4. Input Handling

- [x] 4.1 Remove the `case "n": m.nightMode = !m.nightMode` block from `model.go` `Update()`
- [x] 4.2 Add `"]"` case to `handleKey` in `input.go`: advance `timeScale` through `{1, 2, 5, 10}`, clamp at 10
- [x] 4.3 Add `"["` case to `handleKey` in `input.go`: retreat `timeScale` through `{1, 2, 5, 10}`, clamp at 1

## 5. Rendering — dimFactor

- [x] 5.1 Add `dimFactor(timeOfDay float64) float64` function in `render.go` using `clamp(0.5*(1+math.Cos(2*math.Pi*timeOfDay))*0.85+0.15, 0.15, 1.0)`
- [x] 5.2 Replace `applyColor(hex string, night bool) string` signature with `applyColor(hex string, dim float64) string` in `render.go`; change the body to multiply each channel by `dim` instead of 0.35
- [x] 5.3 Update all call sites of `applyColor` in `render.go` to pass `dimFactor(m.timeOfDay)` instead of `m.nightMode`
- [x] 5.4 Remove `dimColor(hex string) string` helper (its logic is now inlined in `applyColor`)

## 6. Rendering — local illumination

- [x] 6.1 In `renderLocalMap`, compute `globalDim := dimFactor(m.timeOfDay)` once before the loop
- [x] 6.2 For each cell `(x, y)`, determine `cellDim`: if `lm.LitMap[x][y]` is true, use `1.0`; else use `globalDim`
- [x] 6.3 Pass `cellDim` to all `applyColor` calls within the cell rendering block
- [x] 6.4 When `Ground.HasFire` is true and no animal/player occupies the cell, render it as `♨` in color `#ff8800` (using `cellDim`)

## 7. HUD

- [x] 7.1 Add a `formatTime(timeOfDay float64) string` helper in `render.go` that converts `timeOfDay` to `HH:MM` (e.g., `0.5 → "12:00"`)
- [x] 7.2 Update `renderHUD` to include `formatTime(m.timeOfDay)` and `fmt.Sprintf("%d×", m.timeScale)` in the status bar
- [x] 7.3 Update the key bar hint in `renderKeyBar` to replace `n night` with `[/] speed` and show the current speed

## 8. Sidebar

- [x] 8.1 Update `renderSidebar` world-map section to remove any night-mode references
- [x] 8.2 Add fire glyph `♨` / "Campfire" entry to the local map legend in `renderSidebar`

## 9. Tests

- [x] 9.1 Update `TestApplyColor_NightDims` and `TestApplyColor_DayUnchanged` in `render_test.go` to use `float64` dim factor instead of bool
- [x] 9.2 Write `TestDimFactor_Noon` verifying `dimFactor(0.5) ≈ 1.0`
- [x] 9.3 Write `TestDimFactor_Midnight` verifying `dimFactor(0.0) ≈ 0.15`
- [x] 9.4 Write `TestFormatTime` covering `0.0 → "00:00"`, `0.5 → "12:00"`, `0.75 → "18:00"`
- [x] 9.5 Write `TestBuildLitMap` verifying cells within radius 4 of a fire are lit and cells at distance 5 are not
- [x] 9.6 Update any test that references `nightMode` to remove/replace it

## 10. Verification

- [x] 10.1 Run the binary, observe colors dimming and brightening over time on the world map
- [x] 10.2 Descend to a local map; confirm fire glyphs appear in biomes with fires (Plains, Beach, Mountains)
- [x] 10.3 Walk near a fire at midnight; confirm the surrounding radius is brighter than the rest of the map
- [x] 10.4 Press `]` repeatedly; confirm the day/night cycle visibly speeds up, HUD shows `2×`, `5×`, `10×`
- [x] 10.5 Confirm `n` key no longer toggles anything
