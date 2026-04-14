## ADDED Requirements

### Requirement: [ and ] keys control combat playback speed during ScreenCombat
When `m.screenMode == ScreenCombat`, pressing `[` SHALL decrease `m.combatSpeed` by 1 (minimum `CombatSpeedSlow`). Pressing `]` SHALL increase `m.combatSpeed` by 1 (maximum `CombatSpeedFast`). Both keys SHALL be handled in the `ScreenCombat` early-return block alongside the existing `enter`/`space` and `q`/`ctrl+c` handlers.

#### Scenario: ] increases speed from Slow to Normal
- **WHEN** `m.screenMode == ScreenCombat`, `m.combatSpeed == CombatSpeedSlow`, and `]` is pressed
- **THEN** `m.combatSpeed == CombatSpeedNormal`

#### Scenario: ] clamps at CombatSpeedFast
- **WHEN** `m.screenMode == ScreenCombat`, `m.combatSpeed == CombatSpeedFast`, and `]` is pressed
- **THEN** `m.combatSpeed == CombatSpeedFast` (unchanged)

#### Scenario: [ decreases speed from Fast to Normal
- **WHEN** `m.screenMode == ScreenCombat`, `m.combatSpeed == CombatSpeedFast`, and `[` is pressed
- **THEN** `m.combatSpeed == CombatSpeedNormal`

#### Scenario: [ clamps at CombatSpeedSlow
- **WHEN** `m.screenMode == ScreenCombat`, `m.combatSpeed == CombatSpeedSlow`, and `[` is pressed
- **THEN** `m.combatSpeed == CombatSpeedSlow` (unchanged)

#### Scenario: [ and ] are no-ops outside ScreenCombat
- **WHEN** `m.screenMode == ScreenNormal` and `[` or `]` is pressed
- **THEN** `m.combatSpeed` is unchanged (handled by existing time-scale logic or ignored)
