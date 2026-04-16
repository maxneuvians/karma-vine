package game

import (
	"fmt"
	"math/rand"
	"time"

	tea "charm.land/bubbletea/v2"
)

// lootMsg returns a human-readable string for a loot item, or empty string when
// no item was dropped.
func lootMsg(item Item) string {
	if item.Name == "" {
		return ""
	}
	if item.Count > 1 {
		return fmt.Sprintf("Looted: %s x%d", item.Name, item.Count)
	}
	return fmt.Sprintf("Looted: %s", item.Name)
}

// handleKey processes a key event and returns an updated Model and command.
func handleKey(msg tea.KeyPressMsg, m Model) (Model, tea.Cmd) {
	// While the map picker is open, arrow keys navigate the list instead of moving the player.
	if m.showMapPicker {
		switch msg.String() {
		case "up", "w":
			if m.mapPickerCursor > 0 {
				m.mapPickerCursor--
			}
		case "down", "s":
			if m.mapPickerCursor < len(mapModeNames)-1 {
				m.mapPickerCursor++
			}
		case "enter":
			m.mapMode = MapMode(m.mapPickerCursor)
			m.showMapPicker = false
		case "esc", "m":
			m.showMapPicker = false
		case "?":
			m.showMapPicker = false
			m.showHelpPanel = true
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		return m, nil
	}

	// While the combat screen is shown, suppress all keys except dismiss/quit.
	if m.screenMode == ScreenCombat {
		switch msg.String() {
		case "enter", "space":
			// Unpause combat on first Space/Enter press.
			if m.combatPaused {
				m.combatPaused = false
				return m, tea.Tick(combatSpeedDuration(m.combatSpeed), func(t time.Time) tea.Msg { return CombatTickMsg{} })
			}
			// During active playback, space/enter is a no-op.
			if m.combatState != nil && m.combatLogIndex < m.combatState.Round {
				return m, nil
			}
			if m.combatState != nil && m.combatState.PlayerWon {
				// Victory: remove defeated animal, return to normal.
				if m.combatEnemy != nil && m.localMap != nil {
					for i, a := range m.localMap.Animals {
						if a == m.combatEnemy {
							m.localMap.Animals = append(m.localMap.Animals[:i], m.localMap.Animals[i+1:]...)
							break
						}
					}
				}
				// Victory: remove defeated dungeon enemy, apply pre-rolled loot.
				if m.combatDungeonEnemy != nil && m.currentDungeon != nil {
					loot := m.combatState.PendingLoot
					if loot.Name != "" && len(m.inventory.Items) < InventoryMaxSlots {
						// Stack if same name exists.
						stacked := false
						for i := range m.inventory.Items {
							if m.inventory.Items[i].Name == loot.Name {
								m.inventory.Items[i].Count += loot.Count
								stacked = true
								break
							}
						}
						if !stacked {
							m.inventory.Items = append(m.inventory.Items, loot)
						}
					}
					for i, e := range m.currentDungeon.Enemies {
						if e == m.combatDungeonEnemy {
							m.currentDungeon.Enemies = append(m.currentDungeon.Enemies[:i], m.currentDungeon.Enemies[i+1:]...)
							break
						}
					}
				}
				m.playerHP = m.combatState.Player.HP
				m.screenMode = ScreenNormal
				m.paused = false
				m.combatState = nil
				m.combatEnemy = nil
				m.combatDungeonEnemy = nil
			} else {
				// Defeat: show death screen instead of quitting.
				m.deathKiller = m.combatState.Enemy.Name
				m.screenMode = ScreenDeath
				return m, nil
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "]":
			m.combatSpeed = min(CombatSpeedFast, m.combatSpeed+1)
			return m, nil
		case "[":
			m.combatSpeed = max(CombatSpeedSlow, m.combatSpeed-1)
			return m, nil
		}
		return m, nil
	}

	// While the death screen is shown, only r (restart) and q (quit) are active.
	if m.screenMode == ScreenDeath {
		switch msg.String() {
		case "r":
			fresh := NewModel()
			fresh.viewportW = m.viewportW
			fresh.viewportH = m.viewportH
			return fresh, tickCmd()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		return m, nil
	}

	// While the inventory panel is open, arrow keys navigate the cursor.
	if m.screenMode == ScreenInventory {
		switch msg.String() {
		case "space":
			m.paused = !m.paused
			return m, nil
		case "tab":
			m.equipFocused = !m.equipFocused
		case "up", "w":
			if m.equipFocused {
				if m.equipCursor > 0 {
					m.equipCursor--
				}
			} else {
				if m.inventoryCursor > 0 {
					m.inventoryCursor--
				}
			}
		case "down", "s":
			if m.equipFocused {
				if m.equipCursor < NumBodySlots-1 {
					m.equipCursor++
				}
			} else {
				if m.inventoryCursor < len(m.inventory.Items)-1 {
					m.inventoryCursor++
				}
			}
		case "e":
			if m.equipFocused {
				m = unequipSlot(m)
			} else {
				m = equipItem(m)
			}
		case "i":
			m.screenMode = ScreenNormal
			m.paused = m.pausedBeforeInventory
			m = clampInventoryCursor(m)
		case "d":
			m = handleDrop(m)
		case "u":
			if m.mode == ModeLocal || m.mode == ModeDungeon {
				m = handleUse(m)
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.screenMode = ScreenNormal
			m.paused = m.pausedBeforeInventory
			m = clampInventoryCursor(m)
		}
		return m, nil
	}

	switch msg.String() {
	case "space":
		m.paused = !m.paused
		return m, nil

	case "q", "ctrl+c":
		return m, tea.Quit

	case "?":
		m.showHelpPanel = !m.showHelpPanel
		if m.showHelpPanel {
			m.showMapPicker = false
		}

	case "\\":
		m.showSidebar = !m.showSidebar
		if m.showSidebar {
			m.showMapPicker = false
		}

	case "m":
		if m.mode == ModeWorld {
			m.showMapPicker = true
			m.showSidebar = false
			m.mapPickerCursor = int(m.mapMode)
		}

	// Time speed
	case "]":
		m.timeScale = nextTimeScale(m.timeScale)
	case "[":
		m.timeScale = prevTimeScale(m.timeScale)

	// World zoom (world map only)
	case "+", "=":
		if m.mode == ModeWorld {
			m.worldZoom = prevWorldZoom(m.worldZoom)
		}
	case "-":
		if m.mode == ModeWorld {
			m.worldZoom = nextWorldZoom(m.worldZoom)
		}

	// Movement
	case "up", "w":
		if m.paused && m.screenMode == ScreenNormal {
			break
		}
		m = applyDelta(0, -1, m)
	case "down", "s":
		if m.paused && m.screenMode == ScreenNormal {
			break
		}
		m = applyDelta(0, 1, m)
	case "left", "a":
		if m.paused && m.screenMode == ScreenNormal {
			break
		}
		m = applyDelta(-1, 0, m)
	case "right", "d":
		if m.paused && m.screenMode == ScreenNormal {
			break
		}
		m = applyDelta(1, 0, m)

	// Descend to local map / dungeon
	case "enter", ">":
		if m.mode == ModeWorld {
			m.localMap = LocalMapFor(m.worldPos.X, m.worldPos.Y, &m)
			m.mode = ModeLocal
			m.playerPos = findSpawnPoint(m.localMap)
		} else if m.mode == ModeLocal {
			// Check if player is on a dungeon entrance.
			if m.localMap != nil {
				obj := m.localMap.Objects[m.playerPos.X][m.playerPos.Y]
				if obj != nil && obj.Char == '>' {
					meta := DungeonMetaFor(m.worldPos.X, m.worldPos.Y, &m)
					// Record biome on first entry.
					if meta.Biome == 0 {
						tile := TileAt(m.worldPos.X, m.worldPos.Y, &m)
						meta.Biome = tile.Biome
						m.dungeonMeta[m.worldPos] = meta
					}
					level := DungeonLevelFor(m.worldPos.X, m.worldPos.Y, 1, &m)
					m.currentDungeon = level
					m.dungeonDepth = 1
					m.dungeonEntryPos = m.playerPos
					m.mode = ModeDungeon
					m.playerPos = level.UpStair
				}
			}
		} else if m.mode == ModeDungeon {
			// Descend to next dungeon level.
			if m.currentDungeon != nil && m.currentDungeon.HasDownStair &&
				m.playerPos.X == m.currentDungeon.DownStair.X &&
				m.playerPos.Y == m.currentDungeon.DownStair.Y {
				nextDepth := m.dungeonDepth + 1
				level := DungeonLevelFor(m.worldPos.X, m.worldPos.Y, nextDepth, &m)
				m.currentDungeon = level
				m.dungeonDepth = nextDepth
				m.playerPos = level.UpStair
			}
		}

	// Toggle torch/brazier (dungeon only)
	case "f":
		if m.mode == ModeDungeon && m.currentDungeon != nil {
			// Check the player's own cell and all 4 neighbours.
			candidates := []LocalCoord{
				m.playerPos,
				{X: m.playerPos.X, Y: m.playerPos.Y - 1},
				{X: m.playerPos.X, Y: m.playerPos.Y + 1},
				{X: m.playerPos.X - 1, Y: m.playerPos.Y},
				{X: m.playerPos.X + 1, Y: m.playerPos.Y},
			}
			for _, c := range candidates {
				if c.X < 0 || c.X >= DungeonW || c.Y < 0 || c.Y >= DungeonH {
					continue
				}
				obj := m.currentDungeon.Cells[c.X][c.Y].Object
				if obj != nil && (obj.Char == '†' || obj.Char == 'Ω') {
					obj.Lit = !obj.Lit
					// Lighting a torch makes it non-pickupable; unlighting restores it.
					obj.Pickupable = !obj.Lit
				}
			}
		}

	// Toggle inventory panel
	case "i":
		if m.screenMode == ScreenInventory {
			m.screenMode = ScreenNormal
			m.paused = m.pausedBeforeInventory
			m = clampInventoryCursor(m)
		} else {
			m.pausedBeforeInventory = m.paused
			m.paused = true
			m.screenMode = ScreenInventory
		}

	// Pick up item / initiate combat
	case "g":
		if m.mode == ModeLocal && m.localMap != nil {
			// Check for animal at player position → initiate combat.
			for _, a := range m.localMap.Animals {
				if a.X == m.playerPos.X && a.Y == m.playerPos.Y {
					player := buildPlayerCombatant(m)
					enemy := buildEnemyCombatant(*a)
					hooks := buildCombatHooks(m)
					state := resolveCombat(player, enemy, hooks, rand.New(rand.NewSource(rand.Int63())))
					m.combatState = &state
					m.combatEnemy = a
					m.screenMode = ScreenCombat
					m.combatLogIndex = 0
					m.combatPaused = true
					m.paused = true
					return m, nil
				}
			}
		}
		if m.mode == ModeLocal || m.mode == ModeDungeon {
			m = handlePickup(m)
		}

	// Use item
	case "u":
		if m.mode == ModeLocal || m.mode == ModeDungeon {
			m = handleUse(m)
		}

	// Rest at campfire
	case "r":
		if m.mode == ModeLocal && m.localMap != nil {
			px, py := m.playerPos.X, m.playerPos.Y
			if m.localMap.Ground[px][py].HasFire && m.restCooldown == 0 {
				m.playerHP = min(m.playerHP+5, m.playerMaxHP)
				m.restCooldown = 60
			}
		}

	// Ascend from dungeon / local map
	case "esc", "<":
		if m.mode == ModeDungeon {
			if m.dungeonDepth > 1 {
				prevDepth := m.dungeonDepth - 1
				level := DungeonLevelFor(m.worldPos.X, m.worldPos.Y, prevDepth, &m)
				m.currentDungeon = level
				m.dungeonDepth = prevDepth
				m.playerPos = level.DownStair
			} else {
				m.mode = ModeLocal
				m.playerPos = m.dungeonEntryPos
			}
		} else if m.mode == ModeLocal {
			m.mode = ModeWorld
			// Do NOT nil-out m.localMap — it stays in localCache
		}
	}
	// If movement triggered combat, schedule the first combat tick (skip when paused).
	if m.screenMode == ScreenCombat && m.combatState != nil && m.combatLogIndex == 0 && !m.combatPaused {
		return m, tea.Tick(combatSpeedDuration(m.combatSpeed), func(t time.Time) tea.Msg { return CombatTickMsg{} })
	}
	return m, nil
}

// nextTimeScale returns the next higher discrete time scale (1, 2, 5, 10).
func nextTimeScale(s int) int {
	switch s {
	case 1:
		return 2
	case 2:
		return 5
	case 5:
		return 10
	default:
		return 10
	}
}

// prevTimeScale returns the next lower discrete time scale (1, 2, 5, 10).
func prevTimeScale(s int) int {
	switch s {
	case 10:
		return 5
	case 5:
		return 2
	case 2:
		return 1
	default:
		return 1
	}
}

// nextWorldZoom returns the next higher zoom level (1→2→4→8, clamped).
func nextWorldZoom(z int) int {
	switch z {
	case 1:
		return 2
	case 2:
		return 4
	case 4:
		return 8
	default:
		return 8
	}
}

// prevWorldZoom returns the next lower zoom level (8→4→2→1, clamped).
func prevWorldZoom(z int) int {
	switch z {
	case 8:
		return 4
	case 4:
		return 2
	case 2:
		return 1
	default:
		return 1
	}
}

// In ModeWorld, world position is unbounded.
// In ModeLocal, movement is clamped to [0,41]×[0,17] and blocked by Object.Blocking.
func applyDelta(dx, dy int, m Model) Model {
	switch m.mode {
	case ModeWorld:
		m.worldPos.X += dx
		m.worldPos.Y += dy

	case ModeLocal:
		if m.localMap == nil {
			return m
		}
		newX := m.playerPos.X + dx
		newY := m.playerPos.Y + dy
		// Bounds check
		if newX < 0 || newX >= LocalMapW || newY < 0 || newY >= LocalMapH {
			return m
		}
		// Collision check
		obj := m.localMap.Objects[newX][newY]
		if obj != nil && obj.Blocking {
			return m
		}
		m.playerPos.X = newX
		m.playerPos.Y = newY

	case ModeDungeon:
		if m.currentDungeon == nil {
			return m
		}
		newX := m.playerPos.X + dx
		newY := m.playerPos.Y + dy
		// Bounds check
		if newX < 0 || newX >= DungeonW || newY < 0 || newY >= DungeonH {
			return m
		}
		// Wall collision
		if m.currentDungeon.Cells[newX][newY].Kind == CellWall {
			return m
		}
		// Object blocking check
		if obj := m.currentDungeon.Cells[newX][newY].Object; obj != nil && obj.Blocking {
			return m
		}
		// Enemy collision — triggers combat without moving
		for _, e := range m.currentDungeon.Enemies {
			if e.X == newX && e.Y == newY {
				player := buildPlayerCombatant(m)
				enemy := buildDungeonEnemyCombatant(e)
				hooks := buildCombatHooks(m)
				state := resolveCombat(player, enemy, hooks, rand.New(rand.NewSource(rand.Int63())))
				m.combatState = &state
				m.combatDungeonEnemy = e
				m.screenMode = ScreenCombat
				m.combatLogIndex = 0
				m.combatPaused = true
				m.paused = true
				return m
			}
		}
		m.playerPos.X = newX
		m.playerPos.Y = newY
	}
	return m
}

// findSpawnPoint returns the nearest unblocked cell to the map centre,
// scanning outward in a spiral until a non-blocking cell is found.
func findSpawnPoint(lm *LocalMap) LocalCoord {
	const cx, cy = LocalMapW / 2, LocalMapH / 2
	// Scan outward from centre using expanding rings
	for radius := 0; radius <= LocalMapW/2; radius++ {
		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {
				// Only visit the perimeter of each ring
				if abs(dx) != radius && abs(dy) != radius {
					continue
				}
				x, y := cx+dx, cy+dy
				if x < 0 || x >= LocalMapW || y < 0 || y >= LocalMapH {
					continue
				}
				obj := lm.Objects[x][y]
				if obj == nil || !obj.Blocking {
					return LocalCoord{X: x, Y: y}
				}
			}
		}
	}
	// Fallback (should never be reached on a valid map)
	return LocalCoord{X: cx, Y: cy}
}

// ── Inventory helpers ────────────────────────────────────────────────────────

// clampInventoryCursor ensures inventoryCursor is within [0, len-1] (or 0 if empty).
func clampInventoryCursor(m Model) Model {
	max := len(m.inventory.Items) - 1
	if max < 0 {
		m.inventoryCursor = 0
	} else if m.inventoryCursor > max {
		m.inventoryCursor = max
	}
	return m
}

// equipItem equips the inventory item at inventoryCursor to its best available slot.
func equipItem(m Model) Model {
	if len(m.inventory.Items) == 0 || m.inventoryCursor >= len(m.inventory.Items) {
		return m
	}
	item := m.inventory.Items[m.inventoryCursor]
	if len(item.Slots) == 0 {
		return m // not equippable
	}

	// Find the first empty compatible slot.
	targetSlot := item.Slots[0]
	for _, s := range item.Slots {
		if m.inventory.Equipped[s].Name == "" {
			targetSlot = s
			break
		}
	}

	// If slot is occupied, swap — but only if inventory has room for the old item.
	if m.inventory.Equipped[targetSlot].Name != "" {
		old := m.inventory.Equipped[targetSlot]
		// Check if we can fit the old item back (stacking or new slot).
		canStack := false
		for i := range m.inventory.Items {
			if m.inventory.Items[i].Name == old.Name {
				canStack = true
				break
			}
		}
		// After removing the equipped item, we need space: we're removing one item from
		// inventory (the newly equipped) and potentially adding one (the swapped-out).
		// Net change depends on whether the removed item's slot gets freed.
		// Simplify: after removing the item being equipped, can we fit the old?
		futureLen := len(m.inventory.Items) - 1
		if item.Count > 1 {
			futureLen = len(m.inventory.Items) // slot stays (count decremented)
		}
		if !canStack && futureLen >= InventoryMaxSlots {
			return m // can't fit swapped item
		}

		// Return old item to inventory.
		stacked := false
		for i := range m.inventory.Items {
			if m.inventory.Items[i].Name == old.Name {
				m.inventory.Items[i].Count++
				stacked = true
				break
			}
		}
		if !stacked {
			m.inventory.Items = append(m.inventory.Items, old)
		}
	}

	// Place new item in slot.
	equipped := item
	equipped.Count = 1
	m.inventory.Equipped[targetSlot] = equipped

	// Remove from inventory.
	m.inventory.Items[m.inventoryCursor].Count--
	if m.inventory.Items[m.inventoryCursor].Count <= 0 {
		m.inventory.Items = append(m.inventory.Items[:m.inventoryCursor], m.inventory.Items[m.inventoryCursor+1:]...)
	}
	m = clampInventoryCursor(m)
	return m
}

// unequipSlot moves the equipped item at equipCursor back to inventory.
func unequipSlot(m Model) Model {
	slot := BodySlot(m.equipCursor)
	if m.inventory.Equipped[slot].Name == "" {
		return m // empty slot
	}

	old := m.inventory.Equipped[slot]

	// Check if inventory has room.
	canStack := false
	for i := range m.inventory.Items {
		if m.inventory.Items[i].Name == old.Name {
			canStack = true
			break
		}
	}
	if !canStack && len(m.inventory.Items) >= InventoryMaxSlots {
		return m // inventory full
	}

	// Add to inventory.
	stacked := false
	for i := range m.inventory.Items {
		if m.inventory.Items[i].Name == old.Name {
			m.inventory.Items[i].Count++
			stacked = true
			break
		}
	}
	if !stacked {
		m.inventory.Items = append(m.inventory.Items, old)
	}

	// Clear the slot.
	m.inventory.Equipped[slot] = Item{}
	return m
}

// handlePickup picks up a pickupable object from the player's current cell.
func handlePickup(m Model) Model {
	var obj **Object

	if m.mode == ModeLocal && m.localMap != nil {
		obj = &m.localMap.Objects[m.playerPos.X][m.playerPos.Y]
	} else if m.mode == ModeDungeon && m.currentDungeon != nil {
		obj = &m.currentDungeon.Cells[m.playerPos.X][m.playerPos.Y].Object
	}
	if obj == nil || *obj == nil || !(*obj).Pickupable {
		return m
	}
	if len(m.inventory.Items) >= InventoryMaxSlots {
		// Check if we can stack before rejecting.
		found := false
		for i := range m.inventory.Items {
			if m.inventory.Items[i].Name == (*obj).Name {
				found = true
				break
			}
		}
		if !found {
			return m
		}
	}

	o := *obj
	// Try to stack.
	stacked := false
	for i := range m.inventory.Items {
		if m.inventory.Items[i].Name == o.Name {
			m.inventory.Items[i].Count++
			stacked = true
			break
		}
	}
	if !stacked {
		m.inventory.Items = append(m.inventory.Items, Item{
			Char:  o.Char,
			Color: o.Color,
			Name:  o.Name,
			Count: 1,
		})
	}
	// Remove from map.
	*obj = nil
	return m
}

// handleDrop drops the selected inventory item onto the map.
func handleDrop(m Model) Model {
	if m.mode != ModeLocal && m.mode != ModeDungeon {
		return m
	}
	if len(m.inventory.Items) == 0 {
		return m
	}
	if m.inventoryCursor >= len(m.inventory.Items) {
		return m
	}

	item := &m.inventory.Items[m.inventoryCursor]
	dropped := &Object{
		Char:       item.Char,
		Color:      item.Color,
		Name:       item.Name,
		Blocking:   false,
		Pickupable: true,
	}

	placed := false
	if m.mode == ModeLocal && m.localMap != nil {
		x, y, ok := findDropCell(m.playerPos.X, m.playerPos.Y, m.localMap)
		if ok {
			m.localMap.Objects[x][y] = dropped
			placed = true
		}
	} else if m.mode == ModeDungeon && m.currentDungeon != nil {
		x, y, ok := findDungeonDropCell(m.playerPos.X, m.playerPos.Y, m.currentDungeon)
		if ok {
			m.currentDungeon.Cells[x][y].Object = dropped
			placed = true
		}
	}

	if !placed {
		return m
	}

	item.Count--
	if item.Count <= 0 {
		m.inventory.Items = append(m.inventory.Items[:m.inventoryCursor], m.inventory.Items[m.inventoryCursor+1:]...)
	}
	m = clampInventoryCursor(m)
	return m
}

// findDropCell finds a free cell on a local map for dropping an item.
// Checks the player's cell first, then the four cardinal neighbours.
func findDropCell(px, py int, lm *LocalMap) (int, int, bool) {
	candidates := []LocalCoord{
		{X: px, Y: py},
		{X: px, Y: py - 1},
		{X: px + 1, Y: py},
		{X: px, Y: py + 1},
		{X: px - 1, Y: py},
	}
	for _, c := range candidates {
		if c.X < 0 || c.X >= LocalMapW || c.Y < 0 || c.Y >= LocalMapH {
			continue
		}
		if lm.Objects[c.X][c.Y] == nil {
			return c.X, c.Y, true
		}
	}
	return 0, 0, false
}

// findDungeonDropCell finds a free floor cell in a dungeon for dropping an item.
func findDungeonDropCell(px, py int, dl *DungeonLevel) (int, int, bool) {
	candidates := []LocalCoord{
		{X: px, Y: py},
		{X: px, Y: py - 1},
		{X: px + 1, Y: py},
		{X: px, Y: py + 1},
		{X: px - 1, Y: py},
	}
	for _, c := range candidates {
		if c.X < 0 || c.X >= DungeonW || c.Y < 0 || c.Y >= DungeonH {
			continue
		}
		cell := dl.Cells[c.X][c.Y]
		if cell.Kind == CellFloor && cell.Object == nil {
			return c.X, c.Y, true
		}
	}
	return 0, 0, false
}

// ── Item-use dispatch ────────────────────────────────────────────────────────

// treeChars is the set of object chars considered "trees" for chopping.
var treeChars = map[rune]bool{
	'♣': true, // Tree
	'♠': true, // Pine
}

// handleUse triggers a context-sensitive item interaction.
func handleUse(m Model) Model {
	if len(m.inventory.Items) == 0 || m.inventoryCursor >= len(m.inventory.Items) {
		return m
	}
	item := m.inventory.Items[m.inventoryCursor]

	switch item.Name {
	case "Axe":
		m = useAxe(m)
	}
	return m
}

// useAxe chops an adjacent tree object. Checks north, east, south, west in priority order.
func useAxe(m Model) Model {
	dirs := []LocalCoord{
		{X: m.playerPos.X, Y: m.playerPos.Y - 1}, // north
		{X: m.playerPos.X + 1, Y: m.playerPos.Y},  // east
		{X: m.playerPos.X, Y: m.playerPos.Y + 1},  // south
		{X: m.playerPos.X - 1, Y: m.playerPos.Y},  // west
	}

	if m.mode == ModeLocal && m.localMap != nil {
		for _, c := range dirs {
			if c.X < 0 || c.X >= LocalMapW || c.Y < 0 || c.Y >= LocalMapH {
				continue
			}
			obj := m.localMap.Objects[c.X][c.Y]
			if obj != nil && treeChars[obj.Char] {
				m.localMap.Objects[c.X][c.Y] = nil
				return m
			}
		}
	} else if m.mode == ModeDungeon && m.currentDungeon != nil {
		for _, c := range dirs {
			if c.X < 0 || c.X >= DungeonW || c.Y < 0 || c.Y >= DungeonH {
				continue
			}
			obj := m.currentDungeon.Cells[c.X][c.Y].Object
			if obj != nil && treeChars[obj.Char] {
				m.currentDungeon.Cells[c.X][c.Y].Object = nil
				return m
			}
		}
	}
	return m
}

// handleMouseClick processes a mouse click event.
func handleMouseClick(msg tea.MouseClickMsg, m Model) (Model, tea.Cmd) {
	if msg.Button != tea.MouseLeft {
		return m, nil
	}

	mouse := msg.Mouse()

	// In ScreenInventory, click on an item row to select it.
	if m.screenMode == ScreenInventory {
		// The item list starts at row 2 (header + separator).
		row := mouse.Y - 2
		if row >= 0 && row < len(m.inventory.Items) {
			m.inventoryCursor = row
		}
		return m, nil
	}

	// In ScreenNormal, click-to-move on local/dungeon maps.
	if m.screenMode != ScreenNormal {
		return m, nil
	}
	if m.showSidebar || m.showMapPicker {
		return m, nil
	}
	if m.paused {
		return m, nil
	}
	if m.mode != ModeLocal && m.mode != ModeDungeon {
		return m, nil
	}

	mapH := m.viewportH - 2
	if mapH < 1 {
		mapH = 1
	}
	mapW := m.viewportW

	// Compute camera origin (same logic as render functions).
	var maxW, maxH int
	if m.mode == ModeLocal {
		if mapW > LocalMapW {
			mapW = LocalMapW
		}
		if mapH > LocalMapH {
			mapH = LocalMapH
		}
		maxW = LocalMapW
		maxH = LocalMapH
	} else {
		if mapW > DungeonW {
			mapW = DungeonW
		}
		if mapH > DungeonH {
			mapH = DungeonH
		}
		maxW = DungeonW
		maxH = DungeonH
	}

	camX := m.playerPos.X - mapW/2
	camY := m.playerPos.Y - mapH/2
	if camX < 0 {
		camX = 0
	}
	if camY < 0 {
		camY = 0
	}
	if camX > maxW-mapW {
		camX = maxW - mapW
	}
	if camY > maxH-mapH {
		camY = maxH - mapH
	}

	clickX := camX + mouse.X
	clickY := camY + mouse.Y

	dx := clickX - m.playerPos.X
	dy := clickY - m.playerPos.Y
	if dx == 0 && dy == 0 {
		return m, nil
	}

	// Take one cardinal step toward the clicked cell.
	absDx := dx
	if absDx < 0 {
		absDx = -absDx
	}
	absDy := dy
	if absDy < 0 {
		absDy = -absDy
	}
	var stepX, stepY int
	if absDx >= absDy {
		if dx > 0 {
			stepX = 1
		} else {
			stepX = -1
		}
	} else {
		if dy > 0 {
			stepY = 1
		} else {
			stepY = -1
		}
	}

	m = applyDelta(stepX, stepY, m)
	// If movement triggered combat, schedule the first combat tick (skip when paused).
	if m.screenMode == ScreenCombat && m.combatState != nil && m.combatLogIndex == 0 && !m.combatPaused {
		return m, tea.Tick(combatSpeedDuration(m.combatSpeed), func(t time.Time) tea.Msg { return CombatTickMsg{} })
	}
	return m, nil
}

// handleMouseWheel processes a mouse wheel event.
func handleMouseWheel(msg tea.MouseWheelMsg, m Model) (Model, tea.Cmd) {
	if m.screenMode != ScreenInventory {
		return m, nil
	}
	if msg.Button == tea.MouseWheelUp {
		if m.inventoryCursor > 0 {
			m.inventoryCursor--
		}
	} else if msg.Button == tea.MouseWheelDown {
		if m.inventoryCursor < len(m.inventory.Items)-1 {
			m.inventoryCursor++
		}
	}
	return m, nil
}
