## MODIFIED Requirements

### Requirement: NewModel pre-equips a default outfit with real stat values
`NewModel()` SHALL pre-populate `inventory.Equipped` with the following items and stat bonuses:
- `Equipped[SlotChest]`: `Item{Name: "Cloth Tunic", ArmourBonus: 1, Slots: []BodySlot{SlotChest}}`
- `Equipped[SlotLegs]`: `Item{Name: "Cloth Pants", ArmourBonus: 1, Slots: []BodySlot{SlotLegs}}`
- `Equipped[SlotFeet]`: `Item{Name: "Leather Boots", ArmourBonus: 1, Slots: []BodySlot{SlotFeet}}`
- `Equipped[SlotRightHand]`: `Item{Name: "Wooden Sword", DamageBonus: 1, Slots: []BodySlot{SlotRightHand}}`
- `Equipped[SlotLeftHand]`: `Item{Name: "Wooden Shield", ArmourBonus: 1, Slots: []BodySlot{SlotLeftHand}}`

`SlotHead` SHALL remain empty (`Item{}`). None of these items SHALL appear in `inventory.Items`.

#### Scenario: New model has Cloth Tunic with ArmourBonus 1
- **WHEN** `NewModel()` is called
- **THEN** `m.inventory.Equipped[SlotChest].Name == "Cloth Tunic"` and `m.inventory.Equipped[SlotChest].ArmourBonus == 1`

#### Scenario: New model has Cloth Pants with ArmourBonus 1
- **WHEN** `NewModel()` is called
- **THEN** `m.inventory.Equipped[SlotLegs].Name == "Cloth Pants"` and `m.inventory.Equipped[SlotLegs].ArmourBonus == 1`

#### Scenario: New model has Leather Boots with ArmourBonus 1
- **WHEN** `NewModel()` is called
- **THEN** `m.inventory.Equipped[SlotFeet].Name == "Leather Boots"` and `m.inventory.Equipped[SlotFeet].ArmourBonus == 1`

#### Scenario: New model has Wooden Sword with DamageBonus 1
- **WHEN** `NewModel()` is called
- **THEN** `m.inventory.Equipped[SlotRightHand].Name == "Wooden Sword"` and `m.inventory.Equipped[SlotRightHand].DamageBonus == 1`

#### Scenario: New model has Wooden Shield with ArmourBonus 1
- **WHEN** `NewModel()` is called
- **THEN** `m.inventory.Equipped[SlotLeftHand].Name == "Wooden Shield"` and `m.inventory.Equipped[SlotLeftHand].ArmourBonus == 1`

#### Scenario: Starting armour pool is 4 (tunic + pants + boots + shield)
- **WHEN** `NewModel()` is called and combat stats are computed
- **THEN** the player's total armour pool equals 4

### Requirement: Rusty Dagger loot item has DamageBonus 1
`Item{Name: "Rusty Dagger"}` entries in enemy loot tables SHALL have `DamageBonus: 1`. Previously this was 0.

#### Scenario: Rusty Dagger carries DamageBonus 1
- **WHEN** a Rusty Dagger item is defined in a loot table entry
- **THEN** `item.DamageBonus == 1`

### Requirement: Short Sword loot item has DamageBonus 2
`Item{Name: "Short Sword"}` entries in enemy loot tables SHALL have `DamageBonus: 2`. Previously this was 0.

#### Scenario: Short Sword carries DamageBonus 2
- **WHEN** a Short Sword item is defined in a loot table entry
- **THEN** `item.DamageBonus == 2`
