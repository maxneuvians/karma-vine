package game

import (
	"math"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// --- dimFactor ---

func TestDimFactor_Noon(t *testing.T) {
	v := dimFactor(0.5)
	if math.Abs(v-1.0) > 0.01 {
		t.Fatalf("dimFactor(0.5) = %v, want ~1.0", v)
	}
}

func TestDimFactor_Midnight(t *testing.T) {
	v := dimFactor(0.0)
	if math.Abs(v-0.15) > 0.01 {
		t.Fatalf("dimFactor(0.0) = %v, want ~0.15", v)
	}
}

// --- applyColor ---

func TestApplyColor_FullBrightness(t *testing.T) {
	// dim=1.0 should return the original color unchanged
	got := applyColor("#2d7a1f", 1.0)
	if got != "#2d7a1f" {
		t.Fatalf("applyColor dim=1.0: got %q, want %q", got, "#2d7a1f")
	}
}

func TestApplyColor_HalfDim(t *testing.T) {
	got := applyColor("#ff0000", 0.5)
	// #ff = 255, 255*0.5 = 127 = #7f
	if got != "#7f0000" {
		t.Fatalf("applyColor dim=0.5: got %q, want #7f0000", got)
	}
}

func TestApplyColor_InvalidPassthrough(t *testing.T) {
	cases := []string{"", "abc", "#gggggg", "#12345"}
	for _, c := range cases {
		if got := applyColor(c, 0.5); got != c {
			t.Fatalf("applyColor(%q) = %q, want unchanged %q", c, got, c)
		}
	}
}

// --- formatTime ---

func TestFormatTime(t *testing.T) {
	cases := []struct {
		tod  float64
		want string
	}{
		{0.0, "00:00"},
		{0.5, "12:00"},
		{0.75, "18:00"},
		{0.25, "06:00"},
	}
	for _, c := range cases {
		if got := formatTime(c.tod); got != c.want {
			t.Fatalf("formatTime(%v) = %q, want %q", c.tod, got, c.want)
		}
	}
}

// --- WindowSizeMsg ---

func TestUpdate_WindowSizeMsg(t *testing.T) {
	m := NewModel()
	next, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	nm := next.(Model)
	if nm.viewportW != 120 || nm.viewportH != 40 {
		t.Fatalf("WindowSizeMsg: got viewport %dx%d, want 120x40", nm.viewportW, nm.viewportH)
	}
}

// --- timeScale keys ---

func TestUpdate_TimeScaleIncrease(t *testing.T) {
	m := NewModel() // timeScale starts at 1
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("]")})
	nm := next.(Model)
	if nm.timeScale != 2 {
		t.Fatalf("timeScale after ] = %d, want 2", nm.timeScale)
	}
}

func TestUpdate_TimeScaleDecrease(t *testing.T) {
	m := NewModel()
	m.timeScale = 5
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("[")})
	nm := next.(Model)
	if nm.timeScale != 2 {
		t.Fatalf("timeScale after [ = %d, want 2", nm.timeScale)
	}
}

func TestUpdate_TimeScaleClampedMax(t *testing.T) {
	m := NewModel()
	m.timeScale = 10
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("]")})
	nm := next.(Model)
	if nm.timeScale != 10 {
		t.Fatalf("timeScale clamped max: got %d, want 10", nm.timeScale)
	}
}

// --- biomeName ---

func TestBiomeName_AllBiomes(t *testing.T) {
	cases := []struct {
		b    Biome
		want string
	}{
		{DeepOcean, "Deep Ocean"},
		{ShallowWater, "Shallow Water"},
		{Beach, "Beach"},
		{Forest, "Forest"},
		{Plains, "Plains"},
		{DenseForest, "Dense Forest"},
		{Desert, "Desert"},
		{Mountains, "Mountains"},
		{Snow, "Snow"},
		{Jungle, "Jungle"},
		{Savanna, "Savanna"},
		{AridSteppe, "Arid Steppe"},
		{Tundra, "Tundra"},
		{Taiga, "Taiga"},
		{Biome(999), "Unknown"},
	}
	for _, c := range cases {
		if got := biomeName(c.b); got != c.want {
			t.Errorf("biomeName(%d) = %q, want %q", c.b, got, c.want)
		}
	}
}

// --- tempCelsius ---

func TestTempCelsius_EquatorialNoon(t *testing.T) {
	// temperature=1.0, elevation=0.36 (sea level), timeOfDay=0.583 (~14:00 peak)
	// base=40, elevAdj=0, timeAdj≈+5 → ~45
	got := tempCelsius(1.0, 0.36, 0.583)
	if got < 40 || got > 50 {
		t.Errorf("tempCelsius(1.0, 0.36, 0.583) = %d, want ~45", got)
	}
}

