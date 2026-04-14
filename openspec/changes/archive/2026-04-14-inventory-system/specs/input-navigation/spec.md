## ADDED Requirements

### Requirement: `g` picks up item in local and dungeon mode
The system SHALL handle the `g` key in `ModeLocal` and `ModeDungeon` to attempt a pickup on the player's current cell. The key SHALL be ignored in `ModeWorld`.

#### Scenario: `g` in ModeLocal triggers pickup
- **WHEN** the player presses `g` in `ModeLocal`
- **THEN** the input handler calls the pickup logic

#### Scenario: `g` in ModeWorld is silently ignored
- **WHEN** the player presses `g` in `ModeWorld`
- **THEN** no state change occurs

### Requirement: `d` drops selected item in local and dungeon mode
The system SHALL handle the `d` key in `ModeLocal` and `ModeDungeon` to drop the item at `m.inventoryCursor`. The key SHALL be ignored in `ModeWorld`.

#### Scenario: `d` in ModeDungeon triggers drop
- **WHEN** the player presses `d` in `ModeDungeon`
- **THEN** the input handler calls the drop logic

### Requirement: `i` toggles inventory panel in all modes
The system SHALL handle the `i` key in all modes to toggle `m.showInventory`. When `m.showInventory` transitions from `true` to `false`, `m.inventoryCursor` SHALL be clamped to `[0, len(m.inventory.Items)-1]` (or 0 if empty).

#### Scenario: `i` toggles showInventory
- **WHEN** `m.showInventory == false` and the player presses `i`
- **THEN** `m.showInventory == true`

#### Scenario: Second `i` closes inventory
- **WHEN** `m.showInventory == true` and the player presses `i`
- **THEN** `m.showInventory == false`

### Requirement: `u` triggers item use in local and dungeon mode
The system SHALL handle the `u` key in `ModeLocal` and `ModeDungeon` to attempt an item interaction. The key SHALL be ignored in `ModeWorld`.

#### Scenario: `u` in ModeLocal triggers use
- **WHEN** the player presses `u` in `ModeLocal`
- **THEN** the input handler calls the item-use logic

### Requirement: Inventory cursor keys active while inventory is open
When `m.showInventory == true`, the `up`/`w` and `down`/`s` directional keys SHALL move `m.inventoryCursor` instead of moving the player.

#### Scenario: `up` moves cursor up when inventory open
- **WHEN** `showInventory == true` and the player presses `up`
- **THEN** `inventoryCursor` decrements (clamped at 0) and player does not move
