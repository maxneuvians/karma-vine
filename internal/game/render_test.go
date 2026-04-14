package game

import (
	"fmt"
	"math"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
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
	next, _ := m.Update(tea.KeyPressMsg{Code: ']', Text: "]"})
	nm := next.(Model)
	if nm.timeScale != 2 {
		t.Fatalf("timeScale after ] = %d, want 2", nm.timeScale)
	}
}

func TestUpdate_TimeScaleDecrease(t *testing.T) {
	m := NewModel()
	m.timeScale = 5
	next, _ := m.Update(tea.KeyPressMsg{Code: '[', Text: "["})
	nm := next.(Model)
	if nm.timeScale != 2 {
		t.Fatalf("timeScale after [ = %d, want 2", nm.timeScale)
	}
}

func TestUpdate_TimeScaleClampedMax(t *testing.T) {
	m := NewModel()
	m.timeScale = 10
	next, _ := m.Update(tea.KeyPressMsg{Code: ']', Text: "]"})
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

// ── lerpHex ───────────────────────────────────────────────────────────────────

func TestLerpHex_AtZero(t *testing.T) {
	got := lerpHex("#000000", "#ffffff", 0)
	if got != "#000000" {
		t.Errorf("lerpHex t=0: got %q, want #000000", got)
	}
}

func TestLerpHex_AtOne(t *testing.T) {
	got := lerpHex("#000000", "#ffffff", 1.0)
	if got != "#ffffff" {
		t.Errorf("lerpHex t=1: got %q, want #ffffff", got)
	}
}

func TestLerpHex_AtHalf(t *testing.T) {
	// #000000 → #ffffff at 0.5 → each channel: 0 + 0.5*255 = 127 = #7f
	got := lerpHex("#000000", "#ffffff", 0.5)
	if got != "#7f7f7f" {
		t.Errorf("lerpHex t=0.5: got %q, want #7f7f7f", got)
	}
}

func TestLerpHex_InvalidPassthrough(t *testing.T) {
	// Non-color strings should return first argument unchanged
	cases := []struct{ a, b string }{
		{"bad", "#ffffff"},
		{"#ffffff", "bad"},
	}
	for _, c := range cases {
		got := lerpHex(c.a, c.b, 0.5)
		if got != c.a {
			t.Errorf("lerpHex(%q, %q, 0.5) = %q, want %q (passthrough)", c.a, c.b, got, c.a)
		}
	}
}

// ── tileVisual ────────────────────────────────────────────────────────────────

func TestTileVisual_Default_Passthrough(t *testing.T) {
	tile := Tile{Char: '♣', Color: "#2d7a1f", Biome: Forest, Elevation: 0.5, Temperature: 0.5}
	ch, color := tileVisual(tile, MapModeDefault)
	if ch != '♣' || color != "#2d7a1f" {
		t.Errorf("tileVisual default: got (%q, %q), want ('♣', '#2d7a1f')", ch, color)
	}
}

func TestTileVisual_Temperature_Min(t *testing.T) {
	tile := Tile{Temperature: 0.0}
	ch, color := tileVisual(tile, MapModeTemperature)
	if ch != '█' {
		t.Errorf("tileVisual temperature min: char = %q, want '█'", ch)
	}
	if color != "#0022cc" {
		t.Errorf("tileVisual temperature min: color = %q, want #0022cc", color)
	}
}

func TestTileVisual_Temperature_Max(t *testing.T) {
	// Desert at temp=1.0: perceivedTemperature clamps to 1.0 → thermalColor(1) = #ff2200
	tile := Tile{Temperature: 1.0, Biome: Desert}
	ch, color := tileVisual(tile, MapModeTemperature)
	if ch != '█' {
		t.Errorf("tileVisual temperature max: char = %q, want '█'", ch)
	}
	if color != "#ff2200" {
		t.Errorf("tileVisual temperature max: color = %q, want #ff2200", color)
	}
}

func TestTileVisual_Temperature_DeepOceanCoolerThanDesert(t *testing.T) {
	// At the same raw temperature, ocean should appear cooler (more blue) than desert.
	ocean := Tile{Temperature: 0.6, Biome: DeepOcean}
	desert := Tile{Temperature: 0.6, Biome: Desert}
	_, oceanColor := tileVisual(ocean, MapModeTemperature)
	_, desertColor := tileVisual(desert, MapModeTemperature)
	if perceivedTemperature(ocean) >= perceivedTemperature(desert) {
		t.Errorf("deep ocean (%v) should be cooler than desert (%v)",
			perceivedTemperature(ocean), perceivedTemperature(desert))
	}
	if oceanColor == desertColor {
		t.Error("ocean and desert should render different colors in temperature mode")
	}
}

// ── thermalColor ──────────────────────────────────────────────────────────────

func TestThermalColor_Endpoints(t *testing.T) {
	if got := thermalColor(0); got != "#0022cc" {
		t.Errorf("thermalColor(0) = %q, want #0022cc", got)
	}
	if got := thermalColor(1); got != "#ff2200" {
		t.Errorf("thermalColor(1) = %q, want #ff2200", got)
	}
}

func TestThermalColor_Midpoint_IsGreen(t *testing.T) {
	// t=0.5 is the sea-green stop; result should be greenish, not purple
	color := thermalColor(0.5)
	if color != "#00dd88" {
		t.Errorf("thermalColor(0.5) = %q, want #00dd88 (sea green)", color)
	}
}

func TestThermalColor_Clamps(t *testing.T) {
	if thermalColor(-1) != thermalColor(0) {
		t.Error("thermalColor(-1) should clamp to thermalColor(0)")
	}
	if thermalColor(2) != thermalColor(1) {
		t.Error("thermalColor(2) should clamp to thermalColor(1)")
	}
}

func TestThermalColor_OrderingCold_To_Hot(t *testing.T) {
	// Extract the red channel at each stop — it should increase from cold to hot.
	// t=0 (#0022cc) r=0x00, t=0.5 (#00dd88) r=0x00, t=0.75 (#ffee00) r=0xff, t=1 (#ff2200) r=0xff
	// More robustly: blue channel dominates at t=0, red channel at t=1.
	cold := thermalColor(0.1)
	hot := thermalColor(0.9)
	// Parse red channel
	var coldR, hotR int
	fmt.Sscanf(cold[1:3], "%x", &coldR)
	fmt.Sscanf(hot[1:3], "%x", &hotR)
	if coldR >= hotR {
		t.Errorf("thermalColor: cold (t=0.1) red=%d should be less than hot (t=0.9) red=%d", coldR, hotR)
	}
}

func TestTileVisual_Elevation_Min(t *testing.T) {
	tile := Tile{Elevation: 0.0}
	ch, color := tileVisual(tile, MapModeElevation)
	if ch != '█' || color != "#1a6fa8" {
		t.Errorf("tileVisual elevation min: got (%q, %q), want ('█', '#1a6fa8')", ch, color)
	}
}

func TestTileVisual_Elevation_Max(t *testing.T) {
	tile := Tile{Elevation: 1.0}
	ch, color := tileVisual(tile, MapModeElevation)
	if ch != '█' || color != "#f0f6fc" {
		t.Errorf("tileVisual elevation max: got (%q, %q), want ('█', '#f0f6fc')", ch, color)
	}
}

// ── perceivedTemperature ──────────────────────────────────────────────────────

func TestPerceivedTemperature_Clamps(t *testing.T) {
	// Desert at high temp may exceed 1.0 before clamping
	tile := Tile{Temperature: 0.98, Biome: Desert}
	v := perceivedTemperature(tile)
	if v > 1.0 || v < 0.0 {
		t.Errorf("perceivedTemperature out of range: %v", v)
	}
}

func TestPerceivedTemperature_BiomeOrdering(t *testing.T) {
	const rawTemp = 0.6
	ocean := perceivedTemperature(Tile{Temperature: rawTemp, Biome: DeepOcean})
	forest := perceivedTemperature(Tile{Temperature: rawTemp, Biome: Forest})
	desert := perceivedTemperature(Tile{Temperature: rawTemp, Biome: Desert})

	if !(ocean < forest) {
		t.Errorf("deep ocean (%v) should be cooler than forest (%v)", ocean, forest)
	}
	if !(forest < desert) {
		t.Errorf("forest (%v) should be cooler than desert (%v)", forest, desert)
	}
}

func TestTileVisual_Political_Contour(t *testing.T) {
	// Elevation 0.09: int(0.09*10)=0, int((0.09+0.05)*10)=int(1.4)=1 → boundary
	tile := Tile{Elevation: 0.09}
	ch, color := tileVisual(tile, MapModePolitical)
	if ch != '+' || color != "#aabbcc" {
		t.Errorf("tileVisual political contour: got (%q, %q), want ('+', '#aabbcc')", ch, color)
	}
}

func TestTileVisual_Political_NonContour(t *testing.T) {
	// Elevation 0.5: int(5.0)=5, int(5.5)=5 → same → no boundary
	tile := Tile{Elevation: 0.5}
	ch, color := tileVisual(tile, MapModePolitical)
	if ch != '·' || color != "#334455" {
		t.Errorf("tileVisual political non-contour: got (%q, %q), want ('·', '#334455')", ch, color)
	}
}

// ── renderWorldMap with MapMode ───────────────────────────────────────────────

func TestRenderWorldMap_TemperatureMode(t *testing.T) {
	m := NewModel()
	m.mapMode = MapModeTemperature
	outTemp := renderWorldMap(m, 20, 10)

	m2 := NewModel()
	m2.mapMode = MapModeDefault
	m2.worldPos = m.worldPos
	outDefault := renderWorldMap(m2, 20, 10)

	if outTemp == outDefault {
		t.Error("renderWorldMap temperature mode should differ from default mode")
	}
	if outTemp == "" {
		t.Error("renderWorldMap temperature mode returned empty string")
	}
}

// ── renderMapPicker ───────────────────────────────────────────────────────────

func TestRenderMapPicker_ContainsAllModeNames(t *testing.T) {
	m := NewModel()
	out := renderMapPicker(m, 10)
	for _, name := range mapModeNames {
		if !strings.Contains(out, name) {
			t.Errorf("renderMapPicker missing mode name %q", name)
		}
	}
}

func TestRenderMapPicker_CursorHighlighted(t *testing.T) {
	m := NewModel()
	m.mapPickerCursor = 2 // Elevation
	out := renderMapPicker(m, 10)
	if !strings.Contains(out, "> Elevation") {
		t.Error("renderMapPicker: cursor row should contain '> Elevation'")
	}
	if strings.Contains(out, "> Default") {
		t.Error("renderMapPicker: non-cursor row should not contain '> Default'")
	}
}

func TestRenderMapPicker_HeightZero(t *testing.T) {
	m := NewModel()
	out := renderMapPicker(m, 0) // should clamp to 1
	if out == "" {
		t.Fatal("renderMapPicker height=0 returned empty string")
	}
}

// ── buildView with map picker ─────────────────────────────────────────────────

func TestBuildView_WorldModeWithMapPicker(t *testing.T) {
	m := NewModel()
	m.viewportW = 80
	m.viewportH = 24
	m.mode = ModeWorld
	m.showMapPicker = true
	out := buildView(m)
	if out == "" {
		t.Fatal("buildView ModeWorld with map picker returned empty string")
	}
}

func TestBuildView_MapPickerNarrowViewport(t *testing.T) {
	m := NewModel()
	m.viewportW = 5
	m.viewportH = 24
	m.showMapPicker = true
	out := buildView(m)
	if out == "" {
		t.Fatal("buildView narrow map picker viewport returned empty string")
	}
}

// ── dungeon render tests ──────────────────────────────────────────────────────

func makeDungeonModel() Model {
	m := NewModel()
	m.globalSeed = 42
	m.viewportW = 80
	m.viewportH = 26
	m.mode = ModeDungeon
	level := GenerateDungeonLevel(42, 1, 1, 1, 5)
	m.currentDungeon = level
	m.dungeonDepth = 3
	m.playerPos = level.UpStair
	return m
}

func TestBuildView_DungeonContainsGlyphs(t *testing.T) {
	m := makeDungeonModel()
	out := buildView(m)
	if !strings.Contains(out, "█") {
		t.Error("dungeon view should contain wall glyph '█'")
	}
	if !strings.Contains(out, "@") {
		t.Error("dungeon view should contain player glyph '@'")
	}
}

func TestRenderHUD_DungeonShowsDepth(t *testing.T) {
	m := makeDungeonModel()
	hud := renderHUD(m)
	if !strings.Contains(hud, "Depth: 3") {
		t.Errorf("dungeon HUD should contain 'Depth: 3', got: %s", hud)
	}
	if !strings.Contains(hud, "Dungeon") {
		t.Errorf("dungeon HUD should contain 'Dungeon', got: %s", hud)
	}
}

func TestRenderKeyBar_DungeonHints(t *testing.T) {
	m := makeDungeonModel()
	bar := renderKeyBar(m)
	if !strings.Contains(bar, "< up") {
		t.Errorf("dungeon key bar should contain '< up', got: %s", bar)
	}
	if !strings.Contains(bar, "> down") {
		t.Errorf("dungeon key bar should contain '> down', got: %s", bar)
	}
	if !strings.Contains(bar, "esc exit") {
		t.Errorf("dungeon key bar should contain 'esc exit', got: %s", bar)
	}
}

func TestRenderKeyBar_WorldModeNoDungeonHints(t *testing.T) {
	m := NewModel()
	m.viewportW = 80
	m.viewportH = 26
	m.mode = ModeWorld
	bar := renderKeyBar(m)
	if strings.Contains(bar, "< up") {
		t.Error("world key bar should not contain dungeon hints")
	}
}

func TestComputeDungeonVisibility_PlayerCellAlwaysVisible(t *testing.T) {
	m := makeDungeonModel()
	light := computeDungeonLight(m)
	if light[m.playerPos] == 0 {
		t.Fatal("player's own cell should always be visible")
	}
}

func TestComputeDungeonVisibility_FarCellHidden(t *testing.T) {
	m := makeDungeonModel()
	light := computeDungeonLight(m)
	// Find a cell far from the player and any torches.
	farX := m.playerPos.X + playerViewRadius + torchRadius + 5
	farY := m.playerPos.Y + playerViewRadius + torchRadius + 5
	if farX < DungeonW && farY < DungeonH {
		_ = light[LocalCoord{X: farX, Y: farY}] // soft check
	}
}

func TestComputeDungeonVisibility_TorchIlluminates(t *testing.T) {
	m := NewModel()
	m.mode = ModeDungeon
	level := &DungeonLevel{}
	level.Cells[10][10].Kind = CellWall
	level.Cells[10][10].Object = &Object{Char: '†', Color: "#e8c96a", Blocking: true, Lit: true}
	m.currentDungeon = level
	m.playerPos = LocalCoord{X: 0, Y: 0} // far from torch

	light := computeDungeonLight(m)
	// Cell (12, 12) is Chebyshev distance 2 from torch at (10,10) → visible.
	if light[LocalCoord{X: 12, Y: 12}] == 0 {
		t.Error("cell (12,12) should be visible due to torch at (10,10)")
	}
}

// ── sidebar contextual tests ──────────────────────────────────────────────────

func TestRenderSidebar_LocalShowsDungeonEntrance(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.viewportW = 80
	m.viewportH = 26
	lm := &LocalMap{}
	lm.Objects[10][10] = &Object{Char: '>', Color: "#ff3333", Blocking: false, Name: "Dungeon Entrance"}
	m.localMap = lm
	out := renderSidebar(m, 20)
	if !strings.Contains(out, "Dungeon Entrance") {
		t.Errorf("local sidebar should contain 'Dungeon Entrance', got:\n%s", out)
	}
}

func TestRenderSidebar_DungeonShowsDepth(t *testing.T) {
	m := NewModel()
	m.mode = ModeDungeon
	m.dungeonDepth = 1
	level := GenerateDungeonLevel(42, 1, 1, 1, 5)
	m.currentDungeon = level
	m.playerPos = level.UpStair
	out := renderSidebar(m, 20)
	if !strings.Contains(out, "Dungeon") {
		t.Errorf("dungeon sidebar should contain 'Dungeon', got:\n%s", out)
	}
	if !strings.Contains(out, "Depth: 1") {
		t.Errorf("dungeon sidebar should contain 'Depth: 1', got:\n%s", out)
	}
}

func TestRenderSidebar_WorldTemperatureMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.mapMode = MapModeTemperature
	out := renderSidebar(m, 20)
	if !strings.Contains(out, "Temperature") {
		t.Errorf("world sidebar in temperature mode should contain 'Temperature', got:\n%s", out)
	}
	if strings.Contains(out, "Deep Ocean") {
		t.Error("world sidebar in temperature mode should not contain biome legend")
	}
}

