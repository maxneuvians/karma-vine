package game

import "testing"

func TestGenerateDungeonLevel_Determinism(t *testing.T) {
	a := GenerateDungeonLevel(42, 3, 7, 1, 5, Plains)
	b := GenerateDungeonLevel(42, 3, 7, 1, 5, Plains)
	for x := 0; x < DungeonW; x++ {
		for y := 0; y < DungeonH; y++ {
			if a.Cells[x][y].Kind != b.Cells[x][y].Kind {
				t.Fatalf("cell (%d,%d) kind mismatch: %d vs %d", x, y, a.Cells[x][y].Kind, b.Cells[x][y].Kind)
			}
		}
	}
}

func TestGenerateDungeonLevel_HasFloorCell(t *testing.T) {
	level := GenerateDungeonLevel(42, 0, 0, 1, 5, Plains)
	for x := 0; x < DungeonW; x++ {
		for y := 0; y < DungeonH; y++ {
			if level.Cells[x][y].Kind == CellFloor {
				return
			}
		}
	}
	t.Fatal("no floor cell found in generated dungeon level")
}

func TestGenerateDungeonLevel_UpStairOnFloor(t *testing.T) {
	level := GenerateDungeonLevel(42, 1, 2, 1, 5, Plains)
	cell := level.Cells[level.UpStair.X][level.UpStair.Y]
	if cell.Kind != CellFloor {
		t.Fatalf("up-stair at (%d,%d) is on cell kind %d, want CellFloor",
			level.UpStair.X, level.UpStair.Y, cell.Kind)
	}
}

func TestGenerateDungeonLevel_DownStairPresentNonFinal(t *testing.T) {
	level := GenerateDungeonLevel(42, 1, 2, 1, 5, Plains) // depth 1 < maxDepth 5
	if !level.HasDownStair {
		t.Fatal("expected down-stair on non-final level")
	}
	cell := level.Cells[level.DownStair.X][level.DownStair.Y]
	if cell.Kind != CellFloor {
		t.Fatalf("down-stair at (%d,%d) is on cell kind %d, want CellFloor",
			level.DownStair.X, level.DownStair.Y, cell.Kind)
	}
}

func TestGenerateDungeonLevel_NoDownStairFinalLevel(t *testing.T) {
	level := GenerateDungeonLevel(42, 1, 2, 5, 5, Plains) // depth == maxDepth
	if level.HasDownStair {
		t.Fatal("expected no down-stair on final level")
	}
}

func TestGenerateDungeonLevel_TorchOnWallCell(t *testing.T) {
	level := GenerateDungeonLevel(42, 5, 5, 1, 5, Plains)
	found := false
	for x := 0; x < DungeonW; x++ {
		for y := 0; y < DungeonH; y++ {
			obj := level.Cells[x][y].Object
			if obj != nil && obj.Char == '†' {
				found = true
				if level.Cells[x][y].Kind != CellWall {
					t.Fatalf("torch at (%d,%d) is on cell kind %d, want CellWall", x, y, level.Cells[x][y].Kind)
				}
			}
		}
	}
	if !found {
		t.Fatal("no torch found in generated dungeon level")
	}
}

func TestDungeonLevelFor_CacheHit(t *testing.T) {
	m := NewModel()
	m.globalSeed = 42
	m.dungeonMeta[WorldCoord{X: 1, Y: 2}] = DungeonMeta{MaxDepth: 5}

	first := DungeonLevelFor(1, 2, 1, &m)
	second := DungeonLevelFor(1, 2, 1, &m)

	if first != second {
		t.Fatal("expected cache hit to return same pointer")
	}
}

func TestDungeonLevelFor_CacheMiss(t *testing.T) {
	m := NewModel()
	m.globalSeed = 42

	level := DungeonLevelFor(3, 4, 1, &m)
	if level == nil {
		t.Fatal("expected non-nil DungeonLevel on cache miss")
	}
	// Verify it was cached.
	key := dungeonKey{wx: 3, wy: 4, depth: 1}
	if _, ok := m.dungeonCache[key]; !ok {
		t.Fatal("expected level to be stored in cache after miss")
	}
}

func TestDungeonMetaFor_MaxDepthRange(t *testing.T) {
	m := NewModel()
	m.globalSeed = 42

	for wx := 0; wx < 20; wx++ {
		meta := DungeonMetaFor(wx, 0, &m)
		if meta.MaxDepth < 5 || meta.MaxDepth > 10 {
			t.Fatalf("MaxDepth %d out of range [5,10] for wx=%d", meta.MaxDepth, wx)
		}
	}
}

func TestDungeonMetaFor_Stable(t *testing.T) {
	m := NewModel()
	m.globalSeed = 42

	first := DungeonMetaFor(5, 5, &m)
	second := DungeonMetaFor(5, 5, &m)
	if first.MaxDepth != second.MaxDepth {
		t.Fatalf("MaxDepth changed: %d vs %d", first.MaxDepth, second.MaxDepth)
	}
}

func TestGenerateDungeonLevel_AllObjectsHaveNames(t *testing.T) {
	level := GenerateDungeonLevel(42, 1, 1, 1, 5, Plains)
	for x := 0; x < DungeonW; x++ {
		for y := 0; y < DungeonH; y++ {
			if obj := level.Cells[x][y].Object; obj != nil && obj.Name == "" {
				t.Fatalf("dungeon object '%c' at (%d,%d) has empty Name", obj.Char, x, y)
			}
		}
	}
}

// 7.14 All unlit torches and braziers generated in dungeon have Pickupable true.
func TestGenerateDungeonLevel_UnlitTorchesPickupable(t *testing.T) {
	level := GenerateDungeonLevel(42, 1, 1, 1, 5, Plains)
	for x := 0; x < DungeonW; x++ {
		for y := 0; y < DungeonH; y++ {
			obj := level.Cells[x][y].Object
			if obj == nil {
				continue
			}
			if (obj.Char == '†' || obj.Char == 'Ω') && !obj.Lit {
				if !obj.Pickupable {
					t.Fatalf("unlit %q at (%d,%d) should have Pickupable=true", obj.Name, x, y)
				}
			}
		}
	}
}
