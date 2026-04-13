package game

import (
	"math"
	"math/rand"

	opensimplex "github.com/ojrac/opensimplex-go"
)

// Chunk is a 32×32 section of the world map.
type Chunk struct {
	Tiles [32][32]Tile
}

// fbm computes fractional Brownian motion by summing octaves of noise.
// Returns a value in approximately [-1, 1].
func fbm(n opensimplex.Noise, x, y float64, octaves int) float64 {
	v, amp, freq, maxV := 0.0, 1.0, 1.0, 0.0
	for i := 0; i < octaves; i++ {
		v += n.Eval2(x*freq, y*freq) * amp
		maxV += amp
		amp *= 0.5
		freq *= 2.0
	}
	return v / maxV
}

// ridgedFBM computes ridged multi-fractal noise for mountain ranges.
// Returns a value in approximately [0, 1].
func ridgedFBM(n opensimplex.Noise, x, y float64, octaves int) float64 {
	v, amp, freq, maxV := 0.0, 1.0, 1.0, 0.0
	for i := 0; i < octaves; i++ {
		raw := n.Eval2(x*freq, y*freq)
		v += (1.0 - math.Abs(raw)) * amp
		maxV += amp
		amp *= 0.5
		freq *= 2.0
	}
	return v / maxV
}

// norm maps a value from [-1, 1] to [0, 1].
func norm(v float64) float64 {
	return (v + 1.0) * 0.5
}

