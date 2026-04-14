package game

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// 6.1 applyDelta in ModeWorld increments worldPos.
func TestApplyDelta_WorldMovement(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.worldPos = WorldCoord{X: 5, Y: 3}

	m = applyDelta(1, 0, m)
	if m.worldPos.X != 6 || m.worldPos.Y != 3 {
		t.Fatalf("expected worldPos {6,3}, got {%d,%d}", m.worldPos.X, m.worldPos.Y)
	}

	m = applyDelta(0, -1, m)
	if m.worldPos.X != 6 || m.worldPos.Y != 2 {
		t.Fatalf("expected worldPos {6,2}, got {%d,%d}", m.worldPos.X, m.worldPos.Y)
	}
}

// 6.2 applyDelta in ModeLocal is blocked at {0,0} when pressing up (dy=-1).
func TestApplyDelta_BlockedAtTopLeftBoundary(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 0, Y: 0}

	m = applyDelta(0, -1, m) // move up — out of bounds
	if m.playerPos.X != 0 || m.playerPos.Y != 0 {
		t.Fatalf("expected playerPos to stay {0,0} at top-left boundary, got {%d,%d}", m.playerPos.X, m.playerPos.Y)
	}
}

// 6.3 applyDelta in ModeLocal is blocked by a Blocking object.
func TestApplyDelta_BlockedByObject(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}

	// Place a blocking object at {5, 4} (one step up)
	m.localMap.Objects[5][4] = &Object{Char: '#', Color: "white", Blocking: true}

	m = applyDelta(0, -1, m) // move up — blocked
	if m.playerPos.X != 5 || m.playerPos.Y != 5 {
		t.Fatalf("expected playerPos to stay {5,5} when blocked, got {%d,%d}", m.playerPos.X, m.playerPos.Y)
	}

	// Moving in an unblocked direction should succeed
	m = applyDelta(1, 0, m) // move right — clear
	if m.playerPos.X != 6 || m.playerPos.Y != 5 {
		t.Fatalf("expected playerPos {6,5} after unblocked move, got {%d,%d}", m.playerPos.X, m.playerPos.Y)
	}
}

// 6.4 applyDelta in ModeLocal returns early when localMap is nil.
func TestApplyDelta_NilLocalMap(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = nil
	m.playerPos = LocalCoord{X: 5, Y: 5}
	result := applyDelta(1, 0, m)
	if result.playerPos != m.playerPos {
		t.Fatalf("expected playerPos unchanged with nil localMap, got %+v", result.playerPos)
	}
}

// 6.5 applyDelta in ModeLocal is clamped at right and bottom boundaries.
func TestApplyDelta_RightBottomBoundary(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: LocalMapW - 1, Y: LocalMapH - 1}

	m = applyDelta(1, 0, m) // move right — out of bounds
	if m.playerPos.X != LocalMapW-1 {
		t.Fatalf("expected playerPos.X clamped at %d, got %d", LocalMapW-1, m.playerPos.X)
	}
	m = applyDelta(0, 1, m) // move down — out of bounds
	if m.playerPos.Y != LocalMapH-1 {
		t.Fatalf("expected playerPos.Y clamped at %d, got %d", LocalMapH-1, m.playerPos.Y)
	}
}

// --- nextWorldZoom / prevWorldZoom ---

