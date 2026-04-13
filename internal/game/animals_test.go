package game

import "testing"

// helper: Manhattan distance between two points
func manDist(ax, ay, bx, by int) int {
	return abs(ax-bx) + abs(ay-by)
}

// --- randomStep ---

func TestRandomStep_ValidDirection(t *testing.T) {
	for i := 0; i < 50; i++ {
		dx, dy := randomStep()
		if dx == 0 && dy == 0 {
			t.Fatal("randomStep returned (0,0)")
		}
		if dx < -1 || dx > 1 || dy < -1 || dy > 1 {
			t.Fatalf("randomStep out of range: (%d,%d)", dx, dy)
		}
	}
}

// --- fleeStep ---

func TestFleeStep_MovesAway(t *testing.T) {
	// Animal at (5,5), player at (5,6): best flee direction should increase distance.
	dx, dy := fleeStep(5, 5, 5, 6)
	newDist := manDist(5+dx, 5+dy, 5, 6)
	oldDist := manDist(5, 5, 5, 6)
	if newDist <= oldDist {
		t.Fatalf("fleeStep did not increase distance: old=%d new=%d dir=(%d,%d)", oldDist, newDist, dx, dy)
	}
}

// --- moveAnimals ---

// 4.1 Non-flee animal should move at least once across 10 ticks.
func TestMoveAnimals_RandomWalkChangesPosition(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	lm := &LocalMap{}
	lm.Animals = []*Animal{{X: 10, Y: 9, Char: 'd', Color: "#888", Flee: false}}
	m.localMap = lm

	startX, startY := lm.Animals[0].X, lm.Animals[0].Y
	moved := false
	for i := 0; i < 10; i++ {
		next, _ := m.Update(TickMsg{})
		m = next.(Model)
		a := m.localMap.Animals[0]
		if a.X != startX || a.Y != startY {
			moved = true
		}
	}
	if !moved {
		t.Fatal("non-flee animal did not move after 10 ticks")
	}
}

// 4.2 Flee animal adjacent to player moves further away after one tick.
func TestMoveAnimals_FleeMovesAway(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	m.playerPos = LocalCoord{X: 5, Y: 6}
	lm := &LocalMap{}
	lm.Animals = []*Animal{{X: 5, Y: 5, Char: 'r', Color: "#aaa", Flee: true}}
	m.localMap = lm

	initialDist := manDist(5, 5, 5, 6)
	m.Update(TickMsg{})
	a := m.localMap.Animals[0]
	newDist := manDist(a.X, a.Y, 5, 6)
	if newDist <= initialDist {
		t.Fatalf("flee animal did not move away: initial dist=%d new dist=%d pos=(%d,%d)",
			initialDist, newDist, a.X, a.Y)
	}
}

// 4.3 Animal at {0,0} with all valid moves blocked stays put.
func TestMoveAnimals_StaysWhenAllBlocked(t *testing.T) {
	m := NewModel()
	m.mode = ModeLocal
	lm := &LocalMap{}
	// From (0,0) the in-bounds cells are (0,1), (1,0), (1,1) — block them all.
	lm.Objects[0][1] = &Object{Char: '◉', Blocking: true}
	lm.Objects[1][0] = &Object{Char: '◉', Blocking: true}
	lm.Objects[1][1] = &Object{Char: '◉', Blocking: true}
	lm.Animals = []*Animal{{X: 0, Y: 0, Char: 'r', Color: "#aaa", Flee: false}}
	m.localMap = lm

	m.Update(TickMsg{})
	a := m.localMap.Animals[0]
	if a.X != 0 || a.Y != 0 {
		t.Fatalf("animal should stay at (0,0) but moved to (%d,%d)", a.X, a.Y)
	}
}

// --- TickMsg guard ---

func TestTickMsg_SkipsWhenNotLocalMode(t *testing.T) {
	m := NewModel()
	m.mode = ModeWorld
	lm := &LocalMap{}
	lm.Animals = []*Animal{{X: 10, Y: 9, Char: 'd', Color: "#888", Flee: false}}
	m.localMap = lm

	startX, startY := lm.Animals[0].X, lm.Animals[0].Y
	for i := 0; i < 10; i++ {
		next, _ := m.Update(TickMsg{})
		m = next.(Model)
	}
	a := m.localMap.Animals[0]
	if a.X != startX || a.Y != startY {
		t.Fatalf("animals moved in world mode: expected (%d,%d), got (%d,%d)",
			startX, startY, a.X, a.Y)
	}
}

func TestTickMsg_ReschedulesCmd(t *testing.T) {
	m := NewModel()
	_, cmd := m.Update(TickMsg{})
	if cmd == nil {
		t.Fatal("TickMsg handler must return a non-nil reschedule command")
	}
}
