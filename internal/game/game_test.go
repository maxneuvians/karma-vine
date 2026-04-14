package game

import (
	"testing"

	tea "charm.land/bubbletea/v2"
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
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	if cmd == nil {
		t.Fatal("Update 'q': expected a non-nil quit command")
	}
}

func TestModel_Update_QuitOnCtrlC(t *testing.T) {
	m := NewModel()
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	if cmd == nil {
		t.Fatal("Update ctrl+c: expected a non-nil quit command")
	}
}

func TestModel_Update_UnknownKeyNoOp(t *testing.T) {
	m := NewModel()
	next, cmd := m.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
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
	if v := m.View(); v.Content == "" {
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
	if out.Content == "" || out.Content == "World Explorer \u2014 loading..." {
		t.Fatalf("View with non-zero viewport returned unexpected: %q", out.Content)
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

// --- Pause tick tests ---

func TestPause_TickMsg_NoTimeAdvance(t *testing.T) {
	m := NewModel()
	m.paused = true
	m.timeOfDay = 0.5
	m.timeScale = 2
	next, cmd := m.Update(TickMsg{})
	nm := next.(Model)
	if nm.timeOfDay != 0.5 {
		t.Errorf("paused TickMsg: timeOfDay = %v, want 0.5 (unchanged)", nm.timeOfDay)
	}
	if cmd == nil {
		t.Fatal("paused TickMsg should return a non-nil reschedule command")
	}
}

func TestPause_TickMsg_NoAnimalMovement(t *testing.T) {
	m := NewModel()
	m.paused = true
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.localMap.Animals = []*Animal{{X: 10, Y: 10, Char: 'd', Color: "#888", Flee: false}}
	startX, startY := 10, 10
	for i := 0; i < 10; i++ {
		next, _ := m.Update(TickMsg{})
		m = next.(Model)
	}
	a := m.localMap.Animals[0]
	if a.X != startX || a.Y != startY {
		t.Fatalf("paused animals moved: expected (%d,%d), got (%d,%d)", startX, startY, a.X, a.Y)
	}
}

// --- Equipment tests ---

func TestBodySlot_Constants(t *testing.T) {
	if int(SlotFeet) != NumBodySlots-1 {
		t.Fatalf("SlotFeet (%d) should equal NumBodySlots-1 (%d)", SlotFeet, NumBodySlots-1)
	}
}

func TestNewModel_DefaultOutfit(t *testing.T) {
	m := NewModel()
	// Chest should have Cloth Tunic.
	if m.inventory.Equipped[SlotChest].Name != "Cloth Tunic" {
		t.Fatalf("expected Cloth Tunic in Chest, got %q", m.inventory.Equipped[SlotChest].Name)
	}
	// Legs should have Cloth Pants.
	if m.inventory.Equipped[SlotLegs].Name != "Cloth Pants" {
		t.Fatalf("expected Cloth Pants in Legs, got %q", m.inventory.Equipped[SlotLegs].Name)
	}
	// Feet should have Leather Boots.
	if m.inventory.Equipped[SlotFeet].Name != "Leather Boots" {
		t.Fatalf("expected Leather Boots in Feet, got %q", m.inventory.Equipped[SlotFeet].Name)
	}
	// These should NOT be in inventory.Items.
	for _, item := range m.inventory.Items {
		if item.Name == "Cloth Tunic" || item.Name == "Cloth Pants" || item.Name == "Leather Boots" {
			t.Fatalf("default outfit item %q should not be in inventory.Items", item.Name)
		}
	}
}

// 7.1 New model has empty inventory.
func TestNewModel_EmptyInventory(t *testing.T) {
	m := NewModel()
	if m.inventory.Items == nil {
		t.Fatal("NewModel: inventory.Items should be non-nil empty slice")
	}
	if len(m.inventory.Items) != 0 {
		t.Fatalf("NewModel: expected 0 items, got %d", len(m.inventory.Items))
	}
	if m.screenMode == ScreenInventory {
		t.Fatal("NewModel: showInventory should be false")
	}
	if m.inventoryCursor != 0 {
		t.Fatalf("NewModel: inventoryCursor should be 0, got %d", m.inventoryCursor)
	}
}

func TestNewModel_PlayerHP(t *testing.T) {
	m := NewModel()
	if m.playerHP != 20 {
		t.Fatalf("NewModel: expected playerHP == 20, got %d", m.playerHP)
	}
	if m.playerMaxHP != 20 {
		t.Fatalf("NewModel: expected playerMaxHP == 20, got %d", m.playerMaxHP)
	}
	if m.showHelpPanel {
		t.Fatal("NewModel: expected showHelpPanel == false")
	}
}