func TestNextWorldZoom_AllCases(t *testing.T) {
	cases := []struct{ in, want int }{
		{1, 2},
		{2, 4},
		{4, 8},
		{8, 8},   // default → 8
		{99, 8},  // unknown → 8
	}
	for _, c := range cases {
		if got := nextWorldZoom(c.in); got != c.want {
			t.Errorf("nextWorldZoom(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestPrevWorldZoom_AllCases(t *testing.T) {
	cases := []struct{ in, want int }{
		{8, 4},
		{4, 2},
		{2, 1},
		{1, 1},   // default → 1
		{99, 1},  // unknown → 1
	}
	for _, c := range cases {
		if got := prevWorldZoom(c.in); got != c.want {
			t.Errorf("prevWorldZoom(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestNextTimeScale_DefaultCase(t *testing.T) {
	// Any value other than 1, 2, 5 should return 10
	if got := nextTimeScale(99); got != 10 {
		t.Errorf("nextTimeScale(99) = %d, want 10", got)
	}
}

func TestPrevTimeScale_DefaultCase(t *testing.T) {
	// Any value other than 2, 5, 10 should return 1
	if got := prevTimeScale(99); got != 1 {
		t.Errorf("prevTimeScale(99) = %d, want 1", got)
	}
}

// --- findSpawnPoint ---

func TestFindSpawnPoint_UnblockedCentre(t *testing.T) {
	lm := &LocalMap{} // no objects → centre is unblocked
	pos := findSpawnPoint(lm)
	cx, cy := LocalMapW/2, LocalMapH/2
	if pos.X != cx || pos.Y != cy {
		t.Errorf("findSpawnPoint with empty map: got %+v, want {%d,%d}", pos, cx, cy)
	}
}

func TestFindSpawnPoint_BlockedCentre(t *testing.T) {
	lm := &LocalMap{}
	cx, cy := LocalMapW/2, LocalMapH/2
	lm.Objects[cx][cy] = &Object{Char: '#', Color: "white", Blocking: true}
	pos := findSpawnPoint(lm)
	// Should avoid the blocked centre
	if pos.X == cx && pos.Y == cy {
		t.Errorf("findSpawnPoint should skip blocked centre, got {%d,%d}", pos.X, pos.Y)
	}
}

// --- handleKey ---

func TestHandleKey_ToggleSidebar(t *testing.T) {
	m := NewModel()
	m.showSidebar = false
	result, cmd := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}, m)
	if !result.showSidebar {
		t.Fatal("? key should set showSidebar to true")
	}
	if cmd != nil {
		t.Fatal("? key should return nil cmd")
	}
}

func TestHandleKey_WorldZoom_Plus(t *testing.T) {
	// + zooms in → calls prevWorldZoom
	m := NewModel()
	m.mode = ModeWorld
	m.worldZoom = 4
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("+")}, m)
	if result.worldZoom != 2 {
		t.Errorf("+ in ModeWorld: worldZoom = %d, want 2", result.worldZoom)
	}
}

func TestHandleKey_WorldZoom_Equals(t *testing.T) {
	// = zooms in → calls prevWorldZoom
	m := NewModel()
	m.mode = ModeWorld
	m.worldZoom = 4
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("=")}, m)
	if result.worldZoom != 2 {
		t.Errorf("= in ModeWorld: worldZoom = %d, want 2", result.worldZoom)
	}
}

func TestHandleKey_WorldZoom_Minus(t *testing.T) {
	// - zooms out → calls nextWorldZoom
	m := NewModel()
	m.mode = ModeWorld
	m.worldZoom = 2
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("-")}, m)
	if result.worldZoom != 4 {
		t.Errorf("- in ModeWorld: worldZoom = %d, want 4", result.worldZoom)
	}
}

func TestHandleKey_WorldZoom_NotInLocalMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.worldZoom = 1
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("+")}, m)
	if result.worldZoom != 1 {
		t.Errorf("+ in ModeLocal should not change worldZoom: got %d, want 1", result.worldZoom)
	}
}

func TestHandleKey_WorldZoom_MinusNotInLocalMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.worldZoom = 1
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("-")}, m)
	if result.worldZoom != 1 {
		t.Errorf("- in ModeLocal should not change worldZoom: got %d, want 1", result.worldZoom)
	}
}

func TestHandleKey_Movement_ArrowKeys(t *testing.T) {
	keys := []struct {
		msg        tea.KeyMsg
		ddx, ddy   int
	}{
		{tea.KeyMsg{Type: tea.KeyUp}, 0, -1},
		{tea.KeyMsg{Type: tea.KeyDown}, 0, 1},
		{tea.KeyMsg{Type: tea.KeyLeft}, -1, 0},
		{tea.KeyMsg{Type: tea.KeyRight}, 1, 0},
	}
	for _, k := range keys {
		m := NewModel()
		m.mode = ModeWorld
		m.worldPos = WorldCoord{X: 10, Y: 10}
		result, _ := handleKey(k.msg, m)
		wantX := 10 + k.ddx
		wantY := 10 + k.ddy
		if result.worldPos.X != wantX || result.worldPos.Y != wantY {
			t.Errorf("arrow key: got worldPos {%d,%d}, want {%d,%d}",
				result.worldPos.X, result.worldPos.Y, wantX, wantY)
		}
	}
}

