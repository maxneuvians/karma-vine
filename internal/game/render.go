package game

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

// ── Colour dimming ───────────────────────────────────────────────────────────

// dimFactor returns a brightness multiplier in [0.15, 1.0] derived from
// timeOfDay using a cosine curve: 1.0 at noon (0.5), ~0.15 at midnight (0.0/1.0).
func dimFactor(timeOfDay float64) float64 {
	// Shift so the cosine peaks at timeOfDay=0.5 (noon) rather than 0.0.
	v := 0.5*(1+math.Cos(2*math.Pi*(timeOfDay-0.5)))*0.85 + 0.15
	if v < 0.15 {
		return 0.15
	}
	if v > 1.0 {
		return 1.0
	}
	return v
}

// applyColor multiplies each R, G, B channel of a #rrggbb string by dim.
// Non-matching strings are returned unchanged.
func applyColor(hex string, dim float64) string {
	if len(hex) != 7 || hex[0] != '#' {
		return hex
	}
	r, err1 := strconv.ParseUint(hex[1:3], 16, 64)
	g, err2 := strconv.ParseUint(hex[3:5], 16, 64)
	b, err3 := strconv.ParseUint(hex[5:7], 16, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return hex
	}
	return fmt.Sprintf("#%02x%02x%02x",
		uint64(float64(r)*dim),
		uint64(float64(g)*dim),
		uint64(float64(b)*dim),
	)
}

// ── Map mode tile visual ──────────────────────────────────────────────────────

// lerpHex linearly interpolates between two #rrggbb hex colors by factor t ∈ [0,1].
// Returns a unchanged if either string is not a valid #rrggbb color.
func lerpHex(a, b string, t float64) string {
	if len(a) != 7 || a[0] != '#' || len(b) != 7 || b[0] != '#' {
		return a
	}
	ar, e1 := strconv.ParseUint(a[1:3], 16, 64)
	ag, e2 := strconv.ParseUint(a[3:5], 16, 64)
	ab, e3 := strconv.ParseUint(a[5:7], 16, 64)
	br, e4 := strconv.ParseUint(b[1:3], 16, 64)
	bg, e5 := strconv.ParseUint(b[3:5], 16, 64)
	bb, e6 := strconv.ParseUint(b[5:7], 16, 64)
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || e5 != nil || e6 != nil {
		return a
	}
	r := float64(ar) + t*(float64(br)-float64(ar))
	g := float64(ag) + t*(float64(bg)-float64(ag))
	bv := float64(ab) + t*(float64(bb)-float64(ab))
	return fmt.Sprintf("#%02x%02x%02x", uint64(r), uint64(g), uint64(bv))
}