// computeTemperature returns a temperature value in [0, 1] for a given world Y
// coordinate. It applies a cosine latitudinal gradient with a half-period of
// 400 tiles (so temperature cycles hot→cold every 800 tiles) blended with a
// small noise perturbation of amplitude 0.12 to break horizontal symmetry.
// Values near 1.0 are hot (equatorial, Y≈0); values near 0.0 are cold (polar).
func computeTemperature(worldY int, noise opensimplex.Noise) float64 {
	const halfPeriod = 400.0
	base := (math.Cos(math.Pi*math.Abs(float64(worldY))/halfPeriod) + 1.0) * 0.5
	perturbation := noise.Eval2(float64(worldY)*0.003, 0) * 0.12
	t := base + perturbation
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

// generateChunk produces a Chunk for the given chunk coordinates and seed.
// It is a pure function of (cx, cy, globalSeed).
//
// Elevation is built from three layers:
//   - continent (very large scale, domain-warped) — sets ocean vs land
//   - terrain   (medium scale FBM)               — valley and hill detail
//   - ridge     (ridged FBM, land-weighted)       — mountain ranges
//
// Moisture uses its own independent FBM field.
// A thin river-noise band carves shallow-water channels through mid-elevation land.
func generateChunk(cx, cy, globalSeed int) *Chunk {
	continentNoise := opensimplex.New(int64(globalSeed))
	terrainNoise := opensimplex.New(int64(globalSeed + 1))
	ridgeNoise := opensimplex.New(int64(globalSeed + 2))
	moistureNoise := opensimplex.New(int64(globalSeed + 3))
	riverNoise := opensimplex.New(int64(globalSeed + 4))
	warpXNoise := opensimplex.New(int64(globalSeed + 5))
	warpYNoise := opensimplex.New(int64(globalSeed + 6))
	temperatureNoise := opensimplex.New(int64(globalSeed + 7))

	chunk := &Chunk{}
	for lx := 0; lx < 32; lx++ {
		for ly := 0; ly < 32; ly++ {
			worldX := cx*32 + lx
			worldY := cy*32 + ly
			nx := float64(worldX)
			ny := float64(worldY)

			// Domain warp: organically offset the continent sampling point to
			// produce irregular coastlines and natural-looking land masses.
			wx := fbm(warpXNoise, nx*0.006, ny*0.006, 3) * 40.0
			wy := fbm(warpYNoise, nx*0.006, ny*0.006, 3) * 40.0

			// Continental base: very large scale; squared to skew toward ocean.
			continent := norm(fbm(continentNoise, (nx+wx)*0.0045, (ny+wy)*0.0045, 5))
			continent = continent * continent

			// Terrain detail: medium-scale FBM for hills and valleys.
			terrain := norm(fbm(terrainNoise, nx*0.0135, ny*0.0135, 7))

			// Mountain ridges: appear only on elevated land (continent weighting).
			ridge := ridgedFBM(ridgeNoise, nx*0.024, ny*0.024, 5)

			// Combined elevation in [0, 1].
			e := continent*0.55 + terrain*0.28 + ridge*continent*0.17

			// Moisture: independent climate field in [0, 1].
			m := norm(fbm(moistureNoise, nx*0.009, ny*0.009, 4))

			// River carving: interpret thin noise-band valleys as rivers.
			// Only carves through mid-elevation land (avoids ocean and mountains).
			riverV := norm(fbm(riverNoise, nx*0.0165, ny*0.0165, 3))
			if math.Abs(riverV-0.5) < 0.022 && e > 0.41 && e < 0.66 {
				e = 0.33
			}

			temperature := computeTemperature(worldY, temperatureNoise)
			biome, ch, color := classifyBiome(e, m, temperature)
			chunk.Tiles[lx][ly] = Tile{
				Biome:       biome,
				Char:        ch,
				Color:       color,
				Elevation:   e,
				Moisture:    m,
				Temperature: temperature,
			}
		}
	}
	return chunk
}

// classifyBiome returns the Biome, display rune, and hex color for the given
// elevation, moisture, and temperature values. Water/coast cases apply across
// all temperature bands; land biomes are split into hot (≥0.65), temperate
// (0.35–0.65), and cold (<0.35) bands.
func classifyBiome(e, m, temperature float64) (Biome, rune, string) {
	// Water and coast — independent of temperature.
	switch {
	case e < 0.28:
		return DeepOcean, '≋', "#1a6fa8"
	case e < 0.36:
		return ShallowWater, '≈', "#2e9ecf"
	case e < 0.40:
		return Beach, '·', "#e8c96a"
	}

	// Land biomes — split by temperature band.
	switch {
	case temperature >= 0.65: // hot band
		switch {
		case e < 0.50 && m > 0.55:
			return Jungle, '♣', "#1a7a2e"
		case e < 0.50 && m > 0.30:
			return Savanna, 'ˬ', "#b5a04a"
		case e < 0.50:
			return AridSteppe, '·', "#c9a97a"
		case e < 0.62 && m > 0.45:
			return Jungle, '♣', "#1a7a2e"
		case e < 0.62:
			return Desert, '~', "#c8a46a"
		case e < 0.78:
			return Mountains, '▲', "#8fa89c"
		default:
			return Snow, '*', "#ccd9e0"
		}
	case temperature >= 0.35: // temperate band
		switch {
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
	default: // cold band (temperature < 0.35)
		switch {
		case e < 0.50 && m > 0.50:
			return Taiga, '♠', "#3a6b52"
		case e < 0.50:
			return Tundra, '∙', "#8ab08a"
		case e < 0.62:
			return Taiga, '♠', "#3a6b52"
		case e < 0.78:
			return Mountains, '▲', "#8fa89c"
		default:
			return Snow, '*', "#ccd9e0"
		}
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

// isLandBiome returns true for biomes the player can walk on (not ocean/water).
func isLandBiome(b Biome) bool {
	return b != DeepOcean && b != ShallowWater
}

// findWorldSpawn picks a random land biome then searches outward from the
// origin for the first tile of that biome, giving a varied starting location.
func findWorldSpawn(m *Model) WorldCoord {
	landBiomes := []Biome{Beach, Plains, Forest, DenseForest, Desert, Mountains, Snow}
	rng := rand.New(rand.NewSource(int64(m.globalSeed)))
	target := landBiomes[rng.Intn(len(landBiomes))]

	var anyLand *WorldCoord
	for radius := 0; radius <= 512; radius++ {
		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {
				if abs(dx) != radius && abs(dy) != radius {
					continue
				}
				b := TileAt(dx, dy, m).Biome
				if b == target {
					return WorldCoord{X: dx, Y: dy}
				}
				if anyLand == nil && isLandBiome(b) {
					c := WorldCoord{X: dx, Y: dy}
					anyLand = &c
				}
			}
		}
	}
	if anyLand != nil {
		return *anyLand
	}
	return WorldCoord{}
}
