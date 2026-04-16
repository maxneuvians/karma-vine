package game

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

// ── Initiative order ─────────────────────────────────────────────────────────

func TestCombatant_InitiativeOrder(t *testing.T) {
	player := Combatant{Name: "Player", HP: 20, MaxHP: 20, MinDamage: 2, MaxDamage: 3, Initiative: 10}
	enemy := Combatant{Name: "Wolf", HP: 12, MaxHP: 12, MinDamage: 2, MaxDamage: 4, Initiative: 5}
	rng := rand.New(rand.NewSource(42))
	state := resolveCombat(player, enemy, nil, rng)
	if len(state.Log) == 0 {
		t.Fatal("combat log should not be empty")
	}
	if !strings.Contains(state.Log[0], "Player attacks") {
		t.Errorf("first log line should describe player's attack, got: %s", state.Log[0])
	}
}

// ── Armour absorbs damage ───────────────────────────────────────────────────

func TestCombatant_ArmourAbsorbsBeforeHP(t *testing.T) {
	// Attacker deals 3 fixed damage; defender has 5 armour → armour absorbs, no HP loss.
	attacker := Combatant{Name: "Weak", HP: 100, MaxHP: 100, MinDamage: 3, MaxDamage: 3, Initiative: 10}
	defender := Combatant{Name: "Tank", HP: 100, MaxHP: 100, Armour: 5, MinDamage: 3, MaxDamage: 3, Initiative: 1}
	rng := rand.New(rand.NewSource(1))
	state := resolveCombat(attacker, defender, nil, rng)
	// Verify the first attack log line shows armour absorbing.
	found := false
	for _, line := range state.Log {
		if strings.Contains(line, "Weak attacks Tank") && strings.Contains(line, "absorbs") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected an armour-absorb log line; log: %v", state.Log)
	}
}

// ── Player wins ──────────────────────────────────────────────────────────────

func TestResolveCombat_PlayerWins(t *testing.T) {
	player := Combatant{Name: "Player", HP: 100, MaxHP: 100, MinDamage: 10, MaxDamage: 10, Initiative: 10}
	enemy := Combatant{Name: "Rat", HP: 5, MaxHP: 5, MinDamage: 1, MaxDamage: 1, Initiative: 1}
	rng := rand.New(rand.NewSource(0))
	state := resolveCombat(player, enemy, nil, rng)
	if !state.PlayerWon {
		t.Fatal("player should win against weak enemy")
	}
}

// ── Enemy wins ───────────────────────────────────────────────────────────────

func TestResolveCombat_EnemyWins(t *testing.T) {
	player := Combatant{Name: "Player", HP: 5, MaxHP: 5, MinDamage: 1, MaxDamage: 1, Initiative: 1}
	enemy := Combatant{Name: "Dragon", HP: 100, MaxHP: 100, MinDamage: 20, MaxDamage: 20, Initiative: 10}
	rng := rand.New(rand.NewSource(0))
	state := resolveCombat(player, enemy, nil, rng)
	if state.PlayerWon {
		t.Fatal("enemy should win against weak player")
	}
}

// ── No mutation ──────────────────────────────────────────────────────────────

func TestResolveCombat_NoMutation(t *testing.T) {
	player := Combatant{Name: "Player", HP: 20, MaxHP: 20, MinDamage: 3, MaxDamage: 5, Initiative: 5}
	enemy := Combatant{Name: "Wolf", HP: 12, MaxHP: 12, MinDamage: 2, MaxDamage: 4, Initiative: 6}
	origPlayerHP := player.HP
	origEnemyHP := enemy.HP
	rng := rand.New(rand.NewSource(99))
	_ = resolveCombat(player, enemy, nil, rng)
	if player.HP != origPlayerHP {
		t.Fatalf("player HP mutated: %d → %d", origPlayerHP, player.HP)
	}
	if enemy.HP != origEnemyHP {
		t.Fatalf("enemy HP mutated: %d → %d", origEnemyHP, enemy.HP)
	}
}

// ── Round hooks ──────────────────────────────────────────────────────────────

