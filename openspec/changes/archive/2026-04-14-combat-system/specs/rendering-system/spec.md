## ADDED Requirements

### Requirement: buildView dispatches to renderCombatScreen for ScreenCombat
When `m.screenMode == ScreenCombat`, `buildView` SHALL return `renderCombatScreen(m)` immediately, before any map, sidebar, HUD, or key bar rendering. This mirrors the existing `ScreenInventory` early-return pattern.

#### Scenario: buildView returns combat screen for ScreenCombat
- **WHEN** `m.screenMode == ScreenCombat` and `m.combatState` is non-nil
- **THEN** `buildView(m)` returns the same string as `renderCombatScreen(m)`

#### Scenario: buildView does not render HUD in ScreenCombat
- **WHEN** `m.screenMode == ScreenCombat`
- **THEN** `buildView(m)` does not contain the HUD clock or tile info
