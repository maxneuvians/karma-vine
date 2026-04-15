## 1. Types & Data Model

- [x] 1.1 Add `ArmourBonus int` and `DamageBonus int` fields to the `Item` struct in `types.go`
- [x] 1.2 Add `CurrentArmour int` field to the `Combatant` struct in `types.go`

## 2. Starting Equipment

- [x] 2.1 In `NewModel()` (`model.go`), add `Equipped[SlotRightHand]: Item{Name: "Wooden Sword", Slots: []BodySlot{SlotRightHand}, DamageBonus: 1}` to the default outfit
- [x] 2.2 In `NewModel()`, add `Equipped[SlotLeftHand]: Item{Name: "Wooden Shield", Slots: []BodySlot{SlotLeftHand}, ArmourBonus: 1}` to the default outfit

## 3. Item Bonus Application

- [x] 3.1 In `buildPlayerCombatant` (`combat.go`), iterate `m.inventory.Equipped[:]` and accumulate total `ArmourBonus` and `DamageBonus` across all equipped items
- [x] 3.2 Apply the accumulated bonuses: add `totalArmourBonus` to `Combatant.Armour`, add `totalDamageBonus` to both `Combatant.MinDamage` and `Combatant.MaxDamage`

## 4. Combat Engine — Armour Absorption

- [x] 4.1 In `resolveCombat` (`combat.go`), after building working copies of `player` and `enemy`, set `p.CurrentArmour = p.Armour` and `e.CurrentArmour = e.Armour` before the round loop
- [x] 4.2 Replace `rollDamage` with `applyDamage(raw int, defender *Combatant) (armourDrain, hpDrain int)` implementing the absorb-then-overflow logic
- [x] 4.3 Update the round loop in `resolveCombat` to call `applyDamage` for both attack steps (first→second, second→first) and use the returned `(armourDrain, hpDrain)` for log line selection
- [x] 4.4 Emit armour-only log line when `hpDrain == 0 && armourDrain > 0`: `"Round %d: %s attacks %s — absorbs %d (Armour: %d)"`
- [x] 4.5 Emit armour-broken log line when `armourDrain > 0 && hpDrain > 0`: `"Round %d: %s attacks %s — armour broken, %d HP damage (%d HP, 0 Armour)"`
- [x] 4.6 Emit pure HP log line when `armourDrain == 0 && hpDrain > 0`: `"Round %d: %s attacks %s for %d damage (%d HP left)"` (existing format preserved)
- [x] 4.7 Emit no-damage log line when `armourDrain == 0 && hpDrain == 0`: `"Round %d: %s attacks %s but deals no damage"` (replaces old "armour absorbs all damage" line)
- [x] 4.8 Remove the old `rollDamage` function

## 5. Render Layer — Armour Playback

- [x] 5.1 Add `armourAtRound(log []string, roundIndex int, combatantName string, startArmour int) int` in `render.go` that parses the new log formats to reconstruct current armour at a given round
- [x] 5.2 Update `renderHeroPanel` to call `armourAtRound` and display live current armour in the armour stat line (format: `Armour: X/Y` where X is current, Y is max)
- [x] 5.3 Update the enemy panel rendering to do the same for the enemy armour stat line

## 6. Tests

- [x] 6.1 Write unit tests for `applyDamage` in `combat_test.go`: full absorption, partial overflow, zero armour direct damage, zero raw damage no-op
- [x] 6.2 Write unit tests for `resolveCombat` log format: armour-only hit, armour-broken hit, pure HP hit, no-damage hit
- [x] 6.3 Write unit tests for `buildPlayerCombatant` bonus accumulation: Wooden Sword adds +1 to min/max damage, Wooden Shield adds +1 armour, stacking two bonus items
- [x] 6.4 Write unit tests for `NewModel()` starting equipment: Wooden Sword in `SlotRightHand`, Wooden Shield in `SlotLeftHand`, correct bonus fields
- [x] 6.5 Write unit tests for `armourAtRound`: returns `startArmour` at round 0, decrements on armour-only hits, returns 0 after armour-broken hit
- [x] 6.6 Update existing `render_test.go` fixtures that construct `CombatState.Log` manually to use the new log line formats
- [x] 6.7 Run `go test ./internal/game/... -cover` and confirm coverage remains ≥ 90%
