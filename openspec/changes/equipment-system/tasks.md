## 1. Types

- [ ] 1.1 Add `BodySlot int` type and constants (`SlotHead`–`SlotFeet`) and `NumBodySlots = 6` to `types.go`
- [ ] 1.2 Add `Slots []BodySlot` field to `Item` struct in `types.go`
- [ ] 1.3 Add `Equipped [NumBodySlots]Item` field to `Inventory` struct in `types.go`

## 2. Model

- [ ] 2.1 Add `equipFocused bool` and `equipCursor int` fields to `Model` in `model.go`
- [ ] 2.2 Implement `defaultOutfit() [NumBodySlots]Item` helper in `model.go` returning Cloth Tunic (Chest), Cloth Pants (Legs), Leather Boots (Feet)
- [ ] 2.3 Set `inventory.Equipped = defaultOutfit()` in `NewModel()`

## 3. Equip Action

- [ ] 3.1 Implement `equipItem(m Model) Model` in `input.go`: find best empty slot from `item.Slots`, fall back to slot 0 with swap; remove item from `inventory.Items`; guard against non-equippable items and full inventory on swap
- [ ] 3.2 Implement `unequipSlot(m Model) Model` in `input.go`: move `Equipped[equipCursor]` to `inventory.Items` (stack or add slot); guard against empty slot and full inventory

## 4. Input Handlers

- [ ] 4.1 Add `case "tab":` (or `tea.KeyTab`) branch in `handleKey`: toggle `m.equipFocused` when `screenMode == ScreenInventory`, no-op otherwise
- [ ] 4.2 Add `case "e":` branch in `handleKey`: call `equipItem` when `screenMode == ScreenInventory && !equipFocused`; call `unequipSlot` when `screenMode == ScreenInventory && equipFocused`; no-op otherwise
- [ ] 4.3 Update the `up`/`down` routing in `handleKey` for `ScreenInventory`: when `equipFocused == true`, drive `equipCursor` (clamped `[0, NumBodySlots-1]`) instead of `inventoryCursor`

## 5. Rendering

- [ ] 5.1 Update `renderFullscreenInventory` ragdoll slot rows: render `SlotName : [ ItemName ]` when `Equipped[i].Name != ""`, `[ Empty ]` otherwise
- [ ] 5.2 Add `equipCursor` highlight to the active ragdoll slot row when `m.equipFocused == true`
- [ ] 5.3 Differentiate column header styles: active column header (left when `!equipFocused`, right when `equipFocused`) uses a brighter/bold style; inactive header uses a dim style
- [ ] 5.4 Update hint row in the left column to include `e equip  Tab switch`

## 6. Tests

- [ ] 6.1 `TestBodySlot_Constants` in `types_test.go` (or `game_test.go`): assert `int(SlotFeet) == NumBodySlots-1`
- [ ] 6.2 `TestNewModel_DefaultOutfit` in `game_test.go`: assert Cloth Tunic/Cloth Pants/Leather Boots are equipped and not in `inventory.Items`
- [ ] 6.3 `TestEquipItem_EmptySlot` in `input_test.go`: equip item to an empty slot, assert it leaves `inventory.Items` and enters `Equipped`
- [ ] 6.4 `TestEquipItem_Swap` in `input_test.go`: equip to occupied slot, assert old item returns to `inventory.Items`
- [ ] 6.5 `TestEquipItem_NonEquippable` in `input_test.go`: item with no slots, assert no state change
- [ ] 6.6 `TestEquipItem_FullInventorySwapRejected` in `input_test.go`: swap rejected when `inventory.Items` is full
- [ ] 6.7 `TestUnequipSlot_Occupied` in `input_test.go`: unequip occupied slot, assert item appears in `inventory.Items` and slot is cleared
- [ ] 6.8 `TestUnequipSlot_Empty` in `input_test.go`: unequip empty slot, assert no state change
- [ ] 6.9 `TestUnequipSlot_FullInventoryRejected` in `input_test.go`: unequip rejected when inventory full
- [ ] 6.10 `TestTabKey_TogglesEquipFocused` in `input_test.go`: Tab in ScreenInventory toggles `equipFocused`; Tab outside ScreenInventory is no-op
- [ ] 6.11 `TestEquipCursor_Navigation` in `input_test.go`: up/down drive `equipCursor` when `equipFocused == true`; clamp at 0 and `NumBodySlots-1`
- [ ] 6.12 `TestRenderFullscreenInventory_EquippedSlot` in `render_test.go`: assert ragdoll column shows equipped item name
- [ ] 6.13 `TestRenderFullscreenInventory_FocusedHeader` in `render_test.go`: assert active column header style differs when `equipFocused` changes
- [ ] 6.14 Update any existing render tests that assert `[ Empty ]` for default-outfit slots (Chest/Legs/Feet) to expect item names instead
