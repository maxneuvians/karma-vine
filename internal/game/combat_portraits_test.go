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

func TestPlayerPortrait_Dimensions(t *testing.T) {
	if len(playerPortrait) != portraitRows {
		t.Fatalf("playerPortrait should have %d rows, got %d", portraitRows, len(playerPortrait))
	}
	for i, row := range playerPortrait {
		if len(row) != portraitCols {
			t.Errorf("playerPortrait row %d should have %d cols, got %d", i, portraitCols, len(row))
		}
	}
}

func TestHumanoidPortrait_Dimensions(t *testing.T) {
	if len(humanoidPortrait) != portraitRows {
		t.Fatalf("humanoidPortrait should have %d rows, got %d", portraitRows, len(humanoidPortrait))
	}
	for i, row := range humanoidPortrait {
		if len(row) != portraitCols {
			t.Errorf("humanoidPortrait row %d should have %d cols, got %d", i, portraitCols, len(row))
		}
	}
}

func TestBeastPortrait_Dimensions(t *testing.T) {
	if len(beastPortrait) != portraitRows {
		t.Fatalf("beastPortrait should have %d rows, got %d", portraitRows, len(beastPortrait))
	}
	for i, row := range beastPortrait {
		if len(row) != portraitCols {
			t.Errorf("beastPortrait row %d should have %d cols, got %d", i, portraitCols, len(row))
		}
	}
}

func TestUndeadPortrait_Dimensions(t *testing.T) {
	if len(undeadPortrait) != portraitRows {
		t.Fatalf("undeadPortrait should have %d rows, got %d", portraitRows, len(undeadPortrait))
	}
	for i, row := range undeadPortrait {
		if len(row) != portraitCols {
			t.Errorf("undeadPortrait row %d should have %d cols, got %d", i, portraitCols, len(row))
		}
	}
}

func TestFallbackPortrait_Dimensions(t *testing.T) {
	if len(fallbackPortrait) != portraitRows {
		t.Fatalf("fallbackPortrait should have %d rows, got %d", portraitRows, len(fallbackPortrait))
	}
	for i, row := range fallbackPortrait {
		if len(row) != portraitCols {
			t.Errorf("fallbackPortrait row %d should have %d cols, got %d", i, portraitCols, len(row))
		}
	}
}

func TestEnemyPortrait_Humanoid(t *testing.T) {
	p := enemyPortrait('G')
	if p == nil {
		t.Fatal("enemyPortrait('G') should not be nil")
	}
	if len(p) != portraitRows {
		t.Fatalf("expected %d rows, got %d", portraitRows, len(p))
	}
	for _, row := range p {
		if len(row) != portraitCols {
			t.Errorf("expected %d cols, got %d", portraitCols, len(row))
		}
	}
}

func TestEnemyPortrait_Beast(t *testing.T) {
	p := enemyPortrait('W')
	if p == nil {
		t.Fatal("enemyPortrait('W') should not be nil")
	}
	if len(p) != portraitRows {
		t.Fatalf("expected %d rows, got %d", portraitRows, len(p))
	}
}

func TestEnemyPortrait_Undead(t *testing.T) {
	p := enemyPortrait('Z')
	if p == nil {
		t.Fatal("enemyPortrait('Z') should not be nil")
	}
	if len(p) != portraitRows {
		t.Fatalf("expected %d rows, got %d", portraitRows, len(p))
	}
}

func TestEnemyPortrait_Fallback(t *testing.T) {
	p := enemyPortrait('X')
	if p == nil {
		t.Fatal("enemyPortrait('X') should not be nil")
	}
	if len(p) != portraitRows {
		t.Fatalf("expected %d rows, got %d", portraitRows, len(p))
	}
}

func TestRenderPortrait_FullWidth(t *testing.T) {
	out := renderPortrait(playerPortrait, 40)
	lines := strings.Split(out, "\n")
	if len(lines) != portraitRows {
		t.Fatalf("expected %d lines, got %d", portraitRows, len(lines))
	}
	for i, line := range lines {
		visible := stripANSI(line)
		if utf8.RuneCountInString(visible) != portraitCols {
			t.Errorf("row %d: expected %d visible chars, got %d", i, portraitCols, utf8.RuneCountInString(visible))
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
	if len(lines) != portraitRows {
		t.Fatalf("expected %d lines, got %d", portraitRows, len(lines))
	}
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
