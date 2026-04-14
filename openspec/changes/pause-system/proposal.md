## Why

The game has no way to pause world simulation, forcing players to manage their inventory and equipment under time pressure even during menus. Adding a pause state lets players interact with the inventory and equipment screens without the world advancing — time, animal movement, and day/night progression all freeze until they resume.

## What Changes

- Space key toggles a paused state on the model
- While paused, all world-simulation tick commands are suppressed (time progression, animal movement, day/night cycle)
- Player movement input is also blocked while paused
- Inventory navigation, equipment interaction, and other menu actions remain fully functional while paused
- A visible pause indicator is shown in the HUD or as an overlay when the game is paused

## Capabilities

### New Capabilities
- `pause-system`: Defines the `paused bool` field on `Model`, space key toggle, tick suppression logic, and pause indicator in the UI

### Modified Capabilities
- `input-navigation`: Space key handler added; movement keys (`up`/`down`/`left`/`right`/`w`/`a`/`s`/`d`) are no-ops when `m.paused == true`
- `day-night-cycle`: Tick command is not scheduled (or is a no-op) when `m.paused == true`
- `animal-system`: Animal movement tick is suppressed when `m.paused == true`
- `rendering-system`: HUD or overlay shows a "PAUSED" indicator when `m.paused == true`

## Impact

- `internal/game/model.go` — add `paused bool` field
- `internal/game/input.go` — space key handler; guard movement keys behind `!m.paused`
- `internal/game/model.go` (Update) — suppress tick scheduling when paused
- `internal/game/render.go` — render pause indicator
- `internal/game/game_test.go` / `input_test.go` / `render_test.go` — new test coverage for pause behaviour
