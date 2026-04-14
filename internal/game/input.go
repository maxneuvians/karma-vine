package game

import (
	tea "charm.land/bubbletea/v2"
)

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
			m.showSidebar = true
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		return m, nil
	}

	// While the inventory panel is open, arrow keys navigate the cursor.
	if m.screenMode == ScreenInventory {
		switch msg.String() {
		case "up", "w":
			if m.inventoryCursor > 0 {
				m.inventoryCursor--
			}
		case "down", "s":
			if m.inventoryCursor < len(m.inventory.Items)-1 {
				m.inventoryCursor++
			}
		case "i":
			m.screenMode = ScreenNormal
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
			m = clampInventoryCursor(m)
		}
		return m, nil
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "?":
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
		m = applyDelta(0, -1, m)
	case "down", "s":
		m = applyDelta(0, 1, m)
	case "left", "a":
		m = applyDelta(-1, 0, m)
	case "right", "d":
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
					_ = DungeonMetaFor(m.worldPos.X, m.worldPos.Y, &m)
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
			m = clampInventoryCursor(m)
		} else {
			m.screenMode = ScreenInventory
		}

	// Pick up item
	case "g":
		if m.mode == ModeLocal || m.mode == ModeDungeon {
			m = handlePickup(m)
		}

	// Use item
	case "u":
		if m.mode == ModeLocal || m.mode == ModeDungeon {
			m = handleUse(m)
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
