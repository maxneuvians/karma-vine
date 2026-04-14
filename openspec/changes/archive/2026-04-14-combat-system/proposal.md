## Why

The game has exploration and inventory mechanics but no conflict resolution — encounters with entities have no consequence. A combat system gives those encounters meaning and creates a core gameplay loop around risk, preparation, and reward.

## What Changes

- New fullscreen combat screen that activates when the player engages an entity
- Turn-based auto-battler: both sides act in order of initiative with no player input mid-combat
- Combatant stats: HP, armour (damage reduction), damage, initiative
- End conditions: player death ends the game; enemy death returns to exploration
- Between-round hook system for side effects (e.g. item triggers, status conditions) — extensible for future item effects
- Entities (animals) gain combat stats; player derives stats from equipment and base values
- New `ScreenCombat` screen mode added to the screen state machine

## Capabilities

### New Capabilities
- `combat-system`: Core turn-based auto-battle loop, combatant stats, initiative ordering, round resolution, end conditions, and the between-round side-effect hook
- `combat-rendering`: Fullscreen combat screen — combatant stat panels, combat log, round-by-round narration

### Modified Capabilities
- `input-navigation`: `ScreenCombat` added to screen mode enum; key handling during combat (dismiss result screen, quit)
- `animal-system`: Animals gain combat stats (`HP`, `Armour`, `Damage`, `Initiative`) used when combat is initiated
- `rendering-system`: `buildView` routes `ScreenCombat` to the new combat renderer
- `item-interactions`: Document the between-round hook interface that item effects will implement in the future

## Impact

- `internal/game/types.go` — new `CombatStats`, `Combatant`, `CombatState`, `CombatResult` types; `ScreenCombat` constant
- `internal/game/model.go` — `combatState *CombatState` field; combat initiation helpers
- `internal/game/combat.go` — new file: `resolveCombat`, `runRound`, `applyBetweenRoundHooks`
- `internal/game/render.go` — `renderCombatScreen` function; `buildView` dispatch
- `internal/game/input.go` — key handling for `ScreenCombat`
- `internal/game/animal.go` — default combat stats per animal type
- No new dependencies
