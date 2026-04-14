## Why

The current HUD is a single line of key bindings that crowds out player-state information (health, armour) and becomes redundant once players learn the controls. Surfacing combat-relevant stats in the HUD and moving key bindings to a dismissible help panel makes the interface informative at a glance and cleaner for experienced players.

## What Changes

- HUD bar replaced with a structured player-stat bar: HP progress bar, armour value, current tile info, clock, time scale, and a single `? help` hint
- New fullscreen key-bindings help panel (`?` toggles it), listing contextual bindings for the active mode
- Help panel renders within the current viewport dimensions — no clipping on small terminals (fixing a regression from a prior sidebar implementation)
- HP rendered as a visual progress bar (e.g. `[████░░░░░░] 15/20`) and armour as a compact label
- **BREAKING** (display only): the bottom key-bar row is removed; bindings move to the help panel

## Capabilities

### New Capabilities
- `hud-player-stats`: HUD bar showing HP progress bar, armour, tile info, clock, speed, and help hint
- `keybindings-panel`: Fullscreen contextual key-bindings help panel toggled by `?`; content varies by active screen mode; viewport-aware rendering

### Modified Capabilities
- `rendering-system`: `buildView` wires the new HUD renderer and routes `?`-panel rendering; key-bar row removed from normal view
- `input-navigation`: `?` key now toggles `showHelpPanel` rather than `showSidebar` (sidebar toggle moves to a different key or is merged); `ScreenCombat` and `ScreenInventory` suppress the help panel toggle
- `combat-system`: player HP and MaxHP are now read from the model's live HP field (so HUD stays accurate between fights); `buildPlayerCombatant` seeds HP from `m.playerHP`

## Impact

- `internal/game/model.go` — add `playerHP int`, `playerMaxHP int`, `showHelpPanel bool` fields
- `internal/game/render.go` — replace `renderKeyBar` with `renderHUD` (stats bar) and add `renderHelpPanel`; update `buildView`
- `internal/game/input.go` — `?` toggles `showHelpPanel`; remove old key-bar hint string construction
- `internal/game/types.go` — no new types; existing `ScreenMode` unchanged
- No new dependencies
