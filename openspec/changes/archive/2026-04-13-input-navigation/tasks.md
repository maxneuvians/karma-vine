## 1. Key Handler Structure

- [x] 1.1 Create `input.go` and implement `handleKey(msg tea.KeyMsg, m Model) (Model, tea.Cmd)` that switches on `msg.String()`
- [x] 1.2 Replace the stub key handling in `Update()` with a call to `handleKey`

## 2. Movement Delta

- [x] 2.1 Implement `applyDelta(dx, dy int, m Model) Model` in `input.go`
- [x] 2.2 In `ModeWorld`: `m.worldPos.X += dx`, `m.worldPos.Y += dy`, return updated model
- [x] 2.3 In `ModeLocal`: compute `newX = playerPos.X + dx`, `newY = playerPos.Y + dy`; bounds-check against `[0,41]×[0,17]`; collision-check against `localMap.Objects[newX][newY].Blocking`; update `playerPos` only if valid
- [x] 2.4 Add cases in `handleKey` for `"up"`, `"w"`, `"down"`, `"s"`, `"left"`, `"a"`, `"right"`, `"d"` calling `applyDelta`

## 3. Descend (Enter / >)

- [x] 3.1 Add `"enter"` and `">"` cases in `handleKey` (active only when `mode == ModeWorld`)
- [x] 3.2 Implement `findSpawnPoint(lm *LocalMap) LocalCoord` that scans outward from `{21, 9}` for the first non-blocking cell
- [x] 3.3 Set `m.localMap = LocalMapFor(m.worldPos.X, m.worldPos.Y, &m)`, `m.mode = ModeLocal`, `m.playerPos = findSpawnPoint(m.localMap)`

## 4. Ascend (Escape / <)

- [x] 4.1 Add `"esc"` and `"<"` cases in `handleKey` (active only when `mode == ModeLocal`)
- [x] 4.2 Set `m.mode = ModeWorld` — do NOT nil-out `m.localMap`

## 5. Quit

- [x] 5.1 Ensure `"q"` and `"ctrl+c"` cases return `tea.Quit` in all modes

## 6. Tests

- [x] 6.1 Write a unit test for `applyDelta` that verifies world movement increments `worldPos`
- [x] 6.2 Write a unit test that verifies local movement is blocked at `{0, 0}` when pressing `up`
- [x] 6.3 Write a unit test that verifies movement is blocked by a `Blocking` object

## 7. Verification

- [x] 7.1 Run the binary and confirm WASD and arrow keys move the player on the world map
- [x] 7.2 Press Enter to descend; confirm the local map appears and the player is placed near centre
- [x] 7.3 Press Escape to ascend; confirm the world map reappears at the same position
- [x] 7.4 Re-descend to the same tile; confirm the local map is identical to the first visit
- [x] 7.5 Press `q`; confirm the program exits and the terminal is restored
