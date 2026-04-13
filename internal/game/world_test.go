package game

import (
	"testing"

	opensimplex "github.com/ojrac/opensimplex-go"
)

// --- classifyBiome ---

func TestClassifyBiome_DeepOcean(t *testing.T) {
	biome, ch, color := classifyBiome(0.20, 0.5, 0.5)
	if biome != DeepOcean {
		t.Fatalf("expected DeepOcean, got %d", biome)
	}
	if ch != '≋' {
		t.Fatalf("expected '≋', got %c", ch)
	}
	if color != "#1a6fa8" {
		t.Fatalf("expected #1a6fa8, got %s", color)
	}
}

func TestClassifyBiome_Forest(t *testing.T) {
	biome, ch, color := classifyBiome(0.48, 0.60, 0.50)
	if biome != Forest {
		t.Fatalf("expected Forest, got %d", biome)
	}
	if ch != '♣' {
		t.Fatalf("expected '♣', got %c", ch)
	}
	if color != "#2d7a1f" {
		t.Fatalf("expected #2d7a1f, got %s", color)
	}
}

func TestClassifyBiome_Snow(t *testing.T) {
	biome, ch, color := classifyBiome(0.85, 0.5, 0.5)
	if biome != Snow {
		t.Fatalf("expected Snow, got %d", biome)
	}
	if ch != '*' {
		t.Fatalf("expected '*', got %c", ch)
	}
	if color != "#ccd9e0" {
		t.Fatalf("expected #ccd9e0, got %s", color)
	}
}

func TestClassifyBiome_Jungle(t *testing.T) {
	biome, ch, color := classifyBiome(0.45, 0.70, 0.80)
	if biome != Jungle {
		t.Fatalf("expected Jungle, got %d", biome)
	}
	if ch != '♣' {
		t.Fatalf("expected '♣', got %c", ch)
	}
	if color != "#1a7a2e" {
		t.Fatalf("expected #1a7a2e, got %s", color)
	}
}

func TestClassifyBiome_Savanna(t *testing.T) {
	biome, ch, color := classifyBiome(0.45, 0.45, 0.80)
	if biome != Savanna {
		t.Fatalf("expected Savanna, got %d", biome)
	}
	if ch != 'ˬ' {
		t.Fatalf("expected 'ˬ', got %c", ch)
	}
	if color != "#b5a04a" {
		t.Fatalf("expected #b5a04a, got %s", color)
	}
}

func TestClassifyBiome_AridSteppe(t *testing.T) {
	biome, ch, color := classifyBiome(0.45, 0.20, 0.80)
	if biome != AridSteppe {
		t.Fatalf("expected AridSteppe, got %d", biome)
	}
	if ch != '·' {
		t.Fatalf("expected '·', got %c", ch)
	}
	if color != "#c9a97a" {
		t.Fatalf("expected #c9a97a, got %s", color)
	}
}

func TestClassifyBiome_Tundra(t *testing.T) {
	biome, ch, color := classifyBiome(0.45, 0.20, 0.20)
	if biome != Tundra {
		t.Fatalf("expected Tundra, got %d", biome)
	}
	if ch != '∙' {
		t.Fatalf("expected '∙', got %c", ch)
	}
	if color != "#8ab08a" {
		t.Fatalf("expected #8ab08a, got %s", color)
	}
}

func TestClassifyBiome_Taiga(t *testing.T) {
	biome, ch, color := classifyBiome(0.45, 0.70, 0.20)
	if biome != Taiga {
		t.Fatalf("expected Taiga, got %d", biome)
	}
	if ch != '♠' {
		t.Fatalf("expected '♠', got %c", ch)
	}
	if color != "#3a6b52" {
		t.Fatalf("expected #3a6b52, got %s", color)
	}
}

// --- computeTemperature ---

func TestComputeTemperature_Equator(t *testing.T) {
	n := opensimplex.New(7)
	temp := computeTemperature(0, n)
	if temp <= 0.80 {
		t.Fatalf("expected temperature > 0.80 at Y=0, got %f", temp)
	}
}

func TestComputeTemperature_NorthPole(t *testing.T) {
	n := opensimplex.New(7)
	temp := computeTemperature(400, n)
	if temp >= 0.40 {
		t.Fatalf("expected temperature < 0.40 at Y=400, got %f", temp)
	}
}

// --- TileAt smoke test ---

func TestTileAt_Smoke(t *testing.T) {
	validBiomes := map[Biome]bool{
		DeepOcean: true, ShallowWater: true, Beach: true,
		Forest: true, Plains: true, DenseForest: true,
		Desert: true, Mountains: true, Snow: true,
		Jungle: true, Savanna: true, AridSteppe: true,
		Tundra: true, Taiga: true,
	}
	m := NewModel()
	// Sample 1 000 coordinates across a large range
	coords := make([][2]int, 0, 1000)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			for k := 0; k < 10; k++ {
				x := (i-5)*1000 + k*37
				y := (j-5)*1000 + k*53
				coords = append(coords, [2]int{x, y})
			}
		}
	}
	for _, c := range coords {
		tile := TileAt(c[0], c[1], &m)
		if !validBiomes[tile.Biome] {
			t.Fatalf("TileAt(%d,%d): invalid biome %d", c[0], c[1], tile.Biome)
		}
		if tile.Char == 0 {
			t.Fatalf("TileAt(%d,%d): zero Char", c[0], c[1])
		}
		if tile.Color == "" {
			t.Fatalf("TileAt(%d,%d): empty Color", c[0], c[1])
		}
		if tile.Elevation < 0 || tile.Elevation > 1 {
			t.Fatalf("TileAt(%d,%d): Elevation %f out of [0,1]", c[0], c[1], tile.Elevation)
		}
		if tile.Temperature < 0 || tile.Temperature > 1 {
			t.Fatalf("TileAt(%d,%d): Temperature %f out of [0,1]", c[0], c[1], tile.Temperature)
		}
	}
}
