## Context

The game currently runs inside the alternate screen buffer (`tea.WithAltScreen()`) and renders everything as a single string from `buildView`. The inventory is a narrow 22-column panel appended to the right of the map. As a result:

- All UI regions compete for horizontal space — adding an equipment model or rich item details is impractical in the current layout.
- The game does not use mouse input at all; all interaction is keyboard-only.
- There is no concept of a "screen mode" separate from "game mode" (`ModeWorld/Local/Dungeon`). When the inventory should take over the screen, the current approach must awkwardly carve out the map view rather than replacing it.

This change upgrades the project from BubbleTea v1.3.10 to **BubbleTea v2** (`charm.land/bubbletea/v2`) and lipgloss v1 to **lipgloss v2** (`charm.land/lipgloss/v2`). The v2 upgrade is done first, before adding mouse support and fullscreen inventory, so those features are built on the final public API.

Key v2 API changes relevant to this codebase:
- **`View() tea.View`**: `View()` now returns a `tea.View` struct (not `string`). Terminal features like alt-screen and mouse mode are declared as fields on that struct, not as `NewProgram` options.
- **`tea.KeyPressMsg`**: replaces `tea.KeyMsg` struct for key press events. `msg.String()` continues to return `"up"`, `"ctrl+c"`, etc. — switch logic is unchanged; only the type name and test construction differ.
- **Split mouse types**: `tea.MouseClickMsg`, `tea.MouseWheelMsg`, `tea.MouseReleaseMsg`, `tea.MouseMotionMsg` replace the monolithic `tea.MouseMsg` struct. Coordinates accessed via `msg.Mouse().X / .Y`; button via `msg.Button`.
- **Import vanity domain**: `github.com/charmbracelet/bubbletea` → `charm.land/bubbletea/v2`; `github.com/charmbracelet/lipgloss` → `charm.land/lipgloss/v2`.

The proposal adds two related capabilities: **fullscreen inventory** and **mouse support**. Both share the `ScreenMode` concept and are implemented together.

## Goals / Non-Goals

**Goals:**
- Introduce `ScreenMode` (`ScreenNormal` / `ScreenInventory`) on `Model` to cleanly separate game rendering from overlay rendering.
- `renderFullscreenInventory` fills the entire viewport; left column = item list with cursor; right column = ASCII ragdoll body with placeholder equipment slots.
- `i` toggles `screenMode` instead of `showInventory` (migration: `showInventory` removed, replaced by `screenMode == ScreenInventory`).
- Mouse support enabled via `tea.WithMouseCellMotion()` in `main.go`.
- In `ScreenNormal`: left-click at `(x, y)` on the local/dungeon map attempts to pathfind or step the player one cell toward the clicked tile (single step, not full pathfinding).
- In `ScreenInventory`: left-click on an item row sets `inventoryCursor` to that row; scroll-wheel moves cursor up/down.
- Equipment slots in the ragdoll panel are rendered as `[ Empty ]` labels with names only — no equip/unequip logic in this change.

**Non-Goals:**
- Full pathfinding (A\*) for click-to-move — single cardinal step toward click is sufficient.
- Equipping items from inventory to ragdoll slots — reserved for a future change.
- Drag-and-drop between inventory slots — future enhancement.
- Mouse support on the world map (world map uses zoom levels making coordinate math complex) — deferred.
- BubbleTea v2 layer compositor / `OnMouse` region interception — available in v2 but not used in this change (sufficient to use the split message types directly).

## Decisions

### Decision 1: `ScreenMode` replaces `showInventory bool`

`showInventory bool` was a simple toggle but doesn't scale when there could be multiple fullscreen overlays (inventory, character sheet, map legend, etc.). `ScreenMode int` with `ScreenNormal = 0` and `ScreenInventory = 1` is zero-value safe (defaults to `ScreenNormal`), and extending it for future screens requires only a new constant and a new `case` in `buildView`.

**Alternative considered**: Keep `showInventory` and add a separate `showEquipment bool`. Rejected — leads to multiple conflicting bools and complex guard logic.

### Decision 2: Mouse mode declared as `view.MouseMode = tea.MouseModeCellMotion` in `View()`

In v2, terminal feature flags (alt-screen, mouse mode) moved from `tea.NewProgram(...)` options into fields on `tea.View`. `MouseModeCellMotion` reports clicks and scroll events only when the cursor moves between cells (not every pixel), which is appropriate for a character-cell roguelike. `MouseModeAllMotion` would flood the event loop with every mouse movement — unnecessary and CPU-costly.

**Alternative considered**: `MouseModeAllMotion` for hover states. Rejected — adds noise; no hover effect is planned in this change.

### Decision 3: Mouse events handled in a dedicated `handleMouse` function

