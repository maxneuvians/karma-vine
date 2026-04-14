## Context

The game currently has three map tiers (world → local → dungeon) with objects on local and dungeon maps but no mechanism for the player to interact with them beyond collision. The `Object` struct already has `Name`, `Char`, `Color`, `Blocking`, and `Lit` fields. `LocalMap` stores objects in a `[LocalMapW][LocalMapH]*Object` grid; `DungeonLevel` stores them in `[DungeonW][DungeonH]DungeonCell`. Both are cached in `Model`, so any mutation (pickup/drop) persists as long as the cache lives. The BubbleTea model is a value type passed by copy, so all state changes are made on a local copy and returned.

## Goals / Non-Goals

**Goals:**
- Introduce an `Item` type and an `Inventory` held on `Model`
- Mark certain `Object` instances as `Pickupable: true` at map-generation time
- `g` picks up the item on the player's current cell; `d` drops the selected inventory item back onto the cell
- Dropped items re-materialise as `Object` entries in the live (cached) map, persisting across visits
- `i` toggles a small inventory overlay panel accessible from all three modes
- `u` triggers a context-sensitive use action (axe + adjacent tree → chop tree in `ModeLocal`; extensible)
- Inventory display shows item name, glyph, and count (stackable by name)
- Item interactions follow a simple dispatch table: `(itemName, adjacentObjectName) → action`
- Sidebar shows "Inventory: N items" summary line in all modes

**Non-Goals:**
- Item persistence across sessions (save/load not implemented)
- Item durability or degradation
- Crafting or combining items
- Equipment slots or stat bonuses
- NPC trade or economy
- More than one use-interaction type at launch (axe+tree only)
- Items on the world map (only local and dungeon)

## Decisions

### Decision 1: `Item` is a separate struct from `Object`

`Object` is a map-cell occupant with blocking/lighting semantics. `Item` is a carriable entity with a quantity and optional use-effect. They share glyph and name but have different lifecycles. Reusing `Object` would add nullable fields and blur the separation between "world entity" and "carried entity".

**Alternative considered**: Add a `Carriable *Item` pointer to `Object`. Rejected — couples map rendering to inventory logic; a nil check on every render path is fragile.

### Decision 2: `Inventory` as a fixed-capacity slice on `Model`

`Inventory` is `[]Item` with a configurable `MaxSlots` constant (e.g., 8). Items with the same `Name` stack (increment `Count`). Fits naturally in the BubbleTea value-copy model — no pointer indirection needed.

**Alternative considered**: A map keyed by item name. Rejected — map iteration order is non-deterministic (same bug we just fixed in the sidebar); a slice with explicit ordering is simpler.

### Decision 3: `Pickupable bool` field on `Object`

Marking objects at generation time is cleaner than a global lookup table by glyph. Generators in `local.go` and `dungeon.go` set this flag; the input handler checks it before adding to inventory.

**Alternative considered**: A `pickupableChars map[rune]bool` in the input handler. Rejected — same fragility as `localCharNames`; violates the "self-describing at creation" principle we established in the contextual-sidebar change.

### Decision 4: Pickup removes the `Object` from the map; drop restores it

When the player picks up an item, `lm.Objects[x][y]` (or the dungeon cell object) is set to `nil`. When dropped, a new `Object` is created from the `Item` data (with `Pickupable: true`) and placed in the current cell. Since both `LocalMap` and `DungeonLevel` are cached by pointer, this mutation is immediately visible on re-entry.

**Alternative considered**: Keep the object on the map and track "which objects have been picked up" separately. Rejected — more complex state to maintain; nil-on-pickup is the standard roguelike approach.

### Decision 5: Inventory cursor for selection

`Model.inventoryCursor int` tracks which item is selected. Arrow keys navigate when the panel is open. `d` drops the selected item; `u` uses it. The cursor clamps when items are added/removed.

### Decision 6: Use-interactions as a dispatch table

A `var useInteractions = map[useKey]useAction{}` where `useKey` is `struct{ itemName, targetName string }` and `useAction` is a function `func(m *Model, targetX, targetY int)`. This is easy to extend without touching the main input switch. For now, `{"Axe", "Tree"} → chopTree` and `{"Axe", "Pine"} → chopTree`.

The axe item itself is placed on Plains/Forest biome maps at generation time (one per map, on a passable floor cell), similar to the dungeon entrance injection.

### Decision 7: Inventory panel layout

A `renderInventoryPanel(m Model) string` returns a centred overlay (or right-docked panel) of fixed width (22 cols). It lists items with index, glyph, name, and count. When `m.showInventory` is true, `buildView` overlays it. This mirrors the map-picker overlay pattern.

## Risks / Trade-offs

[Dropped item on occupied cell] If the player's cell already has an object, dropping places on the nearest free adjacent cell (same `findSpawnPoint`-style spiral). → Add a `findDropCell(x, y int, lm *LocalMap) (int, int)` helper.

[Stack confusion] Two different items with the same name but different glyphs (unlikely but possible) would stack incorrectly. → Stack by `Name` only; keep glyph of first item in stack (acceptable for now).

[Axe placement in all biomes] The axe must only appear in biomes that have trees. If placed in Desert, use would never find a tree. → Only inject axe in `Forest`, `DenseForest`, `Jungle`, `Taiga`, `Plains` biomes.

[Dungeon pickup of equipped torches] Picking up the last lit torch could leave the player in darkness. → No protection needed at launch; it's a player choice.

[Test coverage] Pickup/drop mutations on cached map pointers need integration-style tests alongside unit tests. → Use `NewModel()` with pre-populated maps in tests, same as current input tests.

## Migration Plan

1. Add `Item`, `Inventory` types and `Pickupable` field to `types.go` (additive, no breakage).
2. Add inventory fields to `Model`, initialise in `NewModel()`.
3. Mark objects `Pickupable: true` in `local.go` (axe injection) and `dungeon.go` (unlit torch/brazier).
4. Implement `g`, `d`, `i`, `u` in `input.go`.
5. Implement `renderInventoryPanel` in `render.go`; extend `buildView` and `renderHUD`.
6. Add tests.

Rollback: all changes are additive until step 4; steps 4–6 can be reverted independently.

## Open Questions

- Should the axe be a rare find (1 per map) or abundant? → 1 per Forest/DenseForest/Jungle/Taiga map; absent on Plains (keeps it as a reward for exploring wooded areas).
- Should chopped trees leave a permanent stump `Object` (impassable) or clear ground? → Clear ground (passable floor cell), so the player gains usable path space.
- Should inventory persist when moving between local tiles? → Yes — inventory is on `Model`, not tied to a map, so it naturally persists.
