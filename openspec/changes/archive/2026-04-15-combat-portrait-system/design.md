## Context

The combat screen (`renderCombatScreen`) currently renders a small ASCII ragdoll (`~O~` etc.) for the hero and a single glyph character centred in a box for the enemy. The playback system (`CombatTickMsg`) auto-starts immediately when `ScreenCombat` is entered: `model.go` schedules the first tick in the same `Update` call that sets `screenMode = ScreenCombat` (`combatLogIndex == 0` branch).

The two changes proposed are:
1. **Portrait upgrade**: replace both panel art sections with 40×20 pixel-grid portraits composed of unicode block characters and lipgloss colour styles.
2. **Pause-on-start**: suppress the auto-tick and show a "Press Space to begin" prompt; only fire the first `CombatTickMsg` when the player presses Space or Enter.

Relevant files: `internal/game/render.go` (hero panel: `renderHeroPanel`, enemy panel inline in `renderCombatScreen`), `internal/game/model.go` (`Update` combat tick scheduling), `internal/game/types.go` (constants and model struct), `internal/game/input.go` (key handling).

## Goals / Non-Goals

**Goals:**
- 40×20 unicode block-character portraits for both player and enemy in their respective combat panels
- Player portrait: a consistent, recognisable heroic humanoid silhouette
- Enemy portrait: a distinct silhouette per enemy archetype (derived from `Template.Char` or `Template.Name`)
- Combat starts paused; a visible prompt ("Press [Space] to begin") replaces the speed hint until unpaused
- Pressing Space or Enter during the paused state unpauses and fires the first tick
- ≥90% test coverage maintained

**Non-Goals:**
- Animated portraits or per-frame transitions
- Unique hand-crafted art for every individual enemy; a small set of archetypes is sufficient
- Changes to combat logic, damage calculation, or playback speed controls
- Multiplayer or networked state

## Decisions

### 1. Portrait representation: hardcoded pixel maps vs. procedural generation

**Decision**: Store each portrait as a `[][]portraitCell` (a 20-row × 40-col grid) where `portraitCell` is a struct with a rune and a `lipgloss.Color`.  Portraits are defined as package-level `var` literals in a new file `combat_portraits.go`.

**Alternatives considered**:
- *Procedural generation at render time*: algorithmic noise + shape primitives. More flexible but complex to tune, hard to test, and output quality is unpredictable for a 40×20 grid.
- *Loading from embedded text assets*: clean separation but adds file I/O, encoding complexity, and a build-time dependency.

**Rationale**: Literal pixel maps are deterministic, testable (snapshot tests), and easy for contributors to tweak visually. The total memory footprint (20 × 40 × ~20 bytes per cell) is negligible.

### 2. Enemy portrait selection

**Decision**: Map `Template.Char` (the single-rune enemy identifier already present on `DungeonEnemy.Template` and `Animal`) to a portrait via a `switch` statement with a catch-all "generic creature" fallback.  A helper `enemyPortrait(char rune) [][]portraitCell` encapsulates the lookup.

**Alternatives considered**:
- *Map keyed by `Template.Name`*: more granular but fragile to name changes.
- *Single generic enemy portrait*: simplest but misses the differentiation goal.

**Rationale**: `Template.Char` is stable, already used for rendering in other contexts, and provides enough variety with ~5–8 archetype portraits (humanoid, beast, undead, dragon, etc.) plus a fallback.

### 3. Rendering portraits as strings

**Decision**: `renderPortrait(p [][]portraitCell, width int) string` iterates rows, renders each cell as `lipgloss.NewStyle().Foreground(cell.color).Render(string(cell.rune))`, joins with `\n`, and centres the 40-column block within `width` using lipgloss `Place`.

**Rationale**: Consistent with existing lipgloss usage in the codebase. The 40-col portrait width is narrower than the ~40% panel width at most viewports, so no clipping occurs at typical terminal widths.

### 4. Pause-on-start: new `combatPaused` model field

**Decision**: Add `combatPaused bool` to `Model`. When `ScreenCombat` is entered, set `combatPaused = true` and do **not** schedule a `CombatTickMsg` (remove the `combatLogIndex == 0` auto-tick branch). In `Update`, handle Space/Enter in `ScreenCombat` when `combatPaused == true`: set `combatPaused = false`, schedule the first `CombatTickMsg`, return.

**Alternatives considered**:
- *Reuse `combatLogIndex == -1` as paused sentinel*: avoids a new field but conflates index semantics with pause state; makes existing `hpAtRound` and log-reveal logic more complex.
- *New `CombatPausedMsg` message type*: more message-passing overhead for a simple boolean toggle.

**Rationale**: Explicit boolean is the clearest representation. One new field, minimal surface area, easy to test.

### 5. Paused UI prompt

**Decision**: In `renderCombatLog` (or a new `renderCombatPausedOverlay`), when `m.combatPaused == true`, replace the speed hint line with a centred `"[ Space ] Begin Combat"` prompt styled with the existing accent colour. The rest of the log panel stays empty. Once unpaused, the existing speed/log UI takes over.

**Rationale**: Reusing the log panel for the prompt avoids layout changes and keeps all combat UI within the existing three-panel structure.

## Risks / Trade-offs

- **Terminal width < 80 columns** → portrait 40 cols + padding may exceed the panel width and wrap. Mitigation: `renderPortrait` clips to `min(40, panelWidth)` columns; portrait cells beyond that are dropped.
- **Colour support varies** → in 256-colour or no-colour terminals, lipgloss degrades gracefully; the block characters remain but shading colours reduce to nearest or no colour. No mitigation needed beyond what lipgloss already does.
- **Test snapshot brittleness** → portrait pixel maps are large. Mitigation: use character-coverage assertions (e.g., "output contains `█`") rather than full string snapshots for unit tests; reserve snapshot tests for integration.
- **combatPaused persists if combat is interrupted** → if `ScreenCombat` is exited early (e.g., game quit), `combatPaused` is reset alongside `combatLogIndex` on next combat entry. Ensure the enter-combat path always resets both fields.

## Migration Plan

1. Add `combatPaused bool` to `Model` in `types.go`; update `NewModel` and the enter-combat path in `model.go` to set `combatPaused = true` and remove the auto-tick
2. Add Space/Enter key handling for the paused state in `input.go`
3. Create `combat_portraits.go` with player portrait, archetype enemy portraits, and `renderPortrait` / `enemyPortrait` helpers
4. Update `renderHeroPanel` to call `renderPortrait(playerPortrait, panelWidth)` instead of the ragdoll loop
5. Update the enemy panel section to call `renderPortrait(enemyPortrait(char), panelWidth)` instead of the glyph box
6. Update `renderCombatLog` to show the pause prompt when `m.combatPaused`
7. Write/update tests for pause flow and portrait rendering
8. No rollback complexity — all changes are confined to the `internal/game` package with no external API surface

## Open Questions

- Should the player portrait change based on equipped gear (e.g., wearing armour vs. not)? For now: no — one fixed player portrait. Can be extended later.
- Should pressing Space during active playback (unpaused) pause it again (like a play/pause toggle)? For now: no — Space only works in the initial paused state; existing speed controls handle pacing.