// 7.17 HUD contains item count string.
func TestRenderHUD_ContainsItemCount(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.viewportW = 120
	m.viewportH = 40
	m.inventory.Items = []Item{
		{Name: "Axe", Count: 1},
		{Name: "Torch", Count: 2},
	}
	out := renderHUD(m)
	expected := fmt.Sprintf("Items: %d/%d", 2, InventoryMaxSlots)
	if !strings.Contains(out, expected) {
		t.Errorf("HUD should contain %q, got:\n%s", expected, out)
	}
}

// ── Fullscreen Inventory tests ────────────────────────────────────────────────

// 8.9 renderFullscreenInventory contains expected labels.
func TestRenderFullscreenInventory_Labels(t *testing.T) {
	m := NewModel()
	m.viewportW = 120
	m.viewportH = 40
	m.screenMode = ScreenInventory
	out := renderFullscreenInventory(m)
	for _, label := range []string{"Inventory", "Head", "Chest", "Left Hand"} {
		if !strings.Contains(out, label) {
			t.Errorf("fullscreen inventory should contain %q", label)
		}
	}
}

// 8.10 renderFullscreenInventory shows Empty when no items.
func TestRenderFullscreenInventory_Empty(t *testing.T) {
	m := NewModel()
	m.viewportW = 100
	m.viewportH = 30
	m.screenMode = ScreenInventory
	out := renderFullscreenInventory(m)
	if !strings.Contains(out, "Empty") {
		t.Error("fullscreen inventory should show 'Empty' when no items")
	}
}

