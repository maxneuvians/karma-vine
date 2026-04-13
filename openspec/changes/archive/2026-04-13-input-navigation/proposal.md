## Why

The explorer has no way to play without input handling. The player must move on both the world map and the local map, descend from the world tier to the local tier, and ascend back. This change wires all key bindings from the brief into `Update()` and implements the tier-transition logic that switches `mode`, loads the correct map, and places the player appropriately.

## What Changes

- Handle arrow keys and WASD for movement in both `ModeWorld` and `ModeLocal`
- In `ModeWorld`, movement updates `worldPos`
- In `ModeLocal`, movement updates `playerPos`; blocked cells (blocking objects) prevent movement
- `Enter` / `>` in `ModeWorld`: descend — call `LocalMapFor`, set `mode = ModeLocal`, place player at map centre `{21, 9}`
- `Escape` / `<` in `ModeLocal`: ascend — set `mode = ModeWorld`, `localMap = nil` is NOT cleared (keeps cache intact)
- `q` / `ctrl+c`: quit

## Capabilities

### New Capabilities
- `input-navigation`: movement key bindings, world-tier and local-tier navigation, descend/ascend tier transitions, collision with blocking objects

### Modified Capabilities
<!-- none — the stub key handler from project-scaffold handled only quit; this replaces it -->

## Impact

- Replaces the stub key handler in the `Update` switch
- No new dependencies or struct fields required
