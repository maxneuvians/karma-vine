## ADDED Requirements

### Requirement: HUD shows PAUSED indicator when game is paused
When `m.paused == true` and `m.screenMode == ScreenNormal`, the HUD status bar SHALL include a visible `[PAUSED]` label. The indicator SHALL appear in all three modes (`ModeWorld`, `ModeLocal`, `ModeDungeon`). When `m.paused == false`, the indicator SHALL NOT appear in the HUD.

#### Scenario: HUD contains PAUSED label when paused
- **WHEN** `m.paused == true` and `m.screenMode == ScreenNormal`
- **THEN** the rendered output contains `[PAUSED]`

#### Scenario: HUD does not contain PAUSED label when unpaused
- **WHEN** `m.paused == false`
- **THEN** the rendered output does NOT contain `[PAUSED]`

#### Scenario: PAUSED indicator appears in ModeWorld HUD
- **WHEN** `m.paused == true` and `m.mode == ModeWorld`
- **THEN** the HUD row contains `[PAUSED]`

#### Scenario: PAUSED indicator appears in ModeDungeon HUD
- **WHEN** `m.paused == true` and `m.mode == ModeDungeon`
- **THEN** the HUD row contains `[PAUSED]`
