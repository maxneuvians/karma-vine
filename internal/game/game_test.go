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
