## Why

The player currently has no way to collect or carry objects found in the world, making exploration purely observational. An inventory system introduces meaningful player agency — picking up items, carrying them across modes, dropping them, and using them to interact with the environment in ways that would otherwise be impossible (e.g., chopping trees with an axe to create paths).

## What Changes

- New `Item` type representing carriable objects (name, glyph, colour, weight/count)
- New `Inventory` structure attached to `Model` holding up to a fixed number of items
- Certain objects placed on local and dungeon maps become pick-up-able items (e.g., axe on Plains/Forest, torch in dungeons)
- Player can press `g` (get/pick up) to pick up an item from their current cell on local and dungeon maps
- Player can press `d` (drop) to drop a selected item back onto the current cell on local and dungeon maps
- Dropped items persist in `LocalMap.Objects` / `DungeonLevel.Cells` for the lifetime of the cached map
- Inventory is viewable via a new `i` key which overlays a small inventory panel anywhere (world, local, dungeon)
- Item interactions: when the player is in `ModeLocal` adjacent to a blocking tree object and holds an axe, pressing `u` (use) removes the tree and places a stump ground tile (passable), creating a traversable path
- The sidebar shows a brief inventory count ("Inventory: N items") in all modes

## Capabilities

### New Capabilities

- `inventory-system`: Core inventory data model, pickup/drop mechanics, persistence in cached maps, and inventory display panel
- `item-interactions`: Context-sensitive item use (axe → chop tree; extensible pattern for future interactions)

### Modified Capabilities

- `input-navigation`: New key handlers — `g` (pick up), `d` (drop), `i` (open/close inventory panel), `u` (use item)
- `rendering-system`: Inventory overlay panel; brief inventory count in HUD or sidebar
- `local-map-generation`: Certain biome objects become pick-up-able (axe on forest/plains tiles); dropped items are re-rendered as objects on the map
- `dungeon-generation`: Torch items become pick-up-able from dungeon walls/floors (unlit only); dropped items persist in the dungeon level cache

## Impact

- `internal/game/types.go`: new `Item` struct; new `Inventory` struct; add `Pickupable bool` field to `Object`
- `internal/game/model.go`: add `inventory Inventory` and `showInventory bool` fields; add `inventoryCursor int` for selection
- `internal/game/local.go`: mark certain biome objects as `Pickupable: true`
- `internal/game/dungeon.go`: mark unlit torches and braziers as `Pickupable: true`
- `internal/game/input.go`: handlers for `g`, `d`, `i`, `u`
- `internal/game/render.go`: `renderInventoryPanel`; extend `buildView` and `renderHUD`
- No new external dependencies
