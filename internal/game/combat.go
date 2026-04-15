package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// combatSpeedDuration returns the tick interval for a given combat playback speed.
func combatSpeedDuration(speed int) time.Duration {
	switch speed {
	case CombatSpeedSlow:
		return 3 * time.Second
	case CombatSpeedNormal:
		return 1 * time.Second
	case CombatSpeedFast:
		return 200 * time.Millisecond
	default:
		return 1 * time.Second
	}
}

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
	var totalArmourBonus, totalDamageBonus int
	for _, item := range m.inventory.Equipped {
		totalArmourBonus += item.ArmourBonus
		totalDamageBonus += item.DamageBonus
	}
	return Combatant{
		Name:       "Player",
		HP:         m.playerHP,
		MaxHP:      m.playerMaxHP,
		Armour:     totalArmourBonus,
		MinDamage:  1 + totalDamageBonus,
		MaxDamage:  3 + totalDamageBonus,
		Initiative: 5,
	}
}

// buildCombatHooks collects between-round hooks from equipped items.
// Returns an empty slice for all current items.
func buildCombatHooks(m Model) []RoundHook {
	return nil
}

// buildDungeonEnemyCombatant constructs a Combatant from a DungeonEnemy.
func buildDungeonEnemyCombatant(e *DungeonEnemy) Combatant {
	return Combatant{
		Name:       e.Template.Name,
		HP:         e.HP,
		MaxHP:      e.MaxHP,
		Armour:     e.Armour,
		MinDamage:  e.MinDamage,
		MaxDamage:  e.MaxDamage,
		Initiative: e.Initiative,
	}
}

// ── Combat resolution ───────────────────────────────────────────────────────

// resolveCombat runs a complete auto-battle between player and enemy.
// The original Combatant values are not mutated — the function works on copies.
func resolveCombat(player, enemy Combatant, hooks []RoundHook, rng *rand.Rand) CombatState {
	// Work on copies so originals are untouched.
	p := player
	e := enemy

	// Initialise current armour to the full pool for this encounter.
	p.CurrentArmour = p.Armour
	e.CurrentArmour = e.Armour

	state := CombatState{
		Player:        p,
		Enemy:         e,
		Hooks:         hooks,
		PlayerStartHP: p.HP,
		EnemyStartHP:  e.HP,
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
		raw := rollRaw(first, rng)
		armourDrain, hpDrain := applyDamage(raw, second)
		state.Log = append(state.Log, combatLogLine(state.Round, firstName, secondName, armourDrain, hpDrain, second))

		if second.HP <= 0 {
			state.Log = append(state.Log, fmt.Sprintf("%s is defeated!", secondName))
			break
		}

		// Second combatant attacks.
		raw = rollRaw(second, rng)
		armourDrain, hpDrain = applyDamage(raw, first)
		state.Log = append(state.Log, combatLogLine(state.Round, secondName, firstName, armourDrain, hpDrain, first))

		if first.HP <= 0 {
			state.Log = append(state.Log, fmt.Sprintf("%s is defeated!", firstName))
			break
		}
	}

	state.PlayerWon = state.Enemy.HP <= 0
	return state
}

// rollRaw rolls the attacker's raw damage before armour is applied.
func rollRaw(attacker *Combatant, rng *rand.Rand) int {
	spread := attacker.MaxDamage - attacker.MinDamage
	if spread < 0 {
		spread = 0
	}
	return rng.Intn(spread+1) + attacker.MinDamage
}

// applyDamage drains defender's CurrentArmour first, then HP with any overflow.
// Returns (armourDrain, hpDrain). Mutates defender in place.
func applyDamage(raw int, defender *Combatant) (armourDrain, hpDrain int) {
	if raw <= 0 {
		return 0, 0
	}
	if defender.CurrentArmour <= 0 {
		defender.HP -= raw
		return 0, raw
	}
	if raw <= defender.CurrentArmour {
		defender.CurrentArmour -= raw
		return raw, 0
	}
	// raw > CurrentArmour — armour broken, overflow goes to HP.
	overflow := raw - defender.CurrentArmour
	armourDrain = defender.CurrentArmour
	defender.CurrentArmour = 0
	defender.HP -= overflow
	return armourDrain, overflow
}