// perceivedTemperature adjusts a tile's raw climate temperature to reflect the
// biome's felt heat using additive offsets, so that same-latitude tiles always
// sort correctly regardless of their absolute temperature.
// The result is clamped to [0, 1].
func perceivedTemperature(t Tile) float64 {
	v := t.Temperature
	switch t.Biome {
	case DeepOcean:
		// Deep water has high thermal mass and strong evaporative cooling;
		// use a combined scale+offset so it's always substantially cooler than land.
		v = v*0.5 - 0.15
	case ShallowWater:
		v = v*0.6 - 0.08
	case Beach:
		v -= 0.05 // wet sand, slightly cooler than open plains
	case Forest:
		v -= 0.12 // canopy shade and transpiration
	case DenseForest, Taiga:
		v -= 0.18
	case Jungle:
		v += 0.08 // hot and humid
	case Savanna:
		v += 0.12
	case Desert, AridSteppe:
		v += 0.20 // bare ground radiates heat strongly
	case Mountains:
		v -= 0.20 // altitude lapse rate
	case Snow, Tundra:
		v -= 0.35 // persistently frozen
	}
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// thermalColor maps t ∈ [0,1] to a 5-stop thermal gradient:
// dark-blue → sky-blue → sea-green → yellow → red.
// Unlike a 2-stop lerp this avoids the confusing purple midpoint.
func thermalColor(t float64) string {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	type rgb struct{ r, g, b float64 }
	stops := []struct {
		at  float64
		col rgb
	}{
		{0.00, rgb{0x00, 0x22, 0xcc}}, // dark blue   (very cold)
		{0.25, rgb{0x00, 0x99, 0xff}}, // sky blue    (cool)
		{0.50, rgb{0x00, 0xdd, 0x88}}, // sea green   (moderate)
		{0.75, rgb{0xff, 0xee, 0x00}}, // yellow      (warm)
		{1.00, rgb{0xff, 0x22, 0x00}}, // red         (hot)
	}
	for i := 1; i < len(stops); i++ {
		if t <= stops[i].at {
			s0, s1 := stops[i-1], stops[i]
			u := (t - s0.at) / (s1.at - s0.at)
			r := s0.col.r + u*(s1.col.r-s0.col.r)
			g := s0.col.g + u*(s1.col.g-s0.col.g)
			b := s0.col.b + u*(s1.col.b-s0.col.b)
			return fmt.Sprintf("#%02x%02x%02x", uint64(r), uint64(g), uint64(b))
		}
	}
	return "#ff2200"
}

// tileVisual returns the character and color to display for a tile in the given MapMode.
// In MapModeDefault the tile's own Char and Color are returned unchanged.
func tileVisual(t Tile, mode MapMode) (ch rune, color string) {
	switch mode {
	case MapModeTemperature:
		return '█', thermalColor(perceivedTemperature(t))
	case MapModeElevation:
		return '█', lerpHex("#1a6fa8", "#f0f6fc", t.Elevation)
	case MapModePolitical:
		// Show a contour character near elevation-band boundaries (every 0.1 unit).
		if int(t.Elevation*10) != int((t.Elevation+0.05)*10) {
			return '+', "#aabbcc"
		}
		return '·', "#334455"
	default: // MapModeDefault
		return t.Char, t.Color
	}
}

// ── Biome name ───────────────────────────────────────────────────────────────

func biomeName(b Biome) string {
	switch b {
	case DeepOcean:
		return "Deep Ocean"
	case ShallowWater:
		return "Shallow Water"
	case Beach:
		return "Beach"
	case Forest:
		return "Forest"
	case Plains:
		return "Plains"
	case DenseForest:
		return "Dense Forest"
	case Desert:
		return "Desert"
	case Mountains:
		return "Mountains"
	case Snow:
		return "Snow"
	case Jungle:
		return "Jungle"
	case Savanna:
		return "Savanna"
	case AridSteppe:
		return "Arid Steppe"
	case Tundra:
		return "Tundra"
	case Taiga:
		return "Taiga"
	default:
		return "Unknown"
	}
}

// ── World map renderer ────────────────────────────────────────────────────────

// renderWorldMap renders mapW × mapH cells centred on worldPos.
// m.worldZoom controls how many world tiles each screen cell represents (1/2/4/8).
func renderWorldMap(m Model, mapW, mapH int) string {
	z := m.worldZoom
	if z < 1 {
		z = 1
	}
	cx := mapW / 2
	cy := mapH / 2
	rows := make([]string, 0, mapH)
	for sy := 0; sy < mapH; sy++ {
		var row strings.Builder
		for sx := 0; sx < mapW; sx++ {
			// Render the player marker at the viewport centre.
			if sx == cx && sy == cy {
				row.WriteString(playerStyle.Render("@"))
				continue
			}
			wx := m.worldPos.X + (sx-cx)*z
			wy := m.worldPos.Y + (sy-cy)*z
			tile := TileAt(wx, wy, &m)
			ch, col := tileVisual(tile, m.mapMode)
			// Data overlay modes show raw values — skip time-of-day dimming so
			// hues remain accurate regardless of the in-game clock.
			var color string
			if m.mapMode == MapModeDefault || m.mapMode == MapModePolitical {
				color = applyColor(col, dimFactor(m.timeOfDay))
			} else {
				color = col
			}
			cell := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(string(ch))
			row.WriteString(cell)
		}
		rows = append(rows, row.String())
	}
	return strings.Join(rows, "\n")
}

// ── Local map renderer ────────────────────────────────────────────────────────

var playerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f0f6fc")).Bold(true)

// renderLocalMap renders a mapW×mapH viewport of the LocalMap centred on playerPos.
func renderLocalMap(m Model, mapW, mapH int) string {
	lm := m.localMap
	if lm == nil {
		if cached, ok := m.localCache[m.worldPos]; ok {
			lm = cached
		}
	}
	if lm == nil {
		return "Local map not loaded."
	}

	// Cap the render dimensions to the local map size so camera clamping never
	// produces a negative origin that would panic on array access.
	if mapW > LocalMapW {
		mapW = LocalMapW
	}
	if mapH > LocalMapH {
		mapH = LocalMapH
	}

	// Compute camera origin so playerPos is centred, clamped to map bounds.
	camX := m.playerPos.X - mapW/2
	camY := m.playerPos.Y - mapH/2
	if camX < 0 {
		camX = 0
	}
	if camY < 0 {
		camY = 0
	}
	if camX > LocalMapW-mapW {
		camX = LocalMapW - mapW
	}
	if camY > LocalMapH-mapH {
		camY = LocalMapH - mapH
	}

	// Build animal lookup by position.
	type pos struct{ x, y int }
	animalAt := make(map[pos]*Animal, len(lm.Animals))
	for _, a := range lm.Animals {
		animalAt[pos{a.X, a.Y}] = a
	}

	rows := make([]string, 0, mapH)
	globalDim := dimFactor(m.timeOfDay)
	for sy := 0; sy < mapH; sy++ {
		y := camY + sy
		var row strings.Builder
		for sx := 0; sx < mapW; sx++ {
			x := camX + sx
			// Per-cell dim: blend global dim with fire illumination intensity.
			// cellDim = max(globalDim, fireIntensity) so fire always adds brightness.
			cellDim := globalDim
			if lm.LitMap[x][y] > globalDim {
				cellDim = lm.LitMap[x][y]
			}
			// Player overrides everything.
			if x == m.playerPos.X && y == m.playerPos.Y {
				row.WriteString(playerStyle.Render("@"))
				continue
			}
			// Animal > fire > object > ground.
			if a, ok := animalAt[pos{x, y}]; ok {
				color := applyColor(a.Color, cellDim)
				row.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(string(a.Char)))
				continue
			}
			if lm.Ground[x][y].HasFire {
				color := applyColor("#ff8800", cellDim)
				row.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render("♨"))
				continue
			}
			if obj := lm.Objects[x][y]; obj != nil {
				color := applyColor(obj.Color, cellDim)
				row.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(string(obj.Char)))
				continue
			}
			g := lm.Ground[x][y]
			color := applyColor(g.Color, cellDim)
			row.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(string(g.Char)))
		}
		rows = append(rows, row.String())
	}
	return strings.Join(rows, "\n")
}

// ── HUD status bar ────────────────────────────────────────────────────────────

