## Context

The current combat engine in `combat.go` uses armour as a flat damage reducer: `rollDamage` computes `max(0, raw - defender.Armour)` and applies the result directly to HP. The `Armour` field on `Combatant` is a static constant throughout a fight — it never changes. There is no concept of armour depletion.

`buildPlayerCombatant` hard-codes `Armour: 0, MinDamage: 1, MaxDamage: 3` and does not consult equipped items. The `Item` struct has no bonus fields.

`NewModel()` pre-equips Cloth Tunic, Cloth Pants, and Leather Boots but leaves both hand slots empty.

The render layer tracks live HP during playback by scraping damage patterns from log lines (`hpAtRound`). Any armour display during playback will need the same treatment.

## Goals / Non-Goals

**Goals:**
- Armour becomes a depletable combat buffer: damage drains `CurrentArmour` first; only overflow reaches HP
- `CurrentArmour` resets to `Armour` (max) at the start of each `resolveCombat` call (i.e., each enemy encounter)
- Enemies' armour also behaves as a depletable buffer within a single encounter
- `Item` gains `ArmourBonus int` and `DamageBonus int` fields
- `buildPlayerCombatant` sums bonuses from all equipped items
- `NewModel()` pre-equips Wooden Sword (`SlotRightHand`, `DamageBonus: 1`) and Wooden Shield (`SlotLeftHand`, `ArmourBonus: 1`)
- Armour bar in the combat UI reflects live current armour during playback
- ≥90% test coverage maintained

**Non-Goals:**
- Per-round armour regeneration or partial recovery
- Enemy item equipment or enemy-side armour bonuses from items
- Armour type distinctions (no light/heavy categories)
- Durability or permanent item degradation

## Decisions

### 1. `CurrentArmour` field on `Combatant`

**Decision**: Add `CurrentArmour int` to `Combatant`. `resolveCombat` sets `CurrentArmour = Armour` (the max) on both combatants at the top of the function, before the round loop. Damage resolution mutates `CurrentArmour` directly on the pointer already in use.

**Alternatives considered**:
- *Separate `currentArmour` map keyed by combatant name*: avoids struct change but is more complex and fragile.
- *Track in `CombatState` as separate slices*: would break the existing `Combatant` pointer pattern in `resolveCombat`.

**Rationale**: `Combatant` is already passed and mutated by pointer in `resolveCombat`. Adding `CurrentArmour` keeps the mutation pattern consistent with the existing `HP` drain.

### 2. Damage resolution: absorb then overflow

**Decision**: Replace `rollDamage` with a new function `applyDamage(raw int, defender *Combatant)` that:
1. If `defender.CurrentArmour > 0`: subtract `raw` from `CurrentArmour`. If `CurrentArmour` goes negative, the absolute value is the HP overflow; clamp `CurrentArmour` to 0 and subtract overflow from `defender.HP`.
2. If `defender.CurrentArmour == 0`: subtract `raw` from `defender.HP` directly.
3. Returns `(armourDrain int, hpDrain int)` for logging.

The old `rollDamage` is replaced entirely. `raw` is still computed as `rng.Intn(attacker.MaxDamage-attacker.MinDamage+1) + attacker.MinDamage` (minimum 0 if `MinDamage == MaxDamage == 0`).

**Note**: Armour no longer reduces raw damage by a flat amount — it absorbs it as a pool instead. This is a **breaking change** to combat balance. Enemies with `Armour: 2` previously reduced every hit by 2; now they absorb the first 2 total damage. This makes early enemies easier. Values may need tuning post-implementation.

**Alternatives considered**:
- *Keep flat reduction and add a separate pool*: hybrid model adds complexity without matching the user's stated intent ("subtract from armour until 0, then from HP").

### 3. Log format

**Decision**: Update log lines to include current armour and HP so the render layer can reconstruct both values during playback:

