package game

import (
	"math/rand"
)

// ── Enemy roster ────────────────────────────────────────────────────────────

// dungeonEnemyRoster maps biomes to their dungeon enemy template.
var dungeonEnemyRoster = map[Biome]*EnemyTemplate{
	Forest: {
		Name: "Goblin", Char: 'g', Color: "#55aa44",
		BaseHP: 8, MaxHP: 20, BaseArmour: 0, MaxArmour: 2,
		BaseMinDamage: 1, MaxMinDamage: 3, BaseMaxDamage: 3, MaxMaxDamage: 7,
		BaseInitiative: 4, MaxInitiative: 7,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '◊', Color: "#88aa66", Name: "Goblin Ear", Count: 1}, Weight: 2},
			{Item: Item{Char: '†', Color: "#aa8866", Name: "Rusty Dagger", Count: 1}, Weight: 1},
		},
	},
	DenseForest: {
		Name: "Goblin", Char: 'g', Color: "#55aa44",
		BaseHP: 8, MaxHP: 20, BaseArmour: 0, MaxArmour: 2,
		BaseMinDamage: 1, MaxMinDamage: 3, BaseMaxDamage: 3, MaxMaxDamage: 7,
		BaseInitiative: 4, MaxInitiative: 7,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '◊', Color: "#88aa66", Name: "Goblin Ear", Count: 1}, Weight: 2},
			{Item: Item{Char: '†', Color: "#aa8866", Name: "Rusty Dagger", Count: 1}, Weight: 1},
		},
	},
	Plains: {
		Name: "Bandit", Char: 'b', Color: "#cc9944",
		BaseHP: 10, MaxHP: 24, BaseArmour: 1, MaxArmour: 3,
		BaseMinDamage: 2, MaxMinDamage: 4, BaseMaxDamage: 4, MaxMaxDamage: 8,
		BaseInitiative: 5, MaxInitiative: 8,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '○', Color: "#ccaa44", Name: "Gold Coin", Count: 1}, Weight: 2},
			{Item: Item{Char: '†', Color: "#aaaaaa", Name: "Short Sword", Count: 1}, Weight: 1},
		},
	},
	Desert: {
		Name: "Sand Wraith", Char: 'W', Color: "#ddcc88",
		BaseHP: 10, MaxHP: 26, BaseArmour: 0, MaxArmour: 1,
		BaseMinDamage: 2, MaxMinDamage: 5, BaseMaxDamage: 5, MaxMaxDamage: 10,
		BaseInitiative: 6, MaxInitiative: 9,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '∞', Color: "#ddcc88", Name: "Wraith Dust", Count: 1}, Weight: 2},
		},
	},
	AridSteppe: {
		Name: "Sand Wraith", Char: 'W', Color: "#ddcc88",
		BaseHP: 10, MaxHP: 26, BaseArmour: 0, MaxArmour: 1,
		BaseMinDamage: 2, MaxMinDamage: 5, BaseMaxDamage: 5, MaxMaxDamage: 10,
		BaseInitiative: 6, MaxInitiative: 9,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '∞', Color: "#ddcc88", Name: "Wraith Dust", Count: 1}, Weight: 2},
		},
	},
	Jungle: {
		Name: "Jungle Troll", Char: 'T', Color: "#44cc44",
		BaseHP: 14, MaxHP: 30, BaseArmour: 1, MaxArmour: 4,
		BaseMinDamage: 3, MaxMinDamage: 5, BaseMaxDamage: 6, MaxMaxDamage: 12,
		BaseInitiative: 3, MaxInitiative: 6,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '♣', Color: "#44cc44", Name: "Troll Hide", Count: 1}, Weight: 2},
			{Item: Item{Char: '♦', Color: "#88ff88", Name: "Jungle Gem", Count: 1}, Weight: 1},
		},
	},
	Savanna: {
		Name: "Jungle Troll", Char: 'T', Color: "#44cc44",
		BaseHP: 14, MaxHP: 30, BaseArmour: 1, MaxArmour: 4,
		BaseMinDamage: 3, MaxMinDamage: 5, BaseMaxDamage: 6, MaxMaxDamage: 12,
		BaseInitiative: 3, MaxInitiative: 6,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '♣', Color: "#44cc44", Name: "Troll Hide", Count: 1}, Weight: 2},
			{Item: Item{Char: '♦', Color: "#88ff88", Name: "Jungle Gem", Count: 1}, Weight: 1},
		},
	},
	Tundra: {
		Name: "Frost Giant", Char: 'G', Color: "#88ccff",
		BaseHP: 16, MaxHP: 35, BaseArmour: 2, MaxArmour: 5,
		BaseMinDamage: 3, MaxMinDamage: 6, BaseMaxDamage: 7, MaxMaxDamage: 14,
		BaseInitiative: 2, MaxInitiative: 5,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '❄', Color: "#aaddff", Name: "Frost Shard", Count: 1}, Weight: 2},
			{Item: Item{Char: '♦', Color: "#ccccff", Name: "Ice Crystal", Count: 1}, Weight: 1},
		},
	},
	Taiga: {
		Name: "Ice Wraith", Char: 'w', Color: "#aaddff",
		BaseHP: 10, MaxHP: 24, BaseArmour: 0, MaxArmour: 2,
		BaseMinDamage: 2, MaxMinDamage: 5, BaseMaxDamage: 5, MaxMaxDamage: 10,
		BaseInitiative: 5, MaxInitiative: 8,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '❄', Color: "#aaddff", Name: "Frost Shard", Count: 1}, Weight: 2},
		},
	},
	Snow: {
		Name: "Ice Wraith", Char: 'w', Color: "#aaddff",
		BaseHP: 10, MaxHP: 24, BaseArmour: 0, MaxArmour: 2,
		BaseMinDamage: 2, MaxMinDamage: 5, BaseMaxDamage: 5, MaxMaxDamage: 10,
		BaseInitiative: 5, MaxInitiative: 8,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '❄', Color: "#aaddff", Name: "Frost Shard", Count: 1}, Weight: 2},
		},
	},
	Mountains: {
		Name: "Stone Golem", Char: 'O', Color: "#888888",
		BaseHP: 18, MaxHP: 40, BaseArmour: 3, MaxArmour: 6,
		BaseMinDamage: 2, MaxMinDamage: 4, BaseMaxDamage: 6, MaxMaxDamage: 12,
		BaseInitiative: 1, MaxInitiative: 3,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '●', Color: "#888888", Name: "Stone Heart", Count: 1}, Weight: 2},
		},
	},
	DeepOcean: {
		Name: "Cave Crustacean", Char: 'c', Color: "#cc5544",
		BaseHP: 6, MaxHP: 16, BaseArmour: 2, MaxArmour: 4,
		BaseMinDamage: 1, MaxMinDamage: 3, BaseMaxDamage: 3, MaxMaxDamage: 6,
		BaseInitiative: 3, MaxInitiative: 5,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '◊', Color: "#cc8866", Name: "Crab Shell", Count: 1}, Weight: 2},
		},
	},
	ShallowWater: {
		Name: "Cave Crustacean", Char: 'c', Color: "#cc5544",
		BaseHP: 6, MaxHP: 16, BaseArmour: 2, MaxArmour: 4,
		BaseMinDamage: 1, MaxMinDamage: 3, BaseMaxDamage: 3, MaxMaxDamage: 6,
		BaseInitiative: 3, MaxInitiative: 5,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '◊', Color: "#cc8866", Name: "Crab Shell", Count: 1}, Weight: 2},
		},
	},
	Beach: {
		Name: "Cave Crustacean", Char: 'c', Color: "#cc5544",
		BaseHP: 6, MaxHP: 16, BaseArmour: 2, MaxArmour: 4,
		BaseMinDamage: 1, MaxMinDamage: 3, BaseMaxDamage: 3, MaxMaxDamage: 6,
		BaseInitiative: 3, MaxInitiative: 5,
		LootTable: []LootEntry{
			{Item: Item{}, Weight: 3},
			{Item: Item{Char: '◊', Color: "#cc8866", Name: "Crab Shell", Count: 1}, Weight: 2},
		},
	},
}

