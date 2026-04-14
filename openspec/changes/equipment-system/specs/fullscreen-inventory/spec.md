## MODIFIED Requirements

### Requirement: renderFullscreenInventory renders a full-viewport inventory
The system SHALL implement `renderFullscreenInventory(m Model) string` that fills the entire viewport (`m.viewportW × m.viewportH`). The layout SHALL be two columns:
- **Left column** (60% of viewport width): item list with header "Inventory" (highlighted when `!m.equipFocused`), separator, item rows (`[glyph] Name  xN`), cursor-highlighted row at `inventoryCursor`, "Empty" placeholder when no items, and hint row at bottom (`i close  e equip  Tab switch`).
- **Right column** (remaining width): header "Equipment" (highlighted when `m.equipFocused`), ASCII ragdoll body outline centred vertically, with named slot labels rendered as `SlotName : [ ItemName ]` when occupied or `SlotName : [ Empty ]` when empty. The row at `equipCursor` SHALL be highlighted when `m.equipFocused == true`.

The panel SHALL be rendered regardless of the current `m.mode` (world/local/dungeon are all valid).

#### Scenario: Fullscreen panel fills viewport
- **WHEN** `renderFullscreenInventory` is called with `viewportW=120, viewportH=40`
- **THEN** the returned string contains exactly 39 newline characters (40 rows, last row has no trailing newline)

#### Scenario: Item list shows glyph and count
- **WHEN** inventory has one item "Axe" with Count=2
- **THEN** the rendered output contains "Axe" and "x2"

#### Scenario: Empty inventory shows placeholder
- **WHEN** `len(m.inventory.Items) == 0`
- **THEN** the rendered output contains "Empty"

#### Scenario: Ragdoll column shows occupied slot with item name
- **WHEN** `Equipped[SlotChest].Name == "Cloth Tunic"`
- **THEN** the ragdoll column contains "Cloth Tunic" and does NOT contain `"Chest  : [ Empty ]"`

#### Scenario: Ragdoll column shows empty slot placeholder
- **WHEN** `Equipped[SlotHead].Name == ""`
- **THEN** the ragdoll column contains `"[ Empty ]"` for the Head slot

#### Scenario: Cursor row is visually distinct in left column
- **WHEN** `equipFocused == false`, `inventoryCursor == 1`, and there are 3 items
- **THEN** the second item row is rendered in a distinct highlight style

#### Scenario: equipCursor row is visually distinct in right column when focused
- **WHEN** `equipFocused == true` and `equipCursor == int(SlotChest)`
- **THEN** the Chest slot row is rendered in a distinct highlight style

#### Scenario: Left header is highlighted when left column is focused
- **WHEN** `equipFocused == false`
- **THEN** the "Inventory" header has a distinct active style and "Equipment" header has an inactive style

#### Scenario: Right header is highlighted when right column is focused
- **WHEN** `equipFocused == true`
- **THEN** the "Equipment" header has a distinct active style and "Inventory" header has an inactive style
