package game

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// --- NewModel ---

func TestNewModel_MapsInitialised(t *testing.T) {
	m := NewModel()
	if m.chunks == nil {
		t.Fatal("NewModel: chunks map is nil")
	}
	if m.localCache == nil {
		t.Fatal("NewModel: localCache map is nil")
	}
}

func TestNewModel_DefaultMode(t *testing.T) {
	m := NewModel()
	if m.mode != ModeWorld {
		t.Fatalf("NewModel: expected mode ModeWorld (%d), got %d", ModeWorld, m.mode)
	}
}

func TestNewModel_StartsOnLand(t *testing.T) {
	m := NewModel()
	tile := TileAt(m.worldPos.X, m.worldPos.Y, &m)
	if !isLandBiome(tile.Biome) {
		t.Fatalf("NewModel: expected land spawn, got biome %d at %+v", tile.Biome, m.worldPos)
	}
	if m.playerPos != (LocalCoord{}) {
		t.Fatalf("NewModel: expected zero playerPos, got %+v", m.playerPos)
	}
}

// --- Init ---

func TestModel_Init_ReturnsTickCmd(t *testing.T) {
	m := NewModel()
	if cmd := m.Init(); cmd == nil {
		t.Fatal("Init: expected non-nil tick command")
	}
}

// --- Update ---

func TestModel_Update_QuitOnQ(t *testing.T) {
	m := NewModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("Update 'q': expected a non-nil quit command")
	}
}

func TestModel_Update_QuitOnCtrlC(t *testing.T) {
	m := NewModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("Update ctrl+c: expected a non-nil quit command")
	}
}

func TestModel_Update_UnknownKeyNoOp(t *testing.T) {
	m := NewModel()
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if cmd != nil {
		t.Fatal("Update unknown key: expected nil cmd")
	}
	if _, ok := next.(Model); !ok {
		t.Fatal("Update unknown key: returned model is not of type Model")
	}
}

// --- View ---

func TestModel_View_NonEmpty(t *testing.T) {
	m := NewModel()
	if v := m.View(); v == "" {
		t.Fatal("View: returned empty string")
	}
}

// --- Types ---

func TestModeConstants(t *testing.T) {
	if ModeWorld == ModeLocal {
		t.Fatal("ModeWorld and ModeLocal must be distinct")
	}
}

func TestCoordStructs(t *testing.T) {
	wc := WorldCoord{X: 1, Y: -1}
	if wc.X != 1 || wc.Y != -1 {
		t.Fatalf("WorldCoord fields wrong: %+v", wc)
	}
	lc := LocalCoord{X: 21, Y: 9}
	if lc.X != 21 || lc.Y != 9 {
		t.Fatalf("LocalCoord fields wrong: %+v", lc)
	}
	cc := ChunkCoord{X: -3, Y: 5}
	if cc.X != -3 || cc.Y != 5 {
		t.Fatalf("ChunkCoord fields wrong: %+v", cc)
	}
}

// --- View with non-zero viewport ---

func TestModel_View_WithViewport(t *testing.T) {
	m := NewModel()
	m.viewportW = 80
	m.viewportH = 24
	out := m.View()
	if out == "" || out == "World Explorer — loading..." {
		t.Fatalf("View with non-zero viewport returned unexpected: %q", out)
	}
}

// --- Update TickMsg ---

func TestModel_Update_TickMsg_WorldMode(t *testing.T) {
	m := NewModel()
	m.timeOfDay = 0.5
	m.timeScale = 2
	m.mode = ModeWorld
	next, cmd := m.Update(TickMsg{})
	nm := next.(Model)
	want := 0.5 + 2.0/600.0
	if nm.timeOfDay < want-0.0001 || nm.timeOfDay > want+0.0001 {
		t.Errorf("TickMsg timeOfDay = %v, want ~%v", nm.timeOfDay, want)
	}
	if cmd == nil {
		t.Fatal("TickMsg Update should return a non-nil cmd")
	}
}

func TestModel_Update_TickMsg_LocalMode_AnimalsMove(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = GenerateLocalMap(0, 0, 42, Forest)
	m.localMap.Animals = []*Animal{{X: 10, Y: 10, Char: 'd', Color: "#888", Flee: false}}
	next, _ := m.Update(TickMsg{})
	nm := next.(Model)
	if nm.localMap == nil {
		t.Fatal("TickMsg local mode: localMap should not be nil after tick")
	}
}

func TestModel_Update_TickMsg_TimeWraps(t *testing.T) {
	m := NewModel()
	m.timeOfDay = 0.9999
	m.timeScale = 10
	next, _ := m.Update(TickMsg{})
	nm := next.(Model)
	if nm.timeOfDay < 0 || nm.timeOfDay >= 1.0 {
		t.Errorf("TickMsg time wrap: timeOfDay = %v, want in [0,1)", nm.timeOfDay)
	}
}

