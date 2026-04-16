package game

import (
	"fmt"
	"testing"

	tea "charm.land/bubbletea/v2"
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
	result, cmd := handleKey(tea.KeyPressMsg{Code: '\\', Text: "\\"}, m)
	if !result.showSidebar {
		t.Fatal("\\ key should set showSidebar to true")
	}
	if cmd != nil {
		t.Fatal("\\ key should return nil cmd")
	}
}

func TestHandleKey_WorldZoom_Plus(t *testing.T) {
	// + zooms in → calls prevWorldZoom
	m := NewModel()
	m.mode = ModeWorld
	m.worldZoom = 4
	result, _ := handleKey(tea.KeyPressMsg{Code: '+', Text: "+"}, m)
	if result.worldZoom != 2 {
		t.Errorf("+ in ModeWorld: worldZoom = %d, want 2", result.worldZoom)
	}
}

func TestHandleKey_WorldZoom_Equals(t *testing.T) {
	// = zooms in → calls prevWorldZoom
	m := NewModel()
	m.mode = ModeWorld
	m.worldZoom = 4
	result, _ := handleKey(tea.KeyPressMsg{Code: '=', Text: "="}, m)
	if result.worldZoom != 2 {
		t.Errorf("= in ModeWorld: worldZoom = %d, want 2", result.worldZoom)
	}
}

func TestHandleKey_WorldZoom_Minus(t *testing.T) {
	// - zooms out → calls nextWorldZoom
	m := NewModel()
	m.mode = ModeWorld
	m.worldZoom = 2
	result, _ := handleKey(tea.KeyPressMsg{Code: '-', Text: "-"}, m)
	if result.worldZoom != 4 {
		t.Errorf("- in ModeWorld: worldZoom = %d, want 4", result.worldZoom)
	}
}

func TestHandleKey_WorldZoom_NotInLocalMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.worldZoom = 1
	result, _ := handleKey(tea.KeyPressMsg{Code: '+', Text: "+"}, m)
	if result.worldZoom != 1 {
		t.Errorf("+ in ModeLocal should not change worldZoom: got %d, want 1", result.worldZoom)
	}
}

func TestHandleKey_WorldZoom_MinusNotInLocalMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.worldZoom = 1
	result, _ := handleKey(tea.KeyPressMsg{Code: '-', Text: "-"}, m)
	if result.worldZoom != 1 {
		t.Errorf("- in ModeLocal should not change worldZoom: got %d, want 1", result.worldZoom)
	}
}

func TestHandleKey_Movement_ArrowKeys(t *testing.T) {
	keys := []struct {
		msg        tea.KeyPressMsg
		ddx, ddy   int
	}{
		{tea.KeyPressMsg{Code: tea.KeyUp}, 0, -1},
		{tea.KeyPressMsg{Code: tea.KeyDown}, 0, 1},
		{tea.KeyPressMsg{Code: tea.KeyLeft}, -1, 0},
		{tea.KeyPressMsg{Code: tea.KeyRight}, 1, 0},
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
		result, _ := handleKey(tea.KeyPressMsg{Code: rune(k.ch[0]), Text: k.ch}, m)
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
	result, cmd := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
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
	result, _ := handleKey(tea.KeyPressMsg{Code: '>', Text: ">"}, m)
	if result.mode != ModeLocal {
		t.Fatalf("> in ModeWorld: mode = %d, want ModeLocal", result.mode)
	}
}

func TestHandleKey_DescendNoopInLocalMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	if result.mode != ModeLocal {
		t.Fatalf("enter in ModeLocal should stay ModeLocal, got %d", result.mode)
	}
}

func TestHandleKey_AscendToWorld_Esc(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEsc}, m)
	if result.mode != ModeWorld {
		t.Fatalf("esc in ModeLocal: mode = %d, want ModeWorld", result.mode)
	}
}

func TestHandleKey_AscendToWorld_LT(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	result, _ := handleKey(tea.KeyPressMsg{Code: '<', Text: "<"}, m)
	if result.mode != ModeWorld {
		t.Fatalf("< in ModeLocal: mode = %d, want ModeWorld", result.mode)
	}
}

func TestHandleKey_AscendNoopInWorldMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEsc}, m)
	if result.mode != ModeWorld {
		t.Fatalf("esc in ModeWorld should stay ModeWorld, got %d", result.mode)
	}
}

// ── map picker input tests ────────────────────────────────────────────────────

func TestHandleKey_M_OpensPickerInWorldMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.mapMode = MapModeElevation
	result, _ := handleKey(tea.KeyPressMsg{Code: 'm', Text: "m"}, m)
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
	result, _ := handleKey(tea.KeyPressMsg{Code: 'm', Text: "m"}, m)
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
	result, _ := handleKey(tea.KeyPressMsg{Code: 'm', Text: "m"}, m)
	if result.showMapPicker {
		t.Fatal("m in ModeLocal should not open picker")
	}
}

func TestHandleKey_Picker_EnterAppliesSelection(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = true
	m.mapPickerCursor = int(MapModeElevation)
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
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
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEsc}, m)
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

	down, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyDown}, m)
	if down.mapPickerCursor != 2 {
		t.Errorf("down in picker: cursor = %d, want 2", down.mapPickerCursor)
	}
	if down.worldPos != m.worldPos {
		t.Error("down in picker: worldPos should not change")
	}

	up, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyUp}, m)
	if up.mapPickerCursor != 0 {
		t.Errorf("up in picker: cursor = %d, want 0", up.mapPickerCursor)
	}
}

func TestHandleKey_Picker_CursorClampsAtBottom(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = true
	m.mapPickerCursor = len(mapModeNames) - 1
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyDown}, m)
	if result.mapPickerCursor != len(mapModeNames)-1 {
		t.Errorf("down at bottom: cursor = %d, want %d", result.mapPickerCursor, len(mapModeNames)-1)
	}
}

func TestHandleKey_Picker_CursorClampsAtTop(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = true
	m.mapPickerCursor = 0
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyUp}, m)
	if result.mapPickerCursor != 0 {
		t.Errorf("up at top: cursor = %d, want 0", result.mapPickerCursor)
	}
}

func TestHandleKey_Picker_MovementBlockedWhileOpen(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = true
	m.worldPos = WorldCoord{X: 5, Y: 5}
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyRight}, m)
	if result.worldPos.X != 5 {
		t.Errorf("right while picker open: worldPos.X = %d, want 5 (blocked)", result.worldPos.X)
	}
}