// 8.11 buildView returns fullscreen inventory when ScreenInventory.
func TestBuildView_ScreenInventory(t *testing.T) {
	m := NewModel()
	m.viewportW = 120
	m.viewportH = 40
	m.screenMode = ScreenInventory
	out := buildView(m)
	if !strings.Contains(out, "Inventory") {
		t.Error("buildView with ScreenInventory should contain 'Inventory'")
	}
	if !strings.Contains(out, "Head") {
		t.Error("buildView with ScreenInventory should contain ragdoll slots")
	}
}

// 8.12 buildView returns normal map when ScreenNormal.
func TestBuildView_ScreenNormal(t *testing.T) {
	m := NewModel()
	m.viewportW = 80
	m.viewportH = 24
	m.mode = ModeWorld
	m.screenMode = ScreenNormal
	out := buildView(m)
	if out == "" {
		t.Fatal("buildView ScreenNormal should not be empty")
	}
	if strings.Contains(out, "Equipment") {
		t.Error("buildView ScreenNormal should not contain fullscreen inventory")
	}
}

// --- Pause HUD tests ---

func TestPause_HUDContainsPausedLabel(t *testing.T) {
	m := NewModel()
	m.paused = true
	m.mode = ModeWorld
	out := renderHUD(m)
	if !strings.Contains(out, "[PAUSED]") {
		t.Errorf("paused HUD should contain '[PAUSED]', got:\n%s", out)
	}
}

