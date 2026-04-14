## 1. Types

- [x] 1.1 Add `ScreenCombat` constant to the `ScreenMode` iota in `types.go`
- [x] 1.2 Add `Combatant` struct (`Name string`, `HP int`, `MaxHP int`, `Armour int`, `MinDamage int`, `MaxDamage int`, `Initiative int`) to `types.go`
- [x] 1.3 Add `type RoundHook func(self, opponent *Combatant)` to `types.go`
- [x] 1.4 Add `CombatState` struct (`Player Combatant`, `Enemy Combatant`, `Hooks []RoundHook`, `Log []string`, `Round int`, `PlayerWon bool`) to `types.go`

## 2. Animal Combat Stats

- [x] 2.1 Add `AnimalCombatStats` struct (`HP`, `Armour`, `MinDamage`, `MaxDamage`, `Initiative int`) in `animal.go` (or new `combat.go`)
- [x] 2.2 Define `animalStatTable map[string]AnimalCombatStats` with entries for all existing animal names (Bear, Wolf, Deer, etc.) and a `"default"` fallback (`HP=5`, `Armour=0`, `MinDamage=1`, `MaxDamage=2`, `Initiative=3`)
- [x] 2.3 Implement `buildEnemyCombatant(a Animal) Combatant` that looks up `a.Name` in `animalStatTable` (with default fallback) and returns a populated `Combatant`

## 3. Player Combatant Builder

- [x] 3.1 Implement `buildPlayerCombatant(m Model) Combatant` in `combat.go`: base stats `HP=MaxHP=20`, `Armour=0`, `MinDamage=1`, `MaxDamage=3`, `Initiative=5`, summing bonus fields from each equipped item (all zero in this iteration)
- [x] 3.2 Implement `buildCombatHooks(m Model) []RoundHook`: iterate equipped items, collect any hooks (returns empty slice for all current items)

## 4. Combat Resolution

- [x] 4.1 Create `internal/game/combat.go` with `resolveCombat(player, enemy Combatant, hooks []RoundHook, rng *rand.Rand) CombatState`
- [x] 4.2 Inside `resolveCombat`: determine initiative order (higher goes first; ties favour player), loop rounds until either HP ≤ 0
- [x] 4.3 Each round: call all hooks with `(attacker, defender)` pointers; check if hooks killed the defender (end combat); resolve attacker attack: `damage = max(0, rng.Intn(MaxDamage-MinDamage+1)+MinDamage - defender.Armour)`; apply to defender HP; append log line; repeat for the other combatant if still alive
- [x] 4.4 Set `CombatState.PlayerWon = enemy.HP <= 0` before returning
- [x] 4.5 Ensure original `player` and `enemy` values passed in are not mutated (work on copies)

## 5. Model

- [x] 5.1 Add `combatState *CombatState` field to `Model` in `model.go`
- [x] 5.2 Add `combatEnemy *Animal` field to `Model` to remember which animal initiated combat (for removal on victory)

## 6. Input — Initiate Combat

- [x] 6.1 In `handleKey`, locate the `g` (pick up) branch for `ModeLocal`
- [x] 6.2 Before the existing item-pickup logic, check if `localMap.Animals` contains an animal at `m.playerPos`; if so, build combatants, call `resolveCombat`, store result in `m.combatState`, store animal reference in `m.combatEnemy`, set `m.screenMode = ScreenCombat`, `m.paused = true`, and return immediately
- [x] 6.3 If no animal is present, fall through to existing item-pickup logic unchanged

## 7. Input — Dismiss Combat Result

- [x] 7.1 In `handleKey`, add an early branch for `m.screenMode == ScreenCombat` that suppresses all keys except `enter`, `space`, `q`, `ctrl+c`
- [x] 7.2 On `enter`/`space` with `m.combatState.PlayerWon == true`: set `m.screenMode = ScreenNormal`, `m.paused = false`, remove `m.combatEnemy` from `localMap.Animals`, clear `m.combatState` and `m.combatEnemy`
- [x] 7.3 On `enter`/`space` with `m.combatState.PlayerWon == false`: return `m, tea.Quit`