func TestRoundHook_FiresBeforeAttack(t *testing.T) {
	var hookCalled bool
	hook := func(self, opponent *Combatant) {
		hookCalled = true
	}
	player := Combatant{Name: "Player", HP: 100, MaxHP: 100, MinDamage: 50, MaxDamage: 50, Initiative: 10}
	enemy := Combatant{Name: "Rat", HP: 5, MaxHP: 5, MinDamage: 1, MaxDamage: 1, Initiative: 1}
	rng := rand.New(rand.NewSource(0))
	_ = resolveCombat(player, enemy, []RoundHook{hook}, rng)
	if !hookCalled {
		t.Fatal("hook should have been called")
	}
}

func TestRoundHook_KillsDefenderEndsRound(t *testing.T) {
	// Hook kills the defender before attacks happen.
	hook := func(self, opponent *Combatant) {
		opponent.HP = 0
	}
	player := Combatant{Name: "Player", HP: 20, MaxHP: 20, MinDamage: 1, MaxDamage: 1, Initiative: 10}
	enemy := Combatant{Name: "Rat", HP: 10, MaxHP: 10, MinDamage: 1, MaxDamage: 1, Initiative: 1}
	rng := rand.New(rand.NewSource(0))
	state := resolveCombat(player, enemy, []RoundHook{hook}, rng)
	// Should end immediately; no attack log lines for round 1.
	for _, line := range state.Log {
		if strings.Contains(line, "attacks") {
			t.Fatalf("no attack should happen when hook kills defender; got: %s", line)
		}
	}
	if !state.PlayerWon {
		t.Fatal("player should win when hook kills enemy")
	}
}

// ── Build enemy combatant ────────────────────────────────────────────────────

func TestBuildEnemyCombatant_KnownAnimal(t *testing.T) {
	a := Animal{Name: "Wolf", Char: 'w', Color: "#555"}
	c := buildEnemyCombatant(a)
	if c.Name != "Wolf" {
		t.Fatalf("expected name Wolf, got %q", c.Name)
	}
	if c.HP != 12 {
		t.Fatalf("Wolf HP should be 12, got %d", c.HP)
	}
	if c.Initiative != 6 {
		t.Fatalf("Wolf Initiative should be 6, got %d", c.Initiative)
	}
}

func TestBuildEnemyCombatant_UnknownAnimal(t *testing.T) {
	a := Animal{Name: "UnknownBeast", Char: '?', Color: "#000"}
	c := buildEnemyCombatant(a)
	if c.Name != "UnknownBeast" {
		t.Fatalf("expected name UnknownBeast, got %q", c.Name)
	}
	if c.HP != 5 || c.Armour != 0 || c.MinDamage != 1 || c.MaxDamage != 2 || c.Initiative != 3 {
		t.Fatalf("fallback stats wrong: %+v", c)
	}
}

// ── Build player combatant ───────────────────────────────────────────────────

func TestBuildPlayerCombatant_BaseStats(t *testing.T) {
	m := NewModel()
	c := buildPlayerCombatant(m)
	if c.HP != 20 || c.MaxHP != 20 {
		t.Fatalf("player HP should be 20, got %d/%d", c.HP, c.MaxHP)
	}
	// Default outfit: Cloth Tunic (+1), Cloth Pants (+1), Leather Boots (+1), Wooden Shield (+1) = 4 armour.
	if c.Armour != 4 {
		t.Fatalf("player Armour should be 4 (tunic+pants+boots+shield), got %d", c.Armour)
	}
	if c.MinDamage != 2 || c.MaxDamage != 4 {
		t.Fatalf("player damage should be 2-4 (sword bonus), got %d-%d", c.MinDamage, c.MaxDamage)
	}
	if c.Initiative != 5 {
		t.Fatalf("player Initiative should be 5, got %d", c.Initiative)
	}
}

func TestBuildPlayerCombatant_UsesModelHP(t *testing.T) {
	m := NewModel()
	m.playerHP = 15
	m.playerMaxHP = 20
	c := buildPlayerCombatant(m)
	if c.HP != 15 {
		t.Fatalf("buildPlayerCombatant should use m.playerHP: expected 15, got %d", c.HP)
	}
	if c.MaxHP != 20 {
		t.Fatalf("buildPlayerCombatant should use m.playerMaxHP: expected 20, got %d", c.MaxHP)
	}
}

