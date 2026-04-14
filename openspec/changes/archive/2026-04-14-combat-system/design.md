## Context

The game currently has exploration (world/local/dungeon modes), an inventory, and equipment slots, but no combat. Animals exist on local maps but serve only as visual elements. This design introduces a self-contained combat system that activates when the player engages an animal, resolves the fight automatically (auto-battler), and then returns to exploration — or ends the game on player death.

A key forward-looking constraint from the proposal is the **between-round hook** system: items (and future status effects) must be able to modify combatant stats between rounds. This must be designed into the core loop now even though no concrete item effects ship in this iteration.

## Goals / Non-Goals

**Goals:**
- Turn-based auto-battle that requires no player input once started
- Deterministic given a fixed RNG seed (testable)
- Combatant stats: HP, Armour (flat damage reduction), Damage (per-hit range), Initiative (determines who acts first)
- Between-round hook interface extensible for items and status effects
- Fullscreen `ScreenCombat` view with stat panels and a scrollable combat log
- Game-over on player death; return to exploration on enemy death
- Animals gain default combat stats (varies by type)
- Player stats derived from base values + equipped items

**Non-Goals:**
- Player choices or abilities during combat (future work)
- Multi-enemy encounters (single opponent per fight)
- Persistent injury / carry-over HP between fights (HP resets after combat)
- Networked or async combat
- Concrete item effects (the hook interface is built; no implementations ship yet)

## Decisions

### 1. Pure function combat resolution (`resolveCombat`)

**Decision:** Combat is resolved in a single pure function call that returns a `CombatLog` slice and a `CombatResult`. The model stores this result; the render layer reads it.

**Rationale:** Keeps the combat loop fully testable without needing to drive a bubbletea `Update` cycle. Matches the existing pattern of pure helpers (`applyDelta`, `equipItem`, etc.).

**Alternative considered:** Step-by-step resolution driven by a new `CombatTickMsg`. Rejected — adds async complexity and makes unit-testing harder with no gameplay benefit since it's an auto-battler anyway.

### 2. Between-round hook as a function slice

**Decision:** Define `type RoundHook func(attacker, defender *Combatant)`. `CombatState` carries `[]RoundHook`. Before each round, all hooks are called with the current combatants (order: player first, then enemy). Hooks may mutate stats in-place.

**Rationale:** Lightweight, no interface boilerplate, easy to append at combatant construction time. Items that have conditional effects (e.g. "if HP < 50%, +5 armour") can register a hook closure that captures the item reference.

**Alternative considered:** An `Effect` interface with `Apply(attacker, defender *Combatant)`. More formal but unnecessary ceremony for a hook that has one job.

### 3. Stats live on `Combatant`, not on `Item` or `Animal` directly

**Decision:** Introduce a `Combatant` struct (`Name`, `HP`, `MaxHP`, `Armour`, `MinDamage`, `MaxDamage`, `Initiative`) constructed at combat start from the entity's source data. Combat operates solely on `Combatant` values.

**Rationale:** Decouples combat resolution from the inventory/animal type system. The player combatant is built by summing base stats + equipment bonuses. Animal combatants are built from a stat table keyed on `Animal.Name`. Neither the `Inventory` nor `Animal` structs need to know about combat.

**Alternative considered:** Adding combat fields directly to `Animal` and deriving player stats inline. Tighter coupling; harder to test and extend.

### 4. HP does not carry over after combat

**Decision:** Player HP resets to `MaxHP` after each combat (win or loss). Enemy HP is irrelevant after combat ends.

**Rationale:** Simplest correct behaviour for the first iteration. Persistent injury requires a model field for current HP, UI for HP display on the HUD, and healing mechanics — all separate concerns.

### 5. `ScreenCombat` blocks normal input

**Decision:** While `m.screenMode == ScreenCombat`, all movement, inventory, and world-interaction keys are suppressed. Only `enter`/`space` (dismiss result) and `q`/`ctrl+c` (quit) are handled.

**Rationale:** Auto-battler by design — player has nothing to do mid-fight. Consistent with how `ScreenInventory` blocks movement.

### 6. Combat initiated by `g` (pick up / interact) on an animal cell

**Decision:** When the player presses `g` (pick up) on a cell occupied by an animal in `ModeLocal`, combat is initiated instead of item pickup.

**Rationale:** Reuses an existing, natural "interact with cell" affordance. Avoids adding a new keybind for a mechanic the player will discover organically. Can be revisited (e.g. automatic on movement into animal cell) later.

**Alternative considered:** Automatic combat on movement into animal cell. More realistic but removes player agency about when to fight; deferred.

## Risks / Trade-offs

- **RNG non-determinism in tests** → Mitigation: `resolveCombat` accepts a `*rand.Rand` parameter; tests pass a seeded source.
- **Hook slice mutation safety** → Mitigation: hooks are read-only after construction; no hooks add/remove other hooks mid-combat.
- **Stat balance is a placeholder** → Mitigation: animal stat tables are defined in one place (`combat.go` or `animal.go`) and trivially tunable.
- **No HP persistence means no tension between fights** → Accepted for now; persistent HP is a follow-up change.

## Open Questions

- Should combat be triggerable by walking into an animal, or only by explicit `g` press? (Deferred — `g` press for now.)
- Should the combat log be limited in length or scroll? (Limit to last N lines that fit the viewport; no scroll in v1.)