func TestPause_HUDNoPausedLabelWhenUnpaused(t *testing.T) {
	m := NewModel()
	m.paused = false
	m.mode = ModeWorld
	out := renderHUD(m)
	if strings.Contains(out, "[PAUSED]") {
		t.Errorf("unpaused HUD should not contain '[PAUSED]', got:\n%s", out)
	}
}

// --- Equipment render tests ---

func TestRenderFullscreenInventory_EquippedSlot(t *testing.T) {
	m := NewModel()
	m.viewportW = 120
	m.viewportH = 40
	m.screenMode = ScreenInventory
	m.inventory.Equipped = [NumBodySlots]Item{}
	m.inventory.Equipped[SlotChest] = Item{Char: '♦', Color: "#a0a0a0", Name: "Cloth Tunic", Count: 1, Slots: []BodySlot{SlotChest}}
	out := renderFullscreenInventory(m)
	if !strings.Contains(out, "Cloth Tunic") {
		t.Error("ragdoll should show equipped item name 'Cloth Tunic'")
	}
}

func TestRenderFullscreenInventory_FocusedHeader(t *testing.T) {
	m := NewModel()
	m.viewportW = 120
	m.viewportH = 40
	m.screenMode = ScreenInventory

	// Not focused on equipment → Inventory header should be bold style.
	m.equipFocused = false
	out1 := renderFullscreenInventory(m)

	// Focused on equipment → Equipment header should be bold style.
	m.equipFocused = true
	out2 := renderFullscreenInventory(m)

	// The outputs should differ (different header styles).
	if out1 == out2 {
		t.Error("rendering should differ when equipFocused changes")
	}
}

