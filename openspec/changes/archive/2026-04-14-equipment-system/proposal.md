## Why

The ragdoll panel in the fullscreen inventory already renders six named body slots, but they are always empty — there is no way to equip, swap, or unequip items. This change makes the ragdoll functional: players can assign wearable items (weapons, tools, clothing) to body slots, see what they have equipped, and start the game with a basic default outfit.

## What Changes

- `BodySlot` type defined with constants for the six existing slots: `Head`, `Chest`, `LeftHand`, `RightHand`, `Legs`, `Feet`
- `Item` struct gains a `Slots []BodySlot` field declaring which slots the item is compatible with; empty means not equippable
- `Inventory` gains an `Equipped [NumBodySlots]Item` array (indexed by `BodySlot`) holding one item per slot
- `e` key in `ScreenInventory` equips the cursor item from the inventory list to its best available slot; if the slot is occupied the old item swaps back into the inventory list
- Ragdoll panel becomes navigable: `Tab` toggles focus between the left (inventory list) and right (ragdoll slots) column; arrow keys in the ragdoll column move a `equipCursor`; `e` when the ragdoll column is focused unequips the selected slot back to inventory
- Ragdoll slot rows show the equipped item name and glyph instead of `[ Empty ]` when occupied
- `NewModel()` pre-equips a default outfit (Cloth Tunic → Chest, Cloth Pants → Legs, Leather Boots → Feet)

## Capabilities

### New Capabilities
- `equipment-system`: `BodySlot` type and constants, `Item.Slots` field, `Inventory.Equipped` array, equip/unequip actions, `equipCursor`, `equipFocused` focus flag, default outfit in `NewModel()`

### Modified Capabilities
- `inventory-system`: `Item` struct gains `Slots []BodySlot`; `Inventory` struct gains `Equipped [NumBodySlots]Item`; pickup logic excludes slots that are considered "equipped" from the carry count
- `fullscreen-inventory`: Ragdoll column renders equipped item name/glyph per slot; adds `equipCursor` highlight on the active slot; `Tab` focus toggle visualised (active column has a highlighted header)
- `input-navigation`: `e` key handler in `ScreenInventory` (equip from left / unequip from right); `Tab` key handler in `ScreenInventory` toggles `equipFocused`; arrow keys in `ScreenInventory` route to `inventoryCursor` or `equipCursor` depending on `equipFocused`

## Impact

- `internal/game/types.go` — `BodySlot` type, constants, `NumBodySlots`; update `Item` and `Inventory` structs
- `internal/game/model.go` — `equipCursor int` and `equipFocused bool` fields on `Model`; default outfit in `NewModel()`
- `internal/game/input.go` — `e` key and `Tab` key handlers; update arrow key routing in `ScreenInventory`
- `internal/game/render.go` — update `renderFullscreenInventory` ragdoll column; active-column header highlight
- `internal/game/input_test.go`, `render_test.go`, `game_test.go` — new test coverage for equip/unequip, default outfit, ragdoll rendering