func TestHandleKey_Movement_WASD(t *testing.T) {
	keys := []struct {
		ch       string
		ddx, ddy int
	}{
		{"w", 0, -1},
		{"s", 0, 1},
		{"a", -1, 0},
		{"d", 1, 0},
	}
	for _, k := range keys {
		m := NewModel()
		m.mode = ModeWorld
		m.worldPos = WorldCoord{X: 10, Y: 10}
		result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k.ch)}, m)
		wantX := 10 + k.ddx
		wantY := 10 + k.ddy
		if result.worldPos.X != wantX || result.worldPos.Y != wantY {
			t.Errorf("key %q: got worldPos {%d,%d}, want {%d,%d}",
				k.ch, result.worldPos.X, result.worldPos.Y, wantX, wantY)
		}
	}
}

func TestHandleKey_DescendToLocal_Enter(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	result, cmd := handleKey(tea.KeyMsg{Type: tea.KeyEnter}, m)
	if result.mode != ModeLocal {
		t.Fatalf("enter in ModeWorld: mode = %d, want ModeLocal", result.mode)
	}
	if result.localMap == nil {
		t.Fatal("enter in ModeWorld: localMap should not be nil")
	}
	if cmd != nil {
		t.Fatal("enter key should return nil cmd")
	}
}

func TestHandleKey_DescendToLocal_GT(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(">")}, m)
	if result.mode != ModeLocal {
		t.Fatalf("> in ModeWorld: mode = %d, want ModeLocal", result.mode)
	}
}

func TestHandleKey_DescendNoopInLocalMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyEnter}, m)
	if result.mode != ModeLocal {
		t.Fatalf("enter in ModeLocal should stay ModeLocal, got %d", result.mode)
	}
}

func TestHandleKey_AscendToWorld_Esc(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyEsc}, m)
	if result.mode != ModeWorld {
		t.Fatalf("esc in ModeLocal: mode = %d, want ModeWorld", result.mode)
	}
}

func TestHandleKey_AscendToWorld_LT(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("<")}, m)
	if result.mode != ModeWorld {
		t.Fatalf("< in ModeLocal: mode = %d, want ModeWorld", result.mode)
	}
}

func TestHandleKey_AscendNoopInWorldMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyEsc}, m)
	if result.mode != ModeWorld {
		t.Fatalf("esc in ModeWorld should stay ModeWorld, got %d", result.mode)
	}
}

// ── map picker input tests ────────────────────────────────────────────────────

func TestHandleKey_M_OpensPickerInWorldMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.mapMode = MapModeElevation
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")}, m)
	if !result.showMapPicker {
		t.Fatal("m in ModeWorld: showMapPicker should be true")
	}
	if result.mapPickerCursor != int(MapModeElevation) {
		t.Errorf("m in ModeWorld: mapPickerCursor = %d, want %d", result.mapPickerCursor, int(MapModeElevation))
	}
	if result.showSidebar {
		t.Fatal("m in ModeWorld: showSidebar should be false")
	}
}

func TestHandleKey_M_ClosesPickerWithoutChangingMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = true
	m.mapMode = MapModeTemperature
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")}, m)
	if result.showMapPicker {
		t.Fatal("m with picker open: showMapPicker should be false")
	}
	if result.mapMode != MapModeTemperature {
		t.Errorf("m closes picker: mapMode changed, got %d want %d", result.mapMode, MapModeTemperature)
	}
}

func TestHandleKey_M_NoopInLocalMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")}, m)
	if result.showMapPicker {
		t.Fatal("m in ModeLocal should not open picker")
	}
}

