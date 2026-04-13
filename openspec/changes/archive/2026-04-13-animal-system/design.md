## Context

BubbleTea provides `tea.Every(d, fn)` for recurring ticks. The tick fires a `TickMsg` into `Update()` every 500 ms. The model is immutable in BubbleTea's Elm pattern: `Update` returns a new model (or the same pointer for mutations) and a follow-up command to re-schedule the next tick. Animals live in `localMap.Animals` (a `[]*Animal` slice); each tick iterates the slice and updates `X`/`Y` in-place (pointer semantics make this safe since the map is already heap-allocated).

## Goals / Non-Goals

**Goals:**
- Tick only fires when `mode == ModeLocal` and `localMap != nil`
- Random direction chosen from `{±1, 0} × {±1, 0}` excluding `(0, 0)` — 8-directional movement
- Flee direction is the unit step (dx, dy) that maximises `distance(newPos, playerPos)` among the 8 candidates
- Bounds: clamp to `[0, 41] × [0, 17]`; if a move would land on a blocking object cell, skip that direction and try the next best; if all blocked, stay put
- The re-schedule command (`tea.Every(500ms, …)`) is returned from `Update` alongside the updated model

**Non-Goals:**
- Pathfinding (A* or similar) — brief specifies random walk; flee is a one-step greedy move
- Animal spawning / despawning during a session
- Animal interactions with each other

## Decisions

**`tea.Every` re-scheduled each tick** — BubbleTea's canonical pattern for repeating events is to return a new `tea.Every` command from `Update`. This is safe, predictable, and avoids goroutine leaks. Alternative: a goroutine with `time.Sleep` sending on a channel — rejected because it bypasses BubbleTea's message loop.

**8-directional movement** — Matches typical rogue-like feel and is consistent with the brief's "move one step randomly". Alternative: 4-directional — removes diagonal movement, making flee look unnatural.

**In-place mutation of `*Animal`** — Since `LocalMap` is heap-allocated and all animals are pointers, updating `X`/`Y` directly is safe and avoids copying the entire animal slice each tick. The BubbleTea model itself still gets a new `Model` value (copy), but the `LocalMap` pointer inside it points to the same structure.

## Risks / Trade-offs

- **Concurrent tick and player move** → BubbleTea processes messages sequentially on one goroutine; no concurrency issue.
- **Animals accumulate in cache indefinitely** → Same as local map cache — a v2 concern.
- **Fleeing animal gets cornered** → If all 8 moves are blocked or out-of-bounds, the animal stays put. This is correct and intentional.
