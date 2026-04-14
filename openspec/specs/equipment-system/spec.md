## ADDED Requirements

### Requirement: BodySlot type and constants are defined
The system SHALL define `type BodySlot int` in `types.go` with constants:
- `SlotHead BodySlot = 0`
- `SlotChest BodySlot = 1`
- `SlotLeftHand BodySlot = 2`
- `SlotRightHand BodySlot = 3`
- `SlotLegs BodySlot = 4`
- `SlotFeet BodySlot = 5`

A constant `NumBodySlots = 6` SHALL also be defined. The constants SHALL correspond positionally to the existing `equipSlots` string slice in `render.go` (`"Head"`, `"Chest"`, `"Left Hand"`, `"Right Hand"`, `"Legs"`, `"Feet"`).

#### Scenario: BodySlot constants are sequential from zero
- **WHEN** `SlotHead`, `SlotChest`, `SlotLeftHand`, `SlotRightHand`, `SlotLegs`, `SlotFeet` are defined
- **THEN** `int(SlotFeet) == NumBodySlots-1`

### Requirement: Model has equipFocused and equipCursor fields
`Model` SHALL have an `equipFocused bool` field and an `equipCursor int` field. `NewModel()` SHALL leave both at their zero values (`false` and `0`). These fields control which column in the `ScreenInventory` panel has keyboard focus.

#### Scenario: New model starts with left column focused
- **WHEN** `NewModel()` is called
- **THEN** `m.equipFocused == false` and `m.equipCursor == 0`

### Requirement: NewModel pre-equips a default outfit
`NewModel()` SHALL pre-populate `inventory.Equipped` with:
- `Equipped[SlotChest]`: `Item{Name: "Cloth Tunic", Slots: []BodySlot{SlotChest}}`
- `Equipped[SlotLegs]`: `Item{Name: "Cloth Pants", Slots: []BodySlot{SlotLegs}}`
- `Equipped[SlotFeet]`: `Item{Name: "Leather Boots", Slots: []BodySlot{SlotFeet}}`

All other slots (`SlotHead`, `SlotLeftHand`, `SlotRightHand`) SHALL be empty (`Item{}`).

#### Scenario: New model has default outfit equipped
- **WHEN** `NewModel()` is called
- **THEN** `m.inventory.Equipped[SlotChest].Name == "Cloth Tunic"` and `m.inventory.Equipped[SlotLegs].Name == "Cloth Pants"` and `m.inventory.Equipped[SlotFeet].Name == "Leather Boots"`

#### Scenario: New model has empty head and hand slots
- **WHEN** `NewModel()` is called
- **THEN** `m.inventory.Equipped[SlotHead].Name == ""` and `m.inventory.Equipped[SlotLeftHand].Name == ""` and `m.inventory.Equipped[SlotRightHand].Name == ""`

#### Scenario: Default outfit items are not in inventory.Items
- **WHEN** `NewModel()` is called
- **THEN** no item in `m.inventory.Items` has `Name == "Cloth Tunic"` or `Name == "Cloth Pants"` or `Name == "Leather Boots"`

### Requirement: e key equips the selected inventory item
When `m.screenMode == ScreenInventory` and `m.equipFocused == false` and the player presses `e`:
1. If `len(m.inventory.Items) == 0` or the selected item's `Slots` is empty, the key SHALL be a no-op.
2. Otherwise, iterate the item's `Slots` in order to find the first empty slot (`Equipped[slot].Name == ""`). If found, equip there.
3. If no empty slot is found among the item's `Slots`, equip to `item.Slots[0]`, moving the currently equipped item at that slot back into `inventory.Items` (stacking if same name, else adding a new slot, respecting `InventoryMaxSlots` cap — if inventory is full the swap is rejected as a no-op).
4. Remove one count of the equipped item from `inventory.Items`; if `Count` drops to 0, remove the slot and clamp `inventoryCursor`.

#### Scenario: e equips item to empty slot
- **WHEN** `equipFocused == false`, cursor is on an item with `Slots: []BodySlot{SlotHead}`, and `Equipped[SlotHead].Name == ""`
- **THEN** `Equipped[SlotHead]` holds that item and it is removed from `inventory.Items`

#### Scenario: e swaps when target slot is occupied
- **WHEN** `equipFocused == false`, cursor is on a "Leather Hat" with `Slots: []BodySlot{SlotHead}`, and `Equipped[SlotHead].Name == "Cloth Cap"`
- **THEN** `Equipped[SlotHead].Name == "Leather Hat"` and `inventory.Items` contains "Cloth Cap"

#### Scenario: e is a no-op when item has no slots
- **WHEN** `equipFocused == false` and the cursor item has `len(Slots) == 0`
- **THEN** inventory and equipped are unchanged

#### Scenario: e is a no-op when inventory is empty
- **WHEN** `equipFocused == false` and `len(inventory.Items) == 0`
- **THEN** no state changes

### Requirement: e key unequips the selected ragdoll slot
When `m.screenMode == ScreenInventory` and `m.equipFocused == true` and the player presses `e`:
1. If `Equipped[equipCursor].Name == ""`, the key SHALL be a no-op.
2. Otherwise, move `Equipped[equipCursor]` to `inventory.Items` (stacking if same name, else adding a new slot — if inventory is full, the unequip is rejected as a no-op).
3. Clear `Equipped[equipCursor]` to `Item{}`.

#### Scenario: e unequips occupied slot to inventory
- **WHEN** `equipFocused == true`, `equipCursor == int(SlotChest)`, and `Equipped[SlotChest].Name == "Cloth Tunic"`
- **THEN** `Equipped[SlotChest].Name == ""` and `inventory.Items` contains "Cloth Tunic"

#### Scenario: e is a no-op when ragdoll slot is empty
- **WHEN** `equipFocused == true` and `Equipped[equipCursor].Name == ""`
- **THEN** no state changes

#### Scenario: e unequip is rejected when inventory is full
- **WHEN** `equipFocused == true`, the selected slot is occupied, and `len(inventory.Items) == InventoryMaxSlots`
- **THEN** inventory and equipped are unchanged
