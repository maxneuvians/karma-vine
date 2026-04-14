## 1. Types

- [x] 1.1 Add `Pickupable bool` field to `Object` struct in `types.go`
- [x] 1.2 Define `Item struct { Char rune; Color string; Name string; Count int }` in `types.go`
- [x] 1.3 Define `Inventory struct { Items []Item }` in `types.go`
- [x] 1.4 Define `InventoryMaxSlots = 8` constant in `types.go`

## 2. Model

- [x] 2.1 Add `inventory Inventory`, `showInventory bool`, `inventoryCursor int` fields to `Model` in `model.go`
- [x] 2.2 Initialise `inventory` with empty `Items` slice in `NewModel()`

## 3. Dungeon Generation

- [x] 3.1 Set `Pickupable: true` on generated torches and braziers in `dungeon.go`

## 4. Local Map Generation

- [x] 4.1 Inject axe object into Forest, DenseForest, Jungle, and Taiga local maps in `local.go`
- [x] 4.2 Add `axeSpawnChance` constant controlling spawn probability; at most one axe per map
- [x] 4.3 Ensure axe is placed only on a passable cell not occupied by another object

## 5. Input Handling

- [x] 5.1 Handle `i` key in all modes to toggle `m.showInventory` in `input.go`
- [x] 5.2 Redirect `up`/`w` and `down`/`s` to move `m.inventoryCursor` when `m.showInventory == true`
- [x] 5.3 Handle `g` key in `ModeLocal` and `ModeDungeon` to pick up a pickupable object from the player's cell
- [x] 5.4 Stack items by name when picking up (increment `Count` on existing slot, add new slot otherwise)
- [x] 5.5 Reject pickup when inventory is at `InventoryMaxSlots`
- [x] 5.6 Handle `d` key in `ModeLocal` and `ModeDungeon` to drop the selected item at `m.inventoryCursor`
- [x] 5.7 Decrement dropped item count; remove slot when count reaches zero; clamp `inventoryCursor`
- [x] 5.8 Place dropped item as `Object{Pickupable: true}` on player's cell; find free adjacent cell if occupied
- [x] 5.9 Handle `u` key in `ModeLocal` and `ModeDungeon` to trigger item-use dispatch
- [x] 5.10 Implement use-dispatch table: `{"Axe", tree glyphs} → chop` (remove tree object, make cell passable)
- [x] 5.11 Chopping a tree does not consume the axe

## 6. Rendering

- [x] 6.1 Implement `renderInventoryPanel(m Model) string` showing title, item list with cursor highlighting, and "Empty" placeholder in `render.go`
- [x] 6.2 Overlay `renderInventoryPanel` in `buildView` when `m.showInventory == true`
- [x] 6.3 Add `Items: N/8` count to `renderHUD` output
- [x] 6.4 Add `i inv`, `g pick`, `d drop`, `u use` hints to `renderKeyBar` (show `g`/`d`/`u` in ModeLocal and ModeDungeon only)

## 7. Tests

- [x] 7.1 Test: new model has empty inventory
- [x] 7.2 Test: pickup removes object from local map and adds item to inventory
- [x] 7.3 Test: picking up same-named item stacks (increments Count)
- [x] 7.4 Test: pickup ignored for non-pickupable objects
- [x] 7.5 Test: pickup rejected when inventory is full
- [x] 7.6 Test: drop places pickupable object on player's cell in local map
- [x] 7.7 Test: drop decrements item count; slot removed when count reaches zero
- [x] 7.8 Test: drop ignored in ModeWorld
- [x] 7.9 Test: `i` toggles showInventory in all modes
- [x] 7.10 Test: cursor keys move inventoryCursor when showInventory is true
- [x] 7.11 Test: axe chops adjacent tree object (tree cell becomes nil)
- [x] 7.12 Test: axe use does not remove axe from inventory
- [x] 7.13 Test: use with no adjacent tree is a no-op
- [x] 7.14 Test: all unlit torches and braziers generated in dungeon have Pickupable true
- [x] 7.15 Test: forest local map may contain an axe with Pickupable true
- [x] 7.16 Test: inventory panel renders "Inventory" title and "Empty" when no items
- [x] 7.17 Test: HUD contains item count string