// renderProgressBar returns a styled bar of filled (█) and empty (░) runes.
// Width is the total number of characters. fillColor and emptyColor are lipgloss hex colours.
func renderProgressBar(current, max, width int, fillColor, emptyColor string) string {
	if width <= 0 {
		return ""
	}
	if max <= 0 {
		empty := strings.Repeat("░", width)
		return lipgloss.NewStyle().Foreground(lipgloss.Color(emptyColor)).Render(empty)
	}
	filled := current * width / max
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	var b strings.Builder
	if filled > 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(fillColor)).Render(strings.Repeat("█", filled)))
	}
	if width-filled > 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(emptyColor)).Render(strings.Repeat("░", width-filled)))
	}
	return b.String()
}

// tempCelsius converts tile climate values to a display temperature in °C.
//   - Base range: 0.0 → -20 °C (polar), 1.0 → 40 °C (equatorial).
//   - Elevation lapse: ~6 °C cooler per 0.15 elevation above sea level (0.36).
//   - Diurnal swing: ±5 °C cosine curve peaking at 14:00 (timeOfDay ≈ 0.583).
func tempCelsius(temperature, elevation, timeOfDay float64) int {
	base := temperature*60 - 20
	elevAdj := -(elevation - 0.36) * 40
	timeAdj := math.Cos(2*math.Pi*(timeOfDay-0.583)) * 5
	return int(math.Round(base + elevAdj + timeAdj))
}

