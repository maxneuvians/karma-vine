## Why

The world map currently renders only one view: biome-colored tiles. The underlying tile data already carries `Temperature`, `Elevation`, and `Moisture` values, but players have no way to visualize them directly. Adding togglable map modes makes the world's geography legible and introduces a sub-navigation UI pattern that can be reused elsewhere.

## What Changes

- Add a `MapMode` type (`Default`, `Temperature`, `Elevation`, `Political`) to the model.
- When a non-default mode is active, override the tile `Char`/`Color` rendered on the world map with mode-specific visuals (gradient colors or contour characters).
- Add a map-mode picker panel — a small overlay menu toggled by a key (e.g. `m`) — that lets the player cycle through or select map modes. The panel follows the same open/close pattern as the `?` sidebar.
- Define the sub-navigation pattern: arrow keys navigate items within the open panel; `enter` or the same toggle key confirms and closes; `esc` cancels without changing mode.

## Capabilities

### New Capabilities

- `world-map-modes`: Defines the `MapMode` enum, per-mode tile rendering overrides, and the map-mode picker panel including its sub-navigation interaction model.

### Modified Capabilities

- `rendering-system`: World map render path must branch on `MapMode` to apply per-mode color and character overrides.
- `input-navigation`: Key handler must route `m` to toggle the mode picker, and while the picker is open, route `up`/`down` and `enter`/`esc` to panel navigation instead of player movement.

## Impact

- `internal/game/types.go` — add `MapMode` type and constants.
- `internal/game/model.go` — add `mapMode MapMode` and `showMapPicker bool` fields; add `mapPickerCursor int` for sub-navigation state.
- `internal/game/render.go` — extend `buildView` / world-map render loop to call mode-specific tile override; add `renderMapPicker` function.
- `internal/game/input.go` — handle `m` key; handle sub-navigation keys when `showMapPicker` is true.
- No new dependencies; no breaking changes to existing key bindings or rendering contracts.