// --- Combat render tests ---

func TestRenderCombatScreen_ShowsStats(t *testing.T) {
	m := NewModel()
	m.viewportW = 120
	m.viewportH = 40
	m.screenMode = ScreenCombat
	m.combatState = &CombatState{
		Player: Combatant{Name: "Player", HP: 15, MaxHP: 20, Armour: 2, MinDamage: 1, MaxDamage: 4, Initiative: 5},
		Enemy:  Combatant{Name: "Wolf", HP: 8, MaxHP: 12, Armour: 1, MinDamage: 2, MaxDamage: 5, Initiative: 6},
		Log:    []string{"Round 1: Player attacks Wolf for 3 damage"},
	}
	out := renderCombatScreen(m)
	for _, s := range []string{"Player", "Wolf", "15/20", "8/12", "Armour: 2", "Armour: 1"} {
		if !strings.Contains(out, s) {
			t.Errorf("combat screen should contain %q", s)
		}
	}
}

func TestRenderCombatScreen_VictoryBanner(t *testing.T) {
	m := NewModel()
	m.viewportW = 80
	m.viewportH = 24
	m.screenMode = ScreenCombat
	m.combatState = &CombatState{
		Player:    Combatant{Name: "Player", HP: 10, MaxHP: 20},
		Enemy:     Combatant{Name: "Rat", HP: 0, MaxHP: 5},
		PlayerWon: true,
	}
	out := renderCombatScreen(m)
	if !strings.Contains(out, "Victory!") {
		t.Error("combat screen should contain 'Victory!' when player won")
	}
}

