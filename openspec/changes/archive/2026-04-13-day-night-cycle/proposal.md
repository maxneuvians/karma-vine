## Why

The world explorer currently lacks any concept of time, making every moment feel identical regardless of how long the player has been exploring. A day/night cycle with automatic color dimming and local fire-lit illumination adds atmosphere, pacing, and meaningful environmental storytelling without requiring new mechanics.

## What Changes

- Introduce a continuous time model: one 24-hour in-game day elapses every 30 real-world seconds
- Replace the manual `n` night-mode toggle with automatic color dimming derived from the current time-of-day
- Add a `hasFire` flag to local map cells; fire cells emit a radius of illumination that cuts through the darkness at night
- The HUD shows the current in-game time (e.g., `06:30`)
- Players can press `[` / `]` to slow/speed up time (default 1×, max 10×) so a full day can pass in as little as 3 seconds
- **BREAKING**: Remove `nightMode bool` toggle and the `n` key binding

## Capabilities

### New Capabilities
- `day-night-cycle`: time-of-day model, automatic color dimming based on sun position, fire illumination radius on local maps, time speed control, HUD clock display

### Modified Capabilities
- `rendering-system`: color rendering now factors in time-of-day dim multiplier and per-cell fire illumination instead of a global boolean flag
- `input-navigation`: replace `n` toggle with `[`/`]` time-speed keys; remove night toggle case

## Impact

- `model.go`: add `timeOfDay float64` (0.0–1.0, 0=midnight, 0.5=noon) and `timeScale int` fields; remove `nightMode bool`
- `types.go`: add `HasFire bool` to `Ground`; add fire cell variants to biome tables
- `local.go`: generate fire spots (campfires, torches) in biome content; populate `HasFire` during generation
- `render.go`: replace night-mode dim logic with time-based `dimFactor`; add local illumination pass
- `model.go` `Update()`: advance `timeOfDay` each tick using `timeScale`
- `input.go`: remove `n` case; add `[`/`]` cases