func TestHandleKey_QuestionMark_ClosesPicker(t *testing.T) {
	m := NewModel()
	m.showMapPicker = true
	m.showHelpPanel = false
	result, _ := handleKey(tea.KeyPressMsg{Code: '?', Text: "?"}, m)
	if result.showMapPicker {
		t.Fatal("?: showMapPicker should be false when help opened")
	}
	if !result.showHelpPanel {
		t.Fatal("?: showHelpPanel should be true")
	}
}

func TestHandleKey_Enter_DescendsWhenPickerClosed(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.showMapPicker = false
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
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
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
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
	m, _ = handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	if m.mode != ModeDungeon {
		t.Fatalf("precondition: mode = %d, want ModeDungeon", m.mode)
	}
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEsc}, m)
	if result.mode != ModeLocal {
		t.Fatalf("esc at depth 1: mode = %d, want ModeLocal", result.mode)
	}
	if result.playerPos.X != 10 || result.playerPos.Y != 10 {
		t.Fatalf("playerPos = {%d,%d}, want {10,10}", result.playerPos.X, result.playerPos.Y)
	}
}

func TestHandleKey_AscendFromDeepLevel(t *testing.T) {
	m := dungeonReadyModel()
	m, _ = handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	if !m.currentDungeon.HasDownStair {
		t.Skip("generated level has no down-stair")
	}
	m.playerPos = m.currentDungeon.DownStair
	m, _ = handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	if m.dungeonDepth != 2 {
		t.Fatalf("dungeonDepth = %d, want 2", m.dungeonDepth)
	}
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEsc}, m)
	if result.dungeonDepth != 1 {
		t.Fatalf("dungeonDepth after ascend = %d, want 1", result.dungeonDepth)
	}
	if result.mode != ModeDungeon {
		t.Fatalf("mode after ascend from depth 2 = %d, want ModeDungeon", result.mode)
	}
}

func TestHandleKey_EscInDungeonNeverSetsWorldMode(t *testing.T) {
	m := dungeonReadyModel()
	m, _ = handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEsc}, m)
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

// ── Inventory input tests ────────────────────────────────────────────────────

// 7.2 Pickup removes object from local map and adds item to inventory.
func TestPickup_LocalMap(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	m.localMap.Objects[5][5] = &Object{Char: '⚒', Color: "#a0a0a0", Name: "Axe", Pickupable: true}

	result, _ := handleKey(tea.KeyPressMsg{Code: 'g', Text: "g"}, m)
	if result.localMap.Objects[5][5] != nil {
		t.Fatal("pickup should remove object from local map")
	}
	if len(result.inventory.Items) != 1 {
		t.Fatalf("expected 1 item in inventory, got %d", len(result.inventory.Items))
	}
	if result.inventory.Items[0].Name != "Axe" || result.inventory.Items[0].Count != 1 {
		t.Fatalf("expected Axe x1, got %+v", result.inventory.Items[0])
	}
}

// 7.3 Picking up same-named item stacks (increments Count).
func TestPickup_Stacking(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	m.inventory.Items = []Item{{Char: '†', Color: "#e8c96a", Name: "Torch", Count: 1}}
	m.localMap.Objects[5][5] = &Object{Char: '†', Color: "#e8c96a", Name: "Torch", Pickupable: true}

	result, _ := handleKey(tea.KeyPressMsg{Code: 'g', Text: "g"}, m)
	if len(result.inventory.Items) != 1 {
		t.Fatalf("expected 1 slot (stacked), got %d", len(result.inventory.Items))
	}
	if result.inventory.Items[0].Count != 2 {
		t.Fatalf("expected Torch count 2, got %d", result.inventory.Items[0].Count)
	}
}

// 7.4 Pickup ignored for non-pickupable objects.
func TestPickup_NonPickupable(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	m.localMap.Objects[5][5] = &Object{Char: '♣', Color: "#2d7a1f", Name: "Tree", Blocking: true, Pickupable: false}

	result, _ := handleKey(tea.KeyPressMsg{Code: 'g', Text: "g"}, m)
	if len(result.inventory.Items) != 0 {
		t.Fatal("pickup should not pick up non-pickupable objects")
	}
	if result.localMap.Objects[5][5] == nil {
		t.Fatal("non-pickupable object should remain on map")
	}
}

// 7.5 Pickup rejected when inventory is full.
func TestPickup_FullInventory(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	// Fill inventory with 8 different items.
	for i := 0; i < InventoryMaxSlots; i++ {
		m.inventory.Items = append(m.inventory.Items, Item{
			Char:  rune('a' + i),
			Color: "#ffffff",
			Name:  fmt.Sprintf("Item%d", i),
			Count: 1,
		})
	}
	m.localMap.Objects[5][5] = &Object{Char: '⚒', Color: "#a0a0a0", Name: "Axe", Pickupable: true}

	result, _ := handleKey(tea.KeyPressMsg{Code: 'g', Text: "g"}, m)
	if len(result.inventory.Items) != InventoryMaxSlots {
		t.Fatalf("inventory should stay at %d, got %d", InventoryMaxSlots, len(result.inventory.Items))
	}
	if result.localMap.Objects[5][5] == nil {
		t.Fatal("object should remain on map when inventory full")
	}
}

// 7.6 Drop places pickupable object on player's cell in local map.
func TestDrop_LocalMap(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	m.inventory.Items = []Item{{Char: '⚒', Color: "#a0a0a0", Name: "Axe", Count: 1}}
	m.screenMode = ScreenInventory
	m.inventoryCursor = 0

	result, _ := handleKey(tea.KeyPressMsg{Code: 'd', Text: "d"}, m)
	obj := result.localMap.Objects[5][5]
	if obj == nil {
		t.Fatal("drop should place object on player's cell")
	}
	if obj.Name != "Axe" || !obj.Pickupable {
		t.Fatalf("dropped object: Name=%q Pickupable=%v, want Axe true", obj.Name, obj.Pickupable)
	}
	if len(result.inventory.Items) != 0 {
		t.Fatalf("expected 0 items after drop, got %d", len(result.inventory.Items))
	}
}