func TestBuildDungeonEnemyCombatant(t *testing.T) {
	tmpl := &EnemyTemplate{
		Name: "Goblin", Char: 'g', Color: "#55aa44",
		BaseHP: 10, MaxHP: 20,
	}
	e := &DungeonEnemy{
		X: 5, Y: 5, Template: tmpl,
		HP: 12, MaxHP: 15, Armour: 2,
		MinDamage: 3, MaxDamage: 6, Initiative: 5,
	}
	c := buildDungeonEnemyCombatant(e)
	if c.Name != "Goblin" {
		t.Fatalf("expected Name 'Goblin', got %q", c.Name)
	}
	if c.HP != 12 {
		t.Fatalf("expected HP 12, got %d", c.HP)
	}
	if c.MaxHP != 15 {
		t.Fatalf("expected MaxHP 15, got %d", c.MaxHP)
	}
	if c.Armour != 2 {
		t.Fatalf("expected Armour 2, got %d", c.Armour)
	}
	if c.MinDamage != 3 {
		t.Fatalf("expected MinDamage 3, got %d", c.MinDamage)
	}
	if c.MaxDamage != 6 {
		t.Fatalf("expected MaxDamage 6, got %d", c.MaxDamage)
	}
	if c.Initiative != 5 {
		t.Fatalf("expected Initiative 5, got %d", c.Initiative)
	}
}

// ── Combat speed duration tests ─────────────────────────────────────────────

func TestCombatSpeedDuration_Slow(t *testing.T) {
	if d := combatSpeedDuration(CombatSpeedSlow); d != 3*time.Second {
		t.Fatalf("expected 3s, got %v", d)
	}
}

func TestCombatSpeedDuration_Fast(t *testing.T) {
	if d := combatSpeedDuration(CombatSpeedFast); d != 200*time.Millisecond {
		t.Fatalf("expected 200ms, got %v", d)
	}
}

func TestCombatSpeedDuration_OutOfRange(t *testing.T) {
	if d := combatSpeedDuration(99); d != 1*time.Second {
		t.Fatalf("expected 1s, got %v", d)
	}
}

// ── Combat log grouping tests ───────────────────────────────────────────────

func TestCombatLogLinesUpTo_ReturnsRound1Only(t *testing.T) {
	log := []string{
		"Round 1: Player attacks Wolf for 3 damage (9 HP left)",
		"Round 1: Wolf attacks Player for 2 damage (18 HP left)",
		"Round 2: Player attacks Wolf for 4 damage (5 HP left)",
		"Round 2: Wolf attacks Player for 1 damage (17 HP left)",
	}
	result := combatLogLinesUpTo(log, 1)
	if len(result) != 2 {
		t.Fatalf("expected 2 lines for round 1, got %d: %v", len(result), result)
	}
}

func TestCombatLogLinesUpTo_ZeroIndex(t *testing.T) {
	log := []string{
		"Round 1: Player attacks Wolf for 3 damage (9 HP left)",
	}
	result := combatLogLinesUpTo(log, 0)
	if len(result) != 0 {
		t.Fatalf("expected 0 lines for roundIndex=0, got %d", len(result))
	}
}

// ── hpAtRound tests ─────────────────────────────────────────────────────────

func TestHpAtRound_NoDamage(t *testing.T) {
	log := []string{
		"Round 1: Player attacks Wolf for 3 damage (9 HP left)",
	}
	hp := hpAtRound(20, log, 1, "Player")
	if hp != 20 {
		t.Fatalf("expected 20 (no damage to Player), got %d", hp)
	}
}

func TestHpAtRound_TakesOneDamage(t *testing.T) {
	log := []string{
		"Round 1: Wolf attacks Hero for 5 damage (15 HP left)",
	}
	hp := hpAtRound(20, log, 1, "Hero")
	if hp != 15 {
		t.Fatalf("expected 15, got %d", hp)
	}
}

