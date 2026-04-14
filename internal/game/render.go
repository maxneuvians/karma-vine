package game

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

// renderHUD renders the single-row status bar.
func renderHUD(m Model) string {
	tile := TileAt(m.worldPos.X, m.worldPos.Y, &m)
	clock := formatTime(m.timeOfDay)
	speed := fmt.Sprintf("%d×", m.timeScale)
	var text string
	// In temperature-overlay mode use the same perceived value that drives the
	// map colour, so the °C reading is consistent with what the player sees.
	displayTemp := tile.Temperature
	if m.mapMode == MapModeTemperature {
		displayTemp = perceivedTemperature(tile)
	}
	celsius := tempCelsius(displayTemp, tile.Elevation, m.timeOfDay)
	if m.mode == ModeLocal {
		text = fmt.Sprintf(" %s  %d°C  local (%d, %d)  world (%d, %d)  %s  %s",
			biomeName(tile.Biome),
			celsius,
			m.playerPos.X, m.playerPos.Y,
			m.worldPos.X, m.worldPos.Y,
			clock, speed,
		)
	} else {
		chunkX := m.worldPos.X / 32
		chunkY := m.worldPos.Y / 32
		text = fmt.Sprintf(" %s  elev: %.2f  %d°C  (%d, %d)  chunk (%d, %d)  %s  %s",
			biomeName(tile.Biome),
			tile.Elevation,
			celsius,
			m.worldPos.X, m.worldPos.Y,
			chunkX, chunkY,
			clock, speed,
		)
	}
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

var localCharNames = map[rune]string{
	'♣': "Tree",
	'♠': "Pine",
	'ψ': "Cactus",
	'○': "Rock",
	'⌂': "Shelter",
	'✿': "Flower",
	'd': "Deer",
	'r': "Rabbit",
	'b': "Bird",
	's': "Snake",
	'l': "Lizard",
	'w': "Wolf",
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

// renderSidebar builds a height-row sidebar with a trailing │ on each line.
func renderSidebar(m Model, height int) string {
	if height < 1 {
		height = 1
	}

	var lines []string
	if m.mode == ModeWorld {
		lines = append(lines,
			sbText(sidebarHeaderStyle.Render(" Biomes")),
			sbText(sidebarSubStyle.Render(" "+strings.Repeat("─", 18))),
		)
		for _, e := range biomeLegend {
			lines = append(lines, sbCell(e.char, applyColor(e.color, dimFactor(m.timeOfDay)), e.name))
		}
	} else {
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
			}
			seenObj := make(map[entry]bool)
			seenAni := make(map[entry]bool)
			hasFire := false
			for x := 0; x < LocalMapW; x++ {
				for y := 0; y < LocalMapH; y++ {
					if obj := lm.Objects[x][y]; obj != nil {
						seenObj[entry{obj.Char, obj.Color}] = true
					}
					if lm.Ground[x][y].HasFire {
						hasFire = true
					}
				}
			}
			for _, a := range lm.Animals {
				seenAni[entry{a.Char, a.Color}] = true
			}
			if hasFire {
				lines = append(lines, sbCell('♨', "#ff8800", "Campfire"))
			}
			if len(seenObj) > 0 {
				lines = append(lines, sbText(sidebarSubStyle.Render(" Objects")))
				for e := range seenObj {
					name := localCharNames[e.char]
					if name == "" {
						name = string(e.char)
					}
					lines = append(lines, sbCell(e.char, applyColor(e.color, dimFactor(m.timeOfDay)), name))
				}
			}
			if len(seenAni) > 0 {
				lines = append(lines, sbText(sidebarSubStyle.Render(" Wildlife")))
				for e := range seenAni {
					name := localCharNames[e.char]
					if name == "" {
						name = string(e.char)
					}
					lines = append(lines, sbCell(e.char, applyColor(e.color, dimFactor(m.timeOfDay)), name))
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

// ── Key bar ───────────────────────────────────────────────────────────────────

// renderKeyBar returns a single row of context-sensitive key binding hints.
func renderKeyBar(m Model) string {
	speed := fmt.Sprintf("%d×", m.timeScale)
	var hints string
	if m.mode == ModeLocal {
		hints = fmt.Sprintf(" ↑↓←→/wasd move  esc/< ascend  [/] speed (%s)  ? sidebar  q quit", speed)
	} else {
		hints = fmt.Sprintf(" ↑↓←→/wasd move  enter/> descend  +/- zoom (%d×)  [/] speed (%s)  m map  ? sidebar  q quit", m.worldZoom, speed)
	}
	return keyBarStyle.Render(hints)
}

// ── View composition ──────────────────────────────────────────────────────────

// buildView composes the full terminal view: optional sidebar | map, HUD, key bar.
func buildView(m Model) string {
	// Reserve 2 rows for HUD + key bar.
	mapH := m.viewportH - 2
	if mapH < 1 {
		mapH = 1
	}

	var mapView string
	if m.showSidebar {
		mapW := m.viewportW - sidebarContentW - 1 // -1 for the │ separator column
		if mapW < 10 {
			mapW = 10
		}
		sidebar := renderSidebar(m, mapH)
		if m.mode == ModeLocal {
			mapView = renderLocalMap(m, mapW, mapH)
		} else {
			mapView = renderWorldMap(m, mapW, mapH)
		}
		mapView = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mapView)
	} else if m.showMapPicker {
		mapW := m.viewportW - pickerContentW - 1 // -1 for the │ separator column
		if mapW < 10 {
			mapW = 10
		}
		picker := renderMapPicker(m, mapH)
		if m.mode == ModeLocal {
			mapView = renderLocalMap(m, mapW, mapH)
		} else {
			mapView = renderWorldMap(m, mapW, mapH)
		}
		mapView = lipgloss.JoinHorizontal(lipgloss.Top, mapView, picker)
	} else {
		if m.mode == ModeLocal {
			mapView = renderLocalMap(m, m.viewportW, mapH)
		} else {
			mapView = renderWorldMap(m, m.viewportW, mapH)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, mapView, renderHUD(m), renderKeyBar(m))
}
