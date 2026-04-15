## MODIFIED Requirements

### Requirement: NewModel pre-equips a default outfit
`NewModel()` SHALL pre-populate `inventory.Equipped` with:
- `Equipped[SlotChest]`: `Item{Name: "Cloth Tunic", Slots: []BodySlot{SlotChest}}`
- `Equipped[SlotLegs]`: `Item{Name: "Cloth Pants", Slots: []BodySlot{SlotLegs}}`
- `Equipped[SlotFeet]`: `Item{Name: "Leather Boots", Slots: []BodySlot{SlotFeet}}`
- `Equipped[SlotRightHand]`: `Item{Name: "Wooden Sword", Slots: []BodySlot{SlotRightHand}, DamageBonus: 1}`
- `Equipped[SlotLeftHand]`: `Item{Name: "Wooden Shield", Slots: []BodySlot{SlotLeftHand}, ArmourBonus: 1}`

`SlotHead` SHALL remain empty (`Item{}`).

None of these items SHALL appear in `inventory.Items`.

#### Scenario: New model has default outfit equipped
- **WHEN** `NewModel()` is called
- **THEN** `m.inventory.Equipped[SlotChest].Name == "Cloth Tunic"` and `m.inventory.Equipped[SlotLegs].Name == "Cloth Pants"` and `m.inventory.Equipped[SlotFeet].Name == "Leather Boots"`

#### Scenario: New model has Wooden Sword in right hand
- **WHEN** `NewModel()` is called
- **THEN** `m.inventory.Equipped[SlotRightHand].Name == "Wooden Sword"` and `m.inventory.Equipped[SlotRightHand].DamageBonus == 1`

#### Scenario: New model has Wooden Shield in left hand
- **WHEN** `NewModel()` is called
- **THEN** `m.inventory.Equipped[SlotLeftHand].Name == "Wooden Shield"` and `m.inventory.Equipped[SlotLeftHand].ArmourBonus == 1`

#### Scenario: New model has empty head slot
- **WHEN** `NewModel()` is called
- **THEN** `m.inventory.Equipped[SlotHead].Name == ""`

#### Scenario: Default outfit items are not in inventory.Items
- **WHEN** `NewModel()` is called
- **THEN** no item in `m.inventory.Items` has `Name` matching any default equipped item name

## ADDED Requirements

### Requirement: Item has ArmourBonus and DamageBonus fields
The `Item` struct SHALL include `ArmourBonus int` and `DamageBonus int` fields. Both default to `0` for items that provide no combat bonuses. An item with `ArmourBonus: 1` adds 1 to the wielder's armour pool; an item with `DamageBonus: 1` adds 1 to both `MinDamage` and `MaxDamage`.

#### Scenario: Item zero-value has no bonuses
- **WHEN** an `Item` is allocated with `Item{}`
- **THEN** `item.ArmourBonus == 0` and `item.DamageBonus == 0`

#### Scenario: Shield item carries ArmourBonus
- **WHEN** `Item{Name: "Wooden Shield", ArmourBonus: 1}` is defined
- **THEN** `item.ArmourBonus == 1`

#### Scenario: Sword item carries DamageBonus
- **WHEN** `Item{Name: "Wooden Sword", DamageBonus: 1}` is defined
- **THEN** `item.DamageBonus == 1`
