package game

import (
	"math"
	"math/rand"

	opensimplex "github.com/ojrac/opensimplex-go"
)

const localNoiseScale = 0.12

// hash derives a deterministic float64 in [0, 1) from integer coordinates and
// a seed. It uses the exact bit-manipulation formula from the brief.
func hash(x, y, seed int) float64 {
	h := x*1619 + y*31337 + seed*6971
	h = (h ^ (h >> 16)) * 0x45d9f3b
	h = (h ^ (h >> 16)) * 0x45d9f3b
	h = h ^ (h >> 16)
	return float64(uint32(h)) / float64(0xffffffff)
}

// ── Biome content tables ────────────────────────────────────────────────────

type groundVariant struct {
	Char     rune
	Color    string
	Passable bool
}

type objectEntry struct {
	Char     rune
	Color    string
	Blocking bool
	Name     string
}

type animalEntry struct {
	Char  rune
	Color string
	Flee  bool
	Name  string
}

type biomeContent struct {
	ground          []groundVariant
	objects         []objectEntry
	animals         []animalEntry
	objectThreshold float64
	maxAnimals      int
	fireThreshold   float64 // noise threshold above which a cell has HasFire; 0 = no fires
}

var biomeTable = map[Biome]biomeContent{
	Forest: {
		ground: []groundVariant{
			{'.', "#5aad3f", true},
			{',', "#4a9a30", true},
			{'\'', "#3d8728", true},
			{';', "#6cbf50", true},
		},
		objects: []objectEntry{
			{'♣', "#2d7a1f", true, "Tree"},
			{'♠', "#3d6b3a", true, "Pine"},
		},
		animals: []animalEntry{
			{'d', "#8B4513", true, "Deer"},
			{'r', "#aaaaaa", true, "Rabbit"},
		},
		objectThreshold: 0.40,
		maxAnimals:      3,
	},
	Desert: {
		ground: []groundVariant{
			{'.', "#c8a46a", true},
			{',', "#d4b47a", true},
			{'·', "#baa060", true},
		},
		objects: []objectEntry{
			{'ψ', "#5aad3f", true, "Cactus"},
			{'○', "#8fa89c", true, "Rock"},
		},
		animals: []animalEntry{
			{'s', "#c8a46a", false, "Snake"},
			{'l', "#5aad3f", false, "Lizard"},
		},
		objectThreshold: 0.50,
		maxAnimals:      2,
	},
	Plains: {
		ground: []groundVariant{
			{'.', "#6cbf50", true},
			{',', "#5aad3f", true},
			{'\'', "#78cc5c", true},
			{';', "#88d96c", true},
		},
		objects: []objectEntry{
			{'⌂', "#7bc96f", false, "Shelter"},
			{'✿', "#f0d060", false, "Flower"},
		},
		animals: []animalEntry{
			{'r', "#aaaaaa", true, "Rabbit"},
			{'b', "#2e9ecf", true, "Bird"},
		},
		objectThreshold: 0.65,
		maxAnimals:      4,		fireThreshold:   0.90,	},
	DenseForest: {
		ground: []groundVariant{
			{'.', "#3d6b3a", true},
			{',', "#2d5a28", true},
			{'\'', "#4a7a44", true},
		},
		objects: []objectEntry{
			{'♠', "#3d6b3a", true, "Pine"},
			{'♣', "#2d7a1f", true, "Tree"},
		},
		animals: []animalEntry{
			{'d', "#8B4513", true, "Deer"},
			{'w', "#555555", false, "Wolf"},
		},
		objectThreshold: 0.30,
		maxAnimals:      2,
	},
	Mountains: {
		ground: []groundVariant{
			{'.', "#8fa89c", true},
			{'·', "#7a9490", true},
			{'°', "#6a8480", true},
		},
		objects: []objectEntry{
			{'◉', "#8fa89c", true, "Boulder"},
			{'▲', "#4a6060", true, "Peak"},
		},
		animals: []animalEntry{
			{'g', "#dddddd", true, "Goat"},
			{'e', "#8B4513", true, "Eagle"},
		},
		objectThreshold: 0.50,
		maxAnimals:      2,
		fireThreshold:   0.92,
	},
	Snow: {
		ground: []groundVariant{
			{'.', "#ccd9e0", true},
			{'*', "#dde8f0", true},
			{'·', "#c0d0dc", true},
			{'°', "#e0ecf4", true},
		},
		objects: []objectEntry{
			{'◆', "#ccd9e0", true, "Ice Rock"},
			{'❄', "#ffffff", false, "Snowflake"},
		},
		animals: []animalEntry{
			{'B', "#eeeeee", false, "Bear"},
			{'r', "#ffffff", true, "Rabbit"},
		},
		objectThreshold: 0.55,
		maxAnimals:      2,
	},
	Beach: {
		ground: []groundVariant{
			{'.', "#e8c96a", true},
			{'·', "#d4b55a", true},
			{',', "#f0d47a", true},
		},
		objects: []objectEntry{
			{'○', "#c8a46a", false, "Shell"},
			{'⊙', "#e8c96a", false, "Pebble"},
		},
		animals: []animalEntry{
			{'c', "#c0392b", true, "Crab"},
			{'s', "#ffffff", true, "Seagull"},
		},
		objectThreshold: 0.65,
		maxAnimals:      3,
		fireThreshold:   0.88,
	},
	Jungle: {
		ground: []groundVariant{
			{'.', "#1a7a2e", true},
			{',', "#15692a", true},
			{'\'', "#228b38", true},
			{';', "#2a9940", true},
		},
		objects: []objectEntry{
			{'♣', "#1a7a2e", true, "Tree"},
			{'♠', "#145520", true, "Pine"},
		},
		animals: []animalEntry{
			{'b', "#2e9ecf", true, "Bird"},
			{'s', "#228b38", false, "Snake"},
		},
		objectThreshold: 0.28,
		maxAnimals:      4,
	},
	Savanna: {
		ground: []groundVariant{
			{'.', "#b5a04a", true},
			{',', "#c4ae55", true},
			{'\'', "#a8943f", true},
			{';', "#cabb60", true},
		},
		objects: []objectEntry{
			{'♣', "#8b6914", true, "Acacia"},
			{'○', "#8fa89c", false, "Rock"},
		},
		animals: []animalEntry{
			{'d', "#c8a020", true, "Antelope"},
			{'b', "#c0a030", true, "Bird"},
		},
		objectThreshold: 0.60,
		maxAnimals:      4,
		fireThreshold:   0.87,
	},
	AridSteppe: {
		ground: []groundVariant{
			{'.', "#c9a97a", true},
			{'·', "#bfa070", true},
			{',', "#d4b484", true},
		},
		objects: []objectEntry{
			{'○', "#8fa89c", false, "Rock"},
			{'ψ', "#7a9060", true, "Scrub"},
		},
		animals: []animalEntry{
			{'l', "#c9a97a", false, "Lizard"},
			{'s', "#d4a060", false, "Snake"},
		},
		objectThreshold: 0.62,
		maxAnimals:      2,
	},
	Tundra: {
		ground: []groundVariant{
			{'.', "#8ab08a", true},
			{'·', "#7a9e7a", true},
			{',', "#96bc96", true},
			{'°', "#6a8e6a", true},
		},
		objects: []objectEntry{
			{'○', "#8fa89c", false, "Rock"},
			{'◆', "#a0b8a0", false, "Stone"},
		},
		animals: []animalEntry{
			{'r', "#dddddd", true, "Rabbit"},
			{'B', "#dddddd", false, "Bear"},
		},
		objectThreshold: 0.68,
		maxAnimals:      2,
	},
	Taiga: {
		ground: []groundVariant{
			{'.', "#3a6b52", true},
			{',', "#2e5a44", true},
			{'\'', "#457a60", true},
		},
		objects: []objectEntry{
			{'♠', "#3a6b52", true, "Pine"},
			{'♣', "#2d5a40", true, "Tree"},
		},
		animals: []animalEntry{
			{'w', "#888888", false, "Wolf"},
			{'d', "#8B4513", true, "Deer"},
		},
		objectThreshold: 0.35,
		maxAnimals:      2,
	},
}

