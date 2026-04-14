## ADDED Requirements

### Requirement: Inventory panel is rendered as an overlay when open
When `m.showInventory == true`, the system SHALL render a `renderInventoryPanel` overlay on top of the current map view. The panel SHALL:
- Display a title line "Inventory"
- List each item slot as `[N] ItemName (xCount)` or similar compact format
- Highlight the row at `m.inventoryCursor` with a distinct background
- Show "Empty" when `len(m.inventory.Items) == 0`
- Be positioned consistently regardless of current game mode

#### Scenario: Inventory panel shows when showInventory is true
- **WHEN** `m.showInventory == true`
- **THEN** the rendered output contains the text "Inventory"

#### Scenario: Inventory panel hidden when showInventory is false
- **WHEN** `m.showInventory == false`
- **THEN** the rendered output does NOT contain the inventory panel

#### Scenario: Empty inventory shows placeholder
- **WHEN** `m.showInventory == true` and `len(m.inventory.Items) == 0`
- **THEN** the rendered output contains "Empty"

### Requirement: HUD shows inventory item count
The HUD row SHALL include the current inventory count in the format `Items: N/8` (where 8 = `InventoryMaxSlots`). This count SHALL be visible in all three modes when the HUD is displayed.

#### Scenario: HUD shows item count
- **WHEN** the player has 2 items and the HUD is rendered
- **THEN** the HUD contains "Items: 2/8" (or equivalent display)

### Requirement: Key bar shows inventory and use hints
The key bar SHALL include the following additional hints:
- `i inv` — toggle inventory
- `g pick` — pick up (shown in ModeLocal and ModeDungeon only)
- `d drop` — drop item (shown in ModeLocal and ModeDungeon only)
- `u use` — use item (shown in ModeLocal and ModeDungeon only)

#### Scenario: Key bar in ModeLocal includes inventory keys
- **WHEN** `m.mode == ModeLocal`
- **THEN** the key bar contains "i", "g", "d", and "u" hints