// 7.7 Drop decrements item count; slot removed when count reaches zero.
func TestDrop_DecrementCount(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	m.inventory.Items = []Item{{Char: '†', Color: "#e8c96a", Name: "Torch", Count: 2}}
	m.screenMode = ScreenInventory
	m.inventoryCursor = 0

	result, _ := handleKey(tea.KeyPressMsg{Code: 'd', Text: "d"}, m)
	if len(result.inventory.Items) != 1 {
		t.Fatalf("expected 1 slot remaining, got %d", len(result.inventory.Items))
	}
	if result.inventory.Items[0].Count != 1 {
		t.Fatalf("expected Torch count 1, got %d", result.inventory.Items[0].Count)
	}
}

// 7.8 Drop ignored in ModeWorld.
func TestDrop_IgnoredInWorldMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.inventory.Items = []Item{{Char: '⚒', Color: "#a0a0a0", Name: "Axe", Count: 1}}
	m.screenMode = ScreenInventory
	m.inventoryCursor = 0

	result, _ := handleKey(tea.KeyPressMsg{Code: 'd', Text: "d"}, m)
	if len(result.inventory.Items) != 1 {
		t.Fatal("drop in ModeWorld should not remove items")
	}
}

// 7.9 `i` toggles showInventory in all modes.
func TestToggleInventory(t *testing.T) {
	for _, mode := range []Mode{ModeWorld, ModeLocal, ModeDungeon} {
		m := NewModel()
		m.mode = mode
		if mode == ModeLocal {
			m.localMap = &LocalMap{}
		} else if mode == ModeDungeon {
			level := &DungeonLevel{}
			level.Cells[5][5].Kind = CellFloor
			m.currentDungeon = level
			m.playerPos = LocalCoord{X: 5, Y: 5}
		}

		// Toggle on.
		result, _ := handleKey(tea.KeyPressMsg{Code: 'i', Text: "i"}, m)
		if result.screenMode == ScreenNormal {
			t.Fatalf("i should open inventory in mode %d", mode)
		}
		// Toggle off.
		result, _ = handleKey(tea.KeyPressMsg{Code: 'i', Text: "i"}, result)
		if result.screenMode == ScreenInventory {
			t.Fatalf("second i should close inventory in mode %d", mode)
		}
	}
}

// 7.10 Cursor keys move inventoryCursor when showInventory is true.
func TestInventoryCursor(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.screenMode = ScreenInventory
	m.inventory.Items = []Item{
		{Name: "Axe", Count: 1},
		{Name: "Torch", Count: 1},
		{Name: "Brazier", Count: 1},
	}
	m.inventoryCursor = 0

	// Move down.
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyDown}, m)
	if result.inventoryCursor != 1 {
		t.Fatalf("down: cursor = %d, want 1", result.inventoryCursor)
	}
	// Player should NOT move.
	if result.playerPos != m.playerPos {
		t.Fatal("down with inventory open should not move player")
	}

	// Move down again.
	result, _ = handleKey(tea.KeyPressMsg{Code: tea.KeyDown}, result)
	if result.inventoryCursor != 2 {
		t.Fatalf("down: cursor = %d, want 2", result.inventoryCursor)
	}

	// Clamp at bottom.
	result, _ = handleKey(tea.KeyPressMsg{Code: tea.KeyDown}, result)
	if result.inventoryCursor != 2 {
		t.Fatalf("down at bottom: cursor = %d, want 2 (clamped)", result.inventoryCursor)
	}

	// Move up.
	result, _ = handleKey(tea.KeyPressMsg{Code: tea.KeyUp}, result)
	if result.inventoryCursor != 1 {
		t.Fatalf("up: cursor = %d, want 1", result.inventoryCursor)
	}
}

// 7.11 Axe chops adjacent tree object (tree cell becomes nil).
func TestUseAxe_ChopsTree(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	m.localMap.Objects[5][4] = &Object{Char: '♣', Color: "#2d7a1f", Name: "Tree", Blocking: true}
	m.inventory.Items = []Item{{Char: '⚒', Color: "#a0a0a0", Name: "Axe", Count: 1}}
	m.inventoryCursor = 0

	result, _ := handleKey(tea.KeyPressMsg{Code: 'u', Text: "u"}, m)
	if result.localMap.Objects[5][4] != nil {
		t.Fatal("axe use should remove adjacent tree")
	}
}

// 7.12 Axe use does not remove axe from inventory.
func TestUseAxe_NotConsumed(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	m.localMap.Objects[5][4] = &Object{Char: '♣', Color: "#2d7a1f", Name: "Tree", Blocking: true}
	m.inventory.Items = []Item{{Char: '⚒', Color: "#a0a0a0", Name: "Axe", Count: 1}}
	m.inventoryCursor = 0

	result, _ := handleKey(tea.KeyPressMsg{Code: 'u', Text: "u"}, m)
	if len(result.inventory.Items) != 1 || result.inventory.Items[0].Name != "Axe" || result.inventory.Items[0].Count != 1 {
		t.Fatalf("axe should not be consumed, got %+v", result.inventory.Items)
	}
}

// 7.13 Use with no adjacent tree is a no-op.
func TestUseAxe_NoTarget(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	m.inventory.Items = []Item{{Char: '⚒', Color: "#a0a0a0", Name: "Axe", Count: 1}}
	m.inventoryCursor = 0

	result, _ := handleKey(tea.KeyPressMsg{Code: 'u', Text: "u"}, m)
	// Nothing should change.
	if len(result.inventory.Items) != 1 {
		t.Fatal("use with no target should not change inventory")
	}
}

// ── ScreenMode tests ─────────────────────────────────────────────────────────

// --- Pause tests ---

func TestPause_SpaceTogglesPaused(t *testing.T) {
	m := NewModel()
	if m.paused {
		t.Fatal("NewModel should start unpaused")
	}
	result, _ := handleKey(tea.KeyPressMsg{Code: ' ', Text: " "}, m)
	if !result.paused {
		t.Fatal("space should set paused to true")
	}
	result, _ = handleKey(tea.KeyPressMsg{Code: ' ', Text: " "}, result)
	if result.paused {
		t.Fatal("second space should set paused to false")
	}
}

func TestPause_MovementBlocked(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	m.paused = true
	m.worldPos = WorldCoord{X: 10, Y: 10}
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyRight}, m)
	if result.worldPos.X != 10 || result.worldPos.Y != 10 {
		t.Fatalf("paused movement: worldPos = {%d,%d}, want {10,10}", result.worldPos.X, result.worldPos.Y)
	}
}

func TestPause_InventoryCursorUnaffected(t *testing.T) {
	m := NewModel()
	m.paused = true
	m.screenMode = ScreenInventory
	m.inventory.Items = []Item{
		{Name: "A", Count: 1},
		{Name: "B", Count: 1},
	}
	m.inventoryCursor = 0
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyDown}, m)
	if result.inventoryCursor != 1 {
		t.Fatalf("paused inventory cursor: got %d, want 1", result.inventoryCursor)
	}
}

