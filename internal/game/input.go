package game

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleKey processes a key event and returns an updated Model and command.
func handleKey(msg tea.KeyMsg, m Model) (Model, tea.Cmd) {
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

	// Descend to local map
	case "enter", ">":
		if m.mode == ModeWorld {
			m.localMap = LocalMapFor(m.worldPos.X, m.worldPos.Y, &m)
			m.mode = ModeLocal
			m.playerPos = findSpawnPoint(m.localMap)
		}

	// Ascend to world map
	case "esc", "<":
		if m.mode == ModeLocal {
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
