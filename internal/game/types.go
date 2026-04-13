package game

// TickMsg is dispatched by the BubbleTea tick loop every 500 ms.
type TickMsg struct{}

// Biome identifies one of the nine biome types.
type Biome int

const (
	DeepOcean   Biome = iota // e < 0.28
	ShallowWater              // e < 0.36
	Beach                     // e < 0.40
	Forest                    // e < 0.50, m > 0.55
	Plains                    // e < 0.50, m <= 0.55
	DenseForest               // e < 0.62, m > 0.45
	Desert                    // e < 0.62, m < 0.35
	Mountains                 // e < 0.78
	Snow                      // e >= 0.78
)

// Tile is a single cell on the world map.
type Tile struct {
	Biome     Biome
	Char      rune
	Color     string
	Elevation float64
	Moisture  float64
}

// Ground is the floor tile in a local map cell.
type Ground struct {
	Char     rune
	Color    string
	Passable bool
	HasFire  bool
}

// Object is a world object occupying a local map cell (tree, rock, cactus, …).
type Object struct {
	Char     rune
	Color    string
	Blocking bool
}

// Animal is a creature on the local map.
type Animal struct {
	X, Y  int
	Char  rune
	Color string
	Flee  bool
}

// LocalMapW and LocalMapH are the dimensions of the local exploration grid.
const (
	LocalMapW = 160
	LocalMapH = 48
)

// LocalMap is the detailed 160×48 local view of a single world-map tile.
type LocalMap struct {
	Ground  [LocalMapW][LocalMapH]Ground
	Objects [LocalMapW][LocalMapH]*Object
	Animals []*Animal
	LitMap  [LocalMapW][LocalMapH]float64
}

// WorldCoord is a position on the infinite world map.
type WorldCoord struct {
	X, Y int
}

// LocalCoord is a position within a local map (42×18 grid).
type LocalCoord struct {
	X, Y int
}

// ChunkCoord identifies a 32×32 chunk of the world map.
type ChunkCoord struct {
	X, Y int
}

// Mode describes which map tier is active.
type Mode int

const (
	ModeWorld Mode = iota
	ModeLocal
)