// combatLogLine builds the appropriate log string based on what was drained.
func combatLogLine(round int, attacker, defender string, armourDrain, hpDrain int, defenderState *Combatant) string {
	switch {
	case armourDrain > 0 && hpDrain == 0:
		return fmt.Sprintf("Round %d: %s attacks %s — absorbs %d (Armour: %d)",
			round, attacker, defender, armourDrain, defenderState.CurrentArmour)
	case armourDrain > 0 && hpDrain > 0:
		return fmt.Sprintf("Round %d: %s attacks %s — armour broken, %d HP damage (%d HP, 0 Armour)",
			round, attacker, defender, hpDrain, defenderState.HP)
	case hpDrain > 0:
		return fmt.Sprintf("Round %d: %s attacks %s for %d damage (%d HP left)",
			round, attacker, defender, hpDrain, defenderState.HP)
	default:
		return fmt.Sprintf("Round %d: %s attacks %s but deals no damage",
			round, attacker, defender)
	}
}

// combatLogLinesUpTo returns log lines up to (but not including) the (roundIndex+1)-th
// distinct round boundary. roundIndex=0 returns empty; roundIndex=1 returns round 1 lines, etc.
func combatLogLinesUpTo(log []string, roundIndex int) []string {
	if roundIndex <= 0 {
		return nil
	}
	distinctRounds := 0
	lastRoundPrefix := ""
	for i, line := range log {
		if strings.HasPrefix(line, "Round ") {
			// Extract round prefix e.g. "Round 1:" to detect distinct rounds.
			colonIdx := strings.Index(line, ":")
			var prefix string
			if colonIdx > 0 {
				prefix = line[:colonIdx]
			} else {
				prefix = line
			}
			if prefix != lastRoundPrefix {
				distinctRounds++
				lastRoundPrefix = prefix
				if distinctRounds > roundIndex {
					return log[:i]
				}
			}
		}
	}
	return log
}

// hpAtRound computes a combatant's HP after damage through the given round index.
func hpAtRound(startHP int, log []string, roundIndex int, combatantName string) int {
	visible := combatLogLinesUpTo(log, roundIndex)
	hp := startHP
	// Pure HP hit: "Round N: <attacker> attacks <target> for D damage (H HP left)"
	needleFor := "attacks " + combatantName + " for "
	// Armour-broken hit: "Round N: <attacker> attacks <target> — armour broken, D HP damage ..."
	needleBroken := "attacks " + combatantName + " — armour broken, "
	for _, line := range visible {
		if idx := strings.Index(line, needleFor); idx >= 0 {
			rest := line[idx+len(needleFor):]
			hp -= parseLeadingInt(rest)
		} else if idx := strings.Index(line, needleBroken); idx >= 0 {
			rest := line[idx+len(needleBroken):]
			hp -= parseLeadingInt(rest)
		}
	}
	if hp < 0 {
		hp = 0
	}
	return hp
}

// armourAtRound computes a combatant's current armour after all hits through the given round index.
func armourAtRound(log []string, roundIndex int, combatantName string, startArmour int) int {
	visible := combatLogLinesUpTo(log, roundIndex)
	// Armour-only hit: "... attacks <target> — absorbs D (Armour: A)" — extract remaining from "(Armour: A)"
	needleAbsorb := "attacks " + combatantName + " — absorbs "
	// Armour-broken: "... attacks <target> — armour broken, ..." — armour goes to 0
	needleBroken := "attacks " + combatantName + " — armour broken,"
	armour := startArmour
	for _, line := range visible {
		if strings.Contains(line, needleBroken) {
			armour = 0
		} else if idx := strings.Index(line, needleAbsorb); idx >= 0 {
			// Parse the remaining armour from "(Armour: A)"
			if ai := strings.Index(line[idx:], "(Armour: "); ai >= 0 {
				rest := line[idx+ai+len("(Armour: "):]
				armour = parseLeadingInt(rest)
			}
		}
	}
	if armour < 0 {
		armour = 0
	}
	return armour
}

// parseLeadingInt reads the leading decimal integer from s, stopping at the first non-digit.
func parseLeadingInt(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			break
		}
	}
	return n
}
