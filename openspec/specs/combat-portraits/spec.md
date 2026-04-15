## ADDED Requirements

### Requirement: Portrait cell type represents a single pixel
The system SHALL define a `portraitCell` struct with a `rune` field and a `lipgloss.Color` field. A portrait SHALL be represented as `[][]portraitCell` with exactly 20 rows and 40 columns.

#### Scenario: Portrait grid has correct dimensions
- **WHEN** `playerPortrait` or any archetype portrait is accessed
- **THEN** it has exactly 20 rows and each row has exactly 40 cells

### Requirement: Player portrait is a fixed humanoid silhouette
The system SHALL define a package-level `playerPortrait [][]portraitCell` that depicts a recognisable heroic humanoid figure using unicode block characters (`█`, `▓`, `▒`, `░`, `▄`, `▀`, `▌`, `▐`, and similar) and appropriate foreground colours (e.g. skin tones, armour greys).

#### Scenario: Player portrait contains block characters
- **WHEN** `renderPortrait(playerPortrait, 40)` is called
- **THEN** the output contains at least one unicode block character (rune ≥ U+2580)

#### Scenario: Player portrait renders in 20 lines
- **WHEN** `renderPortrait(playerPortrait, 40)` is called
- **THEN** the result contains exactly 19 newline characters (20 rows)

### Requirement: Enemy portrait is selected by archetype character
The system SHALL implement `enemyPortrait(char rune) [][]portraitCell` that returns a portrait matching the enemy archetype:
- Humanoid archetypes (e.g. `'H'`, `'K'`, `'G'`, `'T'`): warrior/giant humanoid silhouette
- Beast archetypes (e.g. `'W'`, `'B'`, `'S'`): four-legged creature silhouette
- Undead archetypes (e.g. `'Z'`, `'V'`, `'L'`): skeletal/spectral silhouette
- Any unrecognised char: generic creature fallback portrait

#### Scenario: Known humanoid char returns humanoid portrait
- **WHEN** `enemyPortrait('G')` is called
- **THEN** the returned grid is non-nil and has 20 rows of 40 cells

#### Scenario: Unknown char returns fallback portrait
- **WHEN** `enemyPortrait('X')` is called with an unrecognised rune
- **THEN** the returned grid is non-nil and has 20 rows of 40 cells

### Requirement: renderPortrait renders a portrait as a styled string
The system SHALL implement `renderPortrait(p [][]portraitCell, panelWidth int) string` that:
- Renders each cell as its rune with the cell's foreground colour applied via lipgloss
- Joins cells within a row directly (no separator)
- Joins rows with newline characters
- If `panelWidth < 40`, clips each row to `panelWidth` columns (rightmost cells dropped)

#### Scenario: renderPortrait clips to panelWidth when narrower than 40
- **WHEN** `renderPortrait(playerPortrait, 20)` is called
- **THEN** each rendered row is at most 20 visible characters wide (ignoring ANSI escape sequences)

#### Scenario: renderPortrait full width produces 40-column rows
- **WHEN** `renderPortrait(playerPortrait, 40)` is called
- **THEN** each row (stripped of ANSI escapes) is exactly 40 characters wide