// ── applyDamage tests ────────────────────────────────────────────────────────

func TestApplyDamage_FullAbsorption(t *testing.T) {
	d := &Combatant{HP: 10, CurrentArmour: 3}
	ad, hd := applyDamage(2, d)
	if ad != 2 || hd != 0 {
		t.Fatalf("expected armourDrain=2 hpDrain=0, got %d %d", ad, hd)
	}
	if d.CurrentArmour != 1 {
		t.Fatalf("expected CurrentArmour=1, got %d", d.CurrentArmour)
	}
	if d.HP != 10 {
		t.Fatalf("expected HP unchanged=10, got %d", d.HP)
	}
}

func TestApplyDamage_ArmourBrokenWithOverflow(t *testing.T) {
	d := &Combatant{HP: 10, CurrentArmour: 2}
	ad, hd := applyDamage(5, d)
	if ad != 2 || hd != 3 {
		t.Fatalf("expected armourDrain=2 hpDrain=3, got %d %d", ad, hd)
	}
	if d.CurrentArmour != 0 {
		t.Fatalf("expected CurrentArmour=0, got %d", d.CurrentArmour)
	}
	if d.HP != 7 {
		t.Fatalf("expected HP=7, got %d", d.HP)
	}
}

func TestApplyDamage_DirectHPWhenNoArmour(t *testing.T) {
	d := &Combatant{HP: 10, CurrentArmour: 0}
	ad, hd := applyDamage(4, d)
	if ad != 0 || hd != 4 {
		t.Fatalf("expected armourDrain=0 hpDrain=4, got %d %d", ad, hd)
	}
	if d.HP != 6 {
		t.Fatalf("expected HP=6, got %d", d.HP)
	}
}

func TestApplyDamage_ZeroRawIsNoop(t *testing.T) {
	d := &Combatant{HP: 10, CurrentArmour: 3}
	ad, hd := applyDamage(0, d)
	if ad != 0 || hd != 0 {
		t.Fatalf("expected 0,0 got %d %d", ad, hd)
	}
	if d.HP != 10 || d.CurrentArmour != 3 {
		t.Fatalf("state should be unchanged: HP=%d Armour=%d", d.HP, d.CurrentArmour)
	}
}

// ── resolveCombat log format tests ──────────────────────────────────────────

func TestResolveCombat_ArmourOnlyHitLogLine(t *testing.T) {
	// Player has high enough armour to fully absorb the enemy's fixed 1 damage hit.
	// Enemy has no armour; player always goes first (higher initiative).
	player := Combatant{Name: "Player", HP: 100, MaxHP: 100, Armour: 10, MinDamage: 1, MaxDamage: 1, Initiative: 10}
	enemy := Combatant{Name: "Rat", HP: 5, MaxHP: 5, Armour: 0, MinDamage: 1, MaxDamage: 1, Initiative: 1}
	rng := rand.New(rand.NewSource(0))
	state := resolveCombat(player, enemy, nil, rng)
	// Find a line where Rat attacks Player and armour absorbs.
	found := false
	for _, line := range state.Log {
		if strings.Contains(line, "Rat attacks Player — absorbs") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected armour-absorb log line; log: %v", state.Log)
	}
}

func TestResolveCombat_ArmourBrokenLogLine(t *testing.T) {
	// Player has 1 armour; enemy deals exactly 3 damage → armour broken + 2 HP overflow.
	player := Combatant{Name: "Player", HP: 20, MaxHP: 20, Armour: 1, MinDamage: 10, MaxDamage: 10, Initiative: 1}
	enemy := Combatant{Name: "Orc", HP: 100, MaxHP: 100, Armour: 0, MinDamage: 3, MaxDamage: 3, Initiative: 10}
	rng := rand.New(rand.NewSource(0))
	state := resolveCombat(player, enemy, nil, rng)
	found := false
	for _, line := range state.Log {
		if strings.Contains(line, "Orc attacks Player — armour broken") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected armour-broken log line; log: %v", state.Log)
	}
}

