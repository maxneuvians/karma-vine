## Context

The current bottom of the screen uses two rows: a HUD row (tile info, clock, items) and a key-bar row (raw key binding hints). With combat stats now in the model (`playerHP`, `playerMaxHP`, armour derived from equipment), there is player-state data worth surfacing. The key-bar row is verbose and wasted space for returning players. A prior sidebar implementation had a viewport-clipping bug where the panel could overflow `m.viewportH`; this design must avoid repeating that.

Current `buildView` terminal layout (3 rows reserved):
```
[map rows…]
[HUD row  ]   ← tile/clock/items
[key-bar  ]   ← raw bindings string
```

Target layout (2 rows reserved, 1 row gained for map):
```
[map rows…]
[stats HUD]   ← HP bar, armour, tile, clock, speed, "? help"
```
Help panel is a fullscreen overlay, not a persistent row.

## Goals / Non-Goals

**Goals:**
- HP rendered as a proportional fill bar capped to available terminal width
- Armour shown as a compact badge next to HP
- `[PAUSED]` indicator preserved in HUD
- `?` toggles a fullscreen help panel listing contextual key bindings for the current mode
- Help panel content is mode-specific (world / local / dungeon / inventory / combat)
- Help panel height is strictly clamped to `m.viewportH`; width clamped to `m.viewportW`
- Net: one fewer row consumed by chrome, one more row for the map

**Non-Goals:**
- Animated HP bar transitions
- Multiple simultaneous overlays (help + inventory, etc.)
- Persistent HP loss between combats (out of scope for this change)
- Reworking the sidebar (separate capability)

## Decisions

### 1. HP and MaxHP live on the Model

**Decision:** Add `playerHP int` and `playerMaxHP int` to `Model`. `NewModel()` initialises both to 20 (matching the combat base value). `buildPlayerCombatant` seeds `HP`/`MaxHP` from `m.playerHP`/`m.playerMaxHP` instead of hard-coding 20. On combat victory the HP field is updated to the surviving player combatant's HP (carry-over within the session; still resets on `NewModel`).

**Rationale:** The HUD needs to read HP at render time. Keeping it on the model is the canonical pattern (all other display state lives there). It also makes combat stats accurate between fights within a session.

**Alternative considered:** Derive HP from equipment sum each frame. No way to track damage taken without a mutable field.

### 2. Progress bar is a pure string function

**Decision:** Implement `renderProgressBar(current, max, width int, fillColor, emptyColor string) string` that returns a lipgloss-styled string of `█` (filled) and `░` (empty) runes scaled to `width`. Width is passed in, not hard-coded, so callers can adapt to viewport size.

**Rationale:** Pure function → trivially testable. Width parameter prevents overflow on narrow terminals. Caller (HUD renderer) computes available width after rendering all other fixed-width segments.

**Alternative considered:** Unicode block characters with partial fill (▏▎▍▌▋▊▉█). More precise but much harder to test and renders poorly in many terminal fonts. Full-block approach is robust.

### 3. Key-bar row removed; one row returned to map

**Decision:** `buildView` stops calling `renderKeyBar`. Map height becomes `m.viewportH - 1` (HUD only). The HUD row ends with `  ? help` as the sole binding hint.

**Rationale:** Removing a row of chrome gives the map more space and declutters the screen. The `? help` hint is discoverable enough for new players.

**Risk:** Players who relied on reading the key-bar in-game lose that affordance until they open help. Mitigated by the help panel being one key away and the HUD hint being always visible.

### 4. Help panel is a fullscreen overlay via early return in buildView

**Decision:** `m.showHelpPanel bool` — when true, `buildView` returns `renderHelpPanel(m)` immediately (same pattern as `ScreenInventory` and `ScreenCombat`). `renderHelpPanel` builds a `[]string` of lines and hard-clamps: `lines = lines[:min(len(lines), m.viewportH)]` before joining. Each line is lipgloss-width-constrained to `m.viewportW`.

**Rationale:** Fullscreen overlay avoids the compositor complexity of overlaying panels. The clamp-before-join pattern is the fix for the previous viewport overflow bug — the inventory and combat screens do the same and don't clip.

**Alternative considered:** Help as a sidebar column (like the existing sidebar). Rejected — sidebar is fixed-width and the key binding list needs variable width.

### 5. `?` key toggles showHelpPanel; sidebar `?` binding is replaced

**Decision:** `?` currently toggles `m.showSidebar`. This binding is reassigned to `showHelpPanel`. The sidebar is toggled by a new key (tentatively `\` or left to the sidebar spec). `showHelpPanel` is suppressed (forced false) inside `ScreenCombat` and `ScreenInventory` — `?` is a no-op in those modes.

**Rationale:** `?` is the universal convention for help. The sidebar is an exploration tool with its own niche; it can get a different binding without user-experience loss.

## Risks / Trade-offs

- **[Viewport resize mid-session]** HP bar width is computed from `m.viewportW` at render time → auto-adapts. No risk.
- **[HP carry-over between combats]** Players arrive at next fight with reduced HP. Adds tension but no healing mechanic yet. → Accepted; healing can be added later.
- **[Existing tests check renderKeyBar output]** Tests asserting on key-bar content will need updating. → Handled in task list.
- **[Sidebar key change is breaking for muscle memory]** → Documented in proposal as display-only breaking change.

## Open Questions

- What key replaces `?` for sidebar toggle? (Tentative: `\`; can be decided at implementation time.)
- Should the help panel be scrollable for future bindings growth? (No scroll in v1; clamp is sufficient.)