// ── Equipment input tests ────────────────────────────────────────────────────

func TestEquipItem_EmptySlot(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	// Clear default outfit so Head is empty.
	m.inventory.Equipped = [NumBodySlots]Item{}
	m.inventory.Items = []Item{
		{Char: '⛑', Color: "#ff0000", Name: "Helmet", Count: 1, Slots: []BodySlot{SlotHead}},
	}
	m.inventoryCursor = 0
	result, _ := handleKey(tea.KeyPressMsg{Code: 'e', Text: "e"}, m)
	if result.inventory.Equipped[SlotHead].Name != "Helmet" {
		t.Fatalf("expected Helmet in Head slot, got %q", result.inventory.Equipped[SlotHead].Name)
	}
	if len(result.inventory.Items) != 0 {
		t.Fatalf("expected 0 items after equip, got %d", len(result.inventory.Items))
	}
}

func TestEquipItem_Swap(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	m.inventory.Equipped = [NumBodySlots]Item{}
	m.inventory.Equipped[SlotHead] = Item{Char: '⛑', Color: "#aaa", Name: "Old Hat", Count: 1, Slots: []BodySlot{SlotHead}}
	m.inventory.Items = []Item{
		{Char: '⛑', Color: "#ff0000", Name: "New Helmet", Count: 1, Slots: []BodySlot{SlotHead}},
	}
	m.inventoryCursor = 0
	result, _ := handleKey(tea.KeyPressMsg{Code: 'e', Text: "e"}, m)
	if result.inventory.Equipped[SlotHead].Name != "New Helmet" {
		t.Fatalf("expected New Helmet in Head, got %q", result.inventory.Equipped[SlotHead].Name)
	}
	if len(result.inventory.Items) != 1 || result.inventory.Items[0].Name != "Old Hat" {
		t.Fatalf("expected Old Hat in inventory, got %+v", result.inventory.Items)
	}
}

func TestEquipItem_NonEquippable(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	m.inventory.Equipped = [NumBodySlots]Item{}
	m.inventory.Items = []Item{
		{Char: '†', Color: "#e8c96a", Name: "Torch", Count: 1}, // no Slots
	}
	m.inventoryCursor = 0
	result, _ := handleKey(tea.KeyPressMsg{Code: 'e', Text: "e"}, m)
	if len(result.inventory.Items) != 1 || result.inventory.Items[0].Name != "Torch" {
		t.Fatalf("non-equippable item should remain unchanged, got %+v", result.inventory.Items)
	}
}

func TestEquipItem_FullInventorySwapRejected(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	m.inventory.Equipped = [NumBodySlots]Item{}
	m.inventory.Equipped[SlotHead] = Item{Char: '⛑', Color: "#aaa", Name: "Old Hat", Count: 1, Slots: []BodySlot{SlotHead}}
	// Fill inventory to max.
	for i := 0; i < InventoryMaxSlots; i++ {
		m.inventory.Items = append(m.inventory.Items, Item{
			Char: rune('a' + i), Color: "#fff", Name: fmt.Sprintf("Item%d", i), Count: 1,
			Slots: []BodySlot{SlotHead},
		})
	}
	m.inventoryCursor = 0
	result, _ := handleKey(tea.KeyPressMsg{Code: 'e', Text: "e"}, m)
	// Swap should be rejected because after removing item0, we have 7 slots but adding Old Hat = 8 = max,
	// which is fine. Let me reconsider: future len after removing = 7. 7 < 8 so it fits. Let me adjust.
	// Actually this test should produce a rejection. Let me set count > 1 on the item being equipped
	// so the slot stays (futureLen = 8), and then old hat can't stack and 8 >= InventoryMaxSlots.
	// Let me redo this properly.
	_ = result // rebuild below
	m2 := NewModel()
	m2.screenMode = ScreenInventory
	m2.inventory.Equipped = [NumBodySlots]Item{}
	m2.inventory.Equipped[SlotHead] = Item{Char: '⛑', Color: "#aaa", Name: "Old Hat", Count: 1, Slots: []BodySlot{SlotHead}}
	for i := 0; i < InventoryMaxSlots; i++ {
		m2.inventory.Items = append(m2.inventory.Items, Item{
			Char: rune('a' + i), Color: "#fff", Name: fmt.Sprintf("Filler%d", i), Count: 2,
			Slots: []BodySlot{SlotHead},
		})
	}
	m2.inventoryCursor = 0
	result2, _ := handleKey(tea.KeyPressMsg{Code: 'e', Text: "e"}, m2)
	// futureLen = 8 (count decremented, slot stays). Old Hat can't stack. 8 >= 8, rejected.
	if result2.inventory.Equipped[SlotHead].Name != "Old Hat" {
		t.Fatalf("swap should be rejected when inventory full, but equipped changed to %q", result2.inventory.Equipped[SlotHead].Name)
	}
}

func TestUnequipSlot_Occupied(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	m.equipFocused = true
	m.equipCursor = int(SlotChest)
	m.inventory.Equipped = [NumBodySlots]Item{}
	m.inventory.Equipped[SlotChest] = Item{Char: '♦', Color: "#a0a0a0", Name: "Cloth Tunic", Count: 1, Slots: []BodySlot{SlotChest}}
	m.inventory.Items = []Item{}
	result, _ := handleKey(tea.KeyPressMsg{Code: 'e', Text: "e"}, m)
	if result.inventory.Equipped[SlotChest].Name != "" {
		t.Fatalf("unequip should clear slot, got %q", result.inventory.Equipped[SlotChest].Name)
	}
	if len(result.inventory.Items) != 1 || result.inventory.Items[0].Name != "Cloth Tunic" {
		t.Fatalf("expected Cloth Tunic in inventory, got %+v", result.inventory.Items)
	}
}

func TestUnequipSlot_Empty(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	m.equipFocused = true
	m.equipCursor = int(SlotHead)
	m.inventory.Equipped = [NumBodySlots]Item{}
	m.inventory.Items = []Item{}
	result, _ := handleKey(tea.KeyPressMsg{Code: 'e', Text: "e"}, m)
	if len(result.inventory.Items) != 0 {
		t.Fatal("unequip empty slot should not add items")
	}
}