// fallbackContent is used for biomes with no local content table (water tiles).
var fallbackContent = biomeContent{
	ground: []groundVariant{
		{'~', "#2e9ecf", true},
		{'≈', "#1a6fa8", true},
	},
	objects:         nil,
	animals:         nil,
	objectThreshold: 1.0,
	maxAnimals:      0,
}

// ── GenerateLocalMap ────────────────────────────────────────────────────────

// GenerateLocalMap produces a deterministic 42×18 LocalMap for the given
// world tile coordinates, global seed, and biome.
func GenerateLocalMap(worldX, worldY, globalSeed int, biome Biome) *LocalMap {
	hashVal := hash(worldX, worldY, globalSeed)
	localSeed := int64(hashVal * float64(math.MaxInt64))

	content, ok := biomeTable[biome]
	if !ok {
		content = fallbackContent
	}

	terrainNoise := opensimplex.New(localSeed)
	objectNoise := opensimplex.New(localSeed + 1)
	fireNoise := opensimplex.New(localSeed + 3)
	rng := rand.New(rand.NewSource(localSeed + 2))

	lm := &LocalMap{}

	// 4.4 Populate ground and objects
	for x := 0; x < LocalMapW; x++ {
		for y := 0; y < LocalMapH; y++ {
			nx := float64(x) * localNoiseScale
			ny := float64(y) * localNoiseScale

			// Normalize raw noise [-1,1] → [0,1]
			tn := (terrainNoise.Eval2(nx, ny) + 1) / 2
			on := (objectNoise.Eval2(nx, ny) + 1) / 2
			fn := (fireNoise.Eval2(nx, ny) + 1) / 2

			// Pick ground variant
			idx := int(tn * float64(len(content.ground)))
			if idx >= len(content.ground) {
				idx = len(content.ground) - 1
			}
			g := content.ground[idx]
			hasFire := content.fireThreshold > 0 && fn > content.fireThreshold
			lm.Ground[x][y] = Ground{Char: g.Char, Color: g.Color, Passable: g.Passable, HasFire: hasFire}

			// 4.5 Place object if noise exceeds threshold (fire cells skip objects)
			if !hasFire && len(content.objects) > 0 && on > content.objectThreshold {
				oe := content.objects[rng.Intn(len(content.objects))]
				lm.Objects[x][y] = &Object{Char: oe.Char, Color: oe.Color, Blocking: oe.Blocking, Name: oe.Name}
			}
		}
	}

	// 4.6 Place animals at passable positions
	if len(content.animals) > 0 && content.maxAnimals > 0 {
		// Collect all passable, unblocked positions
		type pos struct{ x, y int }
		var passable []pos
		for x := 0; x < LocalMapW; x++ {
			for y := 0; y < LocalMapH; y++ {
				if lm.Ground[x][y].Passable &&
					(lm.Objects[x][y] == nil || !lm.Objects[x][y].Blocking) {
					passable = append(passable, pos{x, y})
				}
			}
		}
		// Shuffle and place up to maxAnimals
		rng.Shuffle(len(passable), func(i, j int) { passable[i], passable[j] = passable[j], passable[i] })
		count := content.maxAnimals
		if count > len(passable) {
			count = len(passable)
		}
		for i := 0; i < count; i++ {
			ae := content.animals[rng.Intn(len(content.animals))]
			lm.Animals = append(lm.Animals, &Animal{
				X:     passable[i].x,
				Y:     passable[i].y,
				Char:  ae.Char,
				Color: ae.Color,
				Flee:  ae.Flee,
				Name:  ae.Name,
			})
		}
	}

	// Place dungeon entrance on a passable cell not occupied by another object.
	{
		type pos struct{ x, y int }
		var candidates []pos
		for x := 0; x < LocalMapW; x++ {
			for y := 0; y < LocalMapH; y++ {
				if lm.Ground[x][y].Passable && lm.Objects[x][y] == nil && !lm.Ground[x][y].HasFire {
					candidates = append(candidates, pos{x, y})
				}
			}
		}
		if len(candidates) > 0 {
			p := candidates[rng.Intn(len(candidates))]
			lm.Objects[p.x][p.y] = &Object{Char: '>', Color: "#ff3333", Blocking: false, Name: "Dungeon Entrance"}
		}
	}

	buildLitMap(lm)
	return lm
}