// formatTime converts a timeOfDay in [0,1) to a "HH:MM" 24-hour string.
func formatTime(timeOfDay float64) string {
	totalMinutes := int(timeOfDay * 24 * 60)
	h := totalMinutes / 60
	m := totalMinutes % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

// renderHUD renders the single-row status bar with HP bar, armour, tile info, clock, speed, paused indicator, and help hint.
func renderHUD(m Model) string {
	tile := TileAt(m.worldPos.X, m.worldPos.Y, &m)
	clock := formatTime(m.timeOfDay)
	speed := fmt.Sprintf("%d×", m.timeScale)

	// Compute armour from equipped items (currently all items have 0 armour bonus).
	armour := 0

	// Always use perceivedTemperature so the HUD reading is consistent with the
	// temperature map overlay and reflects the biome's felt heat.
	celsius := tempCelsius(perceivedTemperature(tile), tile.Elevation, m.timeOfDay)

	// Mode-contextual tile info.
	var tileInfo string
	if m.mode == ModeDungeon {
		tileInfo = fmt.Sprintf("Dungeon D:%d (%d,%d)", m.dungeonDepth, m.worldPos.X, m.worldPos.Y)
	} else if m.mode == ModeLocal {
		tileInfo = fmt.Sprintf("%s %d°C (%d,%d)", biomeName(tile.Biome), celsius, m.playerPos.X, m.playerPos.Y)
	} else {
		tileInfo = fmt.Sprintf("%s %d°C (%d,%d)", biomeName(tile.Biome), celsius, m.worldPos.X, m.worldPos.Y)
	}

	// Build fixed segments (everything except the HP bar).
	hpLabel := fmt.Sprintf(" HP %d/%d", m.playerHP, m.playerMaxHP)
	armLabel := fmt.Sprintf("ARM:%d", armour)
	pausedStr := ""
	if m.paused {
		pausedStr = "  [PAUSED]"
	}
	helpHint := "? help"

	// fixed = spaces + hpLabel + spaces + armLabel + spaces + tileInfo + spaces + clock + spaces + speed + paused + spaces + helpHint + trailing space
	fixedLen := 1 + len(hpLabel) + 2 + len(armLabel) + 2 + len(tileInfo) + 2 + len(clock) + 2 + len(speed) + len(pausedStr) + 2 + len(helpHint) + 1

	// Compute HP bar width from remaining space, clamped to [5, 20].
	barW := m.viewportW - fixedLen
	if barW < 5 {
		barW = 5
	}
	if barW > 20 {
		barW = 20
	}

	hpBar := renderProgressBar(m.playerHP, m.playerMaxHP, barW, "#22cc55", "#444c56")

	text := fmt.Sprintf(" %s%s  %s  %s  %s  %s%s  %s ",
		hpBar, hpLabel, armLabel, tileInfo, clock, speed, pausedStr, helpHint)

	return lipgloss.NewStyle().Foreground(lipgloss.Color("#ccd9e0")).Render(text)
}

// joinVertical stacks two strings vertically, separated by a newline.
func joinVertical(top, bottom string) string {
	return lipgloss.JoinVertical(lipgloss.Left, top, bottom)
}

// ── Sidebar ───────────────────────────────────────────────────────────────────

const sidebarContentW = 20 // visible character width of the sidebar (separator │ is extra)

var (
	sidebarHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ccd9e0"))
	sidebarSubStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#768390"))
	sidebarSepStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#444c56"))
	keyBarStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#768390"))
)

type legendEntry struct {
	char  rune
	color string
	name  string
}

var biomeLegend = []legendEntry{
	{'≋', "#1a6fa8", "Deep Ocean"},
	{'≈', "#2e9ecf", "Shallow Water"},
	{'·', "#e8c96a", "Beach"},
	{'♣', "#2d7a1f", "Forest"},
	{'░', "#5aad3f", "Plains"},
	{'♠', "#3d6b3a", "Dense Forest"},
	{'~', "#c8a46a", "Desert"},
	{'▲', "#8fa89c", "Mountains"},
	{'*', "#ccd9e0", "Snow"},
}

// sbCell renders a colored icon + name padded to sidebarContentW.
func sbCell(ch rune, color, name string) string {
	icon := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(string(ch))
	return lipgloss.NewStyle().Width(sidebarContentW).Render(" " + icon + " " + name)
}

// sbText renders plain text padded to sidebarContentW.
func sbText(s string) string {
	return lipgloss.NewStyle().Width(sidebarContentW).Render(s)
}

// mapModeOverlayHints provides the sidebar header + hint for each non-default map mode.
var mapModeOverlayHints = map[MapMode][2]string{
	MapModeTemperature: {"Temperature", "blue=cold / red=hot"},
	MapModeElevation:   {"Elevation", "blue=low / white=high"},
	MapModePolitical:   {"Political", "contour lines"},
}

// renderSidebar builds a height-row sidebar with a trailing │ on each line.
func renderSidebar(m Model, height int) string {
	if height < 1 {
		height = 1
	}

	var lines []string

	switch m.mode {
	case ModeWorld:
		if m.mapMode != MapModeDefault {
			hint := mapModeOverlayHints[m.mapMode]
			lines = append(lines,
				sbText(sidebarHeaderStyle.Render(" "+hint[0])),
				sbText(sidebarSubStyle.Render(" "+strings.Repeat("─", 18))),
				sbText(sidebarSubStyle.Render(" "+hint[1])),
			)
		} else {
			lines = append(lines,
				sbText(sidebarHeaderStyle.Render(" Biomes")),
				sbText(sidebarSubStyle.Render(" "+strings.Repeat("─", 18))),
			)
			for _, e := range biomeLegend {
				lines = append(lines, sbCell(e.char, applyColor(e.color, dimFactor(m.timeOfDay)), e.name))
			}
		}

	case ModeLocal:
		lines = append(lines,
			sbText(sidebarHeaderStyle.Render(" Legend")),
			sbText(sidebarSubStyle.Render(" "+strings.Repeat("─", 18))),
			sbCell('@', "#f0f6fc", "You"),
		)
		lm := m.localMap
		if lm == nil {
			if cached, ok := m.localCache[m.worldPos]; ok {
				lm = cached
			}
		}
		if lm != nil {
			type entry struct {
				char  rune
				color string
				name  string
			}
			seenObj := make(map[string]entry)
			seenAni := make(map[string]entry)
			hasFire := false
			for x := 0; x < LocalMapW; x++ {
				for y := 0; y < LocalMapH; y++ {
					if obj := lm.Objects[x][y]; obj != nil {
						key := obj.Name
						if key == "" {
							key = string(obj.Char)
						}
						seenObj[key] = entry{obj.Char, obj.Color, obj.Name}
					}
					if lm.Ground[x][y].HasFire {
						hasFire = true
					}
				}
			}
			for _, a := range lm.Animals {
				key := a.Name
				if key == "" {
					key = string(a.Char)
				}
				seenAni[key] = entry{a.Char, a.Color, a.Name}
			}
			if hasFire {
				lines = append(lines, sbCell('♨', "#ff8800", "Campfire"))
			}
			if len(seenObj) > 0 {
				lines = append(lines, sbText(sidebarSubStyle.Render(" Objects")))
				objKeys := make([]string, 0, len(seenObj))
				for k := range seenObj {
					objKeys = append(objKeys, k)
				}
				sort.Strings(objKeys)
				for _, k := range objKeys {
					e := seenObj[k]
					name := e.name
					if name == "" {
						name = string(e.char)
					}
					lines = append(lines, sbCell(e.char, applyColor(e.color, dimFactor(m.timeOfDay)), name))
				}
			}
			if len(seenAni) > 0 {
				lines = append(lines, sbText(sidebarSubStyle.Render(" Wildlife")))
				aniKeys := make([]string, 0, len(seenAni))
				for k := range seenAni {
					aniKeys = append(aniKeys, k)
				}
				sort.Strings(aniKeys)
				for _, k := range aniKeys {
					e := seenAni[k]
					name := e.name
					if name == "" {
						name = string(e.char)
					}
					lines = append(lines, sbCell(e.char, applyColor(e.color, dimFactor(m.timeOfDay)), name))
				}
			}
		}

	case ModeDungeon:
		lines = append(lines,
			sbText(sidebarHeaderStyle.Render(" Dungeon")),
			sbText(sidebarSubStyle.Render(" "+strings.Repeat("─", 18))),
			sbText(fmt.Sprintf(" Depth: %d", m.dungeonDepth)),
			sbCell('@', "#f0f6fc", "You"),
		)
		if m.currentDungeon != nil {
			type entry struct {
				char  rune
				color string
				name  string
			}
			seen := make(map[string]entry)
			for x := 0; x < DungeonW; x++ {
				for y := 0; y < DungeonH; y++ {
					if obj := m.currentDungeon.Cells[x][y].Object; obj != nil {
						key := obj.Name
						if key == "" {
							key = string(obj.Char)
						}
						seen[key] = entry{obj.Char, obj.Color, obj.Name}
					}
				}
			}
			if len(seen) > 0 {
				lines = append(lines, sbText(sidebarSubStyle.Render(" Contents")))
				seenKeys := make([]string, 0, len(seen))
				for k := range seen {
					seenKeys = append(seenKeys, k)
				}
				sort.Strings(seenKeys)
				for _, k := range seenKeys {
					e := seen[k]
					name := e.name
					if name == "" {
						name = string(e.char)
					}
					lines = append(lines, sbCell(e.char, e.color, name))
				}
			}
		}
	}

	sep := sidebarSepStyle.Render("│")
	hint := sbText(sidebarSubStyle.Render(" ? close"))

	rows := make([]string, height)
	for i := range rows {
		content := sbText("")
		if i < len(lines) {
			content = lines[i]
		}
		rows[i] = content + sep
	}
	rows[height-1] = hint + sep
	return strings.Join(rows, "\n")
}

// ── Help panel ────────────────────────────────────────────────────────────────

// renderHelpPanel builds a fullscreen key-bindings overlay clamped to viewport.
func renderHelpPanel(m Model) string {
	maxW := lipgloss.NewStyle().MaxWidth(m.viewportW)

	var lines []string
	lines = append(lines, maxW.Render(" Key Bindings"))
	lines = append(lines, maxW.Render(" "+strings.Repeat("─", m.viewportW-2)))
	lines = append(lines, maxW.Render(""))

	// Universal bindings.
	lines = append(lines, maxW.Render(" Universal"))
	lines = append(lines, maxW.Render("   q        quit"))
	lines = append(lines, maxW.Render("   space    pause/unpause"))
	lines = append(lines, maxW.Render("   i        inventory"))
	lines = append(lines, maxW.Render("   \\        sidebar"))
	lines = append(lines, maxW.Render("   [/]      time speed"))
	lines = append(lines, maxW.Render("   ?        close help"))
	lines = append(lines, maxW.Render(""))

	switch m.mode {
	case ModeWorld:
		lines = append(lines, maxW.Render(" World Map"))
		lines = append(lines, maxW.Render("   ↑↓←→/wasd  move"))
		lines = append(lines, maxW.Render("   enter/>     descend to local"))
		lines = append(lines, maxW.Render("   +/-         zoom"))
		lines = append(lines, maxW.Render("   m           map mode picker"))
	case ModeLocal:
		lines = append(lines, maxW.Render(" Local Map"))
		lines = append(lines, maxW.Render("   ↑↓←→/wasd  move"))
		lines = append(lines, maxW.Render("   g           pick up / fight"))
		lines = append(lines, maxW.Render("   d           drop"))
		lines = append(lines, maxW.Render("   u           use"))
		lines = append(lines, maxW.Render("   enter/>     enter dungeon"))
		lines = append(lines, maxW.Render("   esc/<       ascend to world"))
	case ModeDungeon:
		lines = append(lines, maxW.Render(" Dungeon"))
		lines = append(lines, maxW.Render("   ↑↓←→/wasd  move"))
		lines = append(lines, maxW.Render("   g           pick up / fight"))
		lines = append(lines, maxW.Render("   d           drop"))
		lines = append(lines, maxW.Render("   u           use"))
		lines = append(lines, maxW.Render("   f           toggle torch"))
		lines = append(lines, maxW.Render("   enter/>     descend deeper"))
		lines = append(lines, maxW.Render("   esc/<       ascend"))
	}

	// Clamp to viewport height.
	if len(lines) > m.viewportH {
		lines = lines[:m.viewportH]
	}

	return strings.Join(lines, "\n")
}

// ── Map mode picker ───────────────────────────────────────────────────────────

const pickerContentW = 22 // visible character width of the picker panel (separator │ is extra)

var (
	pickerHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ccd9e0"))
	pickerCursorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#58a6ff"))
	pickerItemStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ccd9e0"))
)

var mapModeNames = []string{"Default", "Temperature", "Elevation", "Political"}

// pkText renders plain text padded to pickerContentW.
func pkText(s string) string {
	return lipgloss.NewStyle().Width(pickerContentW).Render(s)
}

// renderMapPicker builds a height-row picker panel with a trailing │ on each line.
func renderMapPicker(m Model, height int) string {
	if height < 1 {
		height = 1
	}

	var lines []string
	lines = append(lines,
		pkText(pickerHeaderStyle.Render(" Map Mode")),
		pkText(sidebarSubStyle.Render(" "+strings.Repeat("─", 20))),
	)
	for i, name := range mapModeNames {
		if i == m.mapPickerCursor {
			lines = append(lines, pkText(pickerCursorStyle.Render(" > "+name)))
		} else {
			lines = append(lines, pkText(pickerItemStyle.Render("   "+name)))
		}
	}

	sep := sidebarSepStyle.Render("│")
	hint := pkText(sidebarSubStyle.Render(" m/esc close"))

	rows := make([]string, height)
	for i := range rows {
		content := pkText("")
		if i < len(lines) {
			content = lines[i]
		}
		rows[i] = content + sep
	}
	rows[height-1] = hint + sep
	return strings.Join(rows, "\n")
}



// ── Dungeon renderer ──────────────────────────────────────────────────────────

const (
	playerViewRadius = 6
	torchRadius      = 4
)

// computeDungeonLight returns a per-cell brightness value in [0, 1] based on
// Chebyshev distance to the player and any lit torches/braziers. Cells with a
// brightness of 0 are not rendered (fog of war).
func computeDungeonLight(m Model) map[LocalCoord]float64 {
	light := make(map[LocalCoord]float64)
	if m.currentDungeon == nil {
		return light
	}

	addSource := func(cx, cy, radius int) {
		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {
				x := cx + dx
				y := cy + dy
				if x < 0 || x >= DungeonW || y < 0 || y >= DungeonH {
					continue
				}
				absDx := dx
				if absDx < 0 {
					absDx = -absDx
				}
				absDy := dy
				if absDy < 0 {
					absDy = -absDy
				}
				chebDist := absDx
				if absDy > chebDist {
					chebDist = absDy
				}
				v := float64(radius+1-chebDist) / float64(radius+1)
				coord := LocalCoord{X: x, Y: y}
				if v > light[coord] {
					light[coord] = v
				}
			}
		}
	}

	// Player light source.
	addSource(m.playerPos.X, m.playerPos.Y, playerViewRadius)

	// Lit torch/brazier sources.
	for tx := 0; tx < DungeonW; tx++ {
		for ty := 0; ty < DungeonH; ty++ {
			obj := m.currentDungeon.Cells[tx][ty].Object
			if obj != nil && (obj.Char == '†' || obj.Char == 'Ω') && obj.Lit {
				addSource(tx, ty, torchRadius)
			}
		}
	}

	return light
}