- Armour-only hit: `"Round %d: %s attacks %s — absorbs %d (Armour: %d)"`
- Partial penetration (armour depleted + HP damaged): `"Round %d: %s attacks %s — armour broken, %d HP damage (%d HP, 0 Armour)"`
- Pure HP hit (armour already 0): `"Round %d: %s attacks %s for %d damage (%d HP left)"`
- Zero-damage hit (raw roll was 0): `"Round %d: %s attacks %s but deals no damage"`

The existing scraper `hpAtRound` in `render.go` must be updated or supplemented with `armourAtRound` that parses the new formats.

**Alternatives considered**:
- *Embed JSON structs in log lines*: clean for parsing but ugly in the UI log panel.
- *Store structured events in `CombatState.Log` as `[]LogEvent`*: cleaner long-term but a larger refactor than this change warrants.

### 4. Item bonus fields

**Decision**: Add `ArmourBonus int` and `DamageBonus int` to the `Item` struct in `types.go`. `buildPlayerCombatant` iterates `m.inventory.Equipped[:]` and accumulates totals, adding them to the base `Armour` and `MaxDamage`/`MinDamage` fields of the returned `Combatant`. `DamageBonus` is added to both `MinDamage` and `MaxDamage` to scale the whole damage range.

**Rationale**: Flat bonus fields are the simplest extension that satisfies the stated items. The existing `RoundHook` mechanism remains for future items with complex conditional effects.

### 5. Starting items

**Decision**: In `NewModel()`, add to `inventory.Equipped`:
- `Equipped[SlotRightHand]`: `Item{Name: "Wooden Sword", Slots: []BodySlot{SlotRightHand}, DamageBonus: 1}`
- `Equipped[SlotLeftHand]`: `Item{Name: "Wooden Shield", Slots: []BodySlot{SlotLeftHand}, ArmourBonus: 1}`

These are pre-equipped (not in `inventory.Items`), consistent with how the default outfit is handled.

### 6. Armour bar playback in the render layer

**Decision**: Add `armourAtRound(log []string, roundIndex int, combatantName string, startArmour int) int` alongside the existing `hpAtRound`. The function scans log lines up to `roundIndex` for the new armour-bearing patterns and tracks current armour. The armour bar in `renderHeroPanel` and the enemy panel uses this value when `combatLogIndex > 0`.

## Risks / Trade-offs

- **Combat balance regression** — the absorb-pool model changes effective survivability. [Risk] → Document as intentional; values are data (`animalStats` map, dungeon enemy templates) and can be tuned independently.
- **Log format is now a parsing contract** — `hpAtRound` and `armourAtRound` depend on string patterns. [Risk] → Keep formats in named constants or test them explicitly so regressions are caught.
- **Existing render tests expect old log strings** — `render_test.go` constructs `CombatState.Log` manually. [Risk] → Update all test fixtures to use new log format strings.
- **`CurrentArmour` in `CombatState.Player/Enemy` snapshots** — `CombatState` captures the final `Combatant` state. After combat, `CurrentArmour` will be 0 or whatever remains. This is fine for display (post-combat summary unused) but tests should be aware.

## Migration Plan

1. Add `CurrentArmour int` to `Combatant`; add `ArmourBonus int` / `DamageBonus int` to `Item`
2. Update `buildPlayerCombatant` to sum item bonuses
3. Update `NewModel()` to equip Wooden Sword and Wooden Shield
4. Replace `rollDamage` with `applyDamage`; update `resolveCombat` to initialise `CurrentArmour` and use new log formats
5. Add `armourAtRound` to `render.go`; wire armour bar in `renderHeroPanel` and enemy panel
6. Update all tests (combat logic + render snapshots)
7. No database migrations, no external API changes; entirely internal

## Open Questions

- Should `DamageBonus` increase only `MaxDamage` (increasing variance) or both `MinDamage` and `MaxDamage` (scaling the whole range)? **Decided**: both, to keep the range shape intact.
- Should the Wooden Sword and Shield also appear as obtainable loot, or only as starting equipment for now? **For now**: starting equipment only; loot tables are out of scope.