func TestTempCelsius_PolarMidnight(t *testing.T) {
	// temperature=0.0, elevation=0.36, timeOfDay=0.0 → base=-20, timeAdj≈-5 → ~-25
	got := tempCelsius(0.0, 0.36, 0.0)
	if got > -10 {
		t.Errorf("tempCelsius(0.0, 0.36, 0.0) = %d, expected cold (< -10)", got)
	}
}

func TestTempCelsius_HighElevation(t *testing.T) {
	// High elevation should produce a cooler result than sea level for same temperature
	low := tempCelsius(0.5, 0.36, 0.5)
	high := tempCelsius(0.5, 0.80, 0.5)
	if high >= low {
		t.Errorf("high elevation (%d) should be cooler than sea level (%d)", high, low)
	}
}

// --- renderWorldMap ---

func TestRenderWorldMap_Basic(t *testing.T) {
	m := NewModel()
	out := renderWorldMap(m, 40, 18)
	if out == "" {
		t.Fatal("renderWorldMap returned empty string")
	}
}

func TestRenderWorldMap_ZoomClamped(t *testing.T) {
	m := NewModel()
	m.worldZoom = 0 // should be clamped to 1 inside renderWorldMap
	out := renderWorldMap(m, 10, 5)
	if out == "" {
		t.Fatal("renderWorldMap with worldZoom=0 returned empty string")
	}
}

func TestRenderWorldMap_HighZoom(t *testing.T) {
	m := NewModel()
	m.worldZoom = 8
	out := renderWorldMap(m, 20, 10)
	if out == "" {
		t.Fatal("renderWorldMap with worldZoom=8 returned empty string")
	}
}

// --- renderLocalMap ---

func TestRenderLocalMap_NoMapNoCache(t *testing.T) {
	m := NewModel()
	m.localMap = nil
	got := renderLocalMap(m, 40, 18)
	if got != "Local map not loaded." {
		t.Fatalf("renderLocalMap no map: got %q, want %q", got, "Local map not loaded.")
	}
}

func TestRenderLocalMap_WithLocalMap(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = GenerateLocalMap(0, 0, 42, Forest)
	m.playerPos = LocalCoord{X: 10, Y: 10}
	out := renderLocalMap(m, 40, 18)
	if out == "" {
		t.Fatal("renderLocalMap returned empty string")
	}
}

func TestRenderLocalMap_FromCache(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	lm := GenerateLocalMap(0, 0, 42, Plains)
	m.localMap = nil
	m.localCache[m.worldPos] = lm
	out := renderLocalMap(m, 40, 18)
	if out == "" {
		t.Fatal("renderLocalMap from cache returned empty string")
	}
}

func TestRenderLocalMap_FireAnimalObjectBranches(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.timeOfDay = 0.0 // dark → globalDim=0.15 so LitMap can exceed it

	lm := &LocalMap{}
	for x := 0; x < LocalMapW; x++ {
		for y := 0; y < LocalMapH; y++ {
			lm.Ground[x][y] = Ground{Char: '.', Color: "#5aad3f", Passable: true}
		}
	}
	// Fire at (5,5) — also set LitMap high to exercise the > globalDim branch
	lm.Ground[5][5].HasFire = true
	lm.LitMap[5][5] = 1.0
	// Object at (6,5)
	lm.Objects[6][5] = &Object{Char: '♣', Color: "#2d7a1f", Blocking: false}
	// Animal at (7,5)
	lm.Animals = []*Animal{{X: 7, Y: 5, Char: 'd', Color: "#c8a46a"}}

	m.localMap = lm
	m.playerPos = LocalCoord{X: 10, Y: 10}

	out := renderLocalMap(m, 40, 18)
	if out == "" {
		t.Fatal("renderLocalMap with fire/animal/object returned empty string")
	}
}

func TestRenderLocalMap_CameraClampedRightBottom(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = GenerateLocalMap(0, 0, 1, Plains)
	// Position near bottom-right corner → triggers right/bottom camera clamping
	m.playerPos = LocalCoord{X: LocalMapW - 1, Y: LocalMapH - 1}
	out := renderLocalMap(m, 40, 18)
	if out == "" {
		t.Fatal("renderLocalMap right/bottom clamp returned empty string")
	}
}

// --- renderHUD ---

func TestRenderHUD_WorldMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.timeOfDay = 0.5
	m.timeScale = 1
	out := renderHUD(m)
	if out == "" {
		t.Fatal("renderHUD ModeWorld returned empty string")
	}
}

func TestRenderHUD_LocalMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.timeOfDay = 0.25
	m.timeScale = 2
	m.playerPos = LocalCoord{X: 5, Y: 5}
	out := renderHUD(m)
	if out == "" {
		t.Fatal("renderHUD ModeLocal returned empty string")
	}
}

// --- joinVertical ---

func TestJoinVertical_Basic(t *testing.T) {
	out := joinVertical("top", "bottom")
	if out == "" {
		t.Fatal("joinVertical returned empty string")
	}
}

// --- sbCell / sbText ---