// renderDungeonMap renders a mapW×mapH viewport of the dungeon centred on playerPos.
func renderDungeonMap(m Model, mapW, mapH int) string {
	if m.currentDungeon == nil {
		return "Dungeon not loaded."
	}

	if mapW > DungeonW {
		mapW = DungeonW
	}
	if mapH > DungeonH {
		mapH = DungeonH
	}

	// Camera origin centred on player, clamped to dungeon bounds.
	camX := m.playerPos.X - mapW/2
	camY := m.playerPos.Y - mapH/2
	if camX < 0 {
		camX = 0
	}
	if camY < 0 {
		camY = 0
	}
	if camX > DungeonW-mapW {
		camX = DungeonW - mapW
	}
	if camY > DungeonH-mapH {
		camY = DungeonH - mapH
	}

	vis := computeDungeonLight(m)

	rows := make([]string, 0, mapH)
	for sy := 0; sy < mapH; sy++ {
		y := camY + sy
		var row strings.Builder
		for sx := 0; sx < mapW; sx++ {
			x := camX + sx

			brightness := vis[LocalCoord{X: x, Y: y}]
			if brightness == 0 {
				row.WriteRune(' ')
				continue
			}

			// Player overrides everything.
			if x == m.playerPos.X && y == m.playerPos.Y {
				row.WriteString(playerStyle.Render("@"))
				continue
			}

			cell := m.currentDungeon.Cells[x][y]

			// Enemy layer (before object/floor, after player).
			enemyDrawn := false
			for _, e := range m.currentDungeon.Enemies {
				if e.X == x && e.Y == y {
					row.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(e.Template.Color)).Render(string(e.Template.Char)))
					enemyDrawn = true
					break
				}
			}
			if enemyDrawn {
				continue
			}

			// Object > base cell.
			if cell.Object != nil {
				color := cell.Object.Color
				// Unlit torches and braziers render darker.
				if (cell.Object.Char == '†' || cell.Object.Char == 'Ω') && !cell.Object.Lit {
					color = "#444444"
				}
				row.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(string(cell.Object.Char)))
				continue
			}

			switch cell.Kind {
			case CellWall:
				wallColor := lerpHex("#3d1a08", "#c07a40", brightness)
				row.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(wallColor)).Render("█"))
			case CellFloor:
				row.WriteRune(' ')
			}
		}
		rows = append(rows, row.String())
	}
	return strings.Join(rows, "\n")
}