func TestHandleKey_Picker_EnterAppliesSelection(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = true
	m.mapPickerCursor = int(MapModeElevation)
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyEnter}, m)
	if result.mapMode != MapModeElevation {
		t.Errorf("enter in picker: mapMode = %d, want MapModeElevation", result.mapMode)
	}
	if result.showMapPicker {
		t.Fatal("enter in picker: showMapPicker should be false after confirm")
	}
}

func TestHandleKey_Picker_EscCancels(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = true
	m.mapMode = MapModeTemperature
	m.mapPickerCursor = int(MapModeElevation)
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyEsc}, m)
	if result.showMapPicker {
		t.Fatal("esc in picker: showMapPicker should be false")
	}
	if result.mapMode != MapModeTemperature {
		t.Errorf("esc in picker: mapMode changed, got %d want %d", result.mapMode, MapModeTemperature)
	}
}

func TestHandleKey_Picker_UpDownMoveCursor(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = true
	m.mapPickerCursor = 1

	down, _ := handleKey(tea.KeyMsg{Type: tea.KeyDown}, m)
	if down.mapPickerCursor != 2 {
		t.Errorf("down in picker: cursor = %d, want 2", down.mapPickerCursor)
	}
	if down.worldPos != m.worldPos {
		t.Error("down in picker: worldPos should not change")
	}

	up, _ := handleKey(tea.KeyMsg{Type: tea.KeyUp}, m)
	if up.mapPickerCursor != 0 {
		t.Errorf("up in picker: cursor = %d, want 0", up.mapPickerCursor)
	}
}

func TestHandleKey_Picker_CursorClampsAtBottom(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = true
	m.mapPickerCursor = len(mapModeNames) - 1
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyDown}, m)
	if result.mapPickerCursor != len(mapModeNames)-1 {
		t.Errorf("down at bottom: cursor = %d, want %d", result.mapPickerCursor, len(mapModeNames)-1)
	}
}

func TestHandleKey_Picker_CursorClampsAtTop(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = true
	m.mapPickerCursor = 0
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyUp}, m)
	if result.mapPickerCursor != 0 {
		t.Errorf("up at top: cursor = %d, want 0", result.mapPickerCursor)
	}
}

func TestHandleKey_Picker_MovementBlockedWhileOpen(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = true
	m.worldPos = WorldCoord{X: 5, Y: 5}
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRight}, m)
	if result.worldPos.X != 5 {
		t.Errorf("right while picker open: worldPos.X = %d, want 5 (blocked)", result.worldPos.X)
	}
}

func TestHandleKey_QuestionMark_ClosesPicker(t *testing.T) {
	m := NewModel()
	m.showMapPicker = true
	m.showSidebar = false
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}, m)
	if result.showMapPicker {
		t.Fatal("?: showMapPicker should be false when sidebar opened")
	}
	if !result.showSidebar {
		t.Fatal("?: showSidebar should be true")
	}
}

func TestHandleKey_Enter_DescendsWhenPickerClosed(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = false
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyEnter}, m)
	if result.mode != ModeLocal {
		t.Fatalf("enter with picker closed: mode = %d, want ModeLocal", result.mode)
	}
}

// ── dungeon input tests ───────────────────────────────────────────────────────

// dungeonReadyModel returns a Model in ModeLocal standing on a dungeon entrance.
func dungeonReadyModel() Model {
	m := NewModel()
	m.globalSeed = 42
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 10, Y: 10}
	// Place dungeon entrance under the player.
	m.localMap.Objects[10][10] = &Object{Char: '>', Color: "#e8c96a", Blocking: false}
	return m
}

func TestHandleKey_DescendFromLocalToDungeon(t *testing.T) {
	m := dungeonReadyModel()
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyEnter}, m)
	if result.mode != ModeDungeon {
		t.Fatalf("enter on dungeon entrance: mode = %d, want ModeDungeon", result.mode)
	}
	if result.dungeonDepth != 1 {
		t.Fatalf("dungeonDepth = %d, want 1", result.dungeonDepth)
	}
	if result.currentDungeon == nil {
		t.Fatal("currentDungeon should not be nil")
	}
}

