## Why

The current combat screen uses a simple ASCII ragdoll and a single glyph character to represent combatants, which lacks visual engagement and player identity. Replacing these with rich 40×20 unicode block-character portraits will make combat feel more dramatic and immersive, and starting combat paused gives the player a moment to assess the encounter before it begins.

## What Changes

- Replace the existing ASCII ragdoll hero portrait and single-glyph enemy representation with procedurally generated 40×20 unicode block-character portraits for both the player and enemy
- Portraits use unicode block characters (`█`, `▓`, `▒`, `░`, `▄`, `▀`, etc.) with lipgloss foreground colours to create shaded, realistic-looking pixel-art style depictions
- Player portrait depicts a generic heroic humanoid figure; enemy portrait is derived from enemy type/template to create a distinct silhouette
- Combat begins in a **paused** state: the log is empty, portraits are visible, and a "Press [Space] to begin" prompt is shown
- Pressing Space (or Enter) unpauses combat and starts the playback tick sequence
- The speed controls (`[`, `]`) and existing playback remain unchanged once combat is unpaused

## Capabilities

### New Capabilities
- `combat-portraits`: 40×20 unicode block-character portrait rendering for both player and enemy combatants, displayed in the top-left and top-right panels of the combat screen

### Modified Capabilities
- `combat-rendering`: Portrait panels now render 40×20 block-character art instead of the ASCII ragdoll and enemy glyph; layout and stat sections below the portrait are unchanged
- `combat-playback`: Combat now starts in a paused state (`combatLogIndex == 0`, no tick scheduled); playback only begins after the player presses Space/Enter, which schedules the first `CombatTickMsg`

## Impact

- `combat_render.go` (or equivalent): portrait generation functions added; hero and enemy portrait rendering replaced
- `update.go` / input handling: new `combatPaused bool` model field; Space/Enter key in `ScreenCombat` triggers unpause and first tick
- `types.go`: `combatPaused bool` field added to `Model`
- No external dependencies added; portraits are generated procedurally using existing lipgloss/bubbletea primitives
- Tests for portrait rendering (snapshot or character-coverage checks) and pause/unpause flow required to maintain ≥90% coverage
