## Why

The info sidebar (`?`) shows the same static biome legend or a partial local-map legend regardless of what the player is actually looking at. Animals and objects without entries in the hardcoded `localCharNames` map fall back to their raw glyph character; dungeon entrances show as `>` with no label; and the sidebar shows local-map data even when the player is inside a dungeon. The panel should reflect the current screen at all times.

## What Changes

- Add a `Name string` field to the `Object` and `Animal` structs so every entity is self-describing at creation time — eliminates the fragile `localCharNames` char→string lookup
- Populate `Name` for every object and animal created in `local.go` and `dungeon.go` (dungeon entrance, staircase up/down, torch, brazier, all biome animals and objects)
- Replace sidebar's `localCharNames` fallback with direct `obj.Name` / `animal.Name` reads
- Add a `ModeDungeon` sidebar branch to `renderSidebar` showing: depth, level dimensions, objects visible on the level (torch, brazier, staircases), and a "? close" hint
- Extend the world-map sidebar to show a brief note when a non-default map mode is active (e.g., "Temperature overlay" header instead of the biome colour legend, which is meaningless in that mode)

## Capabilities

### New Capabilities

- `sidebar-contextual-content`: Rules for what the sidebar displays per mode — world (with map-mode awareness), local (with fully named objects/animals), and dungeon (depth, contents)

### Modified Capabilities

- `rendering-system`: `renderSidebar` gains a `ModeDungeon` branch; world-mode sidebar adapts to active `MapMode`
- `local-map-generation`: `Object` and `Animal` structs gain a `Name` field; all creation sites populate it
- `dungeon-generation`: `DungeonCell.Object` entries (torch, brazier, staircases) populated with `Name` field

## Impact

- `internal/game/types.go`: add `Name string` to `Object` and `Animal`
- `internal/game/local.go`: populate `Name` for all objects and animals at generation time
- `internal/game/dungeon.go`: populate `Name` for all dungeon objects (torch, brazier, `<` up, `>` down)
- `internal/game/render.go`: remove `localCharNames` map; update `renderSidebar` for all three modes; add dungeon sidebar section
- No new dependencies
