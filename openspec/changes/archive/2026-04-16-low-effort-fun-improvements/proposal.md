## Why

The game has several low-hanging friction points that make it feel unrewarding and unpolished: loot drops are meaningless (weapons have 0 stat bonuses), HP can never be recovered, losing abruptly quits the process, enemy portraits render as generic blobs, and early dungeon floors feel empty. These five targeted fixes dramatically improve perceived quality with minimal implementation risk.

## What Changes

- **Weapon stat values**: Rusty Dagger gains +1 DMG bonus; Short Sword gains +2 DMG bonus. Starting clothing (Cloth Tunic, Cloth Pants, Leather Boots) each gain +1 armour bonus.
- **Campfire resting**: Player can press `r` while standing on a fire cell to restore 5 HP (up to max HP). A 30-second cooldown (in game-ticks) prevents abuse.
- **Death screen**: On combat defeat, instead of calling `tea.Quit` immediately, the game shows a death screen with the name of the killer and prompts the player to press `r` to restart or `q` to quit.
- **Enemy portrait fix**: Portrait archetype selection currently fails for lowercase enemy chars (all dungeon enemies). Switch to a name-based lookup so Goblin/Bandit/Giant/Troll → humanoid; Wraith/Crustacean/Rat → beast; Sand Wraith/Ice Wraith → undead.
- **Minimum dungeon enemy count**: Floor enemy count is clamped to `max(3, depth)` so the first few floors have at least 3 enemies instead of a single one.

## Capabilities

### New Capabilities

- `campfire-resting`: Player can rest at a fire cell to restore HP, with a cooldown to prevent spam.
- `death-screen`: On defeat, display a death screen instead of quitting immediately, allowing restart or quit.

### Modified Capabilities

- `dungeon-enemy-system`: Minimum enemy count per floor changes from `depth` to `max(3, depth)`.
- `combat-portraits`: Portrait archetype selection switches from char-byte case check to name-based lookup.
- `equipment-system`: Item definitions updated with real stat values for Rusty Dagger, Short Sword, and starting clothing items.

## Impact

- `game/items.go` (or equivalent item definition file): stat value changes for 5 items.
- `game/model.go` / `game/update.go`: campfire rest key handler; death screen state and transitions.
- `game/combat.go` or `game/render.go`: death screen rendering; portrait name-lookup fix.
- `game/dungeon.go`: enemy count clamping.
- No external dependencies added; no API changes; no breaking changes to saved state (game has no persistence).
