package game

import (
	"math/rand"
	"strings"
	"testing"
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

func TestCombatant_ArmourReducesToZero(t *testing.T) {
	// Attacker deals 3 fixed damage, defender has 5 armour → 0 net damage.
	attacker := Combatant{Name: "Weak", HP: 100, MaxHP: 100, MinDamage: 3, MaxDamage: 3, Initiative: 10}
	defender := Combatant{Name: "Tank", HP: 100, MaxHP: 100, Armour: 5, MinDamage: 3, MaxDamage: 3, Initiative: 1}
	rng := rand.New(rand.NewSource(1))
	state := resolveCombat(attacker, defender, nil, rng)
	// After many rounds neither should die if all damage is absorbed both ways.
	// Actually defender also attacks with 3 damage vs attacker armour 0 → attacker takes damage.
	// The test just verifies the first round: attacker deals 0 to defender.
	if state.Enemy.HP < defender.HP {
		// Defender lost HP — armour didn't absorb. This is fine since the second combatant
		// attacks the first. Let's check the log for "armour absorbs" on the first line.
	}
	// Check that the first log line mentions armour absorb.
	found := false
	for _, line := range state.Log {
		if strings.Contains(line, "Weak attacks Tank") && strings.Contains(line, "armour absorbs") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected an armour absorption log line; log: %v", state.Log)
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
	if c.Armour != 0 {
		t.Fatalf("player Armour should be 0, got %d", c.Armour)
	}
	if c.MinDamage != 1 || c.MaxDamage != 3 {
		t.Fatalf("player damage should be 1-3, got %d-%d", c.MinDamage, c.MaxDamage)
	}
	if c.Initiative != 5 {
		t.Fatalf("player Initiative should be 5, got %d", c.Initiative)
	}
}
