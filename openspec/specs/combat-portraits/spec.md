## ADDED Requirements

### Requirement: Portraits are loaded from embedded ANSI files at compile time
The system SHALL embed all files from `internal/game/portraits/*.ansi` into the binary using Go's `//go:embed` directive. Each `.ansi` file contains raw ANSI escape sequences and unicode block characters, one line per row, with lines separated by `\n`.

At startup, `init()` calls `loadPortrait(name, fallback)` for each archetype. If a matching file exists in the embedded FS it is used; otherwise the code-generated fallback string (built by the corresponding `buildXxxPortrait()` helper) is used.

A `portrait` type is defined as `string` (a pre-rendered ANSI string). Portrait variables (`playerPortrait`, `humanoidPortrait`, `beastPortrait`, `undeadPortrait`, `fallbackPortrait`) are package-level `portrait` (string) values.

#### Scenario: Player portrait loads from embedded file when present
- **WHEN** `portraits/player.ansi` exists in the embedded FS
- **THEN** `playerPortrait` is non-empty and contains at least one unicode block character (rune ≥ U+2580)

#### Scenario: Portrait falls back to code-generated version when file absent
- **WHEN** no `.ansi` file exists for an archetype (e.g. `humanoid`, `beast`)
- **THEN** the portrait variable is still non-empty, built from the code-defined `buildXxxPortrait()` grid

### Requirement: Portrait cell type is used only for code-defined fallback grids
The system SHALL retain the `portraitCell` struct (with `r rune` and `color string` fields) and `buildXxxPortrait() [][portraitCols]portraitCell` helpers solely to generate fallback ANSI strings when no `.ansi` file is present. These are not the primary representation; `portrait` (string) is.

### Requirement: Player portrait file depicts a recognisable heroic humanoid
`portraits/player.ansi` SHALL depict a recognisable heroic humanoid figure using unicode block characters (`█`, `▓`, `▒`, `░`, `▄`, `▀`, `▌`, `▐`, and similar) and ANSI 24-bit colour sequences.

#### Scenario: Player portrait contains block characters
- **WHEN** `renderPortrait(playerPortrait, 40)` is called
- **THEN** the output contains at least one unicode block character (rune ≥ U+2580)

### Requirement: Enemy portrait is selected by archetype character
The system SHALL implement `enemyPortrait(char rune) portrait` that returns a portrait string matching the enemy archetype:
- Humanoid archetypes (e.g. `'H'`, `'K'`, `'G'`, `'T'`): warrior/giant humanoid silhouette
- Beast archetypes (e.g. `'W'`, `'B'`, `'S'`): four-legged creature silhouette
- Undead archetypes (e.g. `'Z'`, `'V'`, `'L'`): skeletal/spectral silhouette
- Any unrecognised char: generic creature fallback portrait

#### Scenario: Known humanoid char returns humanoid portrait
- **WHEN** `enemyPortrait('G')` is called
- **THEN** the returned portrait string is non-empty

#### Scenario: Unknown char returns fallback portrait
- **WHEN** `enemyPortrait('X')` is called with an unrecognised rune
- **THEN** the returned portrait string is non-empty

### Requirement: renderPortrait renders a portrait string with line clipping
The system SHALL implement `renderPortrait(p portrait, panelWidth int) string` that:
- Splits the portrait string on `\n` to obtain rows
- Clips each row to at most `panelWidth` visible (non-ANSI-escape) characters via `clipANSILine`
- Joins rows with newline characters
- Returns an empty string when `panelWidth <= 0` or the portrait is empty

`clipANSILine(line string, maxWidth int) string` passes ANSI escape sequences through intact and stops counting after `maxWidth` visible runes.

#### Scenario: renderPortrait clips to panelWidth when narrower than portrait
- **WHEN** `renderPortrait(playerPortrait, 20)` is called
- **THEN** each rendered row is at most 20 visible characters wide (ignoring ANSI escape sequences)

#### Scenario: renderPortrait does not exceed portrait width at full width
- **WHEN** `renderPortrait(playerPortrait, 40)` is called
- **THEN** each row stripped of ANSI escapes is at most 40 characters wide

#### Scenario: renderPortrait returns empty for zero width
- **WHEN** `renderPortrait(playerPortrait, 0)` is called
- **THEN** the result is an empty string

#### Scenario: renderPortrait returns empty for empty portrait
- **WHEN** `renderPortrait("", 40)` is called
- **THEN** the result is an empty string

### Requirement: Portrait height is clamped to leave room for stats
When `renderHeroPanel` or `renderEnemyPanel` renders a portrait taller than `height - statsH` rows, the portrait SHALL be clipped to `height - statsH` rows so that player/enemy stats are always visible.