func TestUnequipSlot_FullInventoryRejected(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	m.equipFocused = true
	m.equipCursor = int(SlotHead)
	m.inventory.Equipped = [NumBodySlots]Item{}
	m.inventory.Equipped[SlotHead] = Item{Char: '⛑', Color: "#ff0000", Name: "Helmet", Count: 1, Slots: []BodySlot{SlotHead}}
	for i := 0; i < InventoryMaxSlots; i++ {
		m.inventory.Items = append(m.inventory.Items, Item{
			Char: rune('a' + i), Color: "#fff", Name: fmt.Sprintf("Filler%d", i), Count: 1,
		})
	}
	result, _ := handleKey(tea.KeyPressMsg{Code: 'e', Text: "e"}, m)
	if result.inventory.Equipped[SlotHead].Name != "Helmet" {
		t.Fatalf("unequip should be rejected when inventory full, got %q", result.inventory.Equipped[SlotHead].Name)
	}
}

func TestTabKey_TogglesEquipFocused(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	if m.equipFocused {
		t.Fatal("should start not focused on equip")
	}
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyTab}, m)
	if !result.equipFocused {
		t.Fatal("Tab should toggle equipFocused to true")
	}
	result, _ = handleKey(tea.KeyPressMsg{Code: tea.KeyTab}, result)
	if result.equipFocused {
		t.Fatal("second Tab should toggle equipFocused to false")
	}
	// Tab outside ScreenInventory is no-op.
	m2 := NewModel()
	m2.screenMode = ScreenNormal
	result2, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyTab}, m2)
	if result2.equipFocused {
		t.Fatal("Tab outside ScreenInventory should not toggle equipFocused")
	}
}

func TestEquipCursor_Navigation(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	m.equipFocused = true
	m.equipCursor = 0

	// Move down.
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyDown}, m)
	if result.equipCursor != 1 {
		t.Fatalf("down: equipCursor = %d, want 1", result.equipCursor)
	}

	// Clamp at top.
	m.equipCursor = 0
	result, _ = handleKey(tea.KeyPressMsg{Code: tea.KeyUp}, m)
	if result.equipCursor != 0 {
		t.Fatalf("up at 0: equipCursor = %d, want 0", result.equipCursor)
	}

	// Clamp at bottom.
	m.equipCursor = NumBodySlots - 1
	result, _ = handleKey(tea.KeyPressMsg{Code: tea.KeyDown}, m)
	if result.equipCursor != NumBodySlots-1 {
		t.Fatalf("down at max: equipCursor = %d, want %d", result.equipCursor, NumBodySlots-1)
	}
}

// ── ScreenMode tests ─────────────────────────────────────────────────────────

// 8.2 NewModel starts in ScreenNormal.
func TestNewModel_ScreenNormal(t *testing.T) {
	m := NewModel()
	if m.screenMode != ScreenNormal {
		t.Fatalf("NewModel: expected ScreenNormal, got %d", m.screenMode)
	}
}

// 8.3 i key toggles screenMode.
func TestScreenMode_IKeyToggle(t *testing.T) {
	m := NewModel()
	m.viewportW = 120
	m.viewportH = 40
	result, _ := handleKey(tea.KeyPressMsg{Code: 'i', Text: "i"}, m)
	if result.screenMode != ScreenInventory {
		t.Fatal("first i should set ScreenInventory")
	}
	result, _ = handleKey(tea.KeyPressMsg{Code: 'i', Text: "i"}, result)
	if result.screenMode != ScreenNormal {
		t.Fatal("second i should return to ScreenNormal")
	}
}

// 8.4 esc closes inventory.
func TestScreenMode_EscClosesInventory(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEsc}, m)
	if result.screenMode != ScreenNormal {
		t.Fatal("esc should set ScreenNormal when in ScreenInventory")
	}
}

// ── Mouse tests ──────────────────────────────────────────────────────────────

// 8.5 Scroll wheel moves inventory cursor.
func TestMouseWheel_InventoryCursor(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	m.inventory.Items = []Item{
		{Name: "A", Count: 1},
		{Name: "B", Count: 1},
		{Name: "C", Count: 1},
	}
	m.inventoryCursor = 0

	// Scroll down.
	result, _ := handleMouseWheel(tea.MouseWheelMsg{Button: tea.MouseWheelDown}, m)
	if result.inventoryCursor != 1 {
		t.Fatalf("scroll down: cursor = %d, want 1", result.inventoryCursor)
	}
	// Scroll up.
	result, _ = handleMouseWheel(tea.MouseWheelMsg{Button: tea.MouseWheelUp}, result)
	if result.inventoryCursor != 0 {
		t.Fatalf("scroll up: cursor = %d, want 0", result.inventoryCursor)
	}
	// Clamp at top.
	result, _ = handleMouseWheel(tea.MouseWheelMsg{Button: tea.MouseWheelUp}, result)
	if result.inventoryCursor != 0 {
		t.Fatalf("scroll up clamped: cursor = %d, want 0", result.inventoryCursor)
	}
}

// 8.6 Click on item row sets inventory cursor.
func TestMouseClick_InventoryRow(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	m.viewportW = 120
	m.viewportH = 40
	m.inventory.Items = []Item{
		{Name: "A", Count: 1},
		{Name: "B", Count: 1},
		{Name: "C", Count: 1},
	}
	m.inventoryCursor = 0

	// Click on row 3 (Y=2 is header+separator, Y=3 is second item → row index 1).
	result, _ := handleMouseClick(tea.MouseClickMsg{X: 5, Y: 3, Button: tea.MouseLeft}, m)
	if result.inventoryCursor != 1 {
		t.Fatalf("click on row 3: cursor = %d, want 1", result.inventoryCursor)
	}
}

// 8.7 Click in ScreenNormal/ModeLocal moves player.
func TestMouseClick_MovePlayer(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.screenMode = ScreenNormal
	m.viewportW = 80
	m.viewportH = 26
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 40, Y: 12}

	// Click to the right of the player (screen centre is mapW/2=40, mapH/2=12).
	// Click at screen X=41 → map X=41 → dx=1, dy=0 → step right.
	result, _ := handleMouseClick(tea.MouseClickMsg{X: 41, Y: 12, Button: tea.MouseLeft}, m)
	if result.playerPos.X != 41 {
		t.Fatalf("click-to-move: playerPos.X = %d, want 41", result.playerPos.X)
	}
}

