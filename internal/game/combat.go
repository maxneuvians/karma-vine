package game

import (
	"fmt"
	"math/rand"
)

// ── Animal combat stats ─────────────────────────────────────────────────────

// AnimalCombatStats holds default combat values for an animal type.
type AnimalCombatStats struct {
	HP        int
	Armour    int
	MinDamage int
	MaxDamage int
	Initiative int
}

// animalStatTable maps animal names to their combat stats.
var animalStatTable = map[string]AnimalCombatStats{
	"Bear":     {HP: 18, Armour: 2, MinDamage: 3, MaxDamage: 6, Initiative: 3},
	"Wolf":     {HP: 12, Armour: 1, MinDamage: 2, MaxDamage: 5, Initiative: 6},
	"Deer":     {HP: 8, Armour: 0, MinDamage: 1, MaxDamage: 2, Initiative: 7},
	"Rabbit":   {HP: 3, Armour: 0, MinDamage: 0, MaxDamage: 1, Initiative: 8},
	"Snake":    {HP: 6, Armour: 0, MinDamage: 2, MaxDamage: 4, Initiative: 5},
	"Lizard":   {HP: 4, Armour: 1, MinDamage: 1, MaxDamage: 2, Initiative: 4},
	"Bird":     {HP: 4, Armour: 0, MinDamage: 1, MaxDamage: 2, Initiative: 9},
	"Crab":     {HP: 5, Armour: 2, MinDamage: 1, MaxDamage: 3, Initiative: 2},
	"Seagull":  {HP: 4, Armour: 0, MinDamage: 1, MaxDamage: 2, Initiative: 7},
	"Goat":     {HP: 10, Armour: 1, MinDamage: 1, MaxDamage: 3, Initiative: 4},
	"Eagle":    {HP: 8, Armour: 0, MinDamage: 2, MaxDamage: 4, Initiative: 8},
	"Antelope": {HP: 9, Armour: 0, MinDamage: 1, MaxDamage: 3, Initiative: 8},
	"default":  {HP: 5, Armour: 0, MinDamage: 1, MaxDamage: 2, Initiative: 3},
}

// buildEnemyCombatant constructs a Combatant from an Animal using the stat table.
func buildEnemyCombatant(a Animal) Combatant {
	stats, ok := animalStatTable[a.Name]
	if !ok {
		stats = animalStatTable["default"]
	}
	return Combatant{
		Name:      a.Name,
		HP:        stats.HP,
		MaxHP:     stats.HP,
		Armour:    stats.Armour,
		MinDamage: stats.MinDamage,
		MaxDamage: stats.MaxDamage,
		Initiative: stats.Initiative,
	}
}

// ── Player combatant builder ────────────────────────────────────────────────

// buildPlayerCombatant constructs the player's Combatant from base stats + equipment.
func buildPlayerCombatant(m Model) Combatant {
	c := Combatant{
		Name:      "Player",
		HP:        m.playerHP,
		MaxHP:     m.playerMaxHP,
		Armour:    0,
		MinDamage: 1,
		MaxDamage: 3,
		Initiative: 5,
	}
	// Sum equipment bonuses (all zero in this iteration; structure in place for future).
	return c
}

// buildCombatHooks collects between-round hooks from equipped items.
// Returns an empty slice for all current items.
func buildCombatHooks(m Model) []RoundHook {
	return nil
}

// ── Combat resolution ───────────────────────────────────────────────────────

// resolveCombat runs a complete auto-battle between player and enemy.
// The original Combatant values are not mutated — the function works on copies.
func resolveCombat(player, enemy Combatant, hooks []RoundHook, rng *rand.Rand) CombatState {
	// Work on copies so originals are untouched.
	p := player
	e := enemy

	state := CombatState{
		Player: p,
		Enemy:  e,
		Hooks:  hooks,
	}

	// Determine initiative order: higher goes first; ties favour player.
	playerFirst := p.Initiative >= e.Initiative

	for {
		state.Round++

		var first, second *Combatant
		var firstName, secondName string
		if playerFirst {
			first, second = &state.Player, &state.Enemy
			firstName, secondName = p.Name, e.Name
		} else {
			first, second = &state.Enemy, &state.Player
			firstName, secondName = e.Name, p.Name
		}

		// Run hooks before attacks.
		for _, hook := range hooks {
			hook(first, second)
		}

		// Check if hooks killed either combatant.
		if second.HP <= 0 {
			state.Log = append(state.Log, fmt.Sprintf("%s is defeated!", secondName))
			break
		}
		if first.HP <= 0 {
			state.Log = append(state.Log, fmt.Sprintf("%s is defeated!", firstName))
			break
		}

		// First combatant attacks.
		damage := rollDamage(first, second, rng)
		second.HP -= damage
		if damage > 0 {
			state.Log = append(state.Log, fmt.Sprintf("Round %d: %s attacks %s for %d damage (%d HP left)",
				state.Round, firstName, secondName, damage, second.HP))
		} else {
			state.Log = append(state.Log, fmt.Sprintf("Round %d: %s attacks %s but armour absorbs all damage (%d HP left)",
				state.Round, firstName, secondName, second.HP))
		}

		if second.HP <= 0 {
			state.Log = append(state.Log, fmt.Sprintf("%s is defeated!", secondName))
			break
		}

		// Second combatant attacks.
		damage = rollDamage(second, first, rng)
		first.HP -= damage
		if damage > 0 {
			state.Log = append(state.Log, fmt.Sprintf("Round %d: %s attacks %s for %d damage (%d HP left)",
				state.Round, secondName, firstName, damage, first.HP))
		} else {
			state.Log = append(state.Log, fmt.Sprintf("Round %d: %s attacks %s but armour absorbs all damage (%d HP left)",
				state.Round, secondName, firstName, first.HP))
		}

		if first.HP <= 0 {
			state.Log = append(state.Log, fmt.Sprintf("%s is defeated!", firstName))
			break
		}
	}

	state.PlayerWon = state.Enemy.HP <= 0
	return state
}

// rollDamage calculates damage from attacker to defender.
func rollDamage(attacker, defender *Combatant, rng *rand.Rand) int {
	raw := rng.Intn(attacker.MaxDamage-attacker.MinDamage+1) + attacker.MinDamage
	damage := raw - defender.Armour
	if damage < 0 {
		damage = 0
	}
	return damage
}
