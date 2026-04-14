package game

// TickMsg is dispatched by the BubbleTea tick loop every 500 ms.
type TickMsg struct{}

// Biome identifies one of the fourteen biome types.
type Biome int

const (
	DeepOcean   Biome = iota // e < 0.28
	ShallowWater              // e < 0.36
	Beach                     // e < 0.40
	Forest                    // temperate, e < 0.50, m > 0.55
	Plains                    // temperate, e < 0.50
	DenseForest               // temperate, e < 0.62, m > 0.45
	Desert                    // e < 0.62, m < 0.35
	Mountains                 // e < 0.78
	Snow                      // e >= 0.78
	Jungle                    // hot, e < 0.50, m > 0.55 or e < 0.62, m > 0.45
	Savanna                   // hot, e < 0.50, m > 0.30
	AridSteppe                // hot, e < 0.50, m <= 0.30
	Tundra                    // cold, e < 0.50, m <= 0.50
	Taiga                     // cold, e < 0.50, m > 0.50 or e < 0.62
)

// Tile is a single cell on the world map.
type Tile struct {
	Biome       Biome
	Char        rune
	Color       string
	Elevation   float64
	Moisture    float64
	Temperature float64
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
	Char       rune
	Color      string
	Blocking   bool
	Lit        bool   // true when the object emits light (torches/braziers in dungeons)
	Name       string // human-readable label (e.g. "Tree", "Torch")
	Pickupable bool   // true when the player can pick up this object
}

// Item is a carriable entity in the player's inventory.
type Item struct {
	Char  rune
	Color string
	Name  string
	Count int
	Slots []BodySlot // which body slots this item can occupy; empty = not equippable
}

// Inventory holds the player's carried items.
type Inventory struct {
	Items    []Item
	Equipped [NumBodySlots]Item
}

// InventoryMaxSlots is the maximum number of distinct item stacks the player can carry.
const InventoryMaxSlots = 8

// Animal is a creature on the local map.
type Animal struct {
	X, Y  int
	Char  rune
	Color string
	Flee  bool
	Name  string // human-readable label (e.g. "Deer", "Wolf")
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
	ModeDungeon
)

// CellKind identifies the type of a dungeon cell.
type CellKind int

const (
	CellWall  CellKind = iota
	CellFloor
)

// DungeonW and DungeonH are the dimensions of a dungeon level.
const (
	DungeonW = 80
	DungeonH = 24
)

// DungeonCell is a single cell in a dungeon level grid.
type DungeonCell struct {
	Kind   CellKind
	Object *Object
}

// DungeonLevel is one floor of a dungeon.
type DungeonLevel struct {
	Cells        [DungeonW][DungeonH]DungeonCell
	UpStair      LocalCoord
	DownStair    LocalCoord
	HasDownStair bool
}

// DungeonMeta stores per-entrance dungeon metadata.
type DungeonMeta struct {
	MaxDepth int
}

// dungeonKey is a cache key for a specific dungeon level.
type dungeonKey struct {
	wx, wy, depth int
}

// MapMode describes the active world-map overlay.
type MapMode int

const (
	MapModeDefault     MapMode = iota // standard biome view
	MapModeTemperature                // temperature gradient overlay
	MapModeElevation                  // elevation gradient overlay
	MapModePolitical                  // contour line overlay
)

// ScreenMode controls which full-screen view is active.
type ScreenMode int

const (
	ScreenNormal    ScreenMode = iota // normal map/HUD view
	ScreenInventory                   // fullscreen inventory overlay
)

// BodySlot identifies one of the six wearable slots on the player's body.
type BodySlot int

const (
	SlotHead      BodySlot = iota
	SlotChest
	SlotLeftHand
	SlotRightHand
	SlotLegs
	SlotFeet
)

// NumBodySlots is the total number of equipment slots.
const NumBodySlots = 6
