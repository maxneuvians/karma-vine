package game

import "testing"

// --- hash ---

func TestHash_Deterministic(t *testing.T) {
	a := hash(5, 10, 42)
	b := hash(5, 10, 42)
	if a != b {
		t.Fatalf("hash(5,10,42) not deterministic: %v != %v", a, b)
	}
}

func TestHash_Range(t *testing.T) {
	// Sample a range of inputs and verify output is in [0, 1)
	for x := -5; x <= 5; x++ {
		for y := -5; y <= 5; y++ {
			v := hash(x, y, 999)
			if v < 0 || v >= 1 {
				t.Fatalf("hash(%d,%d,999)=%v not in [0,1)", x, y, v)
			}
		}
	}
}

// --- GenerateLocalMap ---

func TestGenerateLocalMap_Deterministic(t *testing.T) {
	a := GenerateLocalMap(3, 7, 12345, Forest)
	b := GenerateLocalMap(3, 7, 12345, Forest)
	if a.Ground != b.Ground {
		t.Fatal("GenerateLocalMap not deterministic: Ground arrays differ")
	}
	for x := 0; x < LocalMapW; x++ {
		for y := 0; y < LocalMapH; y++ {
			aObj := a.Objects[x][y]
			bObj := b.Objects[x][y]
			if (aObj == nil) != (bObj == nil) {
				t.Fatalf("GenerateLocalMap not deterministic: Objects[%d][%d] presence differs", x, y)
			}
			if aObj != nil && *aObj != *bObj {
				t.Fatalf("GenerateLocalMap not deterministic: Objects[%d][%d] value differs", x, y)
			}
		}
	}
}

func TestGenerateLocalMap_ForestHasTree(t *testing.T) {
	lm := GenerateLocalMap(0, 0, 1, Forest)
	for x := 0; x < LocalMapW; x++ {
		for y := 0; y < LocalMapH; y++ {
			if obj := lm.Objects[x][y]; obj != nil {
				if obj.Char == '♣' || obj.Char == '♠' {
					return // found a tree
				}
			}
		}
	}
	t.Fatal("Forest local map contains no tree object (♣ or ♠)")
}

func TestGenerateLocalMap_DesertHasCactus(t *testing.T) {
	lm := GenerateLocalMap(0, 0, 1, Desert)
	for x := 0; x < LocalMapW; x++ {
		for y := 0; y < LocalMapH; y++ {
			if obj := lm.Objects[x][y]; obj != nil && obj.Char == 'ψ' {
				return // found a cactus
			}
		}
	}
	t.Fatal("Desert local map contains no cactus object (ψ)")
}

// --- LocalMapFor ---

func TestLocalMapFor_PointerEquality(t *testing.T) {
	m := NewModel()
	p1 := LocalMapFor(10, 20, &m)
	p2 := LocalMapFor(10, 20, &m)
	if p1 != p2 {
		t.Fatal("LocalMapFor returned different pointers for the same coordinate")
	}
}

func TestLocalMapFor_CachesResult(t *testing.T) {
	m := NewModel()
	LocalMapFor(5, 5, &m)
	key := WorldCoord{X: 5, Y: 5}
	if _, ok := m.localCache[key]; !ok {
		t.Fatal("LocalMapFor did not store result in localCache")
	}
}

// --- buildLitMap ---

func TestBuildLitMap_CellsWithinRadiusAreLit(t *testing.T) {
	lm := &LocalMap{}
	// Place a fire at {20, 20}
	lm.Ground[20][20].HasFire = true
	buildLitMap(lm)

	// Fire cell itself: intensity 1.0
	if lm.LitMap[20][20] != 1.0 {
		t.Fatalf("LitMap[20][20] = %v, want 1.0", lm.LitMap[20][20])
	}
	// Distance 2: intensity = 1 - 2/5 = 0.6
	if lm.LitMap[22][20] < 0.5 || lm.LitMap[22][20] > 0.7 {
		t.Fatalf("LitMap[22][20] = %v, want ~0.6 (distance 2)", lm.LitMap[22][20])
	}
	// All cells within radius 4 must have positive intensity.
	for _, delta := range [][2]int{{0, 0}, {4, 0}, {0, 4}, {2, 2}, {3, 1}} {
		x, y := 20+delta[0], 20+delta[1]
		if lm.LitMap[x][y] <= 0 {
			t.Fatalf("LitMap[%d][%d] = %v, want > 0 (distance %d from fire)", x, y, lm.LitMap[x][y], delta[0]+delta[1])
		}
	}
}

func TestBuildLitMap_CellsBeyondRadiusNotLit(t *testing.T) {
	lm := &LocalMap{}
	lm.Ground[20][20].HasFire = true
	buildLitMap(lm)

	// Cells at Manhattan distance 5 must NOT be lit.
	for _, delta := range [][2]int{{5, 0}, {0, 5}, {3, 2}} {
		x, y := 20+delta[0], 20+delta[1]
		if lm.LitMap[x][y] != 0 {
			t.Fatalf("LitMap[%d][%d] = %v, want 0 (distance %d from fire)", x, y, lm.LitMap[x][y], delta[0]+delta[1])
		}
	}
}

