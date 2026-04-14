## ADDED Requirements

### Requirement: Item type is defined
The system SHALL define an `Item` struct with fields `Char rune`, `Color string`, `Name string`, and `Count int`. The `Count` field represents how many of the item are stacked in this slot (minimum 1).

#### Scenario: Item zero-value is safe
- **WHEN** an `Item` is allocated with `Item{}`
- **THEN** it does not panic; `Count` defaults to 0

### Requirement: Inventory is a fixed-capacity collection on Model
The system SHALL define an `Inventory` struct with a `Items []Item` slice. A constant `InventoryMaxSlots = 8` SHALL define the maximum number of distinct item stacks the player can carry. `Model` SHALL have an `inventory Inventory` field, a `showInventory bool` field, and an `inventoryCursor int` field. `NewModel()` SHALL initialise `inventory` with an empty `Items` slice.

#### Scenario: New model has empty inventory
- **WHEN** `NewModel()` is called
- **THEN** `len(m.inventory.Items) == 0`

#### Scenario: Inventory respects max slots
- **WHEN** the player attempts to pick up an item and `len(m.inventory.Items) == InventoryMaxSlots`
- **THEN** the pickup is rejected and inventory remains unchanged

### Requirement: Object has Pickupable field
The system SHALL add a `Pickupable bool` field to the `Object` struct. When `Pickupable == true`, the player may pick up the object from the map.

#### Scenario: Non-pickupable objects are ignored
- **WHEN** the player presses `g` on a cell containing an `Object` with `Pickupable == false`
- **THEN** `m.inventory` is unchanged

### Requirement: Pickup removes object from map and adds Item to inventory
When the player presses `g` in `ModeLocal` or `ModeDungeon` and the player's current cell contains an `Object` with `Pickupable == true`, the system SHALL:
1. Remove the object from the map (set cell to `nil` or clear the cell object)
2. Add an `Item{Char: obj.Char, Color: obj.Color, Name: obj.Name, Count: 1}` to `inventory.Items`; if an item with the same `Name` already exists, increment its `Count` instead of adding a new slot

#### Scenario: Picking up item clears cell
- **WHEN** the player presses `g` on a cell with a pickupable object
- **THEN** the cell no longer contains that object

#### Scenario: Picking up item adds to inventory
- **WHEN** the player picks up an object named "Axe"
- **THEN** `m.inventory` contains an item with `Name == "Axe"` and `Count >= 1`

#### Scenario: Picking up same item stacks
- **WHEN** the player already holds 1 "Torch" and picks up another "Torch"
- **THEN** `m.inventory` has exactly one slot for "Torch" with `Count == 2` rather than two separate slots

#### Scenario: Pickup ignored in ModeWorld
- **WHEN** the player presses `g` in `ModeWorld`
- **THEN** `m.inventory` is unchanged

### Requirement: Drop places item back onto the map
When the player presses `d` in `ModeLocal` or `ModeDungeon` and `len(m.inventory.Items) > 0`, the system SHALL:
1. Take the item at `m.inventoryCursor`
2. Decrement its `Count`; if `Count` reaches 0 remove the slot and clamp `inventoryCursor`
3. Place a new `Object{Char, Color, Name, Blocking: false, Pickupable: true}` on the player's current cell, or on the nearest free adjacent cell if the current cell is occupied

#### Scenario: Drop places object on current cell
- **WHEN** the player drops an item onto an empty cell
- **THEN** the cell now contains an `Object` matching the dropped item with `Pickupable == true`

#### Scenario: Drop clamps around occupied cell
- **WHEN** the player drops an item but the current cell is occupied by a blocking object
- **THEN** the item is placed on the nearest free adjacent cell

#### Scenario: Drop decrements count
- **WHEN** the player has 2 Torches and drops one
- **THEN** the inventory slot for "Torch" has `Count == 1`

#### Scenario: Drop removes slot when count hits zero
- **WHEN** the player has 1 Axe and drops it
- **THEN** the "Axe" slot is removed from inventory

#### Scenario: Drop ignored in ModeWorld
- **WHEN** the player presses `d` in `ModeWorld`
- **THEN** `m.inventory` is unchanged

### Requirement: Inventory cursor navigation
When the inventory panel is open (`m.showInventory == true`), the `up`/`w` and `down`/`s` keys SHALL move `m.inventoryCursor`. The cursor SHALL be clamped to `[0, len(m.inventory.Items)-1]`. When items are removed and the cursor would be out of range, it SHALL be clamped down.

#### Scenario: Cursor wraps at bottom
- **WHEN** `inventoryCursor == len(inventory.Items)-1` and the player presses `down`
- **THEN** `inventoryCursor` stays at `len(inventory.Items)-1` (clamped, no wrap)

#### Scenario: Cursor wraps at top
- **WHEN** `inventoryCursor == 0` and the player presses `up`
- **THEN** `inventoryCursor` stays at 0
