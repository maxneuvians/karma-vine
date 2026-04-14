## Context

The game loop is driven by a BubbleTea `TickMsg` dispatched every 500 ms via `tea.Every`. Each tick advances `timeOfDay`, moves animals in `ModeLocal`, and reschedules itself. Player movement is handled through `tea.KeyPressMsg` routed to `handleKey`. There is currently no mechanism to freeze world simulation while the player interacts with menus.

The model already has a `screenMode ScreenMode` field that controls overlay routing (normal vs. fullscreen inventory). Pausing is a separate concern — the player can be paused while in any `screenMode` or `mode`.

## Goals / Non-Goals

**Goals:**
- Space key toggles a `paused bool` on the model from any game state
- While paused: time stops advancing, animals stop moving, player movement keys are no-ops
- While paused: inventory navigation, equipment slots, `i`/`esc` screen-mode changes, `q`/`ctrl+c` quit all remain functional
- A "PAUSED" indicator is rendered in the HUD when paused
- Tick loop continues running (keeps rescheduling) so unpausing resumes immediately without lag

**Non-Goals:**
- Pausing network or save I/O (not present)
- Separate pause menu or pause-screen overlay (indicator only, no new screen mode)
- Pausing dungeon-specific ticks separately from world ticks (all simulation freezes uniformly)

## Decisions

### 1. Keep ticking, skip effects

**Decision:** The `TickMsg` handler always reschedules the next tick via `tickCmd()`, but when `m.paused == true` it returns early before advancing `timeOfDay` or calling `moveAnimals`.

**Rationale:** Stopping tick scheduling on pause would require re-initialising the loop on unpause (which `Init()` does not re-run). Keeping the loop alive with a no-op tick is simpler and matches BubbleTea's design — the tick cadence is fixed and independent of game state.

**Alternative considered:** Cancel the tick command on pause (return `nil` cmd) and re-issue `tickCmd()` on unpause. Rejected because the `Update` switch has no "just unpaused" event to hook into without adding a new message type.

### 2. Single `paused bool` field, not a new `ScreenMode`

**Decision:** Add `paused bool` to `Model` directly, independent of `screenMode`.

**Rationale:** `screenMode` routes rendering and input to distinct full-screen views. Pause is orthogonal — it affects simulation but not which screen is shown. Conflating the two would force the player out of the inventory screen when unpausing, and would require adding pause-awareness to every screen renderer.

**Alternative considered:** `ScreenPaused` as a third `ScreenMode` constant. Rejected because it breaks the inventory-while-paused use case.

### 3. Guard movement keys in `handleKey`, not in `applyDelta`

**Decision:** Check `m.paused` at the top of the movement key branches in `handleKey`. `applyDelta` itself is not modified.

**Rationale:** `applyDelta` is a pure coordinate-transformation helper; injecting pause logic there spreads the concern. All user-facing key routing lives in `handleKey`, making it the natural place for the pause guard. Mouse-click movement (`handleMouseClick`) must also be guarded.

### 4. Pause indicator in the HUD row

**Decision:** Append `" ⏸ PAUSED"` (or equivalent ASCII `[PAUSED]`) to the existing HUD status bar when `m.paused == true`.

**Rationale:** The HUD is always visible in `ScreenNormal`. In `ScreenInventory` the fullscreen inventory replaces the HUD, so no indicator is needed there (the inventory being open while paused is self-evident from the lack of world movement).

## Risks / Trade-offs

- **Tick drift while paused**: Ticks continue to fire every 500 ms but are no-ops. This is negligible overhead. → No mitigation needed.
- **Space key conflicts**: Space is currently unbound. If a future feature binds space (e.g., action/interact), this will conflict. → Document the binding in the spec.
- **Test helpers**: Existing tests drive simulation by feeding `TickMsg{}` directly. A paused model receiving `TickMsg{}` must not advance time, which existing tests do not expect. → New tests cover the paused `TickMsg` path; existing tests are unaffected because they construct a fresh `NewModel()` where `paused == false`.
