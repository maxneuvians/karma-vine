## Why

The current game has two depth layers (world map → local map), but no vertical exploration underground. Adding a dungeon/cave system gives players a meaningful reason to explore local maps, introduces risk/reward gameplay through deeper levels, and expands the game world into a third dimension.

## What Changes

- New `ModeDungeon` game mode, making the mode system three-tier: world → local → dungeon
- Procedurally generated dungeon levels using BSP (Binary Space Partitioning) rooms + corridor carving
- Dungeon-specific fog-of-war: only the area immediately around the player is visible unless torches or braziers illuminate nearby cells
- Dungeon entrance objects placed on local maps (e.g., a staircase down `>`); player presses `enter`/`>` to descend
- Staircase up `<` objects at the start of each dungeon level to ascend back
- Staircase down `>` objects generated per level to descend further (up to a random max depth of 5–10)
- Wall-mounted items (torches, sconces) placed on wall cells; braziers placed on floor cells — never on the same cell as the player
- Players can occupy floor cells containing items (braziers, dropped loot)
- Dungeon levels are generated once and cached per (worldX, worldY, depth) key
- Max depth randomised per dungeon entrance (between 5 and 10)

## Capabilities

### New Capabilities

- `dungeon-generation`: Procedural level generation — BSP room placement, corridor carving, stair placement, item (torch/brazier) seeding
- `dungeon-rendering`: Dungeon-specific render loop with fog-of-war, lit-cell radius from torches/braziers, and dungeon tile glyphs
- `dungeon-navigation`: Input handling for dungeon mode — movement, descend (`>`/`enter`), ascend (`<`/`esc`), blocked by walls

### Modified Capabilities

- `local-map-generation`: Add dungeon entrance object (`>` staircase) to generated local maps
- `input-navigation`: Extend `handleKey` and `applyDelta` to support `ModeDungeon` movement and stair transitions
- `rendering-system`: Extend `buildView` to dispatch to dungeon render path when `m.mode == ModeDungeon`

## Impact

- `internal/game/types.go`: new `ModeDungeon` constant; new `DungeonLevel`, `DungeonCell`, `DungeonObject` types
- `internal/game/model.go`: new `dungeonMap *DungeonLevel`, `dungeonDepth int`, `dungeonCache map[dungeonKey]*DungeonLevel` fields
- `internal/game/dungeon.go`: new file — generation logic
- `internal/game/render.go`: new `renderDungeonMap` function; `buildView` dispatch
- `internal/game/input.go`: dungeon movement, stair transitions
- `internal/game/local_map.go` (or equivalent): inject dungeon entrance object
- No new dependencies; all procedural generation uses `math/rand` seeded from world coordinates