// 8.8 Click ignored when sidebar open.
func TestMouseClick_IgnoredWhenSidebarOpen(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.screenMode = ScreenNormal
	m.viewportW = 80
	m.viewportH = 26
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 40, Y: 12}
	m.showSidebar = true

	result, _ := handleMouseClick(tea.MouseClickMsg{X: 50, Y: 12, Button: tea.MouseLeft}, m)
	if result.playerPos.X != 40 {
		t.Fatalf("click with sidebar: playerPos.X = %d, want 40 (unchanged)", result.playerPos.X)
	}
}

// ── Combat input tests ───────────────────────────────────────────────────────

func TestHandleKey_GOnAnimalInitiatesCombat(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	m.localMap.Animals = []*Animal{{X: 5, Y: 5, Char: 'w', Color: "#555", Name: "Wolf"}}

	result, _ := handleKey(tea.KeyPressMsg{Code: 'g', Text: "g"}, m)
	if result.screenMode != ScreenCombat {
		t.Fatalf("g on animal: screenMode = %d, want ScreenCombat", result.screenMode)
	}
	if result.combatState == nil {
		t.Fatal("g on animal: combatState should not be nil")
	}
	if !result.paused {
		t.Fatal("g on animal: should be paused during combat")
	}
}

func TestHandleKey_GOnEmptyCellSkipsCombat(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	// No animal at player position.

	result, _ := handleKey(tea.KeyPressMsg{Code: 'g', Text: "g"}, m)
	if result.screenMode != ScreenNormal {
		t.Fatalf("g on empty cell: screenMode = %d, want ScreenNormal", result.screenMode)
	}
}

func TestHandleKey_CombatScreenSuppressesMovement(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.screenMode = ScreenCombat
	m.combatState = &CombatState{PlayerWon: true}
	m.playerPos = LocalCoord{X: 5, Y: 5}

	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyUp}, m)
	if result.playerPos.X != 5 || result.playerPos.Y != 5 {
		t.Fatalf("movement in ScreenCombat: playerPos changed to {%d,%d}", result.playerPos.X, result.playerPos.Y)
	}
}

func TestHandleKey_EnterDismissesVictory(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.screenMode = ScreenCombat
	m.paused = true
	enemy := &Animal{X: 5, Y: 5, Char: 'w', Color: "#555", Name: "Wolf"}
	m.localMap.Animals = []*Animal{enemy}
	m.combatState = &CombatState{PlayerWon: true}
	m.combatEnemy = enemy

	result, cmd := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	if result.screenMode != ScreenNormal {
		t.Fatalf("enter on victory: screenMode = %d, want ScreenNormal", result.screenMode)
	}
	if result.paused {
		t.Fatal("enter on victory: should be unpaused")
	}
	if cmd != nil {
		t.Fatal("enter on victory: cmd should be nil")
	}
	if len(result.localMap.Animals) != 0 {
		t.Fatalf("enter on victory: animal should be removed, got %d", len(result.localMap.Animals))
	}
}

func TestHandleKey_EnterShowsDeathScreenOnDefeat(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenCombat
	m.combatState = &CombatState{
		PlayerWon: false,
		Enemy:     Combatant{Name: "Goblin"},
		Round:     0,
	}

	result, cmd := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	if result.screenMode != ScreenDeath {
		t.Fatalf("enter on defeat: expected ScreenDeath, got %d", result.screenMode)
	}
	if cmd != nil {
		t.Fatal("enter on defeat: should not quit immediately")
	}
}

// --- Help panel toggle ---

func TestHandleKey_QuestionMarkTogglesHelpPanel(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenNormal
	m.showHelpPanel = false
	result, _ := handleKey(tea.KeyPressMsg{Code: '?', Text: "?"}, m)
	if !result.showHelpPanel {
		t.Fatal("? should toggle showHelpPanel to true")
	}
	result2, _ := handleKey(tea.KeyPressMsg{Code: '?', Text: "?"}, result)
	if result2.showHelpPanel {
		t.Fatal("? should toggle showHelpPanel back to false")
	}
}

func TestHandleKey_QuestionMarkSuppressedInInventory(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenInventory
	m.showHelpPanel = false
	result, _ := handleKey(tea.KeyPressMsg{Code: '?', Text: "?"}, m)
	if result.showHelpPanel {
		t.Fatal("? in ScreenInventory should not toggle showHelpPanel")
	}
}

func TestHandleKey_BackslashTogglesSidebar(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenNormal
	m.showSidebar = false
	result, _ := handleKey(tea.KeyPressMsg{Code: '\\', Text: "\\"}, m)
	if !result.showSidebar {
		t.Fatal("\\ should toggle showSidebar to true")
	}
}

// --- Dungeon enemy combat tests ---

func TestHandleKey_MovingIntoEnemyTriggersCombat(t *testing.T) {
	m := NewModel()
	m.mode = ModeDungeon
	m.screenMode = ScreenNormal
	m.viewportW = 80
	m.viewportH = 24
	level := &DungeonLevel{}
	for x := 0; x < DungeonW; x++ {
		for y := 0; y < DungeonH; y++ {
			level.Cells[x][y].Kind = CellFloor
		}
	}
	tmpl := dungeonEnemyRoster[Plains]
	enemy := spawnEnemy(tmpl, 10, 9, 1, 5)
	level.Enemies = []*DungeonEnemy{enemy}
	m.currentDungeon = level
	m.playerPos = LocalCoord{X: 10, Y: 10}

	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyUp}, m)
	if result.screenMode != ScreenCombat {
		t.Fatal("moving into enemy should trigger ScreenCombat")
	}
	if result.playerPos.X != 10 || result.playerPos.Y != 10 {
		t.Fatal("player should not have moved")
	}
}

func TestLootAddedToInventoryAfterVictory(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenCombat
	tmpl := &EnemyTemplate{
		Name: "Test", Char: 'x', Color: "#fff",
		BaseHP: 1, MaxHP: 1,
		LootTable: []LootEntry{
			{Item: Item{Char: '!', Color: "#ff0", Name: "TestLoot", Count: 1}, Weight: 100},
		},
	}
	enemy := &DungeonEnemy{X: 5, Y: 5, Template: tmpl, HP: 0, MaxHP: 1}
	level := &DungeonLevel{}
	level.Enemies = []*DungeonEnemy{enemy}
	m.currentDungeon = level
	m.combatDungeonEnemy = enemy
	m.combatState = &CombatState{
		PlayerWon:   true,
		Player:      Combatant{Name: "Player", HP: 18, MaxHP: 20},
		Enemy:       Combatant{Name: "Test", HP: 0, MaxHP: 1},
		PendingLoot: Item{Char: '!', Color: "#ff0", Name: "TestLoot", Count: 1},
		LootMsg:     "Looted: TestLoot",
	}

	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	found := false
	for _, item := range result.inventory.Items {
		if item.Name == "TestLoot" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected TestLoot in inventory after victory")
	}
	if len(result.currentDungeon.Enemies) != 0 {
		t.Fatalf("expected 0 enemies after defeat, got %d", len(result.currentDungeon.Enemies))
	}
}

