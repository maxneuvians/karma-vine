## MODIFIED Requirements

### Requirement: Unlit torches and braziers are pickupable
Torches and braziers that are generated in the unlit state (`Lit == false`) SHALL have `Pickupable == true` set at generation time. Torches and braziers that start lit (e.g., those adjacent to stairs) SHALL have `Pickupable == false` (or the default) to keep them as fixed light sources.

> **Note:** In the current implementation all torches and braziers start unlit. This requirement simply ensures `Pickupable: true` is set on those objects so the inventory system can pick them up.

#### Scenario: Unlit torch is pickupable
- **WHEN** a dungeon level is generated
- **THEN** every `Object` with `Name == "Torch"` and `Lit == false` has `Pickupable == true`

#### Scenario: Brazier is pickupable when unlit
- **WHEN** a dungeon level is generated
- **THEN** every `Object` with `Name == "Brazier"` and `Lit == false` has `Pickupable == true`

#### Scenario: Lit torch is not pickupable
- **WHEN** a torch has `Lit == true` (manually lit by player)
- **THEN** `Pickupable == false` (or unchanged from generation; lighting a torch removes its pickupability)