// ── Inventory styles ──────────────────────────────────────────────────────────

var (
	invHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ccd9e0"))
	invCursorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#58a6ff"))
	invItemStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ccd9e0"))
	invEmptyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#768390"))
)

// ── Fullscreen Inventory ──────────────────────────────────────────────────────

// ragdoll is the ASCII body outline for the equipment panel.
var ragdoll = []string{
	"    ~O~    ",
	"    /|\\    ",
	"   / | \\   ",
	"     |     ",
	"    / \\    ",
	"   /   \\   ",
}

// equipSlots lists the named equipment positions.
var equipSlots = []string{"Head", "Chest", "Left Hand", "Right Hand", "Legs", "Feet"}

// renderFullscreenInventory fills the entire viewport with a two-column inventory layout.
func renderFullscreenInventory(m Model) string {
	leftW := m.viewportW * 60 / 100
	if leftW < 20 {
		leftW = 20
	}
	rightW := m.viewportW - leftW

	// ── Left column: item list ──
	var leftLines []string
	if !m.equipFocused {
		leftLines = append(leftLines, invHeaderStyle.Render(" Inventory"))
	} else {
		leftLines = append(leftLines, invEmptyStyle.Render(" Inventory"))
	}
	leftLines = append(leftLines,
		sidebarSubStyle.Render(" "+strings.Repeat("─", leftW-2)),
	)

	if len(m.inventory.Items) == 0 {
		leftLines = append(leftLines, invEmptyStyle.Render(" Empty"))
	} else {
		for i, item := range m.inventory.Items {
			icon := lipgloss.NewStyle().Foreground(lipgloss.Color(item.Color)).Render(string(item.Char))
			if i == m.inventoryCursor && !m.equipFocused {
				label := fmt.Sprintf(" > %s %s  x%d", icon, item.Name, item.Count)
				leftLines = append(leftLines, invCursorStyle.Render(label))
			} else {
				label := fmt.Sprintf("   %s %s  x%d", icon, item.Name, item.Count)
				leftLines = append(leftLines, invItemStyle.Render(label))
			}
		}
	}

	// Pad to viewportH-1, then add hint at bottom.
	for len(leftLines) < m.viewportH-1 {
		leftLines = append(leftLines, "")
	}
	hint := " i close  d drop  u use  e equip  Tab switch"
	if m.mode == ModeWorld {
		hint = " i close  e equip  Tab switch"
	}
	leftLines = leftLines[:m.viewportH-1]
	leftLines = append(leftLines, sidebarSubStyle.Render(hint))

	leftCol := lipgloss.NewStyle().Width(leftW).Height(m.viewportH).Render(strings.Join(leftLines, "\n"))

	// ── Right column: ragdoll equipment ──
	var rightLines []string
	if m.equipFocused {
		rightLines = append(rightLines, invHeaderStyle.Render(" Equipment"))
	} else {
		rightLines = append(rightLines, invEmptyStyle.Render(" Equipment"))
	}
	rightLines = append(rightLines,
		sidebarSubStyle.Render(" "+strings.Repeat("─", rightW-2)),
	)

	// Centre ragdoll vertically: place it roughly 1/4 down from the top of the remaining space.
	ragdollStart := (m.viewportH - len(ragdoll) - len(equipSlots) - 2) / 4
	if ragdollStart < 0 {
		ragdollStart = 0
	}
	for i := 0; i < ragdollStart; i++ {
		rightLines = append(rightLines, "")
	}
	for _, line := range ragdoll {
		rightLines = append(rightLines, line)
	}
	rightLines = append(rightLines, "")
	for i, slot := range equipSlots {
		var row string
		if m.equipFocused && i == m.equipCursor {
			if m.inventory.Equipped[i].Name != "" {
				row = fmt.Sprintf(" > %-11s: [ %s ]", slot, m.inventory.Equipped[i].Name)
			} else {
				row = fmt.Sprintf(" > %-11s: [ Empty ]", slot)
			}
			rightLines = append(rightLines, invCursorStyle.Render(row))
		} else {
			if m.inventory.Equipped[i].Name != "" {
				row = fmt.Sprintf("   %-11s: [ %s ]", slot, m.inventory.Equipped[i].Name)
			} else {
				row = fmt.Sprintf("   %-11s: [ Empty ]", slot)
			}
			rightLines = append(rightLines, invItemStyle.Render(row))
		}
	}

	// Pad to viewportH.
	for len(rightLines) < m.viewportH {
		rightLines = append(rightLines, "")
	}
	rightLines = rightLines[:m.viewportH]

	rightCol := lipgloss.NewStyle().Width(rightW).Height(m.viewportH).Render(strings.Join(rightLines, "\n"))

	return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)
}