func TestSbCell_Basic(t *testing.T) {
	out := sbCell('♣', "#2d7a1f", "Forest")
	if out == "" {
		t.Fatal("sbCell returned empty string")
	}
}

func TestSbText_Basic(t *testing.T) {
	out := sbText("hello world")
	if out == "" {
		t.Fatal("sbText returned empty string")
	}
}

// --- renderSidebar ---

func TestRenderSidebar_WorldMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	out := renderSidebar(m, 20)
	if out == "" {
		t.Fatal("renderSidebar ModeWorld returned empty string")
	}
}

func TestRenderSidebar_HeightZero(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	out := renderSidebar(m, 0) // should clamp to 1
	if out == "" {
		t.Fatal("renderSidebar height=0 returned empty string")
	}
}

func TestRenderSidebar_LocalMode_NoMap(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = nil
	out := renderSidebar(m, 20)
	if out == "" {
		t.Fatal("renderSidebar ModeLocal no map returned empty string")
	}
}

func TestRenderSidebar_LocalMode_FromCache(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = nil
	m.localCache[m.worldPos] = GenerateLocalMap(0, 0, 1, Forest)
	out := renderSidebar(m, 30)
	if out == "" {
		t.Fatal("renderSidebar ModeLocal from cache returned empty string")
	}
}

func TestRenderSidebar_LocalMode_WithFireObjectsAnimals(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal

	lm := &LocalMap{}
	// Fire
	lm.Ground[0][0].HasFire = true
	// Known object (in localCharNames)
	lm.Objects[1][0] = &Object{Char: '♣', Color: "#2d7a1f", Blocking: false}
	// Unknown object (not in localCharNames → uses string(char) as name)
	lm.Objects[2][0] = &Object{Char: 'X', Color: "#ffffff", Blocking: false}
	// Known animal
	lm.Animals = []*Animal{
		{X: 5, Y: 5, Char: 'd', Color: "#c8a46a"},
		// Unknown animal glyph
		{X: 6, Y: 5, Char: 'Z', Color: "#aaaaaa"},
	}
	m.localMap = lm

	out := renderSidebar(m, 40)
	if out == "" {
		t.Fatal("renderSidebar with fire/objects/animals returned empty string")
	}
}

// --- renderKeyBar ---

func TestRenderKeyBar_WorldMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.worldZoom = 2
	out := renderKeyBar(m)
	if out == "" {
		t.Fatal("renderKeyBar ModeWorld returned empty string")
	}
}

func TestRenderKeyBar_LocalMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	out := renderKeyBar(m)
	if out == "" {
		t.Fatal("renderKeyBar ModeLocal returned empty string")
	}
}

// --- buildView ---

func TestBuildView_WorldMode(t *testing.T) {
	m := NewModel()
	m.viewportW = 80
	m.viewportH = 24
	m.mode = ModeWorld
	out := buildView(m)
	if out == "" {
		t.Fatal("buildView ModeWorld returned empty string")
	}
}

func TestBuildView_WorldModeWithSidebar(t *testing.T) {
	m := NewModel()
	m.viewportW = 80
	m.viewportH = 24
	m.mode = ModeWorld
	m.showSidebar = true
	out := buildView(m)
	if out == "" {
		t.Fatal("buildView ModeWorld with sidebar returned empty string")
	}
}

func TestBuildView_LocalMode(t *testing.T) {
	m := NewModel()
	m.viewportW = 80
	m.viewportH = 24
	m.mode = ModeLocal
	m.localMap = GenerateLocalMap(0, 0, 42, Forest)
	m.playerPos = LocalCoord{X: 10, Y: 10}
	out := buildView(m)
	if out == "" {
		t.Fatal("buildView ModeLocal returned empty string")
	}
}

func TestBuildView_LocalModeWithSidebar(t *testing.T) {
	m := NewModel()
	m.viewportW = 80
	m.viewportH = 24
	m.mode = ModeLocal
	m.localMap = GenerateLocalMap(0, 0, 42, Forest)
	m.playerPos = LocalCoord{X: 10, Y: 10}
	m.showSidebar = true
	out := buildView(m)
	if out == "" {
		t.Fatal("buildView ModeLocal with sidebar returned empty string")
	}
}

func TestBuildView_SmallViewport_MapHClamped(t *testing.T) {
	// viewportH=2 → mapH = 2-2 = 0 → clamped to 1
	m := NewModel()
	m.viewportW = 40
	m.viewportH = 2
	out := buildView(m)
	if out == "" {
		t.Fatal("buildView small viewport returned empty string")
	}
}

func TestBuildView_SidebarNarrowViewport(t *testing.T) {
	// viewportW=5 → mapW = 5-sidebarContentW-1 = negative → clamped to 10
	m := NewModel()
	m.viewportW = 5
	m.viewportH = 24
	m.showSidebar = true
	out := buildView(m)
	if out == "" {
		t.Fatal("buildView narrow sidebar viewport returned empty string")
	}
}