func TestLootDiscardedWhenInventoryFull(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenCombat
	// Fill inventory.
	for i := 0; i < InventoryMaxSlots; i++ {
		m.inventory.Items = append(m.inventory.Items, Item{Name: fmt.Sprintf("Item%d", i), Count: 1})
	}
	tmpl := &EnemyTemplate{
		Name: "Test", Char: 'x', Color: "#fff",
		BaseHP: 1, MaxHP: 1,
		LootTable: []LootEntry{
			{Item: Item{Name: "TestLoot", Count: 1}, Weight: 100},
		},
	}
	enemy := &DungeonEnemy{X: 5, Y: 5, Template: tmpl, HP: 0, MaxHP: 1}
	level := &DungeonLevel{}
	level.Enemies = []*DungeonEnemy{enemy}
	m.currentDungeon = level
	m.combatDungeonEnemy = enemy
	m.combatState = &CombatState{
		PlayerWon:   true,
		Player:      Combatant{Name: "Player", HP: 20, MaxHP: 20},
		Enemy:       Combatant{Name: "Test", HP: 0, MaxHP: 1},
		PendingLoot: Item{Name: "TestLoot", Count: 1},
	}

	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	if len(result.inventory.Items) != InventoryMaxSlots {
		t.Fatalf("inventory should still be full: got %d", len(result.inventory.Items))
	}
	for _, item := range result.inventory.Items {
		if item.Name == "TestLoot" {
			t.Fatal("TestLoot should not be in full inventory")
		}
	}
}

func TestDungeonMeta_BiomeRecorded(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	m.localMap.Objects[5][5] = &Object{Char: '>', Color: "#e8c96a", Blocking: false, Name: "Staircase Down"}

	_, _ = handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	key := m.worldPos
	meta, ok := m.dungeonMeta[key]
	if !ok {
		t.Fatal("DungeonMeta should be recorded on first descent")
	}
	tile := TileAt(m.worldPos.X, m.worldPos.Y, &m)
	if meta.Biome != tile.Biome {
		t.Fatalf("expected biome %d, got %d", tile.Biome, meta.Biome)
	}
}

// ── Combat speed key tests ──────────────────────────────────────────────────

func TestHandleKey_BracketIncreasesSpeed(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenCombat
	m.combatSpeed = CombatSpeedSlow
	m.combatState = &CombatState{Player: Combatant{Name: "Player", HP: 15, MaxHP: 20}, Enemy: Combatant{Name: "Wolf", HP: 8, MaxHP: 12}}
	result, _ := handleKey(tea.KeyPressMsg{Code: -1, Text: "]"}, m)
	if result.combatSpeed != CombatSpeedNormal {
		t.Fatalf("expected combatSpeed %d, got %d", CombatSpeedNormal, result.combatSpeed)
	}
}

func TestHandleKey_BracketClampsFast(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenCombat
	m.combatSpeed = CombatSpeedFast
	m.combatState = &CombatState{Player: Combatant{Name: "Player", HP: 15, MaxHP: 20}, Enemy: Combatant{Name: "Wolf", HP: 8, MaxHP: 12}}
	result, _ := handleKey(tea.KeyPressMsg{Code: -1, Text: "]"}, m)
	if result.combatSpeed != CombatSpeedFast {
		t.Fatalf("expected combatSpeed %d, got %d", CombatSpeedFast, result.combatSpeed)
	}
}

func TestHandleKey_LeftBracketDecreasesSpeed(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenCombat
	m.combatSpeed = CombatSpeedFast
	m.combatState = &CombatState{Player: Combatant{Name: "Player", HP: 15, MaxHP: 20}, Enemy: Combatant{Name: "Wolf", HP: 8, MaxHP: 12}}
	result, _ := handleKey(tea.KeyPressMsg{Code: -1, Text: "["}, m)
	if result.combatSpeed != CombatSpeedNormal {
		t.Fatalf("expected combatSpeed %d, got %d", CombatSpeedNormal, result.combatSpeed)
	}
}

func TestHandleKey_LeftBracketClampsSlow(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenCombat
	m.combatSpeed = CombatSpeedSlow
	m.combatState = &CombatState{Player: Combatant{Name: "Player", HP: 15, MaxHP: 20}, Enemy: Combatant{Name: "Wolf", HP: 8, MaxHP: 12}}
	result, _ := handleKey(tea.KeyPressMsg{Code: -1, Text: "["}, m)
	if result.combatSpeed != CombatSpeedSlow {
		t.Fatalf("expected combatSpeed %d, got %d", CombatSpeedSlow, result.combatSpeed)
	}
}

// ── Combat pause flow tests ──────────────────────────────────────────────────

func TestCombatEntry_SetsCombatPaused(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.localMap = &LocalMap{}
	m.playerPos = LocalCoord{X: 5, Y: 5}
	m.localMap.Animals = []*Animal{{X: 5, Y: 5, Char: 'w', Color: "#555", Name: "Wolf"}}

	result, cmd := handleKey(tea.KeyPressMsg{Code: 'g', Text: "g"}, m)
	if !result.combatPaused {
		t.Fatal("combat entry should set combatPaused = true")
	}
	if result.combatLogIndex != 0 {
		t.Fatalf("combat entry: combatLogIndex = %d, want 0", result.combatLogIndex)
	}
	if cmd != nil {
		t.Fatal("combat entry should not schedule a tick command (cmd should be nil)")
	}
}