// ── Combat screen renderer ────────────────────────────────────────────────────

var (
	combatHeaderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5555"))
	combatStatStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ccd9e0"))
	combatLogStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#768390"))
	combatVictoryStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#58a6ff"))
	combatDefeatStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5555"))
)

// speedLabel returns "Slow", "Normal", or "Fast" for the given combat speed.
func speedLabel(speed int) string {
	switch speed {
	case CombatSpeedSlow:
		return "Slow"
	case CombatSpeedFast:
		return "Fast"
	default:
		return "Normal"
	}
}

// renderHeroPanel renders the left panel with ragdoll art, name, HP bar, and stats.
func renderHeroPanel(m Model, width, height int) string {
	cs := m.combatState
	var lines []string

	// Centre ragdoll vertically in the top portion.
	artH := len(ragdoll)
	statsH := 5 // name + HP bar + HP label + ARM + DMG
	pad := (height - artH - statsH) / 3
	if pad < 0 {
		pad = 0
	}
	for i := 0; i < pad; i++ {
		lines = append(lines, "")
	}
	for _, line := range ragdoll {
		centered := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(line)
		lines = append(lines, centered)
	}
	lines = append(lines, "")

	// Player stats using HP at current round.
	playerHP := hpAtRound(cs.PlayerStartHP, cs.Log, m.combatLogIndex, cs.Player.Name)
	barW := width - 4
	if barW < 5 {
		barW = 5
	}
	if barW > 20 {
		barW = 20
	}
	lines = append(lines, combatStatStyle.Render(fmt.Sprintf("  %s", cs.Player.Name)))
	lines = append(lines, fmt.Sprintf("  %s HP %d/%d", renderProgressBar(playerHP, cs.Player.MaxHP, barW, "#22cc55", "#444c56"), playerHP, cs.Player.MaxHP))
	lines = append(lines, combatStatStyle.Render(fmt.Sprintf("  ARM:%d  DMG:%d-%d  Init:%d", cs.Player.Armour, cs.Player.MinDamage, cs.Player.MaxDamage, cs.Player.Initiative)))

	// Pad to height
	for len(lines) < height {
		lines = append(lines, "")
	}
	lines = lines[:height]

	// Constrain width
	result := make([]string, len(lines))
	for i, l := range lines {
		result[i] = lipgloss.NewStyle().Width(width).Render(l)
	}
	return strings.Join(result, "\n")
}

// renderEnemyPanel renders the right panel with enemy glyph, name, HP bar, and stats.
func renderEnemyPanel(m Model, width, height int) string {
	cs := m.combatState
	var lines []string

	// Enemy glyph in a bordered box.
	var enemyChar rune
	if m.combatDungeonEnemy != nil {
		enemyChar = m.combatDungeonEnemy.Template.Char
	} else if cs.Enemy.Name != "" {
		enemyChar = rune(cs.Enemy.Name[0])
	} else {
		enemyChar = '?'
	}

	glyphStr := string(enemyChar)
	boxContent := lipgloss.NewStyle().
		Width(5).Height(3).
		Align(lipgloss.Center, lipgloss.Center).
		Render(glyphStr)
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#768390")).
		Render(boxContent)

	// Centre the glyph box similarly to the ragdoll.
	boxLines := strings.Split(box, "\n")
	artH := len(boxLines)
	statsH := 5
	pad := (height - artH - statsH) / 3
	if pad < 0 {
		pad = 0
	}
	for i := 0; i < pad; i++ {
		lines = append(lines, "")
	}
	for _, bl := range boxLines {
		centered := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(bl)
		lines = append(lines, centered)
	}
	lines = append(lines, "")

	// Enemy stats using HP at current round.
	enemyHP := hpAtRound(cs.EnemyStartHP, cs.Log, m.combatLogIndex, cs.Enemy.Name)
	barW := width - 4
	if barW < 5 {
		barW = 5
	}
	if barW > 20 {
		barW = 20
	}
	lines = append(lines, combatStatStyle.Render(fmt.Sprintf("  %s", cs.Enemy.Name)))
	lines = append(lines, fmt.Sprintf("  %s HP %d/%d", renderProgressBar(enemyHP, cs.Enemy.MaxHP, barW, "#ff5555", "#444c56"), enemyHP, cs.Enemy.MaxHP))
	lines = append(lines, combatStatStyle.Render(fmt.Sprintf("  ARM:%d  DMG:%d-%d  Init:%d", cs.Enemy.Armour, cs.Enemy.MinDamage, cs.Enemy.MaxDamage, cs.Enemy.Initiative)))

	// Pad to height
	for len(lines) < height {
		lines = append(lines, "")
	}
	lines = lines[:height]

	// Constrain width
	result := make([]string, len(lines))
	for i, l := range lines {
		result[i] = lipgloss.NewStyle().Width(width).Render(l)
	}
	return strings.Join(result, "\n")
}

