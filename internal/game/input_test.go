package game

import "testing"

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
