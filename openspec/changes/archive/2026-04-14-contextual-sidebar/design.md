## Context

The sidebar currently uses a hardcoded `localCharNames map[rune]string` in `render.go` to translate object/animal glyphs to display names. This lookup fails silently — missing entries fall back to the raw glyph string — and there is no sidebar branch for `ModeDungeon` at all. The dungeon sidebar shows local-map data because the `else` branch in `renderSidebar` unconditionally reads `m.localMap`.

## Goals / Non-Goals

**Goals:**
- Every object and animal carries its own human-readable name at the point of creation
- Sidebar content is fully determined by `m.mode` — world, local, or dungeon — with no cross-mode data leakage
- World sidebar adapts when a non-default map mode is active (the biome colour legend is irrelevant during temperature/elevation/political overlays)
- Dungeon sidebar shows depth, the current level's contents (up/down stairs, torches, braziers), and a close hint
- Local sidebar is complete: all objects and animals show correct names including the dungeon entrance

**Non-Goals:**
- Showing per-cell inspection details on cursor hover (a future "look" command)
- Sidebar scroll for long content
- Showing animal behaviour state (flee/wander) in the sidebar

## Decisions

### Decision 1: Add `Name string` to `Object` and `Animal` — no lookup table

The lookup table approach (`localCharNames[char]`) is a maintenance trap: any new glyph that isn't added to the table silently falls back to the raw character. Adding `Name string` directly to `Object` and `Animal` means the sidebar can call `obj.Name` and fall back to an empty/unknown label only when the generator forgot to set it — which is easily caught by a unit test.

**Alternative considered**: Extend the lookup table with all known glyphs. Rejected — still fragile and requires keeping two locations in sync.

### Decision 2: Populate `Name` at generation sites

`Name` is set when the struct is created in `local.go` (biome objects, animals, dungeon entrance) and `dungeon.go` (torch, brazier, staircase up/down). The render code never needs to derive a name from a glyph.

### Decision 3: Three distinct sidebar sections keyed on `m.mode`

`renderSidebar` switches on `m.mode`:
- `ModeWorld` — biome legend (existing) **or** map-mode note when `m.mapMode != MapModeDefault`
- `ModeLocal` — "Legend" with player, objects by name, wildlife by name (using `obj.Name`, `a.Name`)
- `ModeDungeon` — "Dungeon" header with depth, then contents of the current level (staircases, torches, braziers)

### Decision 4: World map sidebar shows map-mode context

When `m.mapMode != MapModeDefault`, the biome legend is replaced by a short header naming the active overlay and a brief note (e.g., "Temperature", "blue = cold / red = hot"). This is more useful than a biome legend that doesn't match what's on screen.

### Decision 5: Dungeon sidebar scans `m.currentDungeon` directly

The dungeon sidebar iterates `currentDungeon.Cells` to collect unique named objects, then lists them. This mirrors the local-map sidebar which scans `lm.Objects`.

## Risks / Trade-offs

[Breaking change on Object/Animal] Adding `Name string` changes the struct. All tests that construct `Object{}` or `Animal{}` without the new field will still compile (zero value is `""`), but the sidebar would show a blank name for any test fixture that omits it. Mitigated by adding a test that asserts no empty names on a freshly generated local map and dungeon level.

[Dungeon sidebar scan cost] Scanning all 80×24 = 1920 cells on every render to build the sidebar list is negligible (< 1 µs). No caching needed.

## Migration Plan

1. Add `Name string` to `Object` and `Animal` in `types.go`.
2. Update all creation sites in `local.go` and `dungeon.go`.
3. Update `renderSidebar` in `render.go`: remove `localCharNames`, add dungeon branch, add map-mode world branch.
4. Fix any test fixtures that need the `Name` field.
5. Add new tests for sidebar completeness and name coverage.