// renderCombatLog renders the bottom log panel with speed label and visible combat log lines.
func renderCombatLog(m Model, width, height int) string {
	cs := m.combatState
	var lines []string

	// Header with speed label.
	header := combatHeaderStyle.Render(" ⚔ Combat") + "  " +
		combatStatStyle.Render(speedLabel(m.combatSpeed)) + "  " +
		combatLogStyle.Render("[ ] speed")
	lines = append(lines, lipgloss.NewStyle().Width(width).Render(header))
	lines = append(lines, sidebarSubStyle.Render(" "+strings.Repeat("─", width-2)))

	// Visible log lines.
	visibleLines := combatLogLinesUpTo(cs.Log, m.combatLogIndex)
	availRows := height - 2 // subtract header and separator
	if m.combatLogIndex >= cs.Round {
		availRows -= 2 // room for banner + hint
	}
	if availRows < 0 {
		availRows = 0
	}

	// Show last N lines that fit.
	if len(visibleLines) > availRows {
		visibleLines = visibleLines[len(visibleLines)-availRows:]
	}
	for _, l := range visibleLines {
		lines = append(lines, combatLogStyle.Render(" "+l))
	}

	// Victory/Defeated banner after playback completes.
	if m.combatLogIndex >= cs.Round {
		if cs.PlayerWon {
			lines = append(lines, combatVictoryStyle.Render("  Victory!"))
			lines = append(lines, sidebarSubStyle.Render("  press enter to continue"))
		} else {
			lines = append(lines, combatDefeatStyle.Render("  Defeated!"))
			lines = append(lines, sidebarSubStyle.Render("  press enter to quit"))
		}
	}

	// Pad to height.
	for len(lines) < height {
		lines = append(lines, "")
	}
	lines = lines[:height]

	return strings.Join(lines, "\n")
}

// renderCombatScreen fills the viewport with the combat view.
func renderCombatScreen(m Model) string {
	if m.combatState == nil {
		return "No combat active."
	}

	logRows := m.viewportH / 3
	topH := m.viewportH - logRows
	leftW := m.viewportW * 40 / 100
	if leftW < 20 {
		leftW = 20
	}
	rightW := m.viewportW * 40 / 100
	if rightW < 20 {
		rightW = 20
	}

	heroPanel := renderHeroPanel(m, leftW, topH)
	enemyPanel := renderEnemyPanel(m, rightW, topH)
	logPanel := renderCombatLog(m, m.viewportW, logRows)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, heroPanel, enemyPanel)
	return lipgloss.JoinVertical(lipgloss.Left, topRow, logPanel)
}

// ── View composition ──────────────────────────────────────────────────────────

// buildView composes the full terminal view: optional sidebar | map, HUD, key bar.
func buildView(m Model) string {
	// Fullscreen combat takes over the entire viewport.
	if m.screenMode == ScreenCombat {
		return renderCombatScreen(m)
	}

	// Fullscreen help panel takes over the entire viewport.
	if m.showHelpPanel {
		return renderHelpPanel(m)
	}

	// Fullscreen inventory takes over the entire viewport.
	if m.screenMode == ScreenInventory {
		return renderFullscreenInventory(m)
	}

	// Reserve 1 row for HUD.
	mapH := m.viewportH - 1
	if mapH < 1 {
		mapH = 1
	}

	renderMap := func(mapW int) string {
		if m.mode == ModeLocal {
			return renderLocalMap(m, mapW, mapH)
		} else if m.mode == ModeDungeon {
			return renderDungeonMap(m, mapW, mapH)
		}
		return renderWorldMap(m, mapW, mapH)
	}

	var mapView string
	if m.showSidebar {
		mapW := m.viewportW - sidebarContentW - 1 // -1 for the │ separator column
		if mapW < 10 {
			mapW = 10
		}
		mapView = lipgloss.JoinHorizontal(lipgloss.Top, renderSidebar(m, mapH), renderMap(mapW))
	} else if m.showMapPicker {
		mapW := m.viewportW - pickerContentW - 1 // -1 for the │ separator column
		if mapW < 10 {
			mapW = 10
		}
		mapView = lipgloss.JoinHorizontal(lipgloss.Top, renderMap(mapW), renderMapPicker(m, mapH))
	} else {
		mapView = renderMap(m.viewportW)
	}

	return lipgloss.JoinVertical(lipgloss.Left, mapView, renderHUD(m))
}