The existing `handleKey` switch is already long. Adding a `tea.MouseMsg` type-switch branch directly into `Update()` (via a `case tea.MouseMsg:` branch alongside `case tea.KeyMsg:`) keeps routing clean. `handleMouse` is a separate function that mirrors `handleKey`'s signature `func handleMouse(msg tea.MouseMsg, m Model) (Model, tea.Cmd)`.

**Alternative considered**: Add `tea.MouseMsg` case inside `handleKey`. Rejected — `handleKey` receives `tea.KeyMsg` specifically; mixing message types there is semantically wrong.

### Decision 4: Click-to-move is a single cardinal step toward the clicked cell

Full pathfinding is a separate feature. A single step toward the click direction gives responsive feel with minimal complexity: compute `dx = sign(clickX - playerScreenX)`, `dy = sign(clickY - playerScreenY)`, prefer the larger axis, call `applyDelta`. This reuses existing collision logic.

**Alternative considered**: Smooth step-by-step animation toward click. Rejected — requires queuing multiple moves across ticks, which is a larger change.

### Decision 6: BubbleTea v2 upgrade is done as the first step of this change

Migrating to v2 before implementing mouse support and fullscreen inventory ensures the new features are written against the final API (split mouse types, `tea.View` return) rather than v1 and requiring a second migration. The migration scope is bounded: `msg.String()` key matching is unchanged; only the `case` type names in `Update()` and test key construction structs change. Lipgloss v2 has the same surface area used by this project (`lipgloss.NewStyle()`, `JoinHorizontal`, etc.).

**Alternative considered**: Implement ui-enhancements on v1, then migrate to v2 separately. Rejected — mouse event API differences (split types vs. monolithic struct) would require rewriting all the new mouse handling immediately after.

### Decision 5: Ragdoll layout is a static ASCII body outline

The right-hand column of the fullscreen inventory will render a fixed-width ASCII silhouette (`~O~` head, `|H|` chest, etc.) with named slot labels. Slots are read from a `var equipSlots = []string{"Head","Chest","Left Hand","Right Hand","Legs","Feet"}` slice and rendered as `Head   : [ Empty ]`. No `Equipment` type or equip logic is introduced yet — the slots are purely cosmetic/placeholder.

**Alternative considered**: A separate `Equipment` struct on `Model` right now. Rejected — over-engineering before a real equip system is designed; placeholder rendering costs nothing to replace.

## Risks / Trade-offs

[Mouse mode breaks in some terminal emulators] Some SSH environments or multiplexers (tmux) mishandle mouse escape codes. → No mitigation needed for day-to-day development; can add a `--no-mouse` flag later if reported.

[ScreenMode couples rendering and input] `buildView` and `handleKey` both switch on `screenMode`. A future third screen would need changes in both files. → Acceptable for the current scope; a dedicated screen registry can be introduced later when there are ≥3 screens.

[Click coordinate mapping differs per mode] In dungeon mode the map starts at column 0 (full width); with sidebar open it starts after `sidebarContentW`. Mapping a click's terminal column to a map cell requires knowing the current camera offset and sidebar state. → For this change, implement click-to-move only when neither sidebar nor map picker is open; otherwise treat click as a no-op.

[Removing `showInventory` is a breaking rename] Existing tests reference `m.showInventory`. → All callers updated as part of the migration; test assertions updated in the same PR.

## Migration Plan

0. **BubbleTea v2 migration**: update `go.mod`; update all import paths; change `View() string` → `View() tea.View` with `AltScreen = true`; change `case tea.KeyMsg:` → `case tea.KeyPressMsg:`; update `handleKey` signature; update test key construction structs. Run `go build ./...` and `go test ./...` to confirm clean baseline before continuing.
1. Add `ScreenMode` type and constants to `types.go`.
2. Replace `showInventory bool` with `screenMode ScreenMode` on `Model`; update `NewModel()`; update all references.
3. Set `view.MouseMode = tea.MouseModeCellMotion` in `View()` (part of step 0's `View()` change).
4. Implement `handleMouseClick` and `handleMouseWheel` in `input.go`; wire into `Update()` via `case tea.MouseClickMsg:` and `case tea.MouseWheelMsg:`.
5. Implement `renderFullscreenInventory` in `render.go`; update `buildView` to dispatch on `screenMode`.
6. Update all tests that referenced `showInventory → screenMode`.

Rollback: steps 1–6 are additive or straightforward renames; reverting any step is safe.

## Open Questions

- Should the ragdoll silhouette be a hard-coded ASCII art string, or generated from a `[]string` of body rows? (Leaning toward hard-coded for now; easy to replace.)
- Should clicking outside the inventory list (e.g., on the ragdoll column) be a no-op or close the inventory? (Leaning toward no-op; `i`/`esc` closes it.)
- Do we want a visual highlight/border on the selected equipment slot when navigating with arrow keys? (Deferred to the actual equip change.)
