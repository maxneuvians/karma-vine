package game

import (
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

	// Player
	playerPos LocalCoord

	// UI
	viewportW   int
	viewportH   int
	mode        Mode
	showSidebar bool

	// Time
	timeOfDay float64 // [0, 1): 0=midnight, 0.25=6AM, 0.5=noon, 0.75=6PM
	timeScale int    // discrete: 1, 2, 5, 10
}

// NewModel returns an initialised Model with non-nil maps.
func NewModel() Model {
	return Model{
		chunks:     make(map[ChunkCoord]*Chunk),
		localCache: make(map[WorldCoord]*LocalMap),
		timeOfDay:  0.25, // start at 6 AM
		timeScale:  1,
	}
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
	case tea.KeyMsg:
		return handleKey(msg, m)
	case TickMsg:
		// Advance time: at 10× speed one full day takes 30 s (60 ticks).
		// Base rate (1×) is 600 ticks = 5 minutes per cycle.
		delta := float64(m.timeScale) / 600.0
		m.timeOfDay = math.Mod(m.timeOfDay+delta, 1.0)
		if m.mode == ModeLocal && m.localMap != nil {
			moveAnimals(&m)
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m Model) View() string {
	if m.viewportW == 0 || m.viewportH == 0 {
		return "World Explorer — loading..."
	}
	return buildView(m)
}
