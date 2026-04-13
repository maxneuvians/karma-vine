package game

import (
	"math"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// --- dimFactor ---

func TestDimFactor_Noon(t *testing.T) {
	v := dimFactor(0.5)
	if math.Abs(v-1.0) > 0.01 {
		t.Fatalf("dimFactor(0.5) = %v, want ~1.0", v)
	}
}

func TestDimFactor_Midnight(t *testing.T) {
	v := dimFactor(0.0)
	if math.Abs(v-0.15) > 0.01 {
		t.Fatalf("dimFactor(0.0) = %v, want ~0.15", v)
	}
}

// --- applyColor ---

func TestApplyColor_FullBrightness(t *testing.T) {
	// dim=1.0 should return the original color unchanged
	got := applyColor("#2d7a1f", 1.0)
	if got != "#2d7a1f" {
		t.Fatalf("applyColor dim=1.0: got %q, want %q", got, "#2d7a1f")
	}
}

func TestApplyColor_HalfDim(t *testing.T) {
	got := applyColor("#ff0000", 0.5)
	// #ff = 255, 255*0.5 = 127 = #7f
	if got != "#7f0000" {
		t.Fatalf("applyColor dim=0.5: got %q, want #7f0000", got)
	}
}

func TestApplyColor_InvalidPassthrough(t *testing.T) {
	cases := []string{"", "abc", "#gggggg", "#12345"}
	for _, c := range cases {
		if got := applyColor(c, 0.5); got != c {
			t.Fatalf("applyColor(%q) = %q, want unchanged %q", c, got, c)
		}
	}
}

// --- formatTime ---

func TestFormatTime(t *testing.T) {
	cases := []struct {
		tod  float64
		want string
	}{
		{0.0, "00:00"},
		{0.5, "12:00"},
		{0.75, "18:00"},
		{0.25, "06:00"},
	}
	for _, c := range cases {
		if got := formatTime(c.tod); got != c.want {
			t.Fatalf("formatTime(%v) = %q, want %q", c.tod, got, c.want)
		}
	}
}

// --- WindowSizeMsg ---

func TestUpdate_WindowSizeMsg(t *testing.T) {
	m := NewModel()
	next, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	nm := next.(Model)
	if nm.viewportW != 120 || nm.viewportH != 40 {
		t.Fatalf("WindowSizeMsg: got viewport %dx%d, want 120x40", nm.viewportW, nm.viewportH)
	}
}

// --- timeScale keys ---

func TestUpdate_TimeScaleIncrease(t *testing.T) {
	m := NewModel() // timeScale starts at 1
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("]")})
	nm := next.(Model)
	if nm.timeScale != 2 {
		t.Fatalf("timeScale after ] = %d, want 2", nm.timeScale)
	}
}

func TestUpdate_TimeScaleDecrease(t *testing.T) {
	m := NewModel()
	m.timeScale = 5
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("[")})
	nm := next.(Model)
	if nm.timeScale != 2 {
		t.Fatalf("timeScale after [ = %d, want 2", nm.timeScale)
	}
}

func TestUpdate_TimeScaleClampedMax(t *testing.T) {
	m := NewModel()
	m.timeScale = 10
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("]")})
	nm := next.(Model)
	if nm.timeScale != 10 {
		t.Fatalf("timeScale clamped max: got %d, want 10", nm.timeScale)
	}
}

