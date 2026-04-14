## Requirements

### Requirement: Use key triggers context-sensitive item interaction
When the player presses `u` in `ModeLocal` or `ModeDungeon`, the system SHALL determine whether the currently selected inventory item (at `m.inventoryCursor`) has a valid interaction with any object on the player's current cell or in the four cardinal adjacent cells. If no valid interaction exists the key press is a no-op.

#### Scenario: Use is ignored in ModeWorld
- **WHEN** the player presses `u` in `ModeWorld`
- **THEN** no interaction occurs and `m` is unchanged

#### Scenario: Use with empty inventory is a no-op
- **WHEN** the player presses `u` with `len(m.inventory.Items) == 0`
- **THEN** no interaction occurs and `m` is unchanged

### Requirement: Axe can chop adjacent trees
When the player holds an `Axe` (`item.Name == "Axe"`) and presses `u`, the system SHALL scan the four cardinal cells adjacent to the player for an `Object` whose `Char` is a tree glyph (`'♣'` for Tree, `'♠'` for Pine, `'T'` for Palm/Tropical, or as defined by forest biome tables). When a tree object is found:
1. Remove the tree object from that cell (set `Object` to `nil`)
2. Leave the `Axe` in inventory (axe is not consumed)

#### Scenario: Chop removes adjacent tree
- **WHEN** the player holds an Axe and presses `u` adjacent to a Tree object
- **THEN** the tree cell's `Object` is nil

#### Scenario: Chop does not consume axe
- **WHEN** the player chops a tree
- **THEN** the Axe slot in inventory still has `Count >= 1`

#### Scenario: Chop ignores non-tree objects
- **WHEN** the player holds an Axe and presses `u` but no adjacent cell contains a tree object
- **THEN** no cell changes and inventory is unchanged

#### Scenario: Chop targets nearest tree when multiple are adjacent
- **WHEN** multiple adjacent cells contain tree objects
- **THEN** the system removes exactly one tree (northmost first, then east, south, west priority)

### Requirement: Item interaction priority and extensibility
The interaction system SHALL be implemented as a dispatching switch or table keyed on `item.Name`, so that additional item interactions can be added in future without restructuring the handler. Only interactions explicitly defined in the spec are active; all others are no-ops.

#### Scenario: Undefined item has no interaction
- **WHEN** the player holds an item with an unrecognised name (e.g. "Torch") and presses `u`
- **THEN** no cell changes