// buildLitMap precomputes per-cell fire illumination intensity.
// Each cell stores the maximum intensity contributed by nearby fires.
// Intensity = 1.0 at the fire cell, falling off linearly to 0 at radius+1.
func buildLitMap(lm *LocalMap) {
	const fireRadius = 4
	for fx := 0; fx < LocalMapW; fx++ {
		for fy := 0; fy < LocalMapH; fy++ {
			if !lm.Ground[fx][fy].HasFire {
				continue
			}
			x0 := fx - fireRadius
			if x0 < 0 {
				x0 = 0
			}
			x1 := fx + fireRadius
			if x1 >= LocalMapW {
				x1 = LocalMapW - 1
			}
			y0 := fy - fireRadius
			if y0 < 0 {
				y0 = 0
			}
			y1 := fy + fireRadius
			if y1 >= LocalMapH {
				y1 = LocalMapH - 1
			}
			for x := x0; x <= x1; x++ {
				for y := y0; y <= y1; y++ {
					dist := abs(x-fx) + abs(y-fy)
					if dist > fireRadius {
						continue
					}
					// Linear falloff: 1.0 at dist=0, approaches 0 at dist=radius+1.
					intensity := 1.0 - float64(dist)/float64(fireRadius+1)
					if intensity > lm.LitMap[x][y] {
						lm.LitMap[x][y] = intensity
					}
				}
			}
		}
	}
}

// ── LocalMapFor accessor ────────────────────────────────────────────────────

// LocalMapFor returns the LocalMap for the given world-space tile, generating
// and caching it on first access.
func LocalMapFor(worldX, worldY int, m *Model) *LocalMap {
	key := WorldCoord{X: worldX, Y: worldY}
	if lm, ok := m.localCache[key]; ok {
		return lm
	}
	tile := TileAt(worldX, worldY, m)
	lm := GenerateLocalMap(worldX, worldY, m.globalSeed, tile.Biome)
	m.localCache[key] = lm
	return lm
}