func TestResolveCombat_NoDamageLogLine(t *testing.T) {
	// Attacker has MinDamage=MaxDamage=0 so raw is always 0.
	player := Combatant{Name: "Player", HP: 20, MaxHP: 20, MinDamage: 0, MaxDamage: 0, Initiative: 10}
	enemy := Combatant{Name: "Ghost", HP: 5, MaxHP: 5, MinDamage: 10, MaxDamage: 10, Initiative: 1}
	rng := rand.New(rand.NewSource(0))
	state := resolveCombat(player, enemy, nil, rng)
	found := false
	for _, line := range state.Log {
		if strings.Contains(line, "Player attacks Ghost but deals no damage") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected no-damage log line; log: %v", state.Log)
	}
}

// ── buildPlayerCombatant bonus tests ────────────────────────────────────────

func TestBuildPlayerCombatant_WoodenSwordBonus(t *testing.T) {
	m := NewModel()
	// NewModel equips Wooden Sword (DamageBonus: 1) and Wooden Shield (ArmourBonus: 1) by default.
	c := buildPlayerCombatant(m)
	if c.MinDamage != 2 {
		t.Fatalf("expected MinDamage=2 (base 1 + sword +1), got %d", c.MinDamage)
	}
	if c.MaxDamage != 4 {
		t.Fatalf("expected MaxDamage=4 (base 3 + sword +1), got %d", c.MaxDamage)
	}
}

func TestBuildPlayerCombatant_WoodenShieldBonus(t *testing.T) {
	m := NewModel()
	c := buildPlayerCombatant(m)
	// Default outfit: Cloth Tunic (+1), Cloth Pants (+1), Leather Boots (+1), Wooden Shield (+1) = 4.
	if c.Armour != 4 {
		t.Fatalf("expected Armour=4 (tunic+pants+boots+shield), got %d", c.Armour)
	}
}

func TestBuildPlayerCombatant_StackingBonuses(t *testing.T) {
	m := NewModel()
	// Add a second DamageBonus item to head slot.
	m.inventory.Equipped[SlotHead] = Item{Name: "Pointy Hat", DamageBonus: 1}
	c := buildPlayerCombatant(m)
	if c.MaxDamage != 5 {
		t.Fatalf("expected MaxDamage=5 (base 3 + sword +1 + hat +1), got %d", c.MaxDamage)
	}
}

func TestBuildPlayerCombatant_EmptySlotsNoBonus(t *testing.T) {
	m := NewModel()
	// Clear all equipped items.
	m.inventory.Equipped = [NumBodySlots]Item{}
	c := buildPlayerCombatant(m)
	if c.Armour != 0 {
		t.Fatalf("expected Armour=0 with empty slots, got %d", c.Armour)
	}
	if c.MinDamage != 1 || c.MaxDamage != 3 {
		t.Fatalf("expected base damage 1-3, got %d-%d", c.MinDamage, c.MaxDamage)
	}
}

// ── armourAtRound tests ─────────────────────────────────────────────────────

func TestArmourAtRound_RoundZeroReturnsStart(t *testing.T) {
	log := []string{
		"Round 1: Enemy attacks Player — absorbs 1 (Armour: 2)",
	}
	a := armourAtRound(log, 0, "Player", 3)
	if a != 3 {
		t.Fatalf("expected startArmour=3 at round 0, got %d", a)
	}
}

func TestArmourAtRound_DecrementsOnAbsorb(t *testing.T) {
	log := []string{
		"Round 1: Enemy attacks Player — absorbs 1 (Armour: 2)",
	}
	a := armourAtRound(log, 1, "Player", 3)
	if a != 2 {
		t.Fatalf("expected armour=2 after absorb, got %d", a)
	}
}

func TestArmourAtRound_ZeroAfterBroken(t *testing.T) {
	log := []string{
		"Round 1: Enemy attacks Player — armour broken, 3 HP damage (7 HP, 0 Armour)",
	}
	a := armourAtRound(log, 1, "Player", 2)
	if a != 0 {
		t.Fatalf("expected armour=0 after broken, got %d", a)
	}
}
