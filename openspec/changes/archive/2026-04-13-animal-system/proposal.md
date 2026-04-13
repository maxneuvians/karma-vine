## Why

Static local maps feel lifeless. The brief specifies that animals move every 500 ms and those with `Flee: true` move away from the player when within 3 tiles. This change adds the tick loop and movement logic that brings local maps to life, without touching the generation or rendering code.

## What Changes

- Define a `TickMsg` type and register a `tea.Every(500ms, TickMsg{})` command in `Init()`
- Handle `TickMsg` in `Update()`: for each animal in `localMap.Animals`, compute a move using random-direction walk or flee-vector logic
- Flee logic: if `animal.Flee == true` and Manhattan distance to `playerPos` ≤ 3, move one step in the direction that maximises distance; otherwise random step
- Bounds check: animals stay within `[0, 41] × [0, 17]`; blocked moves are skipped (no wrap)
- Animals are preserved in `localCache` so they resume their positions when the player revisits a tile

## Capabilities

### New Capabilities
- `animal-system`: 500 ms tick loop, random-walk movement, flee-from-player behaviour, bounds checking, and cache-persistent animal state

### Modified Capabilities
<!-- none -->

## Impact

- Modifies `Init()` to return the first tick command
- Modifies `Update()` to handle `TickMsg`
- No new dependencies; no changes to type definitions