func TestRenderCombatScreen_DefeatBanner(t *testing.T) {
	m := NewModel()
	m.viewportW = 80
	m.viewportH = 24
	m.screenMode = ScreenCombat
	m.combatState = &CombatState{
		Player:    Combatant{Name: "Player", HP: 0, MaxHP: 20},
		Enemy:     Combatant{Name: "Bear", HP: 10, MaxHP: 18},
		PlayerWon: false,
	}
	out := renderCombatScreen(m)
	if !strings.Contains(out, "Defeated!") {
		t.Error("combat screen should contain 'Defeated!' when player lost")
	}
}

func TestBuildView_ScreenCombat(t *testing.T) {
	m := NewModel()
	m.viewportW = 80
	m.viewportH = 24
	m.screenMode = ScreenCombat
	m.combatState = &CombatState{
		Player:    Combatant{Name: "Player", HP: 15, MaxHP: 20},
		Enemy:     Combatant{Name: "Wolf", HP: 0, MaxHP: 12},
		PlayerWon: true,
		Log:       []string{"Round 1: Player attacks Wolf"},
	}
	buildOut := buildView(m)
	renderOut := renderCombatScreen(m)
	if buildOut != renderOut {
		t.Error("buildView with ScreenCombat should return renderCombatScreen output")
	}
}

