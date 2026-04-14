## 1. Model

- [ ] 1.1 Add `paused bool` field to `Model` struct in `model.go`
- [ ] 1.2 Verify `NewModel()` leaves `paused` at `false` (zero value — no explicit init needed)

## 2. Tick Suppression

- [ ] 2.1 In `Update`'s `TickMsg` handler, add early-return guard: if `m.paused == true`, skip `timeOfDay` advancement and `moveAnimals`, then return `m, tickCmd()`

## 3. Input — Space Key

- [ ] 3.1 Add `case " ":` branch in `handleKey` (before mode-specific branches) that toggles `m.paused` and returns immediately

## 4. Input — Movement Guard

- [ ] 4.1 In `handleKey`, add `m.paused` guard to the directional key branches so they are no-ops when `m.paused == true` and `m.screenMode == ScreenNormal`
- [ ] 4.2 In `handleMouseClick`, add `m.paused` guard alongside the existing sidebar/inventory open checks so left-click movement is suppressed when paused

## 5. Rendering — Pause Indicator

- [ ] 5.1 In `render.go`, locate the HUD status bar construction (used in world, local, and dungeon render paths)
- [ ] 5.2 Append `[PAUSED]` to the HUD string when `m.paused == true`

## 6. Tests

- [ ] 6.1 Add `TestPause_SpaceTogglesPaused` in `input_test.go`: press space → `paused == true`; press space again → `paused == false`
- [ ] 6.2 Add `TestPause_TickMsg_NoTimeAdvance` in `game_test.go`: paused model receives `TickMsg`, assert `timeOfDay` unchanged and cmd is non-nil
- [ ] 6.3 Add `TestPause_TickMsg_NoAnimalMovement` in `game_test.go` or `animals_test.go`: paused model in `ModeLocal`, receive `TickMsg`, assert animal positions unchanged
- [ ] 6.4 Add `TestPause_MovementBlocked` in `input_test.go`: paused model, press directional key in `ScreenNormal`, assert `playerPos` unchanged
- [ ] 6.5 Add `TestPause_InventoryCursorUnaffected` in `input_test.go`: paused model in `ScreenInventory`, press `down`, assert `inventoryCursor` increments
- [ ] 6.6 Add `TestPause_HUDContainsPausedLabel` in `render_test.go`: paused model, assert rendered output contains `[PAUSED]`
- [ ] 6.7 Add `TestPause_HUDNoPausedLabelWhenUnpaused` in `render_test.go`: unpaused model, assert rendered output does NOT contain `[PAUSED]`