func TestCombatPause_SpaceUnpauses(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenCombat
	m.combatPaused = true
	m.combatSpeed = CombatSpeedNormal
	m.combatState = &CombatState{
		Player: Combatant{Name: "Player", HP: 15, MaxHP: 20},
		Enemy:  Combatant{Name: "Wolf", HP: 8, MaxHP: 12},
		Round:  3,
	}

	result, cmd := handleKey(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "}, m)
	if result.combatPaused {
		t.Fatal("Space should set combatPaused = false")
	}
	if cmd == nil {
		t.Fatal("Space on paused combat should return non-nil cmd (tick)")
	}
}

func TestCombatPause_EnterUnpauses(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenCombat
	m.combatPaused = true
	m.combatSpeed = CombatSpeedNormal
	m.combatState = &CombatState{
		Player: Combatant{Name: "Player", HP: 15, MaxHP: 20},
		Enemy:  Combatant{Name: "Wolf", HP: 8, MaxHP: 12},
		Round:  3,
	}

	result, cmd := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	if result.combatPaused {
		t.Fatal("Enter should set combatPaused = false")
	}
	if cmd == nil {
		t.Fatal("Enter on paused combat should return non-nil cmd (tick)")
	}
}

func TestCombatActive_SpaceIsNoOp(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenCombat
	m.combatPaused = false
	m.combatSpeed = CombatSpeedNormal
	m.combatState = &CombatState{
		Player: Combatant{Name: "Player", HP: 15, MaxHP: 20},
		Enemy:  Combatant{Name: "Wolf", HP: 8, MaxHP: 12},
		Round:  3,
	}
	m.combatLogIndex = 1 // mid-playback, not at end

	result, cmd := handleKey(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "}, m)
	if result.combatPaused {
		t.Fatal("Space during active playback should not set combatPaused")
	}
	if cmd != nil {
		t.Fatal("Space during active playback should not schedule additional tick")
	}
}

// ── Campfire resting tests ───────────────────────────────────────────────────

func makeLocalModelWithFire(playerOnFire bool) Model {
	m := NewModel()
	m.mode = ModeLocal
	lm := &LocalMap{}
	// Place player at (5,5)
	m.playerPos = LocalCoord{X: 5, Y: 5}
	if playerOnFire {
		lm.Ground[5][5] = Ground{Char: '.', Color: "#888", Passable: true, HasFire: true}
	} else {
		lm.Ground[5][5] = Ground{Char: '.', Color: "#888", Passable: true}
	}
	m.localMap = lm
	return m
}

func TestRest_OnFireCell_Heals5HP(t *testing.T) {
	m := makeLocalModelWithFire(true)
	m.playerHP = 10
	result, _ := handleKey(tea.KeyPressMsg{Code: 'r', Text: "r"}, m)
	if result.playerHP != 15 {
		t.Fatalf("expected playerHP=15, got %d", result.playerHP)
	}
	if result.restCooldown != 60 {
		t.Fatalf("expected restCooldown=60, got %d", result.restCooldown)
	}
}

func TestRest_CapsAtMaxHP(t *testing.T) {
	m := makeLocalModelWithFire(true)
	m.playerHP = m.playerMaxHP - 2
	result, _ := handleKey(tea.KeyPressMsg{Code: 'r', Text: "r"}, m)
	if result.playerHP != m.playerMaxHP {
		t.Fatalf("expected playerHP=%d (capped), got %d", m.playerMaxHP, result.playerHP)
	}
}

func TestRest_OnNonFireCell_IsNoOp(t *testing.T) {
	m := makeLocalModelWithFire(false)
	m.playerHP = 10
	result, _ := handleKey(tea.KeyPressMsg{Code: 'r', Text: "r"}, m)
	if result.playerHP != 10 {
		t.Fatalf("expected playerHP unchanged=10, got %d", result.playerHP)
	}
	if result.restCooldown != 0 {
		t.Fatalf("expected restCooldown unchanged=0, got %d", result.restCooldown)
	}
}

func TestRest_DuringCooldown_IsNoOp(t *testing.T) {
	m := makeLocalModelWithFire(true)
	m.playerHP = 10
	m.restCooldown = 30
	result, _ := handleKey(tea.KeyPressMsg{Code: 'r', Text: "r"}, m)
	if result.playerHP != 10 {
		t.Fatalf("expected playerHP unchanged=10, got %d", result.playerHP)
	}
}

func TestRest_InDungeonMode_IsNoOp(t *testing.T) {
	m := NewModel()
	m.mode = ModeDungeon
	m.playerHP = 10
	result, _ := handleKey(tea.KeyPressMsg{Code: 'r', Text: "r"}, m)
	if result.playerHP != 10 {
		t.Fatalf("expected playerHP unchanged in dungeon mode, got %d", result.playerHP)
	}
}

// ── Death screen tests ───────────────────────────────────────────────────────

func makeDefeatedModel() Model {
	m := NewModel()
	m.screenMode = ScreenCombat
	m.combatState = &CombatState{
		Player:    Combatant{Name: "Player", HP: 0, MaxHP: 20},
		Enemy:     Combatant{Name: "Goblin", HP: 5, MaxHP: 10},
		PlayerWon: false,
		Round:     2,
		Log:       []string{"Round 1", "Round 2"},
	}
	m.combatLogIndex = 2 // playback complete
	return m
}

func TestDefeat_SetsScreenDeath(t *testing.T) {
	m := makeDefeatedModel()
	result, cmd := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	if result.screenMode != ScreenDeath {
		t.Fatalf("expected ScreenDeath after defeat, got %d", result.screenMode)
	}
	if cmd != nil {
		t.Fatal("defeat should not return tea.Quit command")
	}
}

func TestDefeat_RecordsKillerName(t *testing.T) {
	m := makeDefeatedModel()
	result, _ := handleKey(tea.KeyPressMsg{Code: tea.KeyEnter}, m)
	if result.deathKiller != "Goblin" {
		t.Fatalf("expected deathKiller='Goblin', got %q", result.deathKiller)
	}
}

func TestDeathScreen_RRestarts(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenDeath
	m.deathKiller = "Goblin"
	m.playerHP = 0
	result, cmd := handleKey(tea.KeyPressMsg{Code: 'r', Text: "r"}, m)
	if result.screenMode != ScreenNormal {
		t.Fatalf("expected ScreenNormal after restart, got %d", result.screenMode)
	}
	if result.playerHP != 20 {
		t.Fatalf("expected playerHP=20 after restart, got %d", result.playerHP)
	}
	if cmd == nil {
		t.Fatal("restart should return tickCmd to resume the tick loop")
	}
}

func TestDeathScreen_QQuits(t *testing.T) {
	m := NewModel()
	m.screenMode = ScreenDeath
	_, cmd := handleKey(tea.KeyPressMsg{Code: 'q', Text: "q"}, m)
	if cmd == nil {
		t.Fatal("q on death screen should return tea.Quit command")
	}
}
