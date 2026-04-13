## 1. Tick Infrastructure

- [x] 1.1 Define `TickMsg struct{}` in `types.go`
- [x] 1.2 Update `Init()` to return `tea.Every(500*time.Millisecond, func(t time.Time) tea.Msg { return TickMsg{} })`
- [x] 1.3 Add a `case TickMsg:` branch to `Update()` that calls the animal move function and returns the re-scheduled tick command

## 2. Movement Logic

- [x] 2.1 Implement `moveAnimals(m *Model)` in `animals.go` that iterates `m.localMap.Animals`
- [x] 2.2 For each animal, check `Flee` and Manhattan distance to `m.playerPos`
- [x] 2.3 Implement `randomStep() (int, int)` returning one of the 8 direction pairs at random using `math/rand`
- [x] 2.4 Implement `fleeStep(ax, ay, px, py int) (int, int)` that returns the direction from the 8 candidates maximising Manhattan distance from `(px, py)`
- [x] 2.5 Apply bounds check: clamp candidate `(newX, newY)` to `[0,41]×[0,17]`; skip if resulting cell has a blocking object in `m.localMap.Objects[newX][newY]`
- [x] 2.6 Update `animal.X` and `animal.Y` in-place

## 3. Guard: Only Tick in Local Mode

- [x] 3.1 In the `TickMsg` handler, skip `moveAnimals` if `m.mode != ModeLocal` or `m.localMap == nil`
- [x] 3.2 Still return the re-schedule command so ticks continue when the player descends again

## 4. Tests

- [x] 4.1 Write a unit test that constructs a `LocalMap` with one non-flee animal, dispatches 10 ticks, and asserts the position changes at least once
- [x] 4.2 Write a unit test that places a flee animal adjacent to `playerPos` and asserts the post-tick position is further away
- [x] 4.3 Write a unit test that places an animal at `{0, 0}` with all moves blocked or out-of-bounds and asserts it stays put

## 5. Verification

- [x] 5.1 Run the binary, descend into a local map, and observe animals moving every ~500 ms
- [x] 5.2 Walk toward an animal with `Flee: true` and observe it moving away when within 3 tiles
- [x] 5.3 Ascend and re-descend to the same tile; confirm animals are not reset to initial positions