func TestHandleKey_AscendFromDepth1ReturnsToLocal(t *testing.T) {
	m := dungeonReadyModel()
	m, _ = handleKey(tea.KeyMsg{Type: tea.KeyEnter}, m)
	if m.mode != ModeDungeon {
		t.Fatalf("precondition: mode = %d, want ModeDungeon", m.mode)
	}
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyEsc}, m)
	if result.mode != ModeLocal {
		t.Fatalf("esc at depth 1: mode = %d, want ModeLocal", result.mode)
	}
	if result.playerPos.X != 10 || result.playerPos.Y != 10 {
		t.Fatalf("playerPos = {%d,%d}, want {10,10}", result.playerPos.X, result.playerPos.Y)
	}
}

func TestHandleKey_AscendFromDeepLevel(t *testing.T) {
	m := dungeonReadyModel()
	m, _ = handleKey(tea.KeyMsg{Type: tea.KeyEnter}, m)
	if !m.currentDungeon.HasDownStair {
		t.Skip("generated level has no down-stair")
	}
	m.playerPos = m.currentDungeon.DownStair
	m, _ = handleKey(tea.KeyMsg{Type: tea.KeyEnter}, m)
	if m.dungeonDepth != 2 {
		t.Fatalf("dungeonDepth = %d, want 2", m.dungeonDepth)
	}
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyEsc}, m)
	if result.dungeonDepth != 1 {
		t.Fatalf("dungeonDepth after ascend = %d, want 1", result.dungeonDepth)
	}
	if result.mode != ModeDungeon {
		t.Fatalf("mode after ascend from depth 2 = %d, want ModeDungeon", result.mode)
	}
}

func TestHandleKey_EscInDungeonNeverSetsWorldMode(t *testing.T) {
	m := dungeonReadyModel()
	m, _ = handleKey(tea.KeyMsg{Type: tea.KeyEnter}, m)
	result, _ := handleKey(tea.KeyMsg{Type: tea.KeyEsc}, m)
	if result.mode == ModeWorld {
		t.Fatal("esc in ModeDungeon should never set ModeWorld")
	}
}

func TestApplyDelta_DungeonBlockedByWall(t *testing.T) {
	m := NewModel()
	m.mode = ModeDungeon
	level := &DungeonLevel{}
	level.Cells[5][5].Kind = CellFloor
	m.currentDungeon = level
	m.playerPos = LocalCoord{X: 5, Y: 5}

	result := applyDelta(1, 0, m) // move right into wall
	if result.playerPos.X != 5 {
		t.Fatalf("movement into wall: playerPos.X = %d, want 5", result.playerPos.X)
	}
}

func TestApplyDelta_DungeonBlockedByTorch(t *testing.T) {
	m := NewModel()
	m.mode = ModeDungeon
	level := &DungeonLevel{}
	level.Cells[5][5].Kind = CellFloor
	level.Cells[6][5].Kind = CellWall
	level.Cells[6][5].Object = &Object{Char: '†', Color: "#e8c96a", Blocking: true}
	m.currentDungeon = level
	m.playerPos = LocalCoord{X: 5, Y: 5}

	result := applyDelta(1, 0, m)
	if result.playerPos.X != 5 {
		t.Fatalf("movement into torch: playerPos.X = %d, want 5", result.playerPos.X)
	}
}

func TestApplyDelta_DungeonMovesOverBrazier(t *testing.T) {
	m := NewModel()
	m.mode = ModeDungeon
	level := &DungeonLevel{}
	level.Cells[5][5].Kind = CellFloor
	level.Cells[6][5].Kind = CellFloor
	level.Cells[6][5].Object = &Object{Char: 'Ω', Color: "#e07030", Blocking: false}
	m.currentDungeon = level
	m.playerPos = LocalCoord{X: 5, Y: 5}

	result := applyDelta(1, 0, m)
	if result.playerPos.X != 6 {
		t.Fatalf("movement over brazier: playerPos.X = %d, want 6", result.playerPos.X)
	}
}
