## ADDED Requirements

### Requirement: EnemyTemplate carries a weighted loot table
Each `EnemyTemplate` SHALL have a `LootTable []LootEntry` field. `LootEntry` has `Item Item` and `Weight int`. A `Weight` of 0 means "no drop" (empty item). The total weight is the sum of all entry weights. At least one "no drop" entry (empty `Item`, `Weight >= 1`) SHALL be present in every loot table so that enemies do not always drop something.

#### Scenario: Loot table contains a no-drop entry
- **WHEN** any `EnemyTemplate.LootTable` is inspected
- **THEN** at least one entry has `Item.Name == ""`

### Requirement: resolveEnemyLoot picks one item from the loot table
`resolveEnemyLoot(table []LootEntry, rng *rand.Rand) Item` SHALL compute total weight, pick a random value in `[0, totalWeight)`, and return the `Item` from the entry whose cumulative weight range includes the chosen value. If `table` is empty or total weight is 0, it SHALL return a zero `Item`.

#### Scenario: Single-entry table always returns that entry's item
- **WHEN** `resolveEnemyLoot([]LootEntry{{Item: sword, Weight: 1}}, rng)` is called
- **THEN** `sword` is returned every time

#### Scenario: Zero-weight table returns empty item
- **WHEN** all entries have `Weight == 0`
- **THEN** an empty `Item` is returned

### Requirement: Loot is added to inventory after victorious dungeon combat
In the `ScreenCombat` dismiss handler, when `PlayerWon == true` and the defeated enemy was a `DungeonEnemy`, the system SHALL call `resolveEnemyLoot` with the enemy's template loot table and a fresh RNG. If the result has a non-empty `Name` and `len(m.inventory.Items) < InventoryMaxSlots`, the item SHALL be added to `m.inventory.Items` (stacked if an item with the same `Name` already exists, otherwise appended as a new slot with `Count=1`). If inventory is full or the drop is empty, no change is made to inventory.

#### Scenario: Item added to inventory after victory
- **WHEN** loot resolves to a non-empty item and inventory has space
- **THEN** the item appears in `m.inventory.Items` after dismissing the combat result

#### Scenario: Full inventory silently discards the drop
- **WHEN** `len(m.inventory.Items) == InventoryMaxSlots` and loot resolves to a non-empty item
- **THEN** `m.inventory.Items` is unchanged

#### Scenario: Empty loot drop does not modify inventory
- **WHEN** `resolveEnemyLoot` returns an item with `Name == ""`
- **THEN** `m.inventory.Items` is unchanged
