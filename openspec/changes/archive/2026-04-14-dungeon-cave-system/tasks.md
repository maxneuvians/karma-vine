## 1. Types & Model

- [x] 1.1 Add `ModeDungeon` constant to the `Mode` enum in `internal/game/types.go`
- [x] 1.2 Add `CellKind int` type with `CellWall`, `CellFloor` constants to `internal/game/types.go`
- [x] 1.3 Add `DungeonCell`, `DungeonLevel`, `DungeonMeta`, and `dungeonKey` types to `internal/game/types.go`; add `DungeonW = 80`, `DungeonH = 24` constants
- [x] 1.4 Add `dungeonCache map[dungeonKey]*DungeonLevel`, `dungeonMeta map[WorldCoord]DungeonMeta`, `currentDungeon *DungeonLevel`, `dungeonDepth int`, and `dungeonEntryPos LocalCoord` fields to `Model` in `internal/game/model.go`
- [x] 1.5 Initialise `dungeonCache` and `dungeonMeta` maps in `NewModel()` (parallel to `localCache`)

## 2. Dungeon Generation

- [x] 2.1 Create `internal/game/dungeon.go`; implement BSP helper that recursively splits a rectangle into leaf rects (min leaf size 6×6)
- [x] 2.2 Implement `GenerateDungeonLevel(globalSeed, wx, wy, depth, maxDepth int) *DungeonLevel` using BSP: fill walls, carve rooms (3×3 min interior), connect with L-corridors, place up/down staircases
- [x] 2.3 Add torch placement logic to `GenerateDungeonLevel`: seed torches on wall cells adjacent to floor cells (~1 per 5 rooms)
- [x] 2.4 Add brazier placement logic to `GenerateDungeonLevel`: seed braziers on floor cells inside rooms (~1 per 6 rooms)
- [x] 2.5 Implement `DungeonLevelFor(wx, wy, depth int, m *Model) *DungeonLevel` with cache lookup/store
- [x] 2.6 Implement `DungeonMetaFor(wx, wy int, m *Model) DungeonMeta` that looks up or creates a `DungeonMeta` with randomised `MaxDepth` in `[5,10]`
- [x] 2.7 Write table-driven unit tests in `internal/game/dungeon_test.go`: determinism (same seed → same output), at least one floor cell, up-stair on floor, down-stair absent on final level, down-stair present on non-final level, torch on wall cell
- [x] 2.8 Write unit test for `DungeonLevelFor` cache hit/miss behaviour

## 3. Local Map Dungeon Entrance

- [x] 3.1 Inject a dungeon entrance object (`Char: '>'`, `Color: "#e8c96a"`, `Blocking: false`) into `GenerateLocalMap` in `internal/game/local.go`; pick a passable cell not occupied by another object using the local seed
- [x] 3.2 Add/update unit tests in `internal/game/local_test.go`: exactly one `'>'` object per generated map, entrance is non-blocking, entrance does not overlap other objects

## 4. Input Handling

- [x] 4.1 Add `ModeDungeon` movement branch in `applyDelta` (`internal/game/input.go`): bounds clamp to `[0, DungeonW-1]×[0, DungeonH-1]`, block on `CellWall` or `Object.Blocking`
- [x] 4.2 In `handleKey`, add descent handler: when `ModeLocal` + `enter`/`>` + `playerPos` cell has `Object.Char == '>'`, call `DungeonMetaFor`, call `DungeonLevelFor(wx, wy, 1)`, set `mode = ModeDungeon`, `dungeonDepth = 1`, save `dungeonEntryPos`, place player at `level.UpStair`
- [x] 4.3 In `handleKey`, add dungeon descent handler: when `ModeDungeon` + `enter`/`>` + `playerPos == currentDungeon.DownStair && HasDownStair`, call `DungeonLevelFor` for next depth, update `currentDungeon` and `dungeonDepth`, place player at new level's `UpStair`
- [x] 4.4 In `handleKey`, add dungeon ascent handler: when `ModeDungeon` + `esc`/`<`, if `dungeonDepth > 1` go to previous level placing player at `DownStair`; if `dungeonDepth == 1` set `mode = ModeLocal`, restore `playerPos = dungeonEntryPos`
- [x] 4.5 Guard existing `ModeLocal` ascent (`esc` → `ModeWorld`) so it only fires when `mode == ModeLocal` (not `ModeDungeon`)
- [x] 4.6 Guard existing `ModeLocal` descent (`enter` → `ModeLocal`) so it only fires when player is NOT on a dungeon entrance cell
- [x] 4.7 Write input unit tests: descend from local map opens dungeon, ascend from depth 1 returns to local map, ascend from depth > 1 goes to previous level, descend to next level increments depth, `esc` in dungeon never sets `ModeWorld`, movement blocked by wall, movement blocked by torch (wall object), movement succeeds over brazier

## 5. Rendering

- [x] 5.1 Implement `computeDungeonVisibility(m Model) map[LocalCoord]bool` in `internal/game/render.go`: player Chebyshev radius 6 + torch/brazier Chebyshev radius 4
- [x] 5.2 Implement `renderDungeonMap(m Model) string`: iterate `DungeonH` rows × `DungeonW` cols clipped to viewport, apply visibility mask, render wall/floor/object/player glyphs with correct colors
- [x] 5.3 Extend `buildView` in `internal/game/render.go` with `else if m.mode == ModeDungeon` branch that calls `renderDungeonMap` and composes with HUD
- [x] 5.4 Extend `renderHUD` to show `Dungeon`, `Depth: N`, and `(wx, wy)` when `m.mode == ModeDungeon`
- [x] 5.5 Extend `renderKeyBar` to include `< up`, `> down`, `esc exit` hints when `m.mode == ModeDungeon`
- [x] 5.6 Write render unit tests: `buildView` with `ModeDungeon` contains `#`/`.`/`@` glyphs, HUD shows depth, key bar shows dungeon hints only in `ModeDungeon`, cells outside visibility radius render as blank, torch illuminates nearby cells
