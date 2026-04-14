## 1. Types

- [x] 1.1 Add `Name string` field to `Object` struct in `internal/game/types.go`
- [x] 1.2 Add `Name string` field to `Animal` struct in `internal/game/types.go`

## 2. Local Map Generation — Name Population

- [x] 2.1 In `internal/game/local.go`, add names to all biome object entries in the biome table (tree → `"Tree"`, pine → `"Pine"`, cactus → `"Cactus"`, rock → `"Rock"`, shelter → `"Shelter"`, flower → `"Flower"`, and any other objects in the table)
- [x] 2.2 In `internal/game/local.go`, add names to all biome animal entries in the biome table: `d` → `"Deer"`, `r` → `"Rabbit"`, `s` → `"Snake"`, `l` → `"Lizard"`, `w` → `"Wolf"`, `b` → `"Bird"`, `g` → `"Goat"`, `e` → `"Elk"`, `B` → `"Bear"`, `c` → `"Crab"`, and any others
- [x] 2.3 In `internal/game/local.go`, set `Name: "Dungeon Entrance"` on the dungeon entrance object when it is injected

## 3. Dungeon Generation — Name Population

- [x] 3.1 In `internal/game/dungeon.go`, set `Name: "Staircase Up"` on the up-staircase object
- [x] 3.2 In `internal/game/dungeon.go`, set `Name: "Staircase Down"` on the down-staircase object
- [x] 3.3 In `internal/game/dungeon.go`, set `Name: "Torch"` on torch objects
- [x] 3.4 In `internal/game/dungeon.go`, set `Name: "Brazier"` on brazier objects

## 4. Sidebar Rendering

- [x] 4.1 In `internal/game/render.go`, remove the `localCharNames map[rune]string` declaration
- [x] 4.2 In `renderSidebar`, replace the `else` branch (currently only local mode) with an explicit three-way switch on `m.mode`: `ModeWorld`, `ModeLocal`, `ModeDungeon`
- [x] 4.3 Update the `ModeWorld` branch: when `m.mapMode != MapModeDefault`, replace the biome legend with a header showing the active overlay name and a brief colour-key hint (e.g., `"Temperature"` + `"blue=cold / red=hot"`)
- [x] 4.4 Update the `ModeLocal` branch: replace `localCharNames[e.char]` lookups with `obj.Name` and `a.Name`; ensure the dungeon entrance (`>`) shows as `"Dungeon Entrance"`
- [x] 4.5 Implement the `ModeDungeon` branch: show `"Dungeon"` header, `"Depth: N"` line, then scan `m.currentDungeon.Cells` for unique named objects and list them under a `"Contents"` sub-heading

## 5. Tests

- [x] 5.1 Add test to `internal/game/local_test.go`: assert no object on a generated local map has an empty `Name`
- [x] 5.2 Add test to `internal/game/local_test.go`: assert no animal on a generated local map has an empty `Name`
- [x] 5.3 Add test to `internal/game/local_test.go`: assert the dungeon entrance object has `Name == "Dungeon Entrance"`
- [x] 5.4 Add test to `internal/game/dungeon_test.go`: assert no object in any `DungeonCell` has an empty `Name`
- [x] 5.5 Add test to `internal/game/render_test.go`: `renderSidebar` with `ModeLocal` contains `"Dungeon Entrance"` when local map has an entrance
- [x] 5.6 Add test to `internal/game/render_test.go`: `renderSidebar` with `ModeDungeon` contains `"Dungeon"` and `"Depth: 1"`
- [x] 5.7 Add test to `internal/game/render_test.go`: `renderSidebar` with `ModeWorld` + `MapModeTemperature` contains `"Temperature"` and does not contain `"Deep Ocean"`
- [x] 5.8 Update any existing test fixtures that construct `Object{}` or `Animal{}` directly to include a `Name` field where required by the test assertions
