## ADDED Requirements

### Requirement: Tab key toggles inventory panel focus in ScreenInventory
When `m.screenMode == ScreenInventory`, pressing `Tab` SHALL toggle `m.equipFocused`: `false → true` and `true → false`. Outside `ScreenInventory`, `Tab` SHALL be a no-op.

#### Scenario: Tab switches focus from left to right column
- **WHEN** `screenMode == ScreenInventory`, `equipFocused == false`, and the player presses `Tab`
- **THEN** `m.equipFocused == true`

#### Scenario: Tab switches focus from right to left column
- **WHEN** `screenMode == ScreenInventory`, `equipFocused == true`, and the player presses `Tab`
- **THEN** `m.equipFocused == false`

#### Scenario: Tab is a no-op outside ScreenInventory
- **WHEN** `screenMode == ScreenNormal` and the player presses `Tab`
- **THEN** `m.equipFocused` is unchanged

### Requirement: Arrow keys navigate the focused inventory column in ScreenInventory
When `m.screenMode == ScreenInventory`:
- If `m.equipFocused == false`: `up`/`w` decrements `inventoryCursor` (clamped at 0); `down`/`s` increments `inventoryCursor` (clamped at `len(inventory.Items)-1`). This matches existing behaviour.
- If `m.equipFocused == true`: `up`/`w` decrements `equipCursor` (clamped at 0); `down`/`s` increments `equipCursor` (clamped at `NumBodySlots-1`). `left`/`right`/`a`/`d` are no-ops in both focus states.

#### Scenario: Down moves equipCursor when right column is focused
- **WHEN** `screenMode == ScreenInventory`, `equipFocused == true`, `equipCursor == 1`, and the player presses `down`
- **THEN** `equipCursor == 2`

#### Scenario: equipCursor clamps at bottom
- **WHEN** `screenMode == ScreenInventory`, `equipFocused == true`, `equipCursor == NumBodySlots-1`, and the player presses `down`
- **THEN** `equipCursor` remains `NumBodySlots-1`

#### Scenario: equipCursor clamps at top
- **WHEN** `screenMode == ScreenInventory`, `equipFocused == true`, `equipCursor == 0`, and the player presses `up`
- **THEN** `equipCursor` remains `0`

#### Scenario: Up/down still navigate inventoryCursor when left column is focused
- **WHEN** `screenMode == ScreenInventory`, `equipFocused == false`, and the player presses `down`
- **THEN** `inventoryCursor` increments (clamped) and `equipCursor` is unchanged

### Requirement: e key equips or unequips based on focused column
When `m.screenMode == ScreenInventory`, pressing `e` SHALL:
- If `m.equipFocused == false`: attempt to equip the item at `inventoryCursor` (as defined in the equipment-system spec)
- If `m.equipFocused == true`: attempt to unequip the slot at `equipCursor` (as defined in the equipment-system spec)

Outside `ScreenInventory`, `e` SHALL be a no-op (no existing global binding).

#### Scenario: e equips when left column focused
- **WHEN** `screenMode == ScreenInventory`, `equipFocused == false`, cursor is on an equippable item, and target slot is empty
- **THEN** the item moves from `inventory.Items` to `inventory.Equipped`

#### Scenario: e unequips when right column focused
- **WHEN** `screenMode == ScreenInventory`, `equipFocused == true`, and the selected slot is occupied
- **THEN** the equipped item moves from `inventory.Equipped` to `inventory.Items`

#### Scenario: e is a no-op outside ScreenInventory
- **WHEN** `screenMode == ScreenNormal` and the player presses `e`
- **THEN** no state changes
