package game

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// stripANSI removes ANSI escape sequences from a string for measuring visible width.
func stripANSI(s string) string {
	var b strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEsc = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func TestPlayerPortrait_NonEmpty(t *testing.T) {
	if playerPortrait == "" {
		t.Fatal("playerPortrait should not be empty")
	}
}

func TestPlayerPortrait_HasLines(t *testing.T) {
	lines := strings.Split(playerPortrait, "\n")
	if len(lines) < 10 {
		t.Errorf("playerPortrait should have at least 10 rows, got %d", len(lines))
	}
}

func TestHumanoidPortrait_NonEmpty(t *testing.T) {
	if humanoidPortrait == "" {
		t.Fatal("humanoidPortrait should not be empty")
	}
}

func TestBeastPortrait_NonEmpty(t *testing.T) {
	if beastPortrait == "" {
		t.Fatal("beastPortrait should not be empty")
	}
}

func TestUndeadPortrait_NonEmpty(t *testing.T) {
	if undeadPortrait == "" {
		t.Fatal("undeadPortrait should not be empty")
	}
}

func TestFallbackPortrait_NonEmpty(t *testing.T) {
	if fallbackPortrait == "" {
		t.Fatal("fallbackPortrait should not be empty")
	}
}

func TestEnemyPortrait_Humanoid(t *testing.T) {
	p := enemyPortrait('G')
	if p == "" {
		t.Fatal("enemyPortrait('G') should not be empty")
	}
}

func TestEnemyPortrait_Beast(t *testing.T) {
	p := enemyPortrait('W')
	if p == "" {
		t.Fatal("enemyPortrait('W') should not be empty")
	}
}

func TestEnemyPortrait_Undead(t *testing.T) {
	p := enemyPortrait('Z')
	if p == "" {
		t.Fatal("enemyPortrait('Z') should not be empty")
	}
}

func TestEnemyPortrait_Fallback(t *testing.T) {
	p := enemyPortrait('X')
	if p == "" {
		t.Fatal("enemyPortrait('X') should not be empty")
	}
}

func TestRenderPortrait_FullWidth(t *testing.T) {
	out := renderPortrait(playerPortrait, 40)
	if out == "" {
		t.Fatal("renderPortrait should return a non-empty string")
	}
	lines := strings.Split(out, "\n")
	for i, line := range lines {
		visible := stripANSI(line)
		if utf8.RuneCountInString(visible) > 40 {
			t.Errorf("row %d: expected at most 40 visible chars, got %d", i, utf8.RuneCountInString(visible))
		}
	}
}

func TestRenderPortrait_ContainsBlockChars(t *testing.T) {
	out := renderPortrait(playerPortrait, 40)
	hasBlock := false
	for _, r := range out {
		if r >= 0x2580 {
			hasBlock = true
			break
		}
	}
	if !hasBlock {
		t.Error("player portrait should contain at least one unicode block character (≥ U+2580)")
	}
}

func TestRenderPortrait_ClipToNarrowWidth(t *testing.T) {
	out := renderPortrait(playerPortrait, 20)
	lines := strings.Split(out, "\n")
	for i, line := range lines {
		visible := stripANSI(line)
		w := utf8.RuneCountInString(visible)
		if w > 20 {
			t.Errorf("row %d: expected at most 20 visible chars, got %d", i, w)
		}
	}
}

func TestRenderPortrait_ZeroWidth(t *testing.T) {
	out := renderPortrait(playerPortrait, 0)
	if out != "" {
		t.Errorf("expected empty string for zero width, got %q", out)
	}
}

func TestRenderPortrait_EmptyPortrait(t *testing.T) {
	out := renderPortrait("", 40)
	if out != "" {
		t.Errorf("expected empty string for empty portrait, got %q", out)
	}
}

// ── enemyPortraitByName tests ─────────────────────────────────────────────────

func TestEnemyPortraitByName_Humanoid(t *testing.T) {
	p := enemyPortraitByName("Goblin")
	if p != humanoidPortrait {
		t.Fatal("enemyPortraitByName('Goblin') should return humanoidPortrait")
	}
}

func TestEnemyPortraitByName_Beast(t *testing.T) {
	p := enemyPortraitByName("Cave Crustacean")
	if p != beastPortrait {
		t.Fatal("enemyPortraitByName('Cave Crustacean') should return beastPortrait")
	}
}

func TestEnemyPortraitByName_Undead(t *testing.T) {
	p := enemyPortraitByName("Sand Wraith")
	if p != undeadPortrait {
		t.Fatal("enemyPortraitByName('Sand Wraith') should return undeadPortrait")
	}
}

func TestEnemyPortraitByName_Unknown(t *testing.T) {
	p := enemyPortraitByName("Unknown Monster")
	if p != fallbackPortrait {
		t.Fatal("enemyPortraitByName('Unknown Monster') should return fallbackPortrait")
	}
}