## 8. Rendering

- [x] 8.1 In `buildView` (`render.go`), add an early-return for `m.screenMode == ScreenCombat` that returns `renderCombatScreen(m)`, placed before the `ScreenInventory` check
- [x] 8.2 Implement `renderCombatScreen(m Model) string`: render a `"⚔ Combat"` header, two side-by-side stat panels (player left, enemy right) each showing Name, HP/MaxHP, Armour, Damage range, Initiative
- [x] 8.3 Below the stat panels, render the combat log lines (most-recent-at-bottom, truncated to fit available rows)
- [x] 8.4 When combat is over (`m.combatState.PlayerWon` set), render a `"Victory!"` or `"Defeated!"` banner with the appropriate dismiss hint

## 9. Tests

- [x] 9.1 `TestCombatant_InitiativeOrder` in `combat_test.go`: player Initiative > enemy Initiative → first log line describes player's attack
- [x] 9.2 `TestCombatant_ArmourReducesToZero` in `combat_test.go`: attacker damage ≤ defender Armour → HP unchanged
- [x] 9.3 `TestResolveCombat_PlayerWins` in `combat_test.go`: seeded RNG, player clearly stronger → `PlayerWon == true`
- [x] 9.4 `TestResolveCombat_EnemyWins` in `combat_test.go`: seeded RNG, enemy clearly stronger → `PlayerWon == false`
- [x] 9.5 `TestResolveCombat_NoMutation` in `combat_test.go`: verify original `Combatant` args are unchanged after call
- [x] 9.6 `TestRoundHook_FiresBeforeAttack` in `combat_test.go`: hook that records call order; assert called before damage applied
- [x] 9.7 `TestRoundHook_KillsDefenderEndsRound` in `combat_test.go`: hook sets `defender.HP = 0`; assert no attack log entry for that round
- [x] 9.8 `TestBuildEnemyCombatant_KnownAnimal` in `animal_test.go` (or `combat_test.go`): known name returns correct stats
- [x] 9.9 `TestBuildEnemyCombatant_UnknownAnimal` in same file: unknown name returns fallback stats
- [x] 9.10 `TestBuildPlayerCombatant_BaseStats` in `combat_test.go`: default model → HP=20, Armour=0, MinDamage=1, MaxDamage=3, Initiative=5
- [x] 9.11 `TestHandleKey_GOnAnimalInitiatesCombat` in `input_test.go`: `g` with animal at player pos → `screenMode == ScreenCombat`
- [x] 9.12 `TestHandleKey_GOnEmptyCellSkipsCombat` in `input_test.go`: `g` with no animal → `screenMode` unchanged
- [x] 9.13 `TestHandleKey_CombatScreenSuppressesMovement` in `input_test.go`: `up` in `ScreenCombat` → `playerPos` unchanged
- [x] 9.14 `TestHandleKey_EnterDismissesVictory` in `input_test.go`: `enter` with `PlayerWon == true` → `screenMode == ScreenNormal`
- [x] 9.15 `TestHandleKey_EnterQuitsOnDefeat` in `input_test.go`: `enter` with `PlayerWon == false` → cmd is `tea.Quit`
- [x] 9.16 `TestRenderCombatScreen_ShowsStats` in `render_test.go`: assert player and enemy name, HP, armour present in output
- [x] 9.17 `TestRenderCombatScreen_VictoryBanner` in `render_test.go`: `PlayerWon == true` → output contains `"Victory!"`
- [x] 9.18 `TestRenderCombatScreen_DefeatBanner` in `render_test.go`: `PlayerWon == false` → output contains `"Defeated!"`
- [x] 9.19 `TestBuildView_ScreenCombat` in `render_test.go`: `buildView` with `ScreenCombat` → output equals `renderCombatScreen(m)`
