package game

import (
	"math"
	"math/rand"
	"time"

	tea "charm.land/bubbletea/v2"
)

// Model is the top-level BubbleTea model for the world explorer.
type Model struct {
	// World tier
	globalSeed int
	worldPos   WorldCoord
	chunks     map[ChunkCoord]*Chunk

	// Local tier
	localMap   *LocalMap
	localCache map[WorldCoord]*LocalMap

	// Dungeon tier
	dungeonCache    map[dungeonKey]*DungeonLevel
	dungeonMeta     map[WorldCoord]DungeonMeta
	currentDungeon  *DungeonLevel
	dungeonDepth    int
	dungeonEntryPos LocalCoord

	// Player
	playerPos LocalCoord

	// Inventory
	inventory       Inventory
	screenMode      ScreenMode
	inventoryCursor int
	equipFocused    bool
	equipCursor     int

	// Combat
	combatState        *CombatState
	combatEnemy        *Animal
	combatDungeonEnemy *DungeonEnemy

	// Player stats
	playerHP    int
	playerMaxHP int

	// UI
	viewportW       int
	viewportH       int
	mode            Mode
	showSidebar     bool
	showHelpPanel   bool
	worldZoom       int // 1=normal, 2=2×, 4=4×, 8=8×
	mapMode         MapMode
	showMapPicker   bool
	mapPickerCursor int

	// Time
	timeOfDay float64 // [0, 1): 0=midnight, 0.25=6AM, 0.5=noon, 0.75=6PM
	timeScale int    // discrete: 1, 2, 5, 10

	// Pause
	paused              bool
	pausedBeforeInventory bool // tracks pause state when inventory was opened
}

// defaultOutfit returns the starting equipment for a new character.
func defaultOutfit() [NumBodySlots]Item {
	var eq [NumBodySlots]Item
	eq[SlotChest] = Item{Char: '♦', Color: "#a0a0a0", Name: "Cloth Tunic", Count: 1, Slots: []BodySlot{SlotChest}}
	eq[SlotLegs] = Item{Char: '‖', Color: "#a0a0a0", Name: "Cloth Pants", Count: 1, Slots: []BodySlot{SlotLegs}}
	eq[SlotFeet] = Item{Char: '∩', Color: "#8B4513", Name: "Leather Boots", Count: 1, Slots: []BodySlot{SlotFeet}}
	return eq
}

// NewModel returns an initialised Model with non-nil maps.
func NewModel() Model {
	m := Model{
		globalSeed:   rand.New(rand.NewSource(time.Now().UnixNano())).Int(),
		chunks:       make(map[ChunkCoord]*Chunk),
		localCache:   make(map[WorldCoord]*LocalMap),
		dungeonCache: make(map[dungeonKey]*DungeonLevel),
		dungeonMeta:  make(map[WorldCoord]DungeonMeta),
		inventory:    Inventory{Items: []Item{}, Equipped: defaultOutfit()},
		playerHP:     20,
		playerMaxHP:  20,
		timeOfDay:    0.25, // start at 6 AM
		timeScale:    1,
		worldZoom:    1,
	}
	m.worldPos = findWorldSpawn(&m)
	return m
}

func tickCmd() tea.Cmd {
	return tea.Every(500*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewportW = msg.Width
		m.viewportH = msg.Height
	case tea.KeyPressMsg:
		return handleKey(msg, m)
	case tea.MouseClickMsg:
		return handleMouseClick(msg, m)
	case tea.MouseWheelMsg:
		return handleMouseWheel(msg, m)
	case TickMsg:
		if m.paused {
			return m, tickCmd()
		}
		// Advance time: at 10× speed one full day takes 30 s (60 ticks).
		// Base rate (1×) is 600 ticks = 5 minutes per cycle.
		delta := float64(m.timeScale) / 600.0
		m.timeOfDay = math.Mod(m.timeOfDay+delta, 1.0)
		if m.mode == ModeLocal && m.localMap != nil {
			moveAnimals(&m)
		}
		if m.mode == ModeDungeon {
			m = moveEnemies(m)
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m Model) View() tea.View {
	if m.viewportW == 0 || m.viewportH == 0 {
		return tea.NewView("World Explorer — loading...")
	}
	v := tea.NewView(buildView(m))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
