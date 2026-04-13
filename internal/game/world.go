package game

import (
	"math"

	opensimplex "github.com/ojrac/opensimplex-go"
)

const WorldNoiseScale = 0.07

// Chunk is a 32×32 section of the world map.
type Chunk struct {
	Tiles [32][32]Tile
}

// generateChunk produces a Chunk for the given chunk coordinates and seed.
// It is a pure function of (cx, cy, globalSeed).
func generateChunk(cx, cy, globalSeed int) *Chunk {
	elevNoise := opensimplex.NewNormalized(int64(globalSeed))
	moistNoise := opensimplex.NewNormalized(int64(globalSeed + 1))

	chunk := &Chunk{}
	for lx := 0; lx < 32; lx++ {
		for ly := 0; ly < 32; ly++ {
			worldX := cx*32 + lx
			worldY := cy*32 + ly
			nx := float64(worldX) * WorldNoiseScale
			ny := float64(worldY) * WorldNoiseScale
			e := elevNoise.Eval2(nx, ny)
			m := moistNoise.Eval2(nx, ny)
			biome, ch, color := classifyBiome(e, m)
			chunk.Tiles[lx][ly] = Tile{
				Biome:     biome,
				Char:      ch,
				Color:     color,
				Elevation: e,
				Moisture:  m,
			}
		}
	}
	return chunk
}

// classifyBiome returns the Biome, display rune, and hex color for the given
// elevation and moisture values. Rules are evaluated in the order specified by
// the brief.
func classifyBiome(e, m float64) (Biome, rune, string) {
	switch {
	case e < 0.28:
		return DeepOcean, '≋', "#1a6fa8"
	case e < 0.36:
		return ShallowWater, '≈', "#2e9ecf"
	case e < 0.40:
		return Beach, '·', "#e8c96a"
	case e < 0.50 && m > 0.55:
		return Forest, '♣', "#2d7a1f"
	case e < 0.50:
		return Plains, '░', "#5aad3f"
	case e < 0.62 && m > 0.45:
		return DenseForest, '♠', "#3d6b3a"
	case e < 0.62 && m < 0.35:
		return Desert, '~', "#c8a46a"
	case e < 0.78:
		return Mountains, '▲', "#8fa89c"
	default:
		return Snow, '*', "#ccd9e0"
	}
}

// TileAt returns the Tile at world-space coordinates (worldX, worldY).
// If the containing chunk is not yet cached in m.chunks, it is generated and
// stored before the tile is returned. Negative coordinates are handled
// correctly via math.Floor division.
func TileAt(worldX, worldY int, m *Model) Tile {
	cx := int(math.Floor(float64(worldX) / 32))
	cy := int(math.Floor(float64(worldY) / 32))
	coord := ChunkCoord{X: cx, Y: cy}
	if _, ok := m.chunks[coord]; !ok {
		m.chunks[coord] = generateChunk(cx, cy, m.globalSeed)
	}
	localX := worldX - cx*32
	localY := worldY - cy*32
	return m.chunks[coord].Tiles[localX][localY]
}
