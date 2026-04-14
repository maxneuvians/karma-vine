## Context

`resolveCombat` already pre-computes the full fight and stores every round's narration in `CombatState.Log`. The model enters `ScreenCombat` and currently renders all log lines at once. The change doesn't touch the combat resolution logic — it only controls **how much of the log is revealed** at render time via a `combatLogIndex` cursor that advances on a timer.

The existing `TickMsg` (500 ms) drives world time and animal movement; combat playback needs its own timer at a different rate. The screen layout needs to fit in any terminal that can run the game (minimum ~60 × 18).

## Goals / Non-Goals

**Goals:**
- `combatLogIndex` starts at 0 when combat begins; advances by one "round block" per `CombatTickMsg`
- Three playback speeds: Slow = 3000 ms, Normal = 1000 ms, Fast = 200 ms
- Three-panel layout: left (hero), right (enemy), bottom (log); proportions adapt to viewport
- HP bars in each panel reflect the combatant's HP **at the current log index** (requires storing per-round HP snapshots or recomputing from log)
- Hero portrait: existing ragdoll art, centred in the left panel
- Enemy portrait: large glyph (`Template.Char`) centred in a box in the right panel
- Result banner + dismiss hint appear only after `combatLogIndex` passes the final round

**Non-Goals:**
- Player interaction during playback (no ability choices — still an auto-battler)
- Scroll back through the log
- Animated glyph effects or colours pulsing

## Decisions

### 1. CombatTickMsg is a separate message type with speed-derived interval

**Decision:** Define `type CombatTickMsg struct{}`. When `ScreenCombat` is entered, `Update` schedules a `CombatTickMsg` using `tea.Tick(combatSpeedDuration(m.combatSpeed), ...)`. Each `CombatTickMsg` handler: advances `combatLogIndex` by one round block (all log lines tagged to the current round), then reschedules another `CombatTickMsg` unless playback is complete. `combatSpeed` is an int constant (`CombatSpeedSlow=0`, `CombatSpeedNormal=1`, `CombatSpeedFast=2`) mapping to `{3000, 1000, 200}` ms.

**Rationale:** Decouples combat playback cadence from the world tick. The world `TickMsg` continues running (animals move, time passes) independently. Using `tea.Tick` on `CombatTickMsg` means the interval naturally adapts when speed changes — the next reschedule picks up the new duration.

**Alternative considered:** Using the existing `TickMsg` with a counter. Rejected — the world tick fires at 500 ms; 3 s would need a 6-tick counter, and changing speed mid-combat would drift.

### 2. Log is grouped into round blocks by a sentinel prefix

**Decision:** `resolveCombat` already writes lines like `"Round 1: ..."`. The playback renderer uses `CombatState.Round` (incremented during resolution) to know how many round blocks exist. `combatLogIndex` tracks which round block has been revealed (0 = nothing shown, 1 = round 1 shown, etc.). On each `CombatTickMsg`, `combatLogIndex++`. Lines belonging to round N are lines between the N-th and (N+1)-th `"Round "` prefix in the log slice.

**Rationale:** No structural change to `CombatState.Log`. The grouping is already implicit in the log text.

**Alternative considered:** Storing per-round sub-slices in `CombatState`. More explicit but requires changing the combat resolution return type — risky.

### 3. HP at current round derived by rescanning visible log lines

**Decision:** At render time, the renderer scans the visible log lines (up to `combatLogIndex` round block) for damage patterns (e.g. `"takes N damage"`) and recomputes HP deltas from `CombatState.Player.MaxHP` / `CombatState.Enemy.MaxHP`. This is display-only; the actual resolution result is already in `CombatState`.

**Rationale:** Avoids storing per-round HP snapshots in `CombatState` (would require changing the combat resolution contract). The log already contains all damage numbers — parsing is deterministic and cheap.

**Alternative considered:** Storing `[]int` HP snapshots per round in `CombatState`. Cleaner but changes the existing `CombatState` type and `resolveCombat` signature — higher blast radius.

### 4. Layout: left 40%, right 40%, log 100% bottom rows

**Decision:** Divide `m.viewportW` as `leftW = viewportW*40/100`, `rightW = viewportW*40/100`, with the remaining width unused (or used as a centre gutter). Log panel takes the bottom `logRows = viewportH/3` rows. Top section `topH = viewportH - logRows` is split between left and right panels. Both panels are rendered side-by-side using `lipgloss.JoinHorizontal`.

**Rationale:** 40/40 gives symmetric panels and room for a gutter. The 1/3 bottom for log gives enough room for ~6–8 lines which covers most fights. All measurements use viewport dimensions so it adapts to terminal size.

### 5. Hero portrait reuses ragdoll; enemy portrait is a centred glyph box

**Decision:**
- Hero: the existing `ragdoll` string slice (6 lines, ~11 chars wide) is centred within `leftW`. Below it, HP bar and stats.
- Enemy: a simple box drawn with `lipgloss.Border(lipgloss.RoundedBorder())` containing the enemy's `Template.Char` rendered at a larger perceived size (3×3 block of the same char, or just the single char centred). Below, HP bar and stats.

**Rationale:** The ragdoll already exists and players recognise it from the inventory screen — reuse is consistent. The enemy portrait box makes the glyph prominent without needing actual image data.

### 6. Speed controls: `[` decreases speed (slower), `]` increases speed (faster)

**Decision:** In the `ScreenCombat` key handler, `[` sets `m.combatSpeed = max(0, m.combatSpeed-1)` and `]` sets `m.combatSpeed = min(2, m.combatSpeed+1)`. The speed label (`Slow` / `Normal` / `Fast`) is shown in the log panel header.

**Rationale:** `[`/`]` are unused in combat and bracket-symmetric. They match the existing `[`/`]` speed controls for world time, making the mental model consistent.

## Risks / Trade-offs

- **[Log parsing for HP is fragile]** If log message format changes, HP display breaks. → Mitigation: centralise the damage log format as a named constant or format string.
- **[CombatTickMsg fires during ScreenCombat; if screen changes mid-combat the tick still fires]** → Handler checks `m.screenMode == ScreenCombat` before advancing; stale ticks are no-ops.
- **[Narrow terminals clip panels]** → Minimum widths enforced: `leftW = max(leftW, 20)`, `rightW = max(rightW, 20)`. Log width is full `viewportW`.

## Open Questions

- Should the log panel scroll during playback, or always show the most recent N lines? (Answer: always show most recent N that fit — consistent with existing behaviour.)
- Should speed be persisted between combats? (Answer: yes — `combatSpeed` lives on model, survives across fights.)
