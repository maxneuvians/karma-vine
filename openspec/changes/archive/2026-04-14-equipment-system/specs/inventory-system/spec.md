## MODIFIED Requirements

### Requirement: Item type is defined
The system SHALL define an `Item` struct with fields `Char rune`, `Color string`, `Name string`, `Count int`, and `Slots []BodySlot`. The `Count` field represents how many of the item are stacked in this slot (minimum 1). The `Slots` field lists which `BodySlot` positions the item can be equipped to; an empty slice means the item is not equippable.

#### Scenario: Item zero-value is safe
- **WHEN** an `Item` is allocated with `Item{}`
- **THEN** it does not panic; `Count` defaults to 0 and `Slots` is nil (treated as empty)

#### Scenario: Equippable item declares its slots
- **WHEN** an item is defined with `Slots: []BodySlot{SlotLeftHand, SlotRightHand}`
- **THEN** `len(item.Slots) == 2`

#### Scenario: Non-equippable item has no slots
- **WHEN** an item is defined with `Slots: nil` or `Slots: []BodySlot{}`
- **THEN** `len(item.Slots) == 0`

### Requirement: Inventory is a fixed-capacity collection on Model
The system SHALL define an `Inventory` struct with fields `Items []Item` and `Equipped [NumBodySlots]Item`. A constant `InventoryMaxSlots = 8` SHALL define the maximum number of distinct item stacks the player can carry in `Items`. `Model` SHALL have an `inventory Inventory` field and an `inventoryCursor int` field. `NewModel()` SHALL initialise `inventory.Items` with an empty slice and `inventory.Equipped` with the default outfit.

#### Scenario: New model has empty item list
- **WHEN** `NewModel()` is called
- **THEN** `len(m.inventory.Items) == 0`

#### Scenario: Inventory.Items respects max slots
- **WHEN** the player attempts to pick up an item and `len(m.inventory.Items) == InventoryMaxSlots`
- **THEN** the pickup is rejected and inventory remains unchanged

#### Scenario: Equipped array has NumBodySlots entries
- **WHEN** `Inventory` is declared
- **THEN** `len(m.inventory.Equipped) == NumBodySlots`