// dungeonEnemyFallback is used when the biome has no entry in the roster.
var dungeonEnemyFallback = &EnemyTemplate{
	Name: "Cave Rat", Char: 'r', Color: "#aa8866",
	BaseHP: 5, MaxHP: 12, BaseArmour: 0, MaxArmour: 1,
	BaseMinDamage: 1, MaxMinDamage: 2, BaseMaxDamage: 2, MaxMaxDamage: 4,
	BaseInitiative: 4, MaxInitiative: 7,
	LootTable: []LootEntry{
		{Item: Item{}, Weight: 4},
		{Item: Item{Char: '◊', Color: "#aa8866", Name: "Rat Tail", Count: 1}, Weight: 1},
	},
}

// ── Spawn ───────────────────────────────────────────────────────────────────

// spawnEnemy creates a DungeonEnemy from a template with depth-scaled stats.
func spawnEnemy(tmpl *EnemyTemplate, x, y, depth, maxDepth int) *DungeonEnemy {
	denom := maxDepth - 1
	if denom < 1 {
		denom = 1
	}
	fraction := float64(depth-1) / float64(denom)

	lerp := func(base, max int) int {
		return base + int(fraction*float64(max-base))
	}

	hp := lerp(tmpl.BaseHP, tmpl.MaxHP)
	return &DungeonEnemy{
		X:          x,
		Y:          y,
		Template:   tmpl,
		HP:         hp,
		MaxHP:      hp,
		Armour:     lerp(tmpl.BaseArmour, tmpl.MaxArmour),
		MinDamage:  lerp(tmpl.BaseMinDamage, tmpl.MaxMinDamage),
		MaxDamage:  lerp(tmpl.BaseMaxDamage, tmpl.MaxMaxDamage),
		Initiative: lerp(tmpl.BaseInitiative, tmpl.MaxInitiative),
	}
}

// ── Loot resolution ─────────────────────────────────────────────────────────

// resolveEnemyLoot picks a random item from a loot table using weighted random selection.
func resolveEnemyLoot(table []LootEntry, rng *rand.Rand) Item {
	if len(table) == 0 {
		return Item{}
	}
	total := 0
	for _, e := range table {
		total += e.Weight
	}
	if total <= 0 {
		return Item{}
	}
	roll := rng.Intn(total)
	for _, e := range table {
		roll -= e.Weight
		if roll < 0 {
			return e.Item
		}
	}
	return table[len(table)-1].Item
}

// ── Enemy pathfinding ───────────────────────────────────────────────────────

const enemySightRadius = 8

// chebyshevDist returns the Chebyshev (chessboard) distance between two points.
func chebyshevDist(x1, y1, x2, y2 int) int {
	dx := x1 - x2
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y2
	if dy < 0 {
		dy = -dy
	}
	if dx > dy {
		return dx
	}
	return dy
}

// bfsNextStep runs BFS from (sx, sy) toward (tx, ty) on CellFloor cells,
// skipping cells occupied by other enemies. Returns the next step, or (-1,-1) if no path.
func bfsNextStep(level *DungeonLevel, sx, sy, tx, ty int, occupied map[[2]int]bool) (int, int) {
	type pos [2]int
	start := pos{sx, sy}
	goal := pos{tx, ty}

	if start == goal {
		return -1, -1
	}

	visited := make(map[pos]bool)
	parent := make(map[pos]pos)
	visited[start] = true
	queue := []pos{start}

	dirs := [4]pos{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur == goal {
			// Trace back to find first step from start.
			step := cur
			for parent[step] != start {
				step = parent[step]
			}
			return step[0], step[1]
		}

		for _, d := range dirs {
			nx, ny := cur[0]+d[0], cur[1]+d[1]
			next := pos{nx, ny}
			if nx < 0 || nx >= DungeonW || ny < 0 || ny >= DungeonH {
				continue
			}
			if visited[next] {
				continue
			}
			if level.Cells[nx][ny].Kind != CellFloor {
				continue
			}
			// Block on other enemies (but not the goal cell — that's the player).
			if next != goal && occupied[next] {
				continue
			}
			visited[next] = true
			parent[next] = cur
			queue = append(queue, next)
		}
	}
	return -1, -1
}

// moveEnemies advances each enemy one step toward the player if within sight radius.
// If an enemy reaches the player, combat is initiated.
func moveEnemies(m Model) Model {
	if m.currentDungeon == nil || m.screenMode != ScreenNormal {
		return m
	}

	// Build occupancy map of enemy positions.
	occupied := make(map[[2]int]bool, len(m.currentDungeon.Enemies))
	for _, e := range m.currentDungeon.Enemies {
		occupied[[2]int{e.X, e.Y}] = true
	}

	for _, e := range m.currentDungeon.Enemies {
		dist := chebyshevDist(e.X, e.Y, m.playerPos.X, m.playerPos.Y)
		if dist > enemySightRadius {
			continue
		}

		// Remove from occupancy before moving.
		delete(occupied, [2]int{e.X, e.Y})

		nx, ny := bfsNextStep(m.currentDungeon, e.X, e.Y, m.playerPos.X, m.playerPos.Y, occupied)
		if nx >= 0 && ny >= 0 {
			e.X = nx
			e.Y = ny
		}

		// Re-add to occupancy.
		occupied[[2]int{e.X, e.Y}] = true

		// Check if enemy reached player.
		if e.X == m.playerPos.X && e.Y == m.playerPos.Y {
			player := buildPlayerCombatant(m)
			enemy := buildDungeonEnemyCombatant(e)
			hooks := buildCombatHooks(m)
			state := resolveCombat(player, enemy, hooks, rand.New(rand.NewSource(rand.Int63())))
			if len(e.Template.LootTable) > 0 {
				loot := resolveEnemyLoot(e.Template.LootTable, rand.New(rand.NewSource(rand.Int63())))
				state.PendingLoot = loot
				state.LootMsg = lootMsg(loot)
			}
			m.combatState = &state
			m.combatDungeonEnemy = e
			m.screenMode = ScreenCombat
			m.combatLogIndex = 0
			m.combatPaused = true
			m.paused = true
			break
		}
	}

	return m
}
